package service

import (
	"fmt"
	"strings"

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
