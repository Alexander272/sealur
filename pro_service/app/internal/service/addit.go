package service

import (
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type AdditService struct {
	repo repository.Addit
}

func NewAdditService(repo repository.Addit) *AdditService {
	return &AdditService{repo: repo}
}

func (s *AdditService) GetAll() (addit []*proto.Additional, err error) {
	data, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get additional. error: %w", err)
	}

	for _, d := range data {
		var mats []*proto.AddMaterials
		var mods []*proto.AddMod
		var temps []*proto.AddTemperature
		var mouns []*proto.AddMoun
		var graps []*proto.AddGrap
		var fils []*proto.AddFiller

		tmp := strings.Split(d.Materials, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			mats = append(mats, &proto.AddMaterials{
				Short: parts[0],
				Title: parts[1],
			})
		}
		tmp = strings.Split(d.Mod, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			mods = append(mods, &proto.AddMod{
				Id:          parts[0],
				Short:       parts[1],
				Title:       parts[2],
				Description: parts[3],
			})
		}
		tmp = strings.Split(d.Temperature, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			temps = append(temps, &proto.AddTemperature{
				Id:    parts[0],
				Title: parts[1],
			})
		}
		tmp = strings.Split(d.Mounting, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			mouns = append(mouns, &proto.AddMoun{
				Id:    parts[0],
				Title: parts[1],
			})
		}
		tmp = strings.Split(d.Graphite, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			graps = append(graps, &proto.AddGrap{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
			})
		}
		tmp = strings.Split(d.Fillers, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			fils = append(fils, &proto.AddFiller{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
			})
		}

		addit = append(addit, &proto.Additional{
			Id:          d.Id,
			Materials:   mats,
			Mod:         mods,
			Temperature: temps,
			Mounting:    mouns,
			Graphite:    graps,
			Fillers:     fils,
		})
	}

	return addit, nil
}

func (s *AdditService) Create(addit *proto.CreateAddRequest) (*proto.SuccessResponse, error) {
	err := s.repo.Create(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to create additional. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMat(addit *proto.UpdateAddMatRequest) (*proto.SuccessResponse, error) {
	err := s.repo.UpdateMat(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to update materials. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMod(addit *proto.UpdateAddModRequest) (*proto.SuccessResponse, error) {
	err := s.repo.UpdateMod(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to update mod. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateTemp(addit *proto.UpdateAddTemRequest) (*proto.SuccessResponse, error) {
	err := s.repo.UpdateTemp(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to update temperature. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMoun(addit *proto.UpdateAddMounRequest) (*proto.SuccessResponse, error) {
	err := s.repo.UpdateMoun(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to update mounting. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateGrap(addit *proto.UpdateAddGrapRequest) (*proto.SuccessResponse, error) {
	if err := s.repo.UpdateGrap(addit); err != nil {
		return nil, fmt.Errorf("failed to update graphite. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateFillers(addit *proto.UpdateAddFillersRequest) (*proto.SuccessResponse, error) {
	if err := s.repo.UpdateFillers(addit); err != nil {
		return nil, fmt.Errorf("failed to update fillers. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}
