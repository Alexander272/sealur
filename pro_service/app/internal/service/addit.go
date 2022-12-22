package service

import (
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type AdditService struct {
	repo repository.Addit
}

func NewAdditService(repo repository.Addit) *AdditService {
	return &AdditService{repo: repo}
}

func (s *AdditService) GetAll() (addit []*pro_api.Additional, err error) {
	data, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get additional. error: %w", err)
	}

	for _, d := range data {
		var mats []*pro_api.AddMaterials
		var mods []*pro_api.AddMod
		var temps []*pro_api.AddTemperature
		var mouns []*pro_api.AddMoun
		var graps []*pro_api.AddGrap
		var fils []*pro_api.AddFiller
		var coat []*pro_api.AddCoating
		var constr []*pro_api.AddConstruction
		var obt []*pro_api.AddObturator
		var basis []*pro_api.AddBasis
		var pObt []*pro_api.AddPObturator
		var seal []*pro_api.AddSealant

		tmp := strings.Split(d.Materials, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			mats = append(mats, &pro_api.AddMaterials{
				Short: parts[0],
				Title: parts[1],
			})
		}
		tmp = strings.Split(d.Mod, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			mods = append(mods, &pro_api.AddMod{
				Id:          parts[0],
				Short:       parts[2],
				Title:       parts[1],
				Description: parts[3],
			})
		}
		tmp = strings.Split(d.Temperature, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			temps = append(temps, &pro_api.AddTemperature{
				Id:    parts[0],
				Title: parts[1],
			})
		}
		tmp = strings.Split(d.Mounting, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			mouns = append(mouns, &pro_api.AddMoun{
				Id:    parts[0],
				Title: parts[1],
			})
		}
		tmp = strings.Split(d.Graphite, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			graps = append(graps, &pro_api.AddGrap{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
			})
		}
		tmp = strings.Split(d.Fillers, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			fils = append(fils, &pro_api.AddFiller{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
			})
		}
		tmp = strings.Split(d.Coating, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			coat = append(coat, &pro_api.AddCoating{
				Id:          parts[0],
				Short:       parts[1],
				Title:       parts[2],
				Description: parts[3],
			})
		}
		tmp = strings.Split(d.Construction, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			constr = append(constr, &pro_api.AddConstruction{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
			})
		}
		tmp = strings.Split(d.Obturator, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			obt = append(obt, &pro_api.AddObturator{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
				ForDescr:    parts[3],
			})
		}
		tmp = strings.Split(d.Basis, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			basis = append(basis, &pro_api.AddBasis{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
			})
		}
		tmp = strings.Split(d.Sealant, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			seal = append(seal, &pro_api.AddSealant{
				Id:          parts[0],
				Short:       parts[1],
				Title:       parts[2],
				Description: parts[3],
				ForDescr:    parts[4],
			})
		}
		tmp = strings.Split(d.PObturator, ";")
		for _, v := range tmp {
			parts := strings.Split(v, "@")
			pObt = append(pObt, &pro_api.AddPObturator{
				Short:       parts[0],
				Title:       parts[1],
				Description: parts[2],
				ForDescr:    parts[3],
			})
		}

		addit = append(addit, &pro_api.Additional{
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
			Basis:        basis,
			PObturator:   pObt,
			Sealant:      seal,
		})
	}

	return addit, nil
}

func (s *AdditService) Create(addit *pro_api.CreateAddRequest) (*pro_api.SuccessResponse, error) {
	err := s.repo.Create(addit)
	if err != nil {
		return nil, fmt.Errorf("failed to create additional. error: %w", err)
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMat(addit *pro_api.UpdateAddMatRequest) (*pro_api.SuccessResponse, error) {
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
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMod(addit *pro_api.UpdateAddModRequest) (*pro_api.SuccessResponse, error) {
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
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateTemp(addit *pro_api.UpdateAddTemRequest) (*pro_api.SuccessResponse, error) {
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
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateMoun(addit *pro_api.UpdateAddMounRequest) (*pro_api.SuccessResponse, error) {
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
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateGrap(addit *pro_api.UpdateAddGrapRequest) (*pro_api.SuccessResponse, error) {
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
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateFillers(addit *pro_api.UpdateAddFillersRequest) (*pro_api.SuccessResponse, error) {
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
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateCoating(addit *pro_api.UpdateAddCoatingRequest) (*pro_api.SuccessResponse, error) {
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
		return nil, fmt.Errorf("failed to update coating. error: %w", err)
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateConstruction(addit *pro_api.UpdateAddConstructionRequest) (*pro_api.SuccessResponse, error) {
	var constr string
	for i, c := range addit.Constr {
		if i > 0 {
			constr += ";"
		}
		constr += fmt.Sprintf("%s@%s@%s", c.Short, c.Title, c.Description)
	}

	dto := models.UpdateConstr{
		Id:           addit.Id,
		Construction: constr,
	}
	if err := s.repo.UpdateConstruction(dto); err != nil {
		return nil, fmt.Errorf("failed to update construction. error: %w", err)
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateObturator(addit *pro_api.UpdateAddObturatorRequest) (*pro_api.SuccessResponse, error) {
	var obt string
	for i, o := range addit.Obturator {
		if i > 0 {
			obt += ";"
		}
		obt += fmt.Sprintf("%s@%s@%s@%s", o.Short, o.Title, o.Description, o.ForDescr)
	}

	dto := models.UpdateObturator{
		Id:        addit.Id,
		Obturator: obt,
	}
	if err := s.repo.UpdateObturator(dto); err != nil {
		return nil, fmt.Errorf("failed to update obturator. error: %w", err)
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateBasis(addit *pro_api.UpdateAddBasisRequest) (*pro_api.SuccessResponse, error) {
	var basis string
	for i, b := range addit.Basis {
		if i > 0 {
			basis += ";"
		}
		basis += fmt.Sprintf("%s@%s@%s", b.Short, b.Title, b.Description)
	}

	dto := models.UpdateBasis{
		Id:    addit.Id,
		Basis: basis,
	}
	if err := s.repo.UpdateBasis(dto); err != nil {
		return nil, fmt.Errorf("failed to update basis. error: %w", err)
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdatePObturator(addit *pro_api.UpdateAddPObturatorRequest) (*pro_api.SuccessResponse, error) {
	var pObt string
	for i, p := range addit.PObturator {
		if i > 0 {
			pObt += ";"
		}
		pObt += fmt.Sprintf("%s@%s@%s@%s", p.Short, p.Title, p.Description, p.ForDescr)
	}

	dto := models.UpdatePObturator{
		Id:         addit.Id,
		PObturator: pObt,
	}
	if err := s.repo.UpdatePObturator(dto); err != nil {
		return nil, fmt.Errorf("failed to update p_obturator. error: %w", err)
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}

func (s *AdditService) UpdateSealant(addit *pro_api.UpdateAddSealantRequest) (*pro_api.SuccessResponse, error) {
	var seal string
	for i, s := range addit.Sealant {
		if i > 0 {
			seal += ";"
		}
		seal += fmt.Sprintf("%s@%s@%s@%s@%s", s.Id, s.Short, s.Title, s.Description, s.ForDescr)
	}

	dto := models.UpdateSealant{
		Id:      addit.Id,
		Sealant: seal,
	}
	if err := s.repo.UpdateSealant(dto); err != nil {
		return nil, fmt.Errorf("failed to update sealant. error: %w", err)
	}
	return &pro_api.SuccessResponse{Success: true}, nil
}
