package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/xuri/excelize/v2"
)

type ZipService struct{}

func NewZipService() *ZipService {
	return &ZipService{}
}

func (s *ZipService) Create(fileName string, excel *excelize.File) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	writer := zip.NewWriter(buffer)

	fw, err := writer.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create xlsx in zip. err %w", err)
	}
	_, err = excel.WriteTo(fw)
	if err != nil {
		return nil, fmt.Errorf("failed to write xlsx in zip. err %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer. err %w", err)
	}

	return buffer, nil
}

func (s *ZipService) CreateWithDrawings(excelName string, excel *excelize.File, file bytes.Buffer, drawings []string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)

	writer := zip.NewWriter(buffer)

	fw, err := writer.Create(excelName)
	if err != nil {
		return nil, fmt.Errorf("failed to create xlsx in zip. err %w", err)
	}
	_, err = excel.WriteTo(fw)
	if err != nil {
		return nil, fmt.Errorf("failed to write xlsx in zip. err %w", err)
	}

	reader := bytes.NewReader(file.Bytes())
	zipReader, err := zip.NewReader(reader, reader.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to read zip. err %w", err)
	}

	for _, item := range zipReader.File {
		file, err := s.read(item)
		if err != nil {
			return nil, fmt.Errorf("failed to read file in zip. err %w", err)
		}

		fileName := ""
		for _, v := range drawings {
			if strings.Contains(v, strings.Split(item.FileHeader.Name, "_")[1]) {
				fileName = v
				break
			}
		}

		logger.Debug(fileName)

		fw, err := writer.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create drawing in zip. err %w", err)
		}
		_, err = fw.Write(file)
		if err != nil {
			return nil, fmt.Errorf("failed to write drawing in zip. err %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer. err %w", err)
	}

	return buffer, nil
}

func (s *ZipService) InsertDrawings(file bytes.Buffer, drawings []string, buffer *bytes.Buffer) (*bytes.Buffer, error) {
	reader := bytes.NewReader(file.Bytes())
	zipReader, err := zip.NewReader(reader, reader.Size())
	if err != nil {
		return nil, fmt.Errorf("failed to read zip. err %w", err)
	}

	writer := zip.NewWriter(buffer)

	for _, item := range zipReader.File {
		file, err := s.read(item)
		if err != nil {
			return nil, fmt.Errorf("failed to read file in zip. err %w", err)
		}

		fileName := ""
		for _, v := range drawings {
			if strings.Contains(v, strings.Split(item.FileHeader.Name, "_")[1]) {
				fileName = v
				break
			}
		}

		logger.Debug(fileName)

		fw, err := writer.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create drawing in zip. err %w", err)
		}
		_, err = fw.Write(file)
		if err != nil {
			return nil, fmt.Errorf("failed to write drawing in zip. err %w", err)
		}
		logger.Debug("drawing size ", item.FileInfo().Size())
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer. err %w", err)
	}

	return buffer, nil
}

func (s *ZipService) read(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file in zip. err %w", err)
	}
	defer f.Close()

	return io.ReadAll(f)
}
