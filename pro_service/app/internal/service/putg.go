package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
)

type PutgService struct {
	repo repository.Putg
}

func NewPutgService(repo repository.Putg) *PutgService {
	return &PutgService{repo: repo}
}

func (s *PutgService) Get(req *proto.GetPutgRequest) (putg []*proto.Putg, err error) {
	data, err := s.repo.Get(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get putg. error: %w", err)
	}

	for _, d := range data {
		var construction []*proto.PutgConstructions
		constrs := strings.Split(d.Construction, ";")

		for _, v := range constrs {
			grap := strings.Split(v, "&")[0]
			temp := strings.Split(v, "&")[1]

			var Temps []*proto.ConTemp
			tmp := strings.Split(temp, "@")
			for _, t := range tmp {
				temp := strings.Split(t, ">")[0]
				constr := strings.Split(t, ">")[1]

				var constrs []*proto.Constr
				tmp := strings.Split(constr, "*")
				for _, c := range tmp {
					short := strings.Split(c, "<")[0]
					obts := strings.Split(c, "<")[1]

					var obt []*proto.PutgObt
					tmp := strings.Split(obts, ",")
					for _, o := range tmp {
						short := strings.Split(o, "=")[0]
						url := strings.Split(o, "=")[1]
						obt = append(obt, &proto.PutgObt{Short: short, ImageUrl: url})
					}

					if len(obts) == 0 {
						continue
					}

					constrs = append(constrs, &proto.Constr{Short: short, Obturators: obt})
				}

				if len(constr) == 0 {
					continue
				}

				Temps = append(Temps, &proto.ConTemp{Temp: temp, Constructions: constrs})
			}

			if len(Temps) == 0 {
				continue
			}

			construction = append(construction, &proto.PutgConstructions{
				Grap: grap, Temperatures: Temps,
			})
		}

		var temperatures []*proto.PutgTemp
		tmp := strings.Split(d.Temperatures, ";")

		for _, v := range tmp {
			grap := strings.Split(v, "&")[0]
			tmp := strings.Split(v, "&")[1]

			if tmp == "" {
				continue
			}

			var Temps []*proto.Temp

			temps := strings.Split(tmp, "@")
			for _, t := range temps {
				id := strings.Split(t, ">")[0]
				tmp := strings.Split(t, ">")[1]

				mods := strings.Split(tmp, ",")
				if len(mods) == 0 {
					continue
				}
				Temps = append(Temps, &proto.Temp{Id: id, Mods: mods})
			}

			temperatures = append(temperatures, &proto.PutgTemp{
				Grap: grap, Temps: Temps,
			})
		}

		var reinforce, obturator, iLimiter, oLimiter = &proto.PutgMaterials{}, &proto.PutgMaterials{}, &proto.PutgMaterials{}, &proto.PutgMaterials{}
		tmp = strings.Split(d.Reinforce, "&")
		if len(tmp) > 1 {
			reinforce = &proto.PutgMaterials{Values: strings.Split(tmp[0], ";"), Default: tmp[1], Obturators: strings.Split(tmp[2], ";")}
		}
		tmp = strings.Split(d.Obturator, "&")
		if len(tmp) > 1 {
			obturator = &proto.PutgMaterials{Values: strings.Split(tmp[0], ";"), Default: tmp[1], Obturators: strings.Split(tmp[2], ";")}
		}
		tmp = strings.Split(d.ILimiter, "&")
		if len(tmp) > 1 {
			iLimiter = &proto.PutgMaterials{Values: strings.Split(tmp[0], ";"), Default: tmp[1], Obturators: strings.Split(tmp[2], ";")}
		}
		tmp = strings.Split(d.OLimiter, "&")
		if len(tmp) > 1 {
			oLimiter = &proto.PutgMaterials{Values: strings.Split(tmp[0], ";"), Default: tmp[1], Obturators: strings.Split(tmp[2], ";")}
		}

		p := proto.Putg{
			Id:           d.Id,
			TypeFlId:     d.TypeFlId,
			TypePr:       d.TypePr,
			Form:         d.Form,
			Construction: construction,
			Temperatures: temperatures,
			Reinforce:    reinforce,
			Obturator:    obturator,
			ILimiter:     iLimiter,
			OLimiter:     oLimiter,
			Coating:      strings.Split(d.Coating, ";"),
			Mounting:     strings.Split(d.Mounting, ";"),
			Graphite:     strings.Split(d.Graphite, ";"),
		}
		putg = append(putg, &p)
	}

	return putg, nil
}

