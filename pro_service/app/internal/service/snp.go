package service

import (
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type SNPService struct {
	repo repository.SNP
}

func NewSNPService(repo repository.SNP) *SNPService {
	return &SNPService{repo: repo}
}

func (s *SNPService) Get(req *proto.GetSNPRequest) (snp []*proto.SNP, err error) {
	snp, err = s.repo.Get(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get snp. error: %w", err)
	}
	return snp, nil
}

func (s *SNPService) Create(snp *proto.CreateSNPRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Create(snp)
	if err != nil {
		return nil, fmt.Errorf("failed to create snp. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *SNPService) Update(snp *proto.UpdateSNPRequest) error {
	if err := s.repo.Update(snp); err != nil {
		return fmt.Errorf("failed to update snp. error: %w", err)
	}
	return nil
}

func (s *SNPService) Delete(snp *proto.DeleteSNPRequest) error {
	if err := s.repo.Delete(snp); err != nil {
		return fmt.Errorf("failed to delete snp. error: %w", err)
	}
	return nil
}
