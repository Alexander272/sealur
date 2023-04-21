package flange_model

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/flange_model"
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
	Friction       string         `json:"friction"`
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
	BoltId   string       `json:"boltId"`
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

func (f *CalcFlange) NewFlange() (flange *calc_api.FlangeRequest, err error) {
	pressure, err := strconv.ParseFloat(f.Pressure, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pressure. error: %w", err)
	}
	axialForce, err := strconv.Atoi(f.AxialForce)
	if err != nil {
		return nil, fmt.Errorf("failed to parse axial force. error: %w", err)
	}
	bendingMoment, err := strconv.Atoi(f.BendingMoment)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bending moment. error: %w", err)
	}
	temp, err := strconv.ParseFloat(f.Temp, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse temp. error: %w", err)
	}
	friction := 0.3
	if f.Friction != "" {
		friction, err = strconv.ParseFloat(f.Friction, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse friction. error: %w", err)
		}
	}

	flanges := calc_api.FlangeRequest_Flanges_value[f.Flanges]
	typeB := calc_api.FlangeRequest_Type_value[f.TypeB]
	condition := calc_api.FlangeRequest_Condition_value[f.Condition]
	calculation := calc_api.FlangeRequest_Calcutation_value[f.Calculation]

	flangesData := []*flange_model.FlangeData{}
	flange1, err := f.FlangesData.First.NewFlangeData()
	if err != nil {
		return nil, fmt.Errorf("failed to parse flange 1. error: %w", err)
	}
	flangesData = append(flangesData, flange1)

	if !f.IsSameFlange {
		flange2, err := f.FlangesData.Second.NewFlangeData()
		if err != nil {
			return nil, fmt.Errorf("failed to parse flange 2. error: %w", err)
		}
		flangesData = append(flangesData, flange2)
	}

	bolts, err := f.Bolts.NewBolts()
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolts. error: %w", err)
	}

	gasket, err := f.Gasket.NewGasket()
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket. error: %w", err)
	}

	var embed *flange_model.EmbedData
	if f.IsEmbedded {
		embed, err = f.Embed.NewEmbed()
		if err != nil {
			return nil, fmt.Errorf("failed to parse embed. error: %w", err)
		}
	}

	var washer []*flange_model.WasherData
	if f.IsUseWasher {
		washer, err = f.Washer.NewWasher()
		if err != nil {
			return nil, fmt.Errorf("failed to parse washer. error: %w", err)
		}
	}

	flange = &calc_api.FlangeRequest{
		Pressure:       pressure,
		AxialForce:     int32(axialForce),
		BendingMoment:  int32(bendingMoment),
		Temp:           temp,
		IsWork:         f.IsWork,
		Flanges:        calc_api.FlangeRequest_Flanges(flanges),
		IsSameFlange:   f.IsSameFlange,
		IsEmbedded:     f.IsEmbedded,
		Type:           calc_api.FlangeRequest_Type(typeB),
		Condition:      calc_api.FlangeRequest_Condition(condition),
		Calculation:    calc_api.FlangeRequest_Calcutation(calculation),
		IsUseWasher:    f.IsUseWasher,
		IsNeedFormulas: f.IsNeedFormulas,
		FlangesData:    flangesData,
		Bolts:          bolts,
		Gasket:         gasket,
		Washer:         washer,
		Embed:          embed,
		Friction:       friction,
	}
	return flange, nil
}

