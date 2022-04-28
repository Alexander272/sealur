package service

import (
	"fmt"
	"strings"

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

					constrs = append(constrs, &proto.Constr{Short: short, Obturators: obt})
				}

				Temps = append(Temps, &proto.ConTemp{Temp: temp, Constructions: constrs})
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

		var reinforce, obturator, iLimiter, oLimiter = &proto.Materials{}, &proto.Materials{}, &proto.Materials{}, &proto.Materials{}
		tmp = strings.Split(d.Reinforce, "&")
		if len(tmp) > 1 {
			reinforce = &proto.Materials{Values: strings.Split(tmp[0], ";"), Default: tmp[1]}
		}
		tmp = strings.Split(d.Obturator, "&")
		if len(tmp) > 1 {
			obturator = &proto.Materials{Values: strings.Split(tmp[0], ";"), Default: tmp[1]}
		}
		tmp = strings.Split(d.ILimiter, "&")
		if len(tmp) > 1 {
			iLimiter = &proto.Materials{Values: strings.Split(tmp[0], ";"), Default: tmp[1]}
		}
		tmp = strings.Split(d.OLimiter, "&")
		if len(tmp) > 1 {
			oLimiter = &proto.Materials{Values: strings.Split(tmp[0], ";"), Default: tmp[1]}
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

	reinforce := fmt.Sprintf("%s&%s", strings.Join(dto.Reinforce.Values, ";"), dto.Reinforce.Default)
	if len(dto.Reinforce.Values) == 0 {
		reinforce = ""
	}
	obturator := fmt.Sprintf("%s&%s", strings.Join(dto.Obturator.Values, ";"), dto.Obturator.Default)
	if len(dto.Obturator.Values) == 0 {
		obturator = ""
	}
	iLimiter := fmt.Sprintf("%s&%s", strings.Join(dto.ILimiter.Values, ";"), dto.ILimiter.Default)
	if len(dto.ILimiter.Values) == 0 {
		iLimiter = ""
	}
	oLimiter := fmt.Sprintf("%s&%s", strings.Join(dto.OLimiter.Values, ";"), dto.OLimiter.Default)
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

	reinforce := fmt.Sprintf("%s&%s", strings.Join(dto.Reinforce.Values, ";"), dto.Reinforce.Default)
	if len(dto.Reinforce.Values) == 0 {
		reinforce = ""
	}
	obturator := fmt.Sprintf("%s&%s", strings.Join(dto.Obturator.Values, ";"), dto.Obturator.Default)
	if len(dto.Obturator.Values) == 0 {
		obturator = ""
	}
	iLimiter := fmt.Sprintf("%s&%s", strings.Join(dto.ILimiter.Values, ";"), dto.ILimiter.Default)
	if len(dto.ILimiter.Values) == 0 {
		iLimiter = ""
	}
	oLimiter := fmt.Sprintf("%s&%s", strings.Join(dto.OLimiter.Values, ";"), dto.OLimiter.Default)
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
		return fmt.Errorf("failed to delete snp. error: %w", err)
	}
	return nil
}
