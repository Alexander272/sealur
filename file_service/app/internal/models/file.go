package models

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type File struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Group       string `json:"group"`
	ContentType string `json:"contentType"`
	Bytes       []byte `json:"file"`
}

type CreateFileDTO struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	Group       string `json:"group"`
	ContentType string `json:"contentType"`
	Reader      io.Reader
}

func (d *CreateFileDTO) NormalizeName() {
	d.Name = strings.ReplaceAll(d.Name, " ", "_")

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	newName, _, _ := transform.String(t, d.Name)

	d.Name = newName
}

func NewFile(dto *CreateFileDTO) (*File, error) {
	bytes, err := ioutil.ReadAll(dto.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create file model. err: %w", err)
	}
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate file id. err: %w", err)
	}

	return &File{
		ID:          id.String(),
		Name:        dto.Name,
		Size:        dto.Size,
		Group:       dto.Group,
		ContentType: dto.ContentType,
		Bytes:       bytes,
	}, nil
}
