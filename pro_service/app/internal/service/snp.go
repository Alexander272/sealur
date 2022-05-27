package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Alexander272/sealur/pro_service/internal/models"
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
	data, err := s.repo.Get(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get snp. error: %w", err)
	}

	for _, d := range data {
		var fillers []*proto.Filler
		fil := strings.Split(d.Fillers, ";")

		for _, v := range fil {
			id := strings.Split(v, "&")[0]
			tmp := strings.Split(v, "&")[1]

			var Temps []*proto.Temp

			temps := strings.Split(tmp, "@")
			for _, t := range temps {
				id := strings.Split(t, ">")[0]
				tmp := strings.Split(t, ">")[1]

				mods := strings.Split(tmp, ",")
				Temps = append(Temps, &proto.Temp{Id: id, Mods: mods})
			}

			fillers = append(fillers, &proto.Filler{
				Id: id, Temps: Temps,
			})
		}

		var frame, ir, or = &proto.Materials{}, &proto.Materials{}, &proto.Materials{}
		tmp := strings.Split(d.Frame, "&")
		if len(tmp) > 1 {
			frame = &proto.Materials{Values: strings.Split(tmp[0], ";"), Default: tmp[1]}
		}
		tmp = strings.Split(d.Ir, "&")
		if len(tmp) > 1 {
			ir = &proto.Materials{Values: strings.Split(tmp[0], ";"), Default: tmp[1]}
		}
		tmp = strings.Split(d.Or, "&")
		if len(tmp) > 1 {
			or = &proto.Materials{Values: strings.Split(tmp[0], ";"), Default: tmp[1]}
		}

		s := proto.SNP{
			Id:       d.Id,
			TypeFlId: d.TypeFlId,
			TypePr:   d.TypePr,
			Fillers:  fillers,
			Frame:    frame,
			Ir:       ir,
			Or:       or,
			Mounting: strings.Split(d.Mounting, ";"),
			Graphite: strings.Split(d.Graphite, ";"),
		}
		snp = append(snp, &s)
	}

	return snp, nil
}

