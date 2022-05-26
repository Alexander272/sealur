package service

import (
	"fmt"
	"sync"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
)

type SizeService struct {
	repo repository.Size
}

func NewSizeService(repo repository.Size) *SizeService {
	return &SizeService{repo: repo}
}

func (s *SizeService) Get(req *proto.GetSizesRequest) (sizes []*proto.Size, dn []*proto.Dn, err error) {
	sizes, err = s.repo.Get(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sizes. error: %w", err)
	}

	dn = make([]*proto.Dn, 0, len(sizes))
	for _, s := range sizes {
		if len(dn) > 0 {
			if dn[len(dn)-1].Dn != s.Dn {
				dn = append(dn, &proto.Dn{Dn: s.Dn})
			}
		} else {
			dn = append(dn, &proto.Dn{Dn: s.Dn})
		}
	}

	return sizes, dn, nil
}

func (s *SizeService) GetAll(req *proto.GetSizesRequest) (sizes []*proto.Size, dn []*proto.Dn, err error) {
	sizes, err = s.repo.GetAll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sizes. error: %w", err)
	}

	dn = make([]*proto.Dn, 0, len(sizes))

	return sizes, dn, nil
}

func (s *SizeService) Create(size *proto.CreateSizeRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Create(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create size. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *SizeService) CreateLimit(wg *sync.WaitGroup, limit chan struct{}, size *proto.CreateSizeRequest) {
	defer wg.Done()

	var w sync.WaitGroup
	w.Add(1)
	_, err := s.Create(size)
	if err != nil {
		logger.Error(err)
	}
	w.Done()
	w.Wait()

	<-limit
}

func (s *SizeService) CreateMany(sizes *proto.CreateSizesRequest) error {
	var wg sync.WaitGroup
	limit := make(chan struct{}, 50)
	for _, size := range sizes.Sizes {
		wg.Add(1)
		limit <- struct{}{}
		go s.CreateLimit(&wg, limit, size)
	}

	wg.Wait()
	return nil
}

func (s *SizeService) Update(size *proto.UpdateSizeRequest) error {
	if err := s.repo.Update(size); err != nil {
		return fmt.Errorf("failed to update query. error: %w", err)
	}
	return nil
}

func (s *SizeService) Delete(size *proto.DeleteSizeRequest) error {
	if err := s.repo.Delete(size); err != nil {
		return fmt.Errorf("failed to delete query. error: %w", err)
	}
	return nil
}

func (s *SizeService) DeleteAll(size *proto.DeleteAllSizeRequest) error {
	if err := s.repo.DeleteAll(size); err != nil {
		return fmt.Errorf("failed to delete query. error: %w", err)
	}
	return nil
}
