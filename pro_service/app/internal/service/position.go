package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	proto_file "github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto/proto_file"
)

type PositionService struct {
	repo repository.OrderPosition
	file proto_file.FileServiceClient
}

func NewPositionService(repo repository.OrderPosition, file proto_file.FileServiceClient) *PositionService {
	return &PositionService{
		repo: repo,
		file: file,
	}
}

func (s *PositionService) Get(req *proto.GetPositionsRequest) (positions []*proto.OrderPosition, err error) {
	pos, err := s.repo.Get(req)
	if err != nil {
		return positions, fmt.Errorf("failed to get positions. error: %w", err)
	}

	for _, p := range pos {
		position := &proto.OrderPosition{
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

func (s *PositionService) GetCur(req *proto.GetCurPositionsRequest) (positions []*proto.OrderPosition, err error) {
	pos, err := s.repo.GetCur(req)
	if err != nil {
		return positions, fmt.Errorf("failed to get position for current user. error: %w", err)
	}

	for _, p := range pos {
		position := &proto.OrderPosition{
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

func (s *PositionService) Add(position *proto.AddPositionRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Add(position)
	if err != nil {
		return nil, fmt.Errorf("failed to add position. error: %w", err)
	}

	return &proto.IdResponse{Id: id}, nil
}

func (s *PositionService) Copy(position *proto.CopyPositionRequest) (*proto.IdResponse, error) {
	pos := proto.AddPositionRequest{
		OrderId:     position.OrderId,
		Designation: position.Designation,
		Description: position.Description,
		Count:       position.Count,
		Sizes:       position.Sizes,
		Drawing:     position.Drawing,
	}

	if position.Drawing != "" {
		_, err := s.file.Copy(context.Background(), &proto_file.CopyFileRequest{
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

	return &proto.IdResponse{Id: id}, nil
}

func (s *PositionService) Update(position *proto.UpdatePositionRequest) (*proto.IdResponse, error) {
	if err := s.repo.Update(position); err != nil {
		return nil, fmt.Errorf("failed to update position. error: %w", err)
	}
	return &proto.IdResponse{Id: position.Id}, nil
}

func (s *PositionService) Remove(position *proto.RemovePositionRequest) (*proto.IdResponse, error) {
	image, err := s.repo.Remove(position)
	if err != nil {
		return nil, fmt.Errorf("failed to remove position. error: %w", err)
	}

	if image != "" {
		_, err = s.file.Delete(context.Background(), &proto_file.FileDeleteRequest{
			Id:     strings.Split(image, "_")[0],
			Bucket: "pro",
			Group:  position.OrderId,
			Name:   strings.Split(image, "_")[1],
		})
		if err != nil {
			return nil, fmt.Errorf("failed to remove drawing. error: %w", err)
		}
	}

	return &proto.IdResponse{Id: position.Id}, nil
}
