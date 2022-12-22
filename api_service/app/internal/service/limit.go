package service

import (
	"context"
	"fmt"

	"github.com/Alexander272/sealur/api_service/internal/models"
	"github.com/Alexander272/sealur/api_service/internal/repository"
)

type LimitService struct {
	repo repository.Limit
}

func NewLimitService(repo repository.Limit) *LimitService {
	return &LimitService{
		repo: repo,
	}
}

func (s *LimitService) Create(ctx context.Context, clientIP string) error {
	if err := s.repo.Create(ctx, clientIP); err != nil {
		return fmt.Errorf("failed to create limit. error: %w", err)
	}
	return nil
}

func (s *LimitService) Get(ctx context.Context, clientIP string) (data models.LimitData, err error) {
	data, err = s.repo.Get(ctx, clientIP)
	if err != nil {
		return data, fmt.Errorf("failed to get limit. error: %w", err)
	}

	return data, nil
}

func (s *LimitService) AddAttempt(ctx context.Context, clientIP string) error {
	if err := s.repo.AddAttempt(ctx, clientIP); err != nil {
		return fmt.Errorf("failed to add attempt to limit. error: %w", err)
	}
	return nil
}

func (s *LimitService) Remove(ctx context.Context, clientIP string) error {
	if err := s.repo.Remove(ctx, clientIP); err != nil {
		return fmt.Errorf("failed to remove limit. error: %w", err)
	}
	return nil
}
