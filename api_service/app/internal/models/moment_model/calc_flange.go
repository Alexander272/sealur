package moment_model

import (
	"math"
	"strconv"

	"github.com/Alexander272/sealur_proto/api/moment_api"
)

type CalcFlange struct {
	Pressure       string         `json:"pressure"`
	AxialForce     string         `json:"axialForce"`
	BendingMoment  string         `json:"bendingMoment"`
	Temp           string         `json:"temp"`
	IsWork         bool           `json:"isWork"`
	Flanges        string         `json:"flanges"`
	IsSameFlange   bool           `json:"isSameFlange"`
	IsEmbedded     bool           `json:"isEmbedded"`
	TypeB          string         `json:"type"`
	Condition      string         `json:"condition"`
	Calculation    string         `json:"calculation"`
	Gasket         GasketFullData `json:"gasket"`
	Bolts          BoltsData      `json:"bolts"`
	Embed          EmbedData      `json:"embed"`
	FlangesData    FlangeData     `json:"flangesData"`
	IsUseWasher    bool           `json:"isUseWasher"`
	Washer         WasherData     `json:"washer"`
	IsNeedFormulas bool           `json:"isNeedFormulas"`
}

type GasketFullData struct {
	GasketId  string     `json:"gasketId"`
	EnvId     string     `json:"envId"`
	Thickness string     `json:"thickness"`
	DOut      string     `json:"d_out"`
	DIn       string     `json:"d_in"`
	Data      GasketData `json:"data"`
}

type GasketData struct {
	Title           string `json:"title"`
	TypeG           string `json:"type"`
	Qo              string `json:"qo"`
	M               string `json:"m"`
	Compression     string `json:"compression"`
	Epsilon         string `json:"epsilon"`
	PermissiblePres string `json:"permissiblePres"`
}

type BoltsData struct {
	MarkId   string       `json:"markId"`
	Title    string       `json:"title"`
	Name     string       `json:"name"`
	Diameter string       `json:"diameter"`
	Area     string       `json:"area"`
	Count    string       `json:"count"`
	Temp     string       `json:"temp"`
	Material MaterialData `json:"material"`
}

type EmbedData struct {
	MarkId    string       `json:"markId"`
	Thickness string       `json:"thickness"`
	Material  MaterialData `json:"material"`
}

type FlangeData struct {
	First  Flanges `json:"first"`
	Second Flanges `json:"second"`
}

type WasherData struct {
	First     WasherMaterial `json:"first"`
	Second    WasherMaterial `json:"second"`
	Thickness string         `json:"thickness"`
}

type WasherMaterial struct {
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
}

type MaterialData struct {
	Title       string `json:"title"`
	AlphaF      string `json:"alphaF"`
	EpsilonAt20 string `json:"epsilonAt20"`
	Epsilon     string `json:"epsilon"`
	SigmaAt20   string `json:"sigmaAt20"`
	Sigma       string `json:"sigma"`
}

type Flanges struct {
	TypeF        string       `json:"type"`
	StandartId   string       `json:"standartId"`
	MarkId       string       `json:"markId"`
	Dy           string       `json:"dy"`
	Py           float64      `json:"py"`
	B            string       `json:"b"`
	Row          int32        `json:"row"`
	Corrosion    string       `json:"corrosion"`
	Temp         string       `json:"temp"`
	Size         FlangeSize   `json:"size"`
	Material     MaterialData `json:"material"`
	RingMarkId   string       `json:"ringMarkId"`
	RingMaterial MaterialData `json:"ringMaterial"`
}

type FlangeSize struct {
	DOut string `json:"dOut"`
	D    string `json:"d"`
	H    string `json:"h"`
	S0   string `json:"s0"`
	S1   string `json:"s1"`
	L    string `json:"l"`
	D6   string `json:"d6"`
	Dnk  string `json:"dnk"`
	Dk   string `json:"dk"`
	Ds   string `json:"ds"`
	H0   string `json:"h0"`
	Hk   string `json:"hk"`
}

