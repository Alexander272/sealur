package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type PutgmService struct {
	repo repository.Putgm
}

func NewPutgmService(repo repository.Putgm) *PutgmService {
	return &PutgmService{repo: repo}
}

func (s *PutgmService) Get(req *proto.GetPutgmRequest) (putgm []*proto.Putgm, err error) {
	data, err := s.repo.Get(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg. error: %w", err)
	}

	for _, d := range data {
		var construction []*proto.PutgmConstructions
		constrs := strings.Split(d.Construction, ";")

		for _, v := range constrs {
			grap := strings.Split(v, "&")[0]
			temp := strings.Split(v, "&")[1]

			var bas []*proto.PutgmConstr
			tmp := strings.Split(temp, "@")
			for _, t := range tmp {
				b := strings.Split(t, ">")[0]
				constr := strings.Split(t, ">")[1]

				var obts []*proto.PutgmObt
				tmp := strings.Split(constr, "*")
				for _, c := range tmp {
					short := strings.Split(c, "<")[0]
					seals := strings.Split(c, "<")[1]

					var seal []*proto.PutgmSeal
					tmp := strings.Split(seals, ",")
					for _, o := range tmp {
						short := strings.Split(o, "=")[0]
						url := strings.Split(o, "=")[1]
						seal = append(seal, &proto.PutgmSeal{Seal: short, ImageUrl: url})
					}

					obts = append(obts, &proto.PutgmObt{Obturator: short, Sealant: seal})
				}

				bas = append(bas, &proto.PutgmConstr{Basis: b, Obturator: obts})
			}

			construction = append(construction, &proto.PutgmConstructions{
				Grap: grap, Basis: bas,
			})
		}

		var temperatures []*proto.PutgTemp
		tmp := strings.Split(d.Temperatures, ";")

		for _, v := range tmp {
			grap := strings.Split(v, "&")[0]
			tmp := strings.Split(v, "&")[1]

			var Temps []*proto.Temp

			temps := strings.Split(tmp, "@")
			for _, t := range temps {
				id := strings.Split(t, ">")[0]
				tmp := strings.Split(t, ">")[1]

				mods := strings.Split(tmp, ",")
				Temps = append(Temps, &proto.Temp{Id: id, Mods: mods})
			}

			temperatures = append(temperatures, &proto.PutgTemp{
				Grap: grap, Temps: Temps,
			})
		}

		var basis, obturator = &proto.PutgMaterials{}, &proto.PutgMaterials{}
		tmp = strings.Split(d.Basis, "&")
		if len(tmp) > 1 {
			basis = &proto.PutgMaterials{Values: strings.Split(tmp[0], ";"), Default: tmp[1], Obturators: strings.Split(tmp[2], ";")}
		}
		tmp = strings.Split(d.Obturator, "&")
		if len(tmp) > 1 {
			obturator = &proto.PutgMaterials{Values: strings.Split(tmp[0], ";"), Default: tmp[1], Obturators: strings.Split(tmp[2], ";")}
		}

		p := proto.Putgm{
			Id:           d.Id,
			TypeFlId:     d.TypeFlId,
			TypePr:       d.TypePr,
			Form:         d.Form,
			Construction: construction,
			Temperatures: temperatures,
			Basis:        basis,
			Obturator:    obturator,
			Coating:      strings.Split(d.Coating, ";"),
			Mounting:     strings.Split(d.Mounting, ";"),
			Graphite:     strings.Split(d.Graphite, ";"),
		}
		putgm = append(putgm, &p)
	}

	return putgm, nil
}

