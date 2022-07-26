package service

import (
	"fmt"
	"sync"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type SizeIntService struct {
	repo repository.SizeInt
}

func NewSizeIntService(repo repository.SizeInt) *SizeIntService {
	return &SizeIntService{repo: repo}
}

func (s *SizeIntService) Get(req *pro_api.GetSizesIntRequest) (sizes []*pro_api.SizeInt, dn []*pro_api.Dn, err error) {
	data, err := s.repo.Get(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sizes interview. error: %w", err)
	}

	for _, d := range data {
		s := &pro_api.SizeInt{
			Id:        d.Id,
			Dy:        d.Dy,
			Py:        d.Py,
			D1:        d.D1,
			D2:        d.D2,
			DUp:       d.DUp,
			D:         d.D,
			H1:        d.H1,
			H2:        d.H2,
			Bolt:      d.Bolt,
			CountBolt: d.Count,
		}
		sizes = append(sizes, s)
	}

	dn = make([]*pro_api.Dn, 0, len(sizes))
	for _, s := range sizes {
		if len(dn) > 0 {
			if dn[len(dn)-1].Dn != s.Dy {
				dn = append(dn, &pro_api.Dn{Dn: s.Dy})
			}
		} else {
			dn = append(dn, &pro_api.Dn{Dn: s.Dy})
		}
	}

	return sizes, dn, nil
}

func (s *SizeIntService) GetAll(req *pro_api.GetAllSizeIntRequest) (sizes []*pro_api.SizeInt, dn []*pro_api.Dn, err error) {
	data, err := s.repo.GetAll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get all sizes interview. error: %w", err)
	}

	for _, d := range data {
		s := &pro_api.SizeInt{
			Id:        d.Id,
			Dy:        d.Dy,
			Py:        d.Py,
			D1:        d.D1,
			D2:        d.D2,
			DUp:       d.DUp,
			D:         d.D,
			H1:        d.H1,
			H2:        d.H2,
			Bolt:      d.Bolt,
			CountBolt: d.Count,
			Row:       d.Row,
		}
		sizes = append(sizes, s)
	}

	dn = make([]*pro_api.Dn, 0, len(sizes))
	for _, s := range sizes {
		if len(dn) > 0 {
			if dn[len(dn)-1].Dn != s.Dy {
				dn = append(dn, &pro_api.Dn{Dn: s.Dy})
			}
		} else {
			dn = append(dn, &pro_api.Dn{Dn: s.Dy})
		}
	}

	return sizes, dn, nil
}

func (s *SizeIntService) Create(size *pro_api.CreateSizeIntRequest) (*pro_api.IdResponse, error) {
	id, err := s.repo.Create(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create size interview. error: %w", err)
	}
	return &pro_api.IdResponse{Id: id}, nil
}

func (s *SizeIntService) CreateLimit(wg *sync.WaitGroup, limit chan struct{}, size *pro_api.CreateSizeIntRequest) {
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

func (s *SizeIntService) CreateMany(sizes *pro_api.CreateSizesIntRequest) error {
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

func (s *SizeIntService) Update(size *pro_api.UpdateSizeIntRequest) error {
	if err := s.repo.Update(size); err != nil {
		return fmt.Errorf("failed to update size interview. error: %w", err)
	}
	return nil
}

func (s *SizeIntService) Delete(size *pro_api.DeleteSizeIntRequest) error {
	if err := s.repo.Delete(size); err != nil {
		return fmt.Errorf("failed to delete size interview. error: %w", err)
	}
	return nil
}

func (s *SizeIntService) DeleteAll(size *pro_api.DeleteAllSizeIntRequest) error {
	if err := s.repo.DeleteAll(size); err != nil {
		return fmt.Errorf("failed to delete all size interview. error: %w", err)
	}
	return nil
}