func (f *CalcFlange) NewFlange() (flange *moment_api.CalcFlangeRequest, err error) {
	pressure, err := strconv.ParseFloat(f.Pressure, 64)
	if err != nil {
		return nil, err
	}
	axialForce, err := strconv.Atoi(f.AxialForce)
	if err != nil {
		return nil, err
	}
	bendingMoment, err := strconv.Atoi(f.BendingMoment)
	if err != nil {
		return nil, err
	}
	temp, err := strconv.ParseFloat(f.Temp, 64)
	if err != nil {
		return nil, err
	}

	flanges := moment_api.CalcFlangeRequest_Flanges_value[f.Flanges]
	typeB := moment_api.CalcFlangeRequest_Type_value[f.TypeB]
	condition := moment_api.CalcFlangeRequest_Condition_value[f.Condition]
	calculation := moment_api.CalcFlangeRequest_Calcutation_value[f.Calculation]

	flangesData := []*moment_api.FlangeData{}
	flange1, err := f.FlangesData.First.NewFlangeData()
	if err != nil {
		return nil, err
	}
	flangesData = append(flangesData, flange1)

	if !f.IsSameFlange {
		flange2, err := f.FlangesData.Second.NewFlangeData()
		if err != nil {
			return nil, err
		}
		flangesData = append(flangesData, flange2)
	}

	bolts, err := f.Bolts.NewBolts()
	if err != nil {
		return nil, err
	}

	gasket, err := f.Gasket.NewGasket()
	if err != nil {
		return nil, err
	}

	var embed *moment_api.EmbedData
	if f.IsEmbedded {
		embed, err = f.Embed.NewEmbed()
		if err != nil {
			return nil, err
		}
	}

	var washer []*moment_api.WasherData
	if f.IsUseWasher {
		washer, err = f.Washer.NewWasher()
		if err != nil {
			return nil, err
		}
	}

	flange = &moment_api.CalcFlangeRequest{
		Pressure:       pressure,
		AxialForce:     int32(axialForce),
		BendingMoment:  int32(bendingMoment),
		Temp:           temp,
		IsWork:         f.IsWork,
		Flanges:        moment_api.CalcFlangeRequest_Flanges(flanges),
		IsSameFlange:   f.IsSameFlange,
		IsEmbedded:     f.IsEmbedded,
		Type:           moment_api.CalcFlangeRequest_Type(typeB),
		Condition:      moment_api.CalcFlangeRequest_Condition(condition),
		Calculation:    moment_api.CalcFlangeRequest_Calcutation(calculation),
		IsUseWasher:    f.IsUseWasher,
		IsNeedFormulas: f.IsNeedFormulas,
		FlangesData:    flangesData,
		Bolts:          bolts,
		Gasket:         gasket,
		Washer:         washer,
		Embed:          embed,
	}
	return flange, nil
}

func (f *Flanges) NewFlangeData() (flange *moment_api.FlangeData, err error) {
	var size *moment_api.FlangeData_Size
	if f.StandartId == "another" {
		size, err = f.Size.NewSize()
		if err != nil {
			return nil, err
		}
	}

	var mat, rMat *moment_api.MaterialData
	if f.MarkId == "another" {
		mat, err = f.Material.NewMaterial()
		if err != nil {
			return nil, err
		}
	}
	if f.RingMarkId == "another" {
		rMat, err = f.RingMaterial.NewMaterial()
		if err != nil {
			return nil, err
		}
	}

	typeF := moment_api.FlangeData_Type_value[f.TypeF]
	corrosion, err := strconv.ParseFloat(f.Corrosion, 64)
	if err != nil {
		return nil, err
	}

	var temp, b float64
	if f.Temp != "" {
		temp, err = strconv.ParseFloat(f.Temp, 64)
		if err != nil {
			return nil, err
		}
	}
	if f.B != "" {
		b, err = strconv.ParseFloat(f.B, 64)
		if err != nil {
			return nil, err
		}
	}

	flange = &moment_api.FlangeData{
		Type:         moment_api.FlangeData_Type(typeF),
		StandartId:   f.StandartId,
		MarkId:       f.MarkId,
		Material:     mat,
		Dn:           f.Dy,
		Py:           f.Py,
		B:            b,
		Row:          f.Row,
		Corrosion:    corrosion,
		Size:         size,
		Temp:         temp,
		RingMarkId:   f.RingMarkId,
		RingMaterial: rMat,
	}

	return flange, nil
}

