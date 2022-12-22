package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type PositionService struct {
	repo repository.OrderPosition
	file file_api.FileServiceClient
}

func NewPositionService(repo repository.OrderPosition, file file_api.FileServiceClient) *PositionService {
	return &PositionService{
		repo: repo,
		file: file,
	}
}

func (s *PositionService) Get(req *pro_api.GetPositionsRequest) (positions []*pro_api.OrderPosition, err error) {
	pos, err := s.repo.Get(req)
	if err != nil {
		return positions, fmt.Errorf("failed to get positions. error: %w", err)
	}

	for _, p := range pos {
		position := &pro_api.OrderPosition{
			Id:          p.Id,
			Designation: p.Designation,
			Description: p.Descriprion,
			Count:       p.Count,
			Sizes:       p.Sizes,
			Drawing:     p.Drawing,
			OrderId:     p.OrderId,
		}
		positions = append(positions, position)
	}

	return positions, nil
}

func (s *PositionService) GetCur(req *pro_api.GetCurPositionsRequest) (positions []*pro_api.OrderPosition, err error) {
	pos, err := s.repo.GetCur(req)
	if err != nil {
		return positions, fmt.Errorf("failed to get position for current user. error: %w", err)
	}

	for _, p := range pos {
		position := &pro_api.OrderPosition{
			Id:          p.Id,
			Designation: p.Designation,
			Description: p.Descriprion,
			Count:       p.Count,
			Sizes:       p.Sizes,
			Drawing:     p.Drawing,
			OrderId:     p.OrderId,
		}
		positions = append(positions, position)
	}

	return positions, nil
}

func (s *PositionService) Add(position *pro_api.AddPositionRequest) (*pro_api.IdResponse, error) {
	id, err := s.repo.Add(position)
	if err != nil {
		return nil, fmt.Errorf("failed to add position. error: %w", err)
	}

	return &pro_api.IdResponse{Id: id}, nil
}

func (s *PositionService) Copy(position *pro_api.CopyPositionRequest) (*pro_api.IdResponse, error) {
	pos := pro_api.AddPositionRequest{
		OrderId:     position.OrderId,
		Designation: position.Designation,
		Description: position.Description,
		Count:       position.Count,
		Sizes:       position.Sizes,
		Drawing:     position.Drawing,
	}

	if position.Drawing != "" {
		_, err := s.file.Copy(context.Background(), &file_api.CopyFileRequest{
			Id:       position.Drawing,
			Bucket:   "pro",
			Group:    position.OldOrderId,
			NewGroup: position.OrderId,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to copy drawing. error: %w", err)
		}
	}

	id, err := s.repo.Add(&pos)
	if err != nil {
		return nil, fmt.Errorf("failed to add position. error: %w", err)
	}

	return &pro_api.IdResponse{Id: id}, nil
}

func (s *PositionService) Update(position *pro_api.UpdatePositionRequest) (*pro_api.IdResponse, error) {
	if err := s.repo.Update(position); err != nil {
		return nil, fmt.Errorf("failed to update position. error: %w", err)
	}
	return &pro_api.IdResponse{Id: position.Id}, nil
}

func (s *PositionService) Remove(position *pro_api.RemovePositionRequest) (*pro_api.IdResponse, error) {
	image, err := s.repo.Remove(position)
	if err != nil {
		return nil, fmt.Errorf("failed to remove position. error: %w", err)
	}

	if image != "" {
		_, err = s.file.Delete(context.Background(), &file_api.FileDeleteRequest{
			Id:     strings.Split(image, "_")[0],
			Bucket: "pro",
			Group:  position.OrderId,
			Name:   strings.Split(image, "_")[1],
		})
		if err != nil {
			return nil, fmt.Errorf("failed to remove drawing. error: %w", err)
		}
	}

	return &pro_api.IdResponse{Id: position.Id}, nil
}
