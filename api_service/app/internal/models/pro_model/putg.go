package pro_model

import (
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro_api"
)

type PutgDTO struct {
	FlangeId     string                       `json:"flangeId" binding:"required"`
	TypeFlId     string                       `json:"typeFlId" binding:"required"`
	TypePr       string                       `json:"typePr" binding:"required"`
	Form         string                       `json:"form" binding:"required"`
	Construction []*pro_api.PutgConstructions `json:"construction" binding:"required"`
	Temperatures []*pro_api.PutgTemp          `json:"temperatures" binding:"required"`
	Reinforce    *pro_api.PutgMaterials       `json:"reinforce"`
	Obturator    *pro_api.PutgMaterials       `json:"obturator"`
	ILimiter     *pro_api.PutgMaterials       `json:"iLimiter"`
	OLimiter     *pro_api.PutgMaterials       `json:"oLimiter"`
	Coating      []string                     `json:"coating" binding:"required"`
	Mounting     []string                     `json:"mounting" binding:"required"`
	Graphite     []string                     `json:"graphite" binding:"required"`
}

type PutgData struct {
	Main     MainPutg     `json:"main"`
	Size     SizePutg     `json:"size"`
	Material MaterialPutg `json:"material"`
	Design   DesignPutg   `json:"design"`
}

func (s *PutgData) Parse() *position_model.PositionPutg {
	return &position_model.PositionPutg{
		Main:     s.Main.Parse(),
		Size:     s.Size.Parse(),
		Material: s.Material.Parse(),
		Design:   s.Design.Parse(),
	}
}

type PutgStandard struct {
	Id string `json:"id"`
}

type PutgFlangeType struct {
	Id string `json:"id"`
}

type PutgConfiguration struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

type MainPutg struct {
	Standard      PutgStandard      `json:"standard"`
	FlangeType    PutgFlangeType    `json:"flangeType"`
	Configuration PutgConfiguration `json:"configuration"`
}

func (m *MainPutg) Parse() *position_model.PositionPutg_Main {
	return &position_model.PositionPutg_Main{
		PutgStandardId:    m.Standard.Id,
		FlangeTypeId:      m.FlangeType.Id,
		ConfigurationId:   m.Configuration.Id,
		ConfigurationCode: m.Configuration.Code,
	}
}

type SizePutg struct {
	Dn            string `json:"dn"`
	DnMm          string `json:"dnMm"`
	Pn            Pn     `json:"pn"`
	D4            string `json:"d4"`
	D3            string `json:"d3"`
	D2            string `json:"d2"`
	D1            string `json:"d1"`
	H             string `json:"h"`
	Another       string `json:"another"`
	UseDimensions bool   `json:"useDimensions"`
}

func (s *SizePutg) Parse() *position_model.PositionPutg_Size {
	return &position_model.PositionPutg_Size{
		Dn:            s.Dn,
		DnMm:          s.DnMm,
		Pn:            s.Pn.Parse(),
		D4:            s.D4,
		D3:            s.D3,
		D2:            s.D2,
		D1:            s.D1,
		H:             s.H,
		Another:       s.Another,
		UseDimensions: s.UseDimensions,
	}
}

type PutgFiller struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

type PutgType struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

type PutgConstruction struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

type PutgMaterial struct {
	Id    string `json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
}

type MaterialPutg struct {
	Filler       PutgFiller       `json:"filler"`
	Type         PutgType         `json:"putgType"`
	Construction PutgConstruction `json:"construction"`
	RotaryPlug   PutgMaterial     `json:"rotaryPlug"`
	InnerRing    PutgMaterial     `json:"innerRing"`
	OuterRing    PutgMaterial     `json:"outerRing"`
}

func (s *MaterialPutg) Parse() *position_model.PositionPutg_Material {
	return &position_model.PositionPutg_Material{
		FillerId:         s.Filler.Id,
		FillerCode:       s.Filler.Code,
		TypeId:           s.Type.Id,
		TypeCode:         s.Type.Code,
		ConstructionId:   s.Construction.Id,
		ConstructionCode: s.Construction.Code,
		RotaryPlugId:     s.RotaryPlug.Id,
		RotaryPlugCode:   s.RotaryPlug.Code,
		RotaryPlugTitle:  s.RotaryPlug.Title,
		InnerRingId:      s.InnerRing.Id,
		InnerRindCode:    s.InnerRing.Code,
		InnerRingTitle:   s.InnerRing.Title,
		OuterRingId:      s.OuterRing.Id,
		OuterRingCode:    s.OuterRing.Code,
		OuterRingTitle:   s.OuterRing.Title,
	}
}

type DesignPutg struct {
	HasHole      bool     `json:"hasHole"`
	HasCoating   bool     `json:"hasCoating"`
	HasRemovable bool     `json:"hasRemovable"`
	Jumper       Jumper   `json:"jumper"`
	Mounting     Mounting `json:"mounting"`
	Drawing      string   `json:"drawing"`
}

func (s *DesignPutg) Parse() *position_model.PositionPutg_Design {
	return &position_model.PositionPutg_Design{
		HasJumper:    s.Jumper.HasJumper,
		JumperCode:   s.Jumper.Code,
		JumperWidth:  s.Jumper.Width,
		HasHole:      s.HasHole,
		HasCoating:   s.HasCoating,
		HasRemovable: s.HasRemovable,
		// HasMounting:  s.Mounting.HasMounting,
		// MountingCode: s.Mounting.Code,
		Drawing: s.Drawing,
	}
}