func (m *MaterialData) NewMaterial() (mat *moment_api.MaterialData, err error) {
	alpha, err := strconv.ParseFloat(m.AlphaF, 64)
	if err != nil {
		return nil, err
	}
	alpha = alpha * math.Pow(10, -6)
	eAt20, err := strconv.ParseFloat(m.EpsilonAt20, 64)
	if err != nil {
		return nil, err
	}
	e, err := strconv.ParseFloat(m.Epsilon, 64)
	if err != nil {
		return nil, err
	}
	sAt20, err := strconv.ParseFloat(m.SigmaAt20, 64)
	if err != nil {
		return nil, err
	}
	s, err := strconv.ParseFloat(m.Sigma, 64)
	if err != nil {
		return nil, err
	}

	mat = &moment_api.MaterialData{
		Title:       m.Title,
		AlphaF:      alpha,
		EpsilonAt20: eAt20,
		Epsilon:     e,
		SigmaAt20:   sAt20,
		Sigma:       s,
	}

	return mat, nil
}

func (s *FlangeSize) NewSize() (size *moment_api.FlangeData_Size, err error) {
	dOut, err := strconv.ParseFloat(s.DOut, 64)
	if err != nil {
		return nil, err
	}
	d, err := strconv.ParseFloat(s.D, 64)
	if err != nil {
		return nil, err
	}
	h, err := strconv.ParseFloat(s.H, 64)
	if err != nil {
		return nil, err
	}
	d6, err := strconv.ParseFloat(s.D6, 64)
	if err != nil {
		return nil, err
	}
	s0, err := strconv.ParseFloat(s.S0, 64)
	if err != nil {
		return nil, err
	}
	var s1, l, dnk, dk, ds, h0, hk float64
	if s.S1 != "" {
		s1, err = strconv.ParseFloat(s.S1, 64)
		if err != nil {
			return nil, err
		}
	}
	if s.L != "" {
		l, err = strconv.ParseFloat(s.L, 64)
		if err != nil {
			return nil, err
		}
	}
	if s.Dnk != "" {
		dnk, err = strconv.ParseFloat(s.Dnk, 64)
		if err != nil {
			return nil, err
		}
	}
	if s.Dk != "" {
		dk, err = strconv.ParseFloat(s.Dk, 64)
		if err != nil {
			return nil, err
		}
	}
	if s.Ds != "" {
		ds, err = strconv.ParseFloat(s.Ds, 64)
		if err != nil {
			return nil, err
		}
	}
	if s.H0 != "" {
		h0, err = strconv.ParseFloat(s.H0, 64)
		if err != nil {
			return nil, err
		}
	}
	if s.Hk != "" {
		hk, err = strconv.ParseFloat(s.Hk, 64)
		if err != nil {
			return nil, err
		}
	}

	size = &moment_api.FlangeData_Size{
		DOut: dOut,
		D:    d,
		H:    h,
		S0:   s0,
		S1:   s1,
		L:    l,
		D6:   d6,
		Dnk:  dnk,
		Dk:   dk,
		Ds:   ds,
		H0:   h0,
		Hk:   hk,
	}

	return size, nil
}

