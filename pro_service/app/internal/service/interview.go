package service

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/email_api"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro_api"
	docx "github.com/lukasjarosch/go-docx"
)

type InterviewService struct {
	email email_api.EmailServiceClient
	file  file_api.FileServiceClient
}

func NewInterviewService(email email_api.EmailServiceClient, file file_api.FileServiceClient) *InterviewService {
	return &InterviewService{
		email: email,
		file:  file,
	}
}

func (s *InterviewService) SendInterview(ctx context.Context, req *pro_api.SendInterviewRequest) error {
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
		"_drawing_":      req.DrawingNumber,
		"_info_":         req.Info,
		// "_drawingNumber_": req.DrawingNumber,
		// "_drawing_":       req.Drawing.OrigName,
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
	err = doc.WriteToFile(path.Join("template", "Опрос.docx"))
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to save new docx file. err: %s", err)
	}

	stat, err := os.Stat(path.Join("template", "Опрос.docx"))
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to get stat docx file. err: %s", err)
	}

	ext := filepath.Ext(path.Join("template", "Опрос.docx"))
	names := make([]string, 0, 2)
	names = append(names, stat.Name())
	size := stat.Size()
	buf := new(bytes.Buffer)

	if req.Drawing != nil {
		stream, err := s.file.Download(ctx, &file_api.FileDownloadRequest{
			Id:     req.Drawing.Id,
			Bucket: "pro",
			Group:  req.Drawing.Group,
			Name:   req.Drawing.OrigName,
		})
		if err != nil {
			logger.Errorf("failed to download drawing. err :%s", err.Error())
			return fmt.Errorf("failed to download drawing. err :%w", err)
		}

		_, err = stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to get data. err: %w", err)
		}

		imageData := bytes.Buffer{}

		for {
			logger.Debug("waiting to receive more data")

			req, err := stream.Recv()
			if err == io.EOF {
				logger.Debug("no more data")
				break
			}

			if err != nil {
				return fmt.Errorf("failed to get chunk. err %w", err)
			}

			chunk := req.GetFile().Content

			_, err = imageData.Write(chunk)
			if err != nil {
				return fmt.Errorf("failed to write chunk. err %w", err)
			}
		}

		names = append(names, req.Drawing.OrigName)

		w := zip.NewWriter(buf)
		f, err := w.Create(names[0])
		if err != nil {
			logger.Errorf("failed to create docx in zip. err %s", err.Error())
			return fmt.Errorf("failed to create docx in zip. err %w", err)
		}

		doc, err := os.ReadFile(path.Join("template", "Опрос.docx"))
		if err != nil {
			logger.Error(err)
			return fmt.Errorf("failed to read docx file. err: %s", err)
		}

		_, err = f.Write(doc)
		if err != nil {
			logger.Errorf("failed to write docx in zip. err %s", err.Error())
			return fmt.Errorf("failed to write docx in zip. err %w", err)
		}

		f, err = w.Create(names[1])
		if err != nil {
			logger.Errorf("failed to create image in zip. err %s", err.Error())
			return fmt.Errorf("failed to create image in zip. err %w", err)
		}

		_, err = f.Write(imageData.Bytes())
		if err != nil {
			logger.Errorf("failed to write image in zip. err %s", err.Error())
			return fmt.Errorf("failed to write image in zip. err %w", err)
		}

		if err := w.Close(); err != nil {
			return fmt.Errorf("failed to close writer. err %w", err)
		}

		size = int64(buf.Cap())
	}

	data := &email_api.SendInterviewRequest{
		Request: &email_api.SendInterviewRequest_Data{
			Data: &email_api.InterviewData{
				User: &email_api.User{
					Organization: req.Organization,
					Name:         req.Name,
					Email:        req.Email,
					Phone:        req.Phone,
					Position:     req.Position,
					City:         req.City,
				},
				File: &email_api.FileData{
					Name: names,
					Type: ext,
					Size: size,
				},
			},
		},
	}

	stream, err := s.email.SendInterview(ctx)
	if err != nil {
		return fmt.Errorf("error while connect wuth service. err: %w", err)
	}

	err = stream.Send(data)
	if err != nil {
		logger.Errorf("cannot send docx info to server: %s %s", err.Error(), stream.RecvMsg(nil).Error())
		return fmt.Errorf("cannot send docx info to server. err: %w", err)
	}

	if req.Drawing != nil {
		if err := s.sendZip(stream, buf); err != nil {
			return err
		}
	} else {
		if err := s.sendDoc(stream); err != nil {
			return err
		}
	}

	return nil
}

func (s *InterviewService) sendDoc(stream email_api.EmailService_SendInterviewClient) error {
	file, err := os.Open(path.Join("template", "Опрос.docx"))
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("failed to read docx file. err: %s", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("cannot read chunk to buffer: %s", err.Error())
			return fmt.Errorf("cannot read chunk to buffer: %w", err)
		}

		reqChunk := &email_api.SendInterviewRequest{
			Request: &email_api.SendInterviewRequest_File{
				File: &email_api.File{
					Content: buffer[:n],
				},
			},
		}

		err = stream.Send(reqChunk)
		if err != nil {
			logger.Errorf("cannot send chunk to server: %s", err.Error())
			return fmt.Errorf("cannot send chunk to server: %w", err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("cannot receive response: %s", err.Error())
		return fmt.Errorf("cannot receive response: %w", err)
	}

	return nil
}

func (s *InterviewService) sendZip(stream email_api.EmailService_SendInterviewClient, buf *bytes.Buffer) error {
	reader := bufio.NewReader(buf)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("cannot read chunk to buffer: %s", err.Error())
			return fmt.Errorf("cannot read chunk to buffer: %w", err)
		}

		reqChunk := &email_api.SendInterviewRequest{
			Request: &email_api.SendInterviewRequest_File{
				File: &email_api.File{
					Content: buffer[:n],
				},
			},
		}

		err = stream.Send(reqChunk)
		if err != nil {
			logger.Errorf("cannot send chunk to server: %s", err.Error())
			return fmt.Errorf("cannot send chunk to server: %w", err)
		}
	}

	_, err := stream.CloseAndRecv()
	if err != nil {
		logger.Errorf("cannot receive response: %s", err.Error())
		return fmt.Errorf("cannot receive response: %w", err)
	}

	return nil
}