func (s *PutgService) Create(dto *proto.CreatePutgRequest) (*proto.IdResponse, error) {
	var constructions string
	for i, c := range dto.Construction {
		if i > 0 {
			constructions += ";"
		}
		temps := ""
		for j, t := range c.Temperatures {
			if j > 0 {
				temps += "@"
			}
			constrs := ""
			for k, c := range t.Constructions {
				if k > 0 {
					constrs += "*"
				}
				obts := ""
				for l, o := range c.Obturators {
					if l > 0 {
						obts += ","
					}

					obts += fmt.Sprintf("%s=%s", o.Short, o.ImageUrl)
				}

				constrs += fmt.Sprintf("%s<%s", c.Short, obts)
			}

			temps += fmt.Sprintf("%s>%s", t.Temp, constrs)
		}

		constructions += fmt.Sprintf("%s&%s", c.Grap, temps)
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

	reinforce := fmt.Sprintf("%s&%s&%s", strings.Join(dto.Reinforce.Values, ";"), dto.Reinforce.Default, strings.Join(dto.Reinforce.Obturators, ";"))
	if len(dto.Reinforce.Values) == 0 {
		reinforce = ""
	}
	obturator := fmt.Sprintf("%s&%s&%s", strings.Join(dto.Obturator.Values, ";"), dto.Obturator.Default, strings.Join(dto.Obturator.Obturators, ";"))
	if len(dto.Obturator.Values) == 0 {
		obturator = ""
	}
	iLimiter := fmt.Sprintf("%s&%s&%s", strings.Join(dto.ILimiter.Values, ";"), dto.ILimiter.Default, strings.Join(dto.ILimiter.Obturators, ";"))
	if len(dto.ILimiter.Values) == 0 {
		iLimiter = ""
	}
	oLimiter := fmt.Sprintf("%s&%s&%s", strings.Join(dto.OLimiter.Values, ";"), dto.OLimiter.Default, strings.Join(dto.OLimiter.Obturators, ";"))
	if len(dto.OLimiter.Values) == 0 {
		oLimiter = ""
	}

	putg := models.PutgDTO{
		FlangeId:     dto.FlangeId,
		TypeFlId:     dto.TypeFlId,
		TypePr:       dto.TypePr,
		Form:         dto.Form,
		Construction: constructions,
		Temperatures: temperatures,
		Reinforce:    reinforce,
		Obturator:    obturator,
		ILimiter:     iLimiter,
		OLimiter:     oLimiter,
		Coating:      strings.Join(dto.Coating, ";"),
		Mounting:     strings.Join(dto.Mounting, ";"),
		Graphite:     strings.Join(dto.Graphite, ";"),
	}

	id, err := s.repo.Create(putg)
	if err != nil {
		return nil, fmt.Errorf("failed to create putg. error: %w", err)
	}
	return &proto.IdResponse{Id: id}, nil
}

func (s *PutgService) Update(dto *proto.UpdatePutgRequest) error {
	var constructions string
	for i, c := range dto.Construction {
		if i > 0 {
			constructions += ";"
		}
		temps := ""
		for j, t := range c.Temperatures {
			if j > 0 {
				temps += "@"
			}
			constrs := ""
			for k, c := range t.Constructions {
				if k > 0 {
					constrs += "*"
				}
				obts := ""
				for l, o := range c.Obturators {
					if l > 0 {
						obts += ","
					}

					obts += fmt.Sprintf("%s=%s", o.Short, o.ImageUrl)
				}

				constrs += fmt.Sprintf("%s<%s", c.Short, obts)
			}

			temps += fmt.Sprintf("%s>%s", t.Temp, constrs)
		}

		constructions += fmt.Sprintf("%s&%s", c.Grap, temps)
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

	reinforce := fmt.Sprintf("%s&%s&%s", strings.Join(dto.Reinforce.Values, ";"), dto.Reinforce.Default, strings.Join(dto.Reinforce.Obturators, ";"))
	if len(dto.Reinforce.Values) == 0 {
		reinforce = ""
	}
	obturator := fmt.Sprintf("%s&%s&%s", strings.Join(dto.Obturator.Values, ";"), dto.Obturator.Default, strings.Join(dto.Obturator.Obturators, ";"))
	if len(dto.Obturator.Values) == 0 {
		obturator = ""
	}
	iLimiter := fmt.Sprintf("%s&%s&%s", strings.Join(dto.ILimiter.Values, ";"), dto.ILimiter.Default, strings.Join(dto.ILimiter.Obturators, ";"))
	if len(dto.ILimiter.Values) == 0 {
		iLimiter = ""
	}
	oLimiter := fmt.Sprintf("%s&%s&%s", strings.Join(dto.OLimiter.Values, ";"), dto.OLimiter.Default, strings.Join(dto.OLimiter.Obturators, ";"))
	if len(dto.OLimiter.Values) == 0 {
		oLimiter = ""
	}

	putg := models.PutgDTO{
		Id:           dto.Id,
		FlangeId:     dto.FlangeId,
		TypeFlId:     dto.TypeFlId,
		TypePr:       dto.TypePr,
		Form:         dto.Form,
		Construction: constructions,
		Temperatures: temperatures,
		Reinforce:    reinforce,
		Obturator:    obturator,
		ILimiter:     iLimiter,
		OLimiter:     oLimiter,
		Coating:      strings.Join(dto.Coating, ";"),
		Mounting:     strings.Join(dto.Mounting, ";"),
		Graphite:     strings.Join(dto.Graphite, ";"),
	}

	if err := s.repo.Update(putg); err != nil {
		return fmt.Errorf("failed to update putg. error: %w", err)
	}
	return nil
}

func (s *PutgService) Delete(putg *proto.DeletePutgRequest) error {
	if err := s.repo.Delete(putg); err != nil {
		return fmt.Errorf("failed to delete putg. error: %w", err)
	}
	return nil
}

func (s *PutgService) DeleteGrap(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`construction like '%%%s&%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putg. error: %w", err)
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
			Obturator:    cur.Obturator,
			ILimiter:     cur.ILimiter,
			OLimiter:     cur.OLimiter,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     graphite,
		}

		go s.putgUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgService) DeleteTemp(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`construction like '%%%s>%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putg. error: %w", err)
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

		var newTemp []string
		tmp = strings.Split(cur.Temperatures, ";")
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
			Construction: strings.Join(newConstr, ";"),
			Temperature:  strings.Join(newTemp, ";"),
			Obturator:    cur.Obturator,
			ILimiter:     cur.ILimiter,
			OLimiter:     cur.OLimiter,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgService) DeleteMod(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`temperatures like '%%>%s%%' OR temperatures like '%%,%s%%'`, id, id))
	if err != nil {
		return fmt.Errorf("failed to get putg. error: %w", err)
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
			Obturator:    cur.Obturator,
			ILimiter:     cur.ILimiter,
			OLimiter:     cur.OLimiter,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgService) DeleteMat(id string, materials []*proto.AddMaterials) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`obturator like '%%%s%%' OR i_limiter like '%%%s%%' OR o_limiter like '%%%s%%'`, id, id, id))
	if err != nil {
		return fmt.Errorf("failed to get putg. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

		var obt, il, ol string

		parts := strings.Split(cur.ILimiter, "&")
		if parts[0] == "" {
			il = cur.ILimiter
		} else {
			mats := strings.Split(parts[0], ";")
			mats = filter(mats, id)

			if len(mats) == 0 {
				il = ""
			} else {
				if parts[1] == id {
					parts[1] = materials[0].Short
				}
				if parts[0] == "*" {
					il = fmt.Sprintf("%s&%s&%s", parts[0], parts[1], parts[2])
				} else {
					il = fmt.Sprintf("%s&%s&%s", strings.Join(mats, ";"), parts[1], parts[2])
				}
			}

		}

		parts = strings.Split(cur.OLimiter, "&")
		if parts[0] == "" {
			ol = cur.OLimiter
		} else {
			mats := strings.Split(parts[0], ";")
			mats = filter(mats, id)

			if len(mats) == 0 {
				ol = ""
			} else {
				if parts[1] == id {
					parts[1] = materials[0].Short
				}
				if parts[0] == "*" {
					ol = fmt.Sprintf("%s&%s&%s", parts[0], parts[1], parts[2])
				} else {
					ol = fmt.Sprintf("%s&%s&%s", strings.Join(mats, ";"), parts[1], parts[2])
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
			Obturator:    obt,
			ILimiter:     il,
			OLimiter:     ol,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgService) DeleteCon(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`construction like '%%>%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putg. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

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
			Obturator:    cur.Obturator,
			ILimiter:     cur.ILimiter,
			OLimiter:     cur.OLimiter,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgService) DeleteObt(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`construction like '%%<%s%%' OR obturator like '%%&%s%%' OR i_limiter like '%%&%s%%' 
	OR o_limiter like '%%&%s%%'`, id, id, id, id))
	if err != nil {
		return fmt.Errorf("failed to get putg. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
		wg.Add(1)
		limit <- struct{}{}

		var obt, il, ol string

		parts := strings.Split(cur.ILimiter, "&")
		if strings.Contains(parts[2], id) {
			obts := strings.Split(parts[2], ";")
			obts = filter(obts, id)
			il = fmt.Sprintf("%s&%s&%s", parts[0], parts[1], strings.Join(obts, ";"))
		}

		parts = strings.Split(cur.OLimiter, "&")
		if strings.Contains(parts[2], id) {
			obts := strings.Split(parts[2], ";")
			obts = filter(obts, id)
			ol = fmt.Sprintf("%s&%s&%s", parts[0], parts[1], strings.Join(obts, ";"))
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
			Obturator:    obt,
			ILimiter:     il,
			OLimiter:     ol,
			Coating:      cur.Coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgService) DeleteMoun(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`mounting like '%%%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putg. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
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

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: cur.Construction,
			Temperature:  cur.Temperatures,
			Obturator:    cur.Obturator,
			ILimiter:     cur.ILimiter,
			OLimiter:     cur.OLimiter,
			Coating:      cur.Coating,
			Mounting:     mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgService) DeleteCoating(id string) error {
	var wg sync.WaitGroup
	putg, err := s.repo.GetByCondition(fmt.Sprintf(`coating like '%%%s%%'`, id))
	if err != nil {
		return fmt.Errorf("failed to get putg. error: %w", err)
	}

	limit := make(chan struct{}, 30)
	for _, cur := range putg {
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

		upPutg := models.UpdateAdditDTO{
			Id:           cur.Id,
			Construction: cur.Construction,
			Temperature:  cur.Temperatures,
			Obturator:    cur.Obturator,
			ILimiter:     cur.ILimiter,
			OLimiter:     cur.OLimiter,
			Coating:      coating,
			Mounting:     cur.Mounting,
			Graphite:     cur.Graphite,
		}

		go s.putgUpdate(upPutg, &wg, limit)
	}

	wg.Wait()
	return nil
}

func (s *PutgService) putgUpdate(putg models.UpdateAdditDTO, wg *sync.WaitGroup, limit chan struct{}) error {
	defer wg.Done()
	if err := s.repo.UpdateAddit(putg); err != nil {
		return fmt.Errorf("failed to update putg addit. error: %w", err)
	}

	<-limit
	return nil
}