func (s *SNPService) Create(dto *proto.CreateSNPRequest) (*proto.IdResponse, error) {
	var fillers string
	for i, f := range dto.Fillers {
		if i > 0 {
			fillers += ";"
		}
		temps := ""
		for j, t := range f.Temps {
			if j > 0 {
				temps += "@"
			}
			temps += fmt.Sprintf("%s>%s", t.Id, strings.Join(t.Mods, ","))
		}
		fillers += fmt.Sprintf("%s&%s", f.Id, temps)
	}

	frame := fmt.Sprintf("%s&%s", strings.Join(dto.Frame.Values, ";"), dto.Frame.Default)
	if len(dto.Frame.Values) == 0 {
		frame = ""
	}
	ir := fmt.Sprintf("%s&%s", strings.Join(dto.Ir.Values, ";"), dto.Ir.Default)
	if len(dto.Ir.Values) == 0 {
		ir = ""
	}
	or := fmt.Sprintf("%s&%s", strings.Join(dto.Or.Values, ";"), dto.Or.Default)
	if len(dto.Or.Values) == 0 {
		or = ""
	}

	snp := models.SnpDTO{
		StandId:  dto.StandId,
		FlangeId: dto.FlangeId,
		TypeFlId: dto.TypeFlId,
		TypePr:   dto.TypePr,
		Fillers:  fillers,
		Frame:    frame,
		Ir:       ir,
		Or:       or,
		Mounting: strings.Join(dto.Mounting, ";"),
		Graphite: strings.Join(dto.Graphite, ";"),
	}

	id, err := s.repo.Create(snp)
	if err != nil {
		return nil, fmt.Errorf("failed to create snp. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *SNPService) Update(dto *proto.UpdateSNPRequest) error {
	var fillers string
	for i, f := range dto.Fillers {
		if i > 0 {
			fillers += ";"
		}
		temps := ""
		for j, t := range f.Temps {
			if j > 0 {
				temps += "@"
			}
			temps += fmt.Sprintf("%s>%s", t.Id, strings.Join(t.Mods, ","))
		}
		fillers += fmt.Sprintf("%s&%s", f.Id, temps)
	}

	frame := fmt.Sprintf("%s&%s", strings.Join(dto.Frame.Values, ";"), dto.Frame.Default)
	if len(dto.Frame.Values) == 0 {
		frame = ""
	}
	ir := fmt.Sprintf("%s&%s", strings.Join(dto.Ir.Values, ";"), dto.Ir.Default)
	if len(dto.Ir.Values) == 0 {
		ir = ""
	}
	or := fmt.Sprintf("%s&%s", strings.Join(dto.Or.Values, ";"), dto.Or.Default)
	if len(dto.Or.Values) == 0 {
		or = ""
	}

	snp := models.SnpDTO{
		Id:       dto.Id,
		StandId:  dto.StandId,
		FlangeId: dto.FlangeId,
		TypeFlId: dto.TypeFlId,
		TypePr:   dto.TypePr,
		Fillers:  fillers,
		Frame:    frame,
		Ir:       ir,
		Or:       or,
		Mounting: strings.Join(dto.Mounting, ";"),
		Graphite: strings.Join(dto.Graphite, ";"),
	}

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

func (s *SNPService) AddMat(id string) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(`(NOT(frame = '*') AND NOT(frame = '')) 
		OR (NOT(in_ring = '*') AND NOT(in_ring = '')) OR (NOT(ou_ring = '*') AND NOT(ou_ring = ''))`)
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var frame, ir, or string

		parts := strings.Split(cur.Ir, "&")
		if parts[0] == "*" || parts[0] == "" {
			ir = cur.Ir
		} else {
			ir = fmt.Sprintf("%s;%s&%s", parts[0], id, parts[1])
		}
		parts = strings.Split(cur.Or, "&")
		if parts[0] == "*" || parts[0] == "" {
			or = cur.Or
		} else {
			or = fmt.Sprintf("%s;%s&%s", parts[0], id, parts[1])
		}
		parts = strings.Split(cur.Frame, "&")
		if parts[0] == "*" || parts[0] == "" {
			frame = cur.Frame
		} else {
			frame = fmt.Sprintf("%s;%s&%s", parts[0], id, parts[1])
		}

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  cur.Fillers,
			Mounting: cur.Mounting,
			Graphite: cur.Graphite,
			Frame:    frame,
			Ir:       ir,
			Or:       or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) DeleteMat(id string, materials []*proto.AddMaterials) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(fmt.Sprintf(`frame like '%%%s%%' OR in_ring like '%%%s%%' OR ou_ring like '%%%s%%'`, id, id, id))
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var frame, ir, or string

		parts := strings.Split(cur.Ir, "&")
		if parts[0] == "" {
			ir = cur.Ir
		} else {
			mats := strings.Split(parts[0], ";")
			mats = filter(mats, id)

			if len(mats) == 0 {
				ir = ""
			} else {
				if parts[1] == id {
					parts[1] = materials[0].Short
				}
				if parts[0] == "*" {
					ir = fmt.Sprintf("%s&%s", parts[0], parts[1])
				} else {
					ir = fmt.Sprintf("%s&%s", strings.Join(mats, ";"), parts[1])
				}
			}

		}

		parts = strings.Split(cur.Or, "&")
		if parts[0] == "" {
			or = cur.Or
		} else {
			mats := strings.Split(parts[0], ";")
			mats = filter(mats, id)

			if len(mats) == 0 {
				or = ""
			} else {
				if parts[1] == id {
					parts[1] = materials[0].Short
				}
				if parts[0] == "*" {
					or = fmt.Sprintf("%s&%s", parts[0], parts[1])
				} else {
					or = fmt.Sprintf("%s&%s", strings.Join(mats, ";"), parts[1])
				}
			}
		}

		parts = strings.Split(cur.Frame, "&")
		if parts[0] == "" {
			frame = cur.Frame
		} else {
			mats := strings.Split(parts[0], ";")
			mats = filter(mats, id)

			if len(mats) == 0 {
				frame = ""
			} else {
				if parts[1] == id {
					parts[1] = materials[0].Short
				}
				if parts[0] == "*" {
					frame = fmt.Sprintf("%s&%s", parts[0], parts[1])
				} else {
					frame = fmt.Sprintf("%s&%s", strings.Join(mats, ";"), parts[1])
				}
			}
		}

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  cur.Fillers,
			Mounting: cur.Mounting,
			Graphite: cur.Graphite,
			Frame:    frame,
			Ir:       ir,
			Or:       or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) AddMoun(id string) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(`NOT(mounting = '' OR mounting='*')`)
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var mounting string

		if cur.Mounting == "*" || cur.Mounting == "" {
			mounting = cur.Mounting
		} else {
			mounting = fmt.Sprintf("%s;%s", cur.Mounting, id)
		}

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  cur.Fillers,
			Mounting: mounting,
			Graphite: cur.Graphite,
			Frame:    cur.Frame,
			Ir:       cur.Ir,
			Or:       cur.Or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) DeleteMoun(id string) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(fmt.Sprintf(`mounting like '%%%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var mounting string
		if cur.Mounting == "*" || cur.Mounting == "" {
			mounting = cur.Mounting
		} else {
			mouns := strings.Split(cur.Mounting, ";")
			mouns = filter(mouns, id)
			mounting = strings.Join(mouns, ";")
		}

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  cur.Fillers,
			Mounting: mounting,
			Graphite: cur.Graphite,
			Frame:    cur.Frame,
			Ir:       cur.Ir,
			Or:       cur.Or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) AddGrap(id string) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(`NOT(mounting = '' OR mounting='*')`)
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var graphite string

		if cur.Graphite == "*" || cur.Graphite == "" {
			graphite = cur.Graphite
		} else {
			graphite = fmt.Sprintf("%s;%s", cur.Graphite, id)
		}

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  cur.Fillers,
			Mounting: cur.Mounting,
			Graphite: graphite,
			Frame:    cur.Frame,
			Ir:       cur.Ir,
			Or:       cur.Or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) DeleteGrap(id string) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(fmt.Sprintf(`graphite like '%%%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var graphite string
		if cur.Graphite == "*" || cur.Graphite == "" {
			graphite = cur.Graphite
		} else {
			graps := strings.Split(cur.Graphite, ";")
			graps = filter(graps, id)
			graphite = strings.Join(graps, ";")
		}

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  cur.Fillers,
			Mounting: cur.Mounting,
			Graphite: graphite,
			Frame:    cur.Frame,
			Ir:       cur.Ir,
			Or:       cur.Or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) DeleteFiller(id string) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(fmt.Sprintf(`filler like '%%%s&%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var fillers string
		tmp := strings.Split(cur.Fillers, ";")
		newFil := make([]string, 0, len(tmp))
		for _, fil := range tmp {
			if !strings.Contains(fil, fmt.Sprintf("%s&", id)) {
				newFil = append(newFil, fil)
			}
		}
		fillers = strings.Join(newFil, ";")

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  fillers,
			Mounting: cur.Mounting,
			Graphite: cur.Graphite,
			Frame:    cur.Frame,
			Ir:       cur.Ir,
			Or:       cur.Or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) DeleteTemp(id string) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(fmt.Sprintf(`filler like '%%%s>%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var fillers []string
		tmp := strings.Split(cur.Fillers, ";")
		for _, fil := range tmp {
			if strings.Contains(fil, fmt.Sprintf("%s>", id)) {
				tmp := strings.Split(fil, "&")[0]
				newFil := make([]string, 0, 10)
				arr := strings.Split(strings.Split(fil, "&")[1], "@")

				for _, t := range arr {
					if !strings.Contains(t, fmt.Sprintf("%s>", id)) {
						newFil = append(newFil, t)
					}
				}

				if len(newFil) == 0 {
					continue
				}

				fillers = append(fillers, fmt.Sprintf("%s&%s", tmp, strings.Join(newFil, "@")))
			} else {
				fillers = append(fillers, fil)
			}
		}

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  strings.Join(fillers, ";"),
			Mounting: cur.Mounting,
			Graphite: cur.Graphite,
			Frame:    cur.Frame,
			Ir:       cur.Ir,
			Or:       cur.Or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) DeleteMod(id string) error {
	var wg sync.WaitGroup
	snp, err := s.repo.GetByCondition(fmt.Sprintf(`filler like '%%>%s%%' OR filler like '%%,%s%%'`, id, id))
	if err != nil {
		return fmt.Errorf("failed to get snp. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range snp {
		wg.Add(1)
		limit <- struct{}{}

		var fillers []string
		tmp := strings.Split(cur.Fillers, ";")
		for _, fil := range tmp {
			tmp := strings.Split(fil, "&")[0]
			newFil := make([]string, 0, 10)
			arr := strings.Split(strings.Split(fil, "&")[1], "@")
			for _, t := range arr {
				if strings.Contains(t, fmt.Sprintf(">%s", id)) || strings.Contains(t, fmt.Sprintf(",%s", id)) {
					tmp := strings.Split(t, ">")[0]
					newTmp := make([]string, 0, 10)
					arr := strings.Split(strings.Split(t, ">")[1], ",")

					for _, t := range arr {
						if t != id {
							newTmp = append(newTmp, t)
						}
					}

					if len(newTmp) == 0 {
						continue
					}

					newFil = append(newFil, fmt.Sprintf("%s>%s", tmp, strings.Join(newTmp, ",")))
				} else {
					newFil = append(newFil, t)
				}
			}

			if len(newFil) == 0 {
				continue
			}

			fillers = append(fillers, fmt.Sprintf("%s&%s", tmp, strings.Join(newFil, "@")))
		}

		upSnp := models.UpdateAdditDTO{
			Id:       cur.Id,
			Fillers:  strings.Join(fillers, ";"),
			Mounting: cur.Mounting,
			Graphite: cur.Graphite,
			Frame:    cur.Frame,
			Ir:       cur.Ir,
			Or:       cur.Or,
		}

		go s.snpUpdate(upSnp, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *SNPService) snpUpdate(snp models.UpdateAdditDTO, wg *sync.WaitGroup, limit chan struct{}) error {
	defer wg.Done()
	if err := s.repo.UpdateAddit(snp); err != nil {
		return fmt.Errorf("failed to update snp addit. error: %w", err)
	}

	<-limit
	return nil
}

func filter(arr []string, id string) []string {
	var res []string
	for i := range arr {
		if arr[i] != id {
			res = append(res, arr[i])
		}
	}
	return res
}