func (f *Flanges) NewFlangeData() (flange *flange_model.FlangeData, err error) {
	var size *flange_model.FlangeData_Size
	if f.StandartId == "another" {
		size, err = f.Size.NewSize()
		if err != nil {
			return nil, fmt.Errorf("failed to parse size. error: %w", err)
		}
	}

	var mat, rMat *flange_model.MaterialData
	if f.MarkId == "another" {
		mat, err = f.Material.NewMaterial()
		if err != nil {
			return nil, fmt.Errorf("failed to parse material. error: %w", err)
		}
	}
	if f.RingMarkId == "another" {
		rMat, err = f.RingMaterial.NewMaterial()
		if err != nil {
			return nil, fmt.Errorf("failed to parse ring material. error: %w", err)
		}
	}

	typeF := flange_model.FlangeData_Type_value[f.TypeF]
	corrosion, err := strconv.ParseFloat(f.Corrosion, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse corrosion. error: %w", err)
	}

	var temp, b float64
	if f.Temp != "" {
		temp, err = strconv.ParseFloat(f.Temp, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse temp. error: %w", err)
		}
	}
	if f.B != "" {
		b, err = strconv.ParseFloat(f.B, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse b. error: %w", err)
		}
	}

	flange = &flange_model.FlangeData{
		Type:         flange_model.FlangeData_Type(typeF),
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

func (m *MaterialData) NewMaterial() (mat *flange_model.MaterialData, err error) {
	alpha, err := strconv.ParseFloat(m.AlphaF, 64)
	if err != nil {
		return nil, err
	}
	alpha = alpha * math.Pow(10, -6)

	var eAt20, e, sAt20, s float64
	if m.EpsilonAt20 != "" {
		eAt20, err = strconv.ParseFloat(m.EpsilonAt20, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse epsilon at 20. error: %w", err)
		}
	}

	if m.Epsilon != "" {
		e, err = strconv.ParseFloat(m.Epsilon, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse epsilon. error: %w", err)
		}
	}

	if m.SigmaAt20 != "" {
		sAt20, err = strconv.ParseFloat(m.SigmaAt20, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sigma at 20. error: %w", err)
		}
	}

	if m.Sigma != "" {
		s, err = strconv.ParseFloat(m.Sigma, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sigma. error: %w", err)
		}
	}

	mat = &flange_model.MaterialData{
		Title:       m.Title,
		AlphaF:      alpha,
		EpsilonAt20: eAt20,
		Epsilon:     e,
		SigmaAt20:   sAt20,
		Sigma:       s,
	}

	return mat, nil
}

func (s *FlangeSize) NewSize() (size *flange_model.FlangeData_Size, err error) {
	dOut, err := strconv.ParseFloat(s.DOut, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse size dOut. error: %w", err)
	}
	d, err := strconv.ParseFloat(s.D, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse size d. error: %w", err)
	}
	h, err := strconv.ParseFloat(s.H, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse size h. error: %w", err)
	}
	d6, err := strconv.ParseFloat(s.D6, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse size d6. error: %w", err)
	}
	s0, err := strconv.ParseFloat(s.S0, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse size s0. error: %w", err)
	}
	var s1, l, dnk, dk, ds, h0, hk float64
	if s.S1 != "" {
		s1, err = strconv.ParseFloat(s.S1, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse size s1. error: %w", err)
		}
	}
	if s.L != "" {
		l, err = strconv.ParseFloat(s.L, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse size l. error: %w", err)
		}
	}
	if s.Dnk != "" {
		dnk, err = strconv.ParseFloat(s.Dnk, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse size dnk. error: %w", err)
		}
	}
	if s.Dk != "" {
		dk, err = strconv.ParseFloat(s.Dk, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse size dk. error: %w", err)
		}
	}
	if s.Ds != "" {
		ds, err = strconv.ParseFloat(s.Ds, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse size ds. error: %w", err)
		}
	}
	if s.H0 != "" {
		h0, err = strconv.ParseFloat(s.H0, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse size h0. error: %w", err)
		}
	}
	if s.Hk != "" {
		hk, err = strconv.ParseFloat(s.Hk, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse size hk. error: %w", err)
		}
	}

	size = &flange_model.FlangeData_Size{
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

func (b *BoltsData) NewBolts() (bolts *flange_model.BoltData, err error) {
	bolts = &flange_model.BoltData{
		MarkId: b.MarkId,
		BoltId: b.BoltId,
		Title:  b.Title,
	}
	if b.Count != "" {
		count, err := strconv.Atoi(b.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt count. error: %w", err)
		}
		bolts.Count = int32(count)
	}
	if b.Temp != "" {
		temp, err := strconv.ParseFloat(b.Temp, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt temp. error: %w", err)
		}
		bolts.Temp = temp
	}
	if b.BoltId == "another" {
		diameter, err := strconv.ParseFloat(b.Diameter, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt diameter. error: %w", err)
		}
		area, err := strconv.ParseFloat(b.Area, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt area. error: %w", err)
		}
		bolts.Diameter = diameter
		bolts.Area = area
	}
	if b.MarkId == "another" {
		bolts.Material, err = b.Material.NewMaterial()
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt materials. error: %w", err)
		}
	}

	return bolts, nil
}

func (g *GasketFullData) NewGasket() (gasket *flange_model.GasketData, err error) {
	thickness, err := strconv.ParseFloat(g.Thickness, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket thickness. error: %w", err)
	}
	dOut, err := strconv.ParseFloat(g.DOut, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket dOut. error: %w", err)
	}
	dIn, err := strconv.ParseFloat(g.DIn, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket dIn. error: %w", err)
	}
	gasket = &flange_model.GasketData{
		GasketId:  g.GasketId,
		EnvId:     g.EnvId,
		Thickness: thickness,
		DOut:      dOut,
		DIn:       dIn,
	}
	if g.GasketId == "another" {
		data, err := g.Data.NewGasketData()
		if err != nil {
			return nil, fmt.Errorf("failed to parse gasket data. error: %w", err)
		}

		gasket.Data = data
	}

	return gasket, nil
}

func (g *GasketData) NewGasketData() (data *flange_model.GasketData_Data, err error) {
	q, err := strconv.ParseFloat(g.Qo, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket q. error: %w", err)
	}
	m, err := strconv.ParseFloat(g.M, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket m. error: %w", err)
	}
	compression, err := strconv.ParseFloat(g.Compression, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket compression. error: %w", err)
	}
	e, err := strconv.ParseFloat(g.Epsilon, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket epsilon. error: %w", err)
	}
	permissiblePres, err := strconv.ParseFloat(g.PermissiblePres, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket permissible pres. error: %w", err)
	}

	typeG := flange_model.GasketData_Type_value[g.TypeG]

	data = &flange_model.GasketData_Data{
		Title:           g.Title,
		Type:            flange_model.GasketData_Type(typeG),
		Qo:              q,
		M:               m,
		Compression:     compression,
		Epsilon:         e,
		PermissiblePres: permissiblePres,
	}

	return data, nil
}

func (w *WasherData) NewWasher() (washers []*flange_model.WasherData, err error) {
	thickness, err := strconv.ParseFloat(w.Thickness, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse washer thickness. error: %w", err)
	}

	var material *flange_model.MaterialData
	if w.First.MarkId == "another" {
		material, err = w.First.Material.NewMaterial()
		if err != nil {
			return nil, fmt.Errorf("failed to parse washer materials. error: %w", err)
		}
	}

	washers = append(washers, &flange_model.WasherData{
		MarkId:    w.First.MarkId,
		Thickness: thickness,
		Material:  material,
	})

	if w.Second != (WasherMaterial{}) {
		var material *flange_model.MaterialData
		if w.Second.MarkId == "another" {
			material, err = w.First.Material.NewMaterial()
			if err != nil {
				return nil, fmt.Errorf("failed to parse washer materials. error: %w", err)
			}
		}

		washers = append(washers, &flange_model.WasherData{
			MarkId:    w.Second.MarkId,
			Thickness: thickness,
			Material:  material,
		})
	}

	return washers, nil
}

func (e *EmbedData) NewEmbed() (embed *flange_model.EmbedData, err error) {
	thickness, err := strconv.ParseFloat(e.Thickness, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embed thickness. error: %w", err)
	}

	var material *flange_model.MaterialData
	if e.MarkId == "another" {
		material, err = e.Material.NewMaterial()
		if err != nil {
			return nil, fmt.Errorf("failed to parse embed materials. error: %w", err)
		}
	}

	embed = &flange_model.EmbedData{
		MarkId:    e.MarkId,
		Thickness: thickness,
		Material:  material,
	}
	return embed, nil
}
