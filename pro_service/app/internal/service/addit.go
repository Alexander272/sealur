package service

import (
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
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
		var coat []*proto.AddCoating
		var constr []*proto.AddConstruction
		var obt []*proto.AddObturator

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
				Short:       parts[2],
				Title:       parts[1],
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
		tmp = strings.Split(d.Coating, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			coat = append(coat, &proto.AddCoating{
				Id:          parts[0],
				Short:       parts[1],
				Title:       parts[2],
				Description: parts[3],
			})
		}
		tmp = strings.Split(d.Construction, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			constr = append(constr, &proto.AddConstruction{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
				ForDescr:    parts[3],
			})
		}
		tmp = strings.Split(d.Obturator, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			obt = append(obt, &proto.AddObturator{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
			})
		}

		addit = append(addit, &proto.Additional{
			Id:           d.Id,
			Materials:    mats,
			Mod:          mods,
			Temperature:  temps,
			Mounting:     mouns,
			Graphite:     graps,
			Fillers:      fils,
			Coating:      coat,
			Construction: constr,
			Obturator:    obt,
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
	var mat string
	for i, am := range addit.Materials {
		if i > 0 {
			mat += ";"
		}
		mat += fmt.Sprintf("%s@%s", am.Short, am.Title)
	}

	dto := models.UpdateMat{
		Id:        addit.Id,
		Materials: mat,
	}
	err := s.repo.UpdateMat(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to update materials. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMod(addit *proto.UpdateAddModRequest) (*proto.SuccessResponse, error) {
	var mod string
	for i, am := range addit.Mod {
		if i > 0 {
			mod += ";"
		}
		mod += fmt.Sprintf("%s@%s@%s@%s", am.Id, am.Title, am.Short, am.Description)
	}

	dto := models.UpdateMod{
		Id:  addit.Id,
		Mod: mod,
	}
	err := s.repo.UpdateMod(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to update mod. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateTemp(addit *proto.UpdateAddTemRequest) (*proto.SuccessResponse, error) {
	var temp string
	for i, at := range addit.Temperature {
		if i > 0 {
			temp += ";"
		}
		temp += fmt.Sprintf("%s@%s", at.Id, at.Title)
	}

	dto := models.UpdateTemp{
		Id:          addit.Id,
		Temperature: temp,
	}
	err := s.repo.UpdateTemp(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to update temperature. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMoun(addit *proto.UpdateAddMounRequest) (*proto.SuccessResponse, error) {
	var moun string
	for i, am := range addit.Mounting {
		if i > 0 {
			moun += ";"
		}
		moun += fmt.Sprintf("%s@%s", am.Id, am.Title)
	}

	dto := models.UpdateMoun{
		Id:       addit.Id,
		Mounting: moun,
	}
	err := s.repo.UpdateMoun(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to update mounting. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateGrap(addit *proto.UpdateAddGrapRequest) (*proto.SuccessResponse, error) {
	var grap string
	for i, ag := range addit.Graphite {
		if i > 0 {
			grap += ";"
		}
		grap += fmt.Sprintf("%s@%s@%s", ag.Short, ag.Title, ag.Description)
	}

	dto := models.UpdateGrap{
		Id:       addit.Id,
		Graphite: grap,
	}
	if err := s.repo.UpdateGrap(dto); err != nil {
		return nil, fmt.Errorf("failed to update graphite. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateFillers(addit *proto.UpdateAddFillersRequest) (*proto.SuccessResponse, error) {
	var fil string
	for i, af := range addit.Fillers {
		if i > 0 {
			fil += ";"
		}
		fil += fmt.Sprintf("%s@%s@%s", af.Short, af.Title, af.Description)
	}

	dto := models.UpdateFill{
		Id:      addit.Id,
		Fillers: fil,
	}
	if err := s.repo.UpdateFillers(dto); err != nil {
		return nil, fmt.Errorf("failed to update fillers. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateCoating(addit *proto.UpdateAddCoatingRequest) (*proto.SuccessResponse, error) {
	var coat string
	for i, c := range addit.Coating {
		if i > 0 {
			coat += ";"
		}
		coat += fmt.Sprintf("%s@%s@%s@%s", c.Id, c.Short, c.Title, c.Description)
	}

	dto := models.UpdateCoating{
		Id:      addit.Id,
		Coating: coat,
	}
	if err := s.repo.UpdateCoating(dto); err != nil {
		return nil, fmt.Errorf("failed to update fillers. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateConstruction(addit *proto.UpdateAddConstructionRequest) (*proto.SuccessResponse, error) {
	var constr string
	for i, c := range addit.Constr {
		if i > 0 {
			constr += ";"
		}
		constr += fmt.Sprintf("%s@%s@%s@%s", c.Short, c.Title, c.Description, c.ForDescr)
	}

	dto := models.UpdateConstr{
		Id:           addit.Id,
		Construction: constr,
	}
	if err := s.repo.UpdateConstruction(dto); err != nil {
		return nil, fmt.Errorf("failed to update fillers. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateObturator(addit *proto.UpdateAddObturatorRequest) (*proto.SuccessResponse, error) {
	var obt string
	for i, o := range addit.Obturator {
		if i > 0 {
			obt += ";"
		}
		obt += fmt.Sprintf("%s@%s@%s", o.Short, o.Title, o.Description)
	}

	dto := models.UpdateObturator{
		Id:        addit.Id,
		Obturator: obt,
	}
	if err := s.repo.UpdateObturator(dto); err != nil {
		return nil, fmt.Errorf("failed to update fillers. error: %w", err)
	}
	return &proto.SuccessResponse{Success: true}, nil
}