func (s *PutgmService) Create(dto *proto.CreatePutgmRequest) (*proto.IdResponse, error) {
	var constructions string
	for i, c := range dto.Construction {
		if i > 0 {
			constructions += ";"
		}
		bas := ""
		for j, t := range c.Basis {
			if j > 0 {
				bas += "@"
			}
			constrs := ""
			for k, c := range t.Obturator {
				if k > 0 {
					constrs += "*"
				}
				seal := ""
				for l, o := range c.Sealant {
					if l > 0 {
						seal += ","
					}

					seal += fmt.Sprintf("%s=%s", o.Seal, o.ImageUrl)
				}

				constrs += fmt.Sprintf("%s<%s", c.Obturator, seal)
			}

			bas += fmt.Sprintf("%s>%s", t.Basis, constrs)
		}

		constructions += fmt.Sprintf("%s&%s", c.Grap, bas)
	}

	var temperatures string
	for i, t := range dto.Temperatures {
		if i > 0 {
			temperatures += ";"
		}
		temps := ""
		for j, t := range t.Temps {
			if j > 0 {
				temps += "@"
			}
			temps += fmt.Sprintf("%s>%s", t.Id, strings.Join(t.Mods, ","))
		}
		temperatures += fmt.Sprintf("%s&%s", t.Grap, temps)
	}

	basis := fmt.Sprintf("%s&%s&%s", strings.Join(dto.Basis.Values, ";"), dto.Basis.Default, strings.Join(dto.Basis.Obturators, ";"))
	if len(dto.Basis.Values) == 0 {
		basis = ""
	}
	obturator := fmt.Sprintf("%s&%s&%s", strings.Join(dto.Obturator.Values, ";"), dto.Obturator.Default, strings.Join(dto.Obturator.Obturators, ";"))
	if len(dto.Obturator.Values) == 0 {
		obturator = ""
	}

	putgm := models.PutgmDTO{
		FlangeId:     dto.FlangeId,
		TypeFlId:     dto.TypeFlId,
		TypePr:       dto.TypePr,
		Form:         dto.Form,
		Construction: constructions,
		Temperatures: temperatures,
		Basis:        basis,
		Obturator:    obturator,
		Coating:      strings.Join(dto.Coating, ";"),
		Mounting:     strings.Join(dto.Mounting, ";"),
		Graphite:     strings.Join(dto.Graphite, ";"),
	}

	id, err := s.repo.Create(putgm)
	if err != nil {
		return nil, fmt.Errorf("failed to create putg. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *PutgmService) Update(dto *proto.UpdatePutgmRequest) error {
	var constructions string
	for i, c := range dto.Construction {
		if i > 0 {
			constructions += ";"
		}
		bas := ""
		for j, t := range c.Basis {
			if j > 0 {
				bas += "@"
			}
			constrs := ""
			for k, c := range t.Obturator {
				if k > 0 {
					constrs += "*"
				}
				seal := ""
				for l, o := range c.Sealant {
					if l > 0 {
						seal += ","
					}

					seal += fmt.Sprintf("%s=%s", o.Seal, o.ImageUrl)
				}

				constrs += fmt.Sprintf("%s<%s", c.Obturator, seal)
			}

			bas += fmt.Sprintf("%s>%s", t.Basis, constrs)
		}

		constructions += fmt.Sprintf("%s&%s", c.Grap, bas)
	}

	var temperatures string
	for i, t := range dto.Temperatures {
		if i > 0 {
			temperatures += ";"
		}
		temps := ""
		for j, t := range t.Temps {
			if j > 0 {
				temps += "@"
			}
			temps += fmt.Sprintf("%s>%s", t.Id, strings.Join(t.Mods, ","))
		}
		temperatures += fmt.Sprintf("%s&%s", t.Grap, temps)
	}

	basis := fmt.Sprintf("%s&%s&%s", strings.Join(dto.Basis.Values, ";"), dto.Basis.Default, strings.Join(dto.Basis.Obturators, ";"))
	if len(dto.Basis.Values) == 0 {
		basis = ""
	}
	obturator := fmt.Sprintf("%s&%s&%s", strings.Join(dto.Obturator.Values, ";"), dto.Obturator.Default, strings.Join(dto.Obturator.Obturators, ";"))
	if len(dto.Obturator.Values) == 0 {
		obturator = ""
	}

	putgm := models.PutgmDTO{
		Id:           dto.Id,
		FlangeId:     dto.FlangeId,
		TypeFlId:     dto.TypeFlId,
		TypePr:       dto.TypePr,
		Form:         dto.Form,
		Construction: constructions,
		Temperatures: temperatures,
		Basis:        basis,
		Obturator:    obturator,
		Coating:      strings.Join(dto.Coating, ";"),
		Mounting:     strings.Join(dto.Mounting, ";"),
		Graphite:     strings.Join(dto.Graphite, ";"),
	}

	if err := s.repo.Update(putgm); err != nil {
		return fmt.Errorf("failed to update putg. error: %w", err)
	}
	return nil
}

func (s *PutgmService) Delete(putgm *proto.DeletePutgmRequest) error {
	if err := s.repo.Delete(putgm); err != nil {
		return fmt.Errorf("failed to delete putgm. error: %w", err)
	}
	return nil
}

func (s *PutgmService) DeleteGrap(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`construction like '%%%s&%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
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

		newConstr := make([]string, 0, 10)
		var constructions = strings.Split(cur.Construction, ";")
		for _, v := range constructions {
			if !strings.Contains(v, fmt.Sprintf("%s&", id)) {
				newConstr = append(newConstr, v)
			}
		}

		newTemp := make([]string, 0, 10)
		var graps = strings.Split(cur.Graphite, ";")
		for _, v := range graps {
			if !strings.Contains(v, fmt.Sprintf("%s&", id)) {
				newTemp = append(newTemp, v)
			}
		}

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: strings.Join(newConstr, ";"),
			Temperature:  strings.Join(newTemp, ";"),
			Basis:        cur.Basis,
			PObturator:   cur.Obturator,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     graphite,
		}

		go s.putgmUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) DeleteTemp(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`temperatures like '%%%s>%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

		var newTemp []string
		tmp := strings.Split(cur.Temperatures, ";")
		for _, con := range tmp {
			if strings.Contains(con, fmt.Sprintf("%s>", id)) {
				tmp := strings.Split(con, "&")[0]
				newMod := make([]string, 0, 10)
				arr := strings.Split(strings.Split(con, "&")[1], "@")

				for _, t := range arr {
					if !strings.Contains(t, fmt.Sprintf("%s>", id)) {
						newMod = append(newMod, t)
					}
				}

				if len(newMod) == 0 {
					continue
				}

				newTemp = append(newTemp, fmt.Sprintf("%s&%s", tmp, strings.Join(newMod, "@")))
			} else {
				newTemp = append(newTemp, con)
			}
		}

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: cur.Construction,
			Temperature:  strings.Join(newTemp, ";"),
			Basis:        cur.Basis,
			PObturator:   cur.Obturator,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgmUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) DeleteMod(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`temperatures like '%%>%s%%' OR temperatures like '%%,%s%%'`, id, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

		var newTemp []string
		tmp := strings.Split(cur.Temperatures, ";")
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

			newTemp = append(newTemp, fmt.Sprintf("%s&%s", tmp, strings.Join(newFil, "@")))
		}

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: cur.Construction,
			Temperature:  strings.Join(newTemp, ";"),
			Basis:        cur.Basis,
			PObturator:   cur.Obturator,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgmUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) DeleteMat(id string, materials []*proto.AddMaterials) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`obturator like '%%%s%%' OR basis like '%%%s%%'`, id, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

		var obt, bas string

		parts := strings.Split(cur.Basis, "&")
		if parts[0] == "" {
			bas = cur.Basis
		} else {
			mats := strings.Split(parts[0], ";")
			mats = filter(mats, id)

			if len(mats) == 0 {
				bas = ""
			} else {
				if parts[1] == id {
					parts[1] = materials[0].Short
				}
				if parts[0] == "*" {
					bas = fmt.Sprintf("%s&%s&%s", parts[0], parts[1], parts[2])
				} else {
					bas = fmt.Sprintf("%s&%s&%s", strings.Join(mats, ";"), parts[1], parts[2])
				}
			}

		}

		parts = strings.Split(cur.Obturator, "&")
		if parts[0] == "" {
			obt = cur.Obturator
		} else {
			mats := strings.Split(parts[0], ";")
			mats = filter(mats, id)

			if len(mats) == 0 {
				obt = ""
			} else {
				if parts[1] == id {
					parts[1] = materials[0].Short
				}
				if parts[0] == "*" {
					obt = fmt.Sprintf("%s&%s&%s", parts[0], parts[1], parts[2])
				} else {
					obt = fmt.Sprintf("%s&%s&%s", strings.Join(mats, ";"), parts[1], parts[2])
				}
			}
		}

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: cur.Construction,
			Temperature:  cur.Temperatures,
			Basis:        bas,
			PObturator:   obt,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgmUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) DeleteCon(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`construction like '%%&%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

		var newConstr []string
		tmp := strings.Split(cur.Construction, ";")
		for _, con := range tmp {
			if strings.Contains(con, fmt.Sprintf("%s>", id)) {
				tmp := strings.Split(con, "&")[0]
				newTemp := make([]string, 0, 10)
				arr := strings.Split(strings.Split(con, "&")[1], "@")

				for _, t := range arr {
					if !strings.Contains(t, fmt.Sprintf("%s>", id)) {
						newTemp = append(newTemp, t)
					}
				}

				if len(newTemp) == 0 {
					continue
				}

				newConstr = append(newConstr, fmt.Sprintf("%s&%s", tmp, strings.Join(newTemp, "@")))
			} else {
				newConstr = append(newConstr, con)
			}
		}

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: strings.Join(newConstr, ";"),
			Temperature:  cur.Temperatures,
			Basis:        cur.Basis,
			PObturator:   cur.Obturator,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgmUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) DeleteObt(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`construction like '%%<%s%%' OR obturator like '%%&%s%%' OR basis like '%%&%s%%'`, id, id, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

		var obt, bas string

		parts := strings.Split(cur.Basis, "&")
		if strings.Contains(parts[2], id) {
			obts := strings.Split(parts[2], ";")
			obts = filter(obts, id)
			bas = fmt.Sprintf("%s&%s&%s", parts[0], parts[1], strings.Join(obts, ";"))
		}

		parts = strings.Split(cur.Obturator, "&")
		if strings.Contains(parts[2], id) {
			obts := strings.Split(parts[2], ";")
			obts = filter(obts, id)
			obt = fmt.Sprintf("%s&%s&%s", parts[0], parts[1], strings.Join(obts, ";"))
		}

		var newConstr []string
		tmp := strings.Split(cur.Construction, ";")
		for _, con := range tmp {
			if strings.Contains(con, fmt.Sprintf(">%s", id)) || strings.Contains(con, fmt.Sprintf("*%s", id)) {
				tmp := strings.Split(con, "&")[0]
				newTemp := make([]string, 0, 10)
				arr := strings.Split(strings.Split(con, "&")[1], "@")

				for _, t := range arr {
					if strings.Contains(con, fmt.Sprintf(">%s", id)) || strings.Contains(con, fmt.Sprintf("*%s", id)) {

						tmp := strings.Split(t, ">")[0]
						newCon := make([]string, 0, 10)
						arr := strings.Split(strings.Split(t, ">")[1], "*")

						for _, t := range arr {
							if !strings.Contains(t, fmt.Sprintf("%s<", id)) {
								newCon = append(newCon, t)
							}
						}
						if len(newCon) == 0 {
							continue
						}

						newTemp = append(newTemp, fmt.Sprintf("%s>%s", tmp, strings.Join(newCon, "*")))
					} else {
						newTemp = append(newTemp, t)
					}
				}

				if len(newTemp) == 0 {
					continue
				}

				newConstr = append(newConstr, fmt.Sprintf("%s&%s", tmp, strings.Join(newTemp, "@")))
			} else {
				newConstr = append(newConstr, con)
			}
		}

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: strings.Join(newConstr, ";"),
			Temperature:  cur.Temperatures,
			Basis:        bas,
			PObturator:   obt,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgmUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) DeleteSeal(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`construction like '%%<%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

		var newConstr []string
		tmp := strings.Split(cur.Construction, ";")
		for _, con := range tmp {
			if strings.Contains(con, fmt.Sprintf("<%s", id)) {
				tmp := strings.Split(con, "&")[0]
				newTemp := make([]string, 0, 10)
				arr := strings.Split(strings.Split(con, "&")[1], "@")

				for _, t := range arr {
					if strings.Contains(con, fmt.Sprintf("<%s", id)) {

						tmp := strings.Split(t, ">")[0]
						newCon := make([]string, 0, 10)
						arr := strings.Split(strings.Split(t, ">")[1], "*")

						for _, t := range arr {
							if strings.Contains(t, fmt.Sprintf("<%s", id)) {
								tmp := strings.Split(t, "<")[0]
								newObt := make([]string, 0, 10)
								arr := strings.Split(strings.Split(t, "<")[1], ",")

								for _, t := range arr {
									if !strings.Contains(t, fmt.Sprintf("%s=", id)) {
										newObt = append(newObt, t)
									}
								}
								if len(newObt) == 0 {
									continue
								}

								newCon = append(newCon, fmt.Sprintf("%s<%s", tmp, strings.Join(newObt, ",")))
							} else {
								newCon = append(newCon, t)
							}
						}
						if len(newCon) == 0 {
							continue
						}

						newTemp = append(newTemp, fmt.Sprintf("%s>%s", tmp, strings.Join(newCon, "*")))
					} else {
						newTemp = append(newTemp, t)
					}
				}

				if len(newTemp) == 0 {
					continue
				}

				newConstr = append(newConstr, fmt.Sprintf("%s&%s", tmp, strings.Join(newTemp, "@")))
			} else {
				newConstr = append(newConstr, con)
			}
		}

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: strings.Join(newConstr, ";"),
			Temperature:  cur.Temperatures,
			Basis:        cur.Basis,
			PObturator:   cur.Obturator,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgmUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) DeleteMoun(id string) error {
	var wg sync.WaitGroup
	putgm, err := s.repo.GetByCondition(fmt.Sprintf(`mounting like '%%%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putgm {
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

		upPutgm := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: cur.Construction,
			Temperature:  cur.Temperatures,
			Basis:        cur.Basis,
			PObturator:   cur.Obturator,
			Coating:      cur.Coating,
			Mounting:     mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgmUpdate(upPutgm, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) DeleteCoating(id string) error {
	var wg sync.WaitGroup
	putgm, err := s.repo.GetByCondition(fmt.Sprintf(`coating like '%%%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putgm. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putgm {
		wg.Add(1)
		limit <- struct{}{}

		var coating string
		if cur.Coating == "*" || cur.Coating == "" {
			coating = cur.Mounting
		} else {
			mouns := strings.Split(cur.Coating, ";")
			mouns = filter(mouns, id)
			coating = strings.Join(mouns, ";")
		}

		upPutgm := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: cur.Construction,
			Temperature:  cur.Temperatures,
			Basis:        cur.Basis,
			PObturator:   cur.Obturator,
			Coating:      coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgmUpdate(upPutgm, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgmService) putgmUpdate(putgm models.UpdateAdditDTO, wg *sync.WaitGroup, limit chan struct{}) error {
	defer wg.Done()
	if err := s.repo.UpdateAddit(putgm); err != nil {
		return fmt.Errorf("failed to update putg addit. error: %w", err)
	}

	<-limit
	return nil
}
