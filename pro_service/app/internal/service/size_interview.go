package service

import (
	"fmt"
	"sync"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
)

type SizeIntService struct {
	repo repository.SizeInt
}

func NewSizeIntService(repo repository.SizeInt) *SizeIntService {
	return &SizeIntService{repo: repo}
}

func (s *SizeIntService) Get(req *proto.GetSizesIntRequest) (sizes []*proto.SizeInt, dn []*proto.Dn, err error) {
	data, err := s.repo.Get(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sizes interview. error: %w", err)
	}

	for _, d := range data {
		s := &proto.SizeInt{
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

	dn = make([]*proto.Dn, 0, len(sizes))
	for _, s := range sizes {
		if len(dn) > 0 {
			if dn[len(dn)-1].Dn != s.Dy {
				dn = append(dn, &proto.Dn{Dn: s.Dy})
			}
		} else {
			dn = append(dn, &proto.Dn{Dn: s.Dy})
		}
	}

	return sizes, dn, nil
}

func (s *SizeIntService) GetAll(req *proto.GetAllSizeIntRequest) (sizes []*proto.SizeInt, dn []*proto.Dn, err error) {
	data, err := s.repo.GetAll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get all sizes interview. error: %w", err)
	}

	for _, d := range data {
		s := &proto.SizeInt{
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

	dn = make([]*proto.Dn, 0, len(sizes))
	for _, s := range sizes {
		if len(dn) > 0 {
			if dn[len(dn)-1].Dn != s.Dy {
				dn = append(dn, &proto.Dn{Dn: s.Dy})
			}
		} else {
			dn = append(dn, &proto.Dn{Dn: s.Dy})
		}
	}

	return sizes, dn, nil
}

func (s *SizeIntService) Create(size *proto.CreateSizeIntRequest) (*proto.IdResponse, error) {
	id, err := s.repo.Create(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create size interview. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *SizeIntService) CreateLimit(wg *sync.WaitGroup, limit chan struct{}, size *proto.CreateSizeIntRequest) {
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

func (s *SizeIntService) CreateMany(sizes *proto.CreateSizesIntRequest) error {
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

func (s *SizeIntService) Update(size *proto.UpdateSizeIntRequest) error {
	if err := s.repo.Update(size); err != nil {
		return fmt.Errorf("failed to update size interview. error: %w", err)
	}
	return nil
}

func (s *SizeIntService) Delete(size *proto.DeleteSizeIntRequest) error {
	if err := s.repo.Delete(size); err != nil {
		return fmt.Errorf("failed to delete size interview. error: %w", err)
	}
	return nil
}

func (s *SizeIntService) DeleteAll(size *proto.DeleteAllSizeIntRequest) error {
	if err := s.repo.DeleteAll(size); err != nil {
		return fmt.Errorf("failed to delete all size interview. error: %w", err)
	}
	return nil
}