func (b *BoltsData) NewBolts() (bolts *moment_api.BoltData, err error) {
	bolts = &moment_api.BoltData{
		MarkId: b.MarkId,
		Name:   b.Name,
		Title:  b.Title,
	}
	if b.Count != "" {
		count, err := strconv.Atoi(b.Count)
		if err != nil {
			return nil, err
		}
		bolts.Count = int32(count)
	}
	if b.Temp != "" {
		temp, err := strconv.ParseFloat(b.Temp, 64)
		if err != nil {
			return nil, err
		}
		bolts.Temp = temp
	}
	if b.Name == "another" {
		diameter, err := strconv.ParseFloat(b.Diameter, 64)
		if err != nil {
			return nil, err
		}
		area, err := strconv.ParseFloat(b.Area, 64)
		if err != nil {
			return nil, err
		}
		bolts.Diameter = diameter
		bolts.Area = area
	}
	if b.MarkId == "another" {
		bolts.Material, err = b.Material.NewMaterial()
		if err != nil {
			return nil, err
		}
	}

	return bolts, nil
}

func (g *GasketFullData) NewGasket() (gasket *moment_api.GasketData, err error) {
	thickness, err := strconv.ParseFloat(g.Thickness, 64)
	if err != nil {
		return nil, err
	}
	dOut, err := strconv.ParseFloat(g.DOut, 64)
	if err != nil {
		return nil, err
	}
	dIn, err := strconv.ParseFloat(g.DIn, 64)
	if err != nil {
		return nil, err
	}
	gasket = &moment_api.GasketData{
		GasketId:  g.GasketId,
		EnvId:     g.EnvId,
		Thickness: thickness,
		DOut:      dOut,
		DIn:       dIn,
	}
	if g.GasketId == "another" {
		data, err := g.Data.NewGasketData()
		if err != nil {
			return nil, err
		}

		gasket.Data = data
	}

	return gasket, nil
}

func (g *GasketData) NewGasketData() (data *moment_api.GasketData_Data, err error) {
	q, err := strconv.ParseFloat(g.Qo, 64)
	if err != nil {
		return nil, err
	}
	m, err := strconv.ParseFloat(g.M, 64)
	if err != nil {
		return nil, err
	}
	compression, err := strconv.ParseFloat(g.Compression, 64)
	if err != nil {
		return nil, err
	}
	e, err := strconv.ParseFloat(g.Epsilon, 64)
	if err != nil {
		return nil, err
	}
	permissiblePres, err := strconv.ParseFloat(g.PermissiblePres, 64)
	if err != nil {
		return nil, err
	}

	typeG := moment_api.GasketData_Type_value[g.TypeG]

	data = &moment_api.GasketData_Data{
		Title:           g.Title,
		Type:            moment_api.GasketData_Type(typeG),
		Qo:              q,
		M:               m,
		Compression:     compression,
		Epsilon:         e,
		PermissiblePres: permissiblePres,
	}

	return data, nil
}

func (w *WasherData) NewWasher() (washers []*moment_api.WasherData, err error) {
	thickness, err := strconv.ParseFloat(w.Thickness, 64)
	if err != nil {
		return nil, err
	}

	var material *moment_api.MaterialData
	if w.First.MarkId == "another" {
		material, err = w.First.Material.NewMaterial()
		if err != nil {
			return nil, err
		}
	}

	washers = append(washers, &moment_api.WasherData{
		MarkId:    w.First.MarkId,
		Thickness: thickness,
		Material:  material,
	})

	if w.Second != (WasherMaterial{}) {
		var material *moment_api.MaterialData
		if w.Second.MarkId == "another" {
			material, err = w.First.Material.NewMaterial()
			if err != nil {
				return nil, err
			}
		}

		washers = append(washers, &moment_api.WasherData{
			MarkId:    w.Second.MarkId,
			Thickness: thickness,
			Material:  material,
		})
	}

	return washers, nil
}

func (e *EmbedData) NewEmbed() (embed *moment_api.EmbedData, err error) {
	thickness, err := strconv.ParseFloat(e.Thickness, 64)
	if err != nil {
		return nil, err
	}

	var material *moment_api.MaterialData
	if e.MarkId == "another" {
		material, err = e.Material.NewMaterial()
		if err != nil {
			return nil, err
		}
	}

	embed = &moment_api.EmbedData{
		MarkId:    e.MarkId,
		Thickness: thickness,
		Material:  material,
	}
	return embed, nil
}
