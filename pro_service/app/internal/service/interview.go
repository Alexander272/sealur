package service

import (
	"fmt"
	"path"
	"strings"

	"github.com/Alexander272/sealur/pro_service/internal/transport/grpc/proto"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	docx "github.com/lukasjarosch/go-docx"
)

type InterviewService struct {
}

func NewInterviewService() *InterviewService {
	return &InterviewService{}
}

func (s *InterviewService) SendInterview(req *proto.SendInterviewRequest) error {
	var mounting, lubricant string
	if req.Mounting {
		mounting = "Да"
	} else {
		mounting = "Нет"
	}
	if req.Lubricant {
		lubricant = "Да"
	} else {
		lubricant = "Нет"
	}

	medium := make([]string, 0, 3)
	if req.Abrasive {
		medium = append(medium, "абразивная")
	}
	if req.Crystallized {
		medium = append(medium, "кристализуемая")
	}
	if req.Penetrating {
		medium = append(medium, "с высокой проникающей способностью")
	}

	//TODO отправляться должно менеджеру

	replaceMap := docx.PlaceholderMap{
		"_organization_": req.Organization,
		"_name_":         req.Name,
		"_email_":        req.Email,
		"_phone_":        req.Phone,
		"_position_":     req.Position,
		"_city_":         req.City,
		"_techprocess_":  req.Techprocess,
		"_equipment_":    req.Equipment,
		"_seal_":         req.Seal,
		"_consumer_":     req.Consumer,
		"_factory_":      req.Factory,
		"_developer_":    req.Developer,
		"_typeFl_":       req.TypeFl,
		"_flange_":       req.Flange,
		"_d1_":           req.Sizes.D1,
		"_d2_":           req.Sizes.D2,
		"_dUp_":          req.Sizes.DUp,
		"_d_":            req.Sizes.D,
		"_h1_":           req.Sizes.H1,
		"_h2_":           req.Sizes.H2,
		"_material_":     req.Material,
		"_along_":        req.Along,
		"_across_":       req.Across,
		"_nonFlatness_":  req.NonFlatness,
		"_mounting_":     mounting,
		"_boltMaterial_": req.BoltMaterial,
		"_bolt_":         req.Sizes.Bolt,
		"_countBolt_":    req.Sizes.CountBolt,
		"_lubricant_":    lubricant,
		"_diffFrom_":     req.DiffFrom,
		"_diffTo_":       req.DiffTo,
		"_tempWorkPipe_": req.TempWorkPipe,
		"_tempWorkAnn_":  req.TempWorkAnn,
		"_presWork_":     req.PresWork,
		"_presTest_":     req.PresTest,
		"_pressure_":     req.Pressure,
		"_presWorkPipe_": req.PresWorkPipe,
		"_presWorkAnn_":  req.PresWorkAnn,
		"_environ_":      req.Environ,
		"_environPipe_":  req.EnvironPipe,
		"_environAnn_":   req.EnvironAnn,
		"_medium_":       strings.Join(medium, ", "),
		"_condition_":    req.Condition,
		"_period_":       req.Period,
		"_drawing_":      "",
		"_info_":         req.Info,
	}

	if req.Type == "stand" {
		replaceMap["_dy_"] = req.Sizes.Dy
		replaceMap["_py_"] = req.Sizes.Py
	} else {
		replaceMap["_dIn_"] = req.Sizes.DIn
		replaceMap["_dOut_"] = req.Sizes.DOut
		replaceMap["_h_"] = req.Sizes.H
	}

	pathway := path.Join("template", "template_stand.docx")
	if req.Type == "not_stand" {
		pathway = path.Join("template", "template_notstand.docx")
	}

	// read and parse the template docx
	doc, err := docx.Open(pathway)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to read docx file. err: %s", err)
	}

	// replace the keys with values from replaceMap
	err = doc.ReplaceAll(replaceMap)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to replace values docx file. err: %s", err)
	}

	// write out a new file
	err = doc.WriteToFile("опрос.docx")
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to save new docx file. err: %s", err)
	}

	return nil
}
