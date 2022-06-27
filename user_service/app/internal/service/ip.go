package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/user_service/internal/repo"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
)

type IpService struct {
	repo repo.IP
}

func NewIpService(repo repo.IP) *IpService {
	return &IpService{repo: repo}
}

func (s *IpService) Add(ctx context.Context, ip *proto_user.AddIpRequest) error {
	if err := s.repo.Add(ctx, ip); err != nil {
		return fmt.Errorf("failed to add ip. error: %w", err)
	}
	return nil
}
