package service

import (
	"fmt"
	"sync"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type SizeService struct {
	repo repository.Size
}

func NewSizeService(repo repository.Size) *SizeService {
	return &SizeService{repo: repo}
}

func (s *SizeService) Get(req *pro_api.GetSizesRequest) (sizes []*pro_api.Size, dn []*pro_api.Dn, err error) {
	sizes, err = s.repo.Get(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sizes. error: %w", err)
	}

	dn = make([]*pro_api.Dn, 0, len(sizes))
	for _, s := range sizes {
		if len(dn) > 0 {
			if dn[len(dn)-1].Dn != s.Dn {
				dn = append(dn, &pro_api.Dn{Dn: s.Dn})
			}
		} else {
			dn = append(dn, &pro_api.Dn{Dn: s.Dn})
		}
	}

	return sizes, dn, nil
}

func (s *SizeService) GetAll(req *pro_api.GetSizesRequest) (sizes []*pro_api.Size, dn []*pro_api.Dn, err error) {
	sizes, err = s.repo.GetAll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sizes. error: %w", err)
	}

	dn = make([]*pro_api.Dn, 0, len(sizes))

	return sizes, dn, nil
}

func (s *SizeService) Create(size *pro_api.CreateSizeRequest) (*pro_api.IdResponse, error) {
	id, err := s.repo.Create(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create size. error: %w", err)
	}
	return &pro_api.IdResponse{Id: id}, nil
}

func (s *SizeService) CreateLimit(wg *sync.WaitGroup, limit chan struct{}, size *pro_api.CreateSizeRequest) {
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

func (s *SizeService) CreateMany(sizes *pro_api.CreateSizesRequest) error {
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

func (s *SizeService) Update(size *pro_api.UpdateSizeRequest) error {
	if err := s.repo.Update(size); err != nil {
		return fmt.Errorf("failed to update query. error: %w", err)
	}
	return nil
}

func (s *SizeService) Delete(size *pro_api.DeleteSizeRequest) error {
	if err := s.repo.Delete(size); err != nil {
		return fmt.Errorf("failed to delete query. error: %w", err)
	}
	return nil
}

func (s *SizeService) DeleteAll(size *pro_api.DeleteAllSizeRequest) error {
	if err := s.repo.DeleteAll(size); err != nil {
		return fmt.Errorf("failed to delete query. error: %w", err)
	}
	return nil
}
