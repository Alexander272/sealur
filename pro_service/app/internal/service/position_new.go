package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/google/uuid"
)

type PositionServiceNew struct {
	repo    repository.Position
	snp     PositionSnp
	fileApi file_api.FileServiceClient
}

func NewPositionService_New(repo repository.Position, snp PositionSnp, fileApi file_api.FileServiceClient) *PositionServiceNew {
	return &PositionServiceNew{
		repo:    repo,
		snp:     snp,
		fileApi: fileApi,
	}
}

func (s *PositionServiceNew) Get(ctx context.Context, orderId string) (positions []*position_model.FullPosition, err error) {
	snpPosition, err := s.snp.Get(ctx, orderId)
	if err != nil {
		return nil, err
	}

	positions = append(positions, snpPosition...)

	return positions, nil
}

func (s *PositionServiceNew) GetAll(ctx context.Context, orderId string) ([]*position_model.OrderPosition, error) {
	positions, err := s.repo.Get(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions. error: %w", err)
	}
	return positions, nil
}

func (s *PositionServiceNew) GetFull(ctx context.Context, orderId string) ([]*position_model.OrderPosition, error) {
	positions, err := s.repo.Get(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions. error: %w", err)
	}

	snpId := make([]string, 0, len(positions))
	snpIndex := make(map[string]int, 0)

	for i, op := range positions {
		if op.Type == position_model.PositionType_Snp {
			snpId = append(snpId, op.Id)
			snpIndex[op.Id] = i
		}
	}

	snpPositions, err := s.snp.GetFull(ctx, snpId)
	if err != nil {
		return nil, err
	}

	for _, ops := range snpPositions {
		index := snpIndex[ops.Main.PositionId]
		positions[index].SnpData = ops
	}

	return positions, nil
}

func (s *PositionServiceNew) GetAnalytics(ctx context.Context, req *order_api.GetOrderAnalytics) (*order_api.Analytics, error) {
	data, err := s.repo.GetAnalytics(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get position analytics. error: %w", err)
	}

	return data, nil
}

func (s *PositionServiceNew) Create(ctx context.Context, position *position_model.FullPosition) (string, error) {
	candidate, err := s.repo.GetByTitle(ctx, position.Title, position.OrderId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to get position by title. error: %w", err)
	}

	if candidate != "" {
		return "", fmt.Errorf("position exists")
	}

	id, err := s.repo.Create(ctx, position)
	if err != nil {
		return "", fmt.Errorf("failed to create position. error: %w", err)
	}

	position.Id = id
	if position.Type == position_model.PositionType_Snp {
		if err := s.snp.Create(ctx, position); err != nil {
			return "", err
		}
	}
	return id, nil
}

func (s *PositionServiceNew) CreateSeveral(ctx context.Context, positions []*position_model.FullPosition, orderId string) error {
	var positionSnp []*position_model.FullPosition

	for _, p := range positions {
		id := uuid.New()
		p.Id = id.String()
		p.OrderId = orderId

		if p.Type == position_model.PositionType_Snp {
			positionSnp = append(positionSnp, p)
		}
	}

	if err := s.repo.CreateSeveral(ctx, positions); err != nil {
		return fmt.Errorf("failed to create several positions. error: %w", err)
	}

	if err := s.snp.CreateSeveral(ctx, positionSnp); err != nil {
		return err
	}

	return nil
}

func (s *PositionServiceNew) Update(ctx context.Context, position *position_model.FullPosition) error {
	if err := s.repo.Update(ctx, position); err != nil {
		return fmt.Errorf("failed to update position. error: %w", err)
	}

	if position.Type == position_model.PositionType_Snp {
		if err := s.snp.Update(ctx, position); err != nil {
			return err
		}
	}

	return nil
}

func (s *PositionServiceNew) Copy(ctx context.Context, position *position_api.CopyPosition) (string, error) {
	curPosition, err := s.repo.GetById(ctx, position.Id)
	if err != nil {
		return "", fmt.Errorf("failed to get position. error: %w", err)
	}

	candidate, err := s.repo.GetByTitle(ctx, curPosition.Title, position.OrderId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to get position by title. error: %w", err)
	}

	if candidate != "" {
		return "", fmt.Errorf("position exists")
	}

	curPosition.Count = position.Count
	curPosition.OrderId = position.OrderId
	if position.Amount != "" {
		curPosition.Amount = position.Amount
	}
	id, err := s.repo.Create(ctx, curPosition)
	if err != nil {
		return "", fmt.Errorf("failed to create position. error: %w", err)
	}
	// id, err := s.repo.Copy(ctx, position)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to copy position. error: %w", err)
	// }

	// По сути я могу возвращать drawing и из него вырезать id файла, а id заявки (старой и новой) буду принимать с клиента
	drawing, err := s.snp.Copy(ctx, id, position)
	if err != nil {
		return "", err
	}

	if drawing != "" {
		parts := strings.Split(drawing, "/")

		logger.Debug(position.FromOrderId, ", ", position.OrderId)

		_, err := s.fileApi.Copy(context.Background(), &file_api.CopyFileRequest{
			Id:       fmt.Sprintf("%s_%s", parts[len(parts)-2], parts[len(parts)-1]),
			Bucket:   "pro",
			Group:    position.FromOrderId,
			NewGroup: position.OrderId,
		})
		if err != nil {
			return "", fmt.Errorf("failed to copy drawing. error: %w", err)
		}
	}

	return drawing, nil
}

func (s *PositionServiceNew) Delete(ctx context.Context, positionId string) error {
	if err := s.repo.Delete(ctx, positionId); err != nil {
		return fmt.Errorf("failed to delete position. error: %w", err)
	}

	//? оно удаляется само (каскадом)
	// можно было бы тут получать чертеж, но оно ничего не удалит
	if err := s.snp.Delete(ctx, positionId); err != nil {
		return err
	}

	//TODO удалять чертеж

	return nil
}
