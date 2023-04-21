package dev_cooling_model

import (
	"fmt"
	"strconv"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/dev_cooling_model"
)

type Calc struct {
	// Расчетное давление
	Pressure string `json:"pressure"`
	// Расчетная температура
	Temp string `json:"temp"`
	// Способ крепления труб
	Method string `json:"method"`
	// Тип соединения
	TypeBolt string `json:"typeBolt"`
	// Способ крепления труб в трубной решетке
	Mounting string `json:"mounting"`
	// Тип крепления труб в трубной решетке
	TypeMounting string `json:"typeMounting"`
	// Схема камеры аппарата воздушного охлаждения
	CameraDiagram string `json:"cameraDiagram"`
	// Схема размещения отверстий
	Layout string `json:"layout"`

	Cap            CapData        `json:"cap"`
	TubeSheet      TubeSheetData  `json:"tubeSheet"`
	Tube           TubeData       `json:"tube"`
	Bolts          BoltData       `json:"bolts"`
	Gasket         GasketFullData `json:"gasket"`
	IsNeedFormulas bool           `json:"isNeedFormulas"`
	Friction       string         `json:"friction"`
}

type MaterialData struct {
	Title     string `json:"title"`
	Epsilon   string `json:"epsilon"`
	SigmaAt20 string `json:"sigmaAt20"`
	Sigma     string `json:"sigma"`
}

type CapData struct {
	// s4 - Толщина донышка крышки
	BottomThick string `json:"bottomThick"`
	// s5 - Толщина стенки крышки в месте присоединения к фланцу
	WallThick string `json:"wallThick"`
	// s6 - Толщина фланца крышки
	FlangeThick string `json:"flangeThick"`
	// s7 - Толщина боковой стенки
	SideWallThick string `json:"sideWallThick"`
	// B0 - Внутренний размер камеры в поперечном направлении
	InnerSize string `json:"innerSize"`
	// B4 - Наружный размер камеры в поперечном направлении
	OuterSize string `json:"outerSize"`
	// H - Глубина камеры (крышки)
	Depth string `json:"depth"`
	// L0 - Внутренний размер камеры в продольном направлении
	L string `json:"L"`
	// φ - Коэффициент прочности сварного шва
	Strength string `json:"strength"`
	// cк - Прибавка на коррозию
	Corrosion string `json:"corrosion"`
	// R - Радиус гиба в углу крышки камеры
	Radius string `json:"radius"`
	// Id Материала крышки
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
}

type TubeSheetData struct {
	// s1 - Толщина трубной решетки в пределах зоны перфорации
	ZoneThick string `json:"zoneThick"`
	// s2 - Толщина трубной решетки в месте уплотнения
	PlaceThick string `json:"placeThick"`
	// s3 - Толщина трубной решетки вне зоны уплотнения
	OutZoneThick string `json:"outZoneThick"`
	// B1 - Ширина зоны решетки толщиной s1
	Width string `json:"width"`
	// t1 - Шаг отверстий под трубы в продольном направлении
	StepLong string `json:"stepLong"`
	// t2 - Шаг отверстий под трубы в поперечном направлении
	StepTrans string `json:"stepTrans"`
	// z - Число рядов труб в поперечном направлении
	Count string `json:"count"`
	// d0 - Диаметр трубных отверстий в решетках
	Diameter string `json:"diameter"`
	// ср - Прибавка на коррозию
	Corrosion string `json:"corrosion"`
	// Id Материала трубной решетки
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
}

type TubeData struct {
	// L - Длина труб
	Length string `json:"length"`
	// Lк - Приведенная длина труб при продольном изгибе
	ReducedLength string `json:"reducedLength"`
	// dТ - Наружный диаметр трубы
	Diameter string `json:"diameter"`
	// sT - Толщина стенки трубы
	Thickness string `json:"thickness"`
	// сT - Прибавка на коррозию
	Corrosion string `json:"corrosion"`
	// l0 - Глубина развальцовки
	Depth string `json:"depth"`
	// aT - Размер сварного шва приварки труб
	Size string `json:"size"`
	// Id Материала труб
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
}

type BoltData struct {
	// B3 - Расстояние между осями болтов/шпилек в поперечном направлении
	Distance string `json:"distance"`
	// n - Количество болтов/шпилек
	Count string `json:"count"`
	// Id болта
	BoltId string `json:"boltId"`
	// lб - Длина болта/шпильки между опорными поверхностями
	Lenght string `json:"lenght"`
	// Id Материала болтов
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
	// d - Наружный диаметр болта/шпильки
	Diameter string `json:"diameter"`
	// fб - Площадь болта/шпильки
	Area string `json:"area"`
}

type GasketFullData struct {
	// Id прокладки
	GasketId string `json:"gasketId"`
	// Id среды
	EnvId string `json:"envId"`
	// Толщина прокладки
	Thickness string `json:"thickness"`
	// bp - Ширина прокладки
	Width string `json:"width"`
	// L2 - Размер прокладки в продольном направлении
	SizeLong string `json:"sizeLong"`
	// B2 - Размер прокладки в поперечном направление
	SizeTrans string     `json:"sizeTrans"`
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

func (c *Calc) Parse() (*calc_api.DevCoolingRequest, error) {
	pressure, err := strconv.ParseFloat(c.Pressure, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pressure. error: %w", err)
	}
	temp, err := strconv.ParseFloat(c.Temp, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse temp. error: %w", err)
	}
	friction := 0.3
	if c.Friction != "" {
		friction, err = strconv.ParseFloat(c.Friction, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse friction. error: %w", err)
		}
	}

	typeBolt := calc_api.DevCoolingRequest_TypeBolt_value[c.TypeBolt]
	method := calc_api.DevCoolingRequest_MountingMethod_value[c.Method]
	mounting := calc_api.DevCoolingRequest_Mounting_value[c.Mounting]
	typeMounting := calc_api.DevCoolingRequest_TypeMounting_value[c.TypeMounting]
	cameraDiagram := calc_api.DevCoolingRequest_CameraDiagram_value[c.CameraDiagram]
	layout := calc_api.DevCoolingRequest_Layout_value[c.Layout]

	cap, err := c.Cap.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap. error: %w", err)
	}

	tubeSheet, err := c.TubeSheet.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet. error: %w", err)
	}

	tube, err := c.Tube.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube. error: %w", err)
	}

	bolts, err := c.Bolts.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolts. error: %w", err)
	}

	gasket, err := c.Gasket.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket. error: %w", err)
	}

	result := &calc_api.DevCoolingRequest{
		Pressure:       pressure,
		Temp:           temp,
		TypeBolt:       calc_api.DevCoolingRequest_TypeBolt(typeBolt),
		Method:         calc_api.DevCoolingRequest_MountingMethod(method),
		Mounting:       calc_api.DevCoolingRequest_Mounting(mounting),
		TypeMounting:   calc_api.DevCoolingRequest_TypeMounting(typeMounting),
		CameraDiagram:  calc_api.DevCoolingRequest_CameraDiagram(cameraDiagram),
		Layout:         calc_api.DevCoolingRequest_Layout(layout),
		Cap:            cap,
		TubeSheet:      tubeSheet,
		Tube:           tube,
		Bolts:          bolts,
		Gasket:         gasket,
		IsNeedFormulas: c.IsNeedFormulas,
		Friction:       friction,
	}

	return result, nil
}

func (m *MaterialData) Parse() (mat *dev_cooling_model.MaterialData, err error) {
	var e, sAt20, s float64
	if m.Epsilon != "" {
		e, err = strconv.ParseFloat(m.Epsilon, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse material epsilon. error: %w", err)
		}
	}
	if m.SigmaAt20 != "" {
		sAt20, err = strconv.ParseFloat(m.SigmaAt20, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse material sigma at 20. error: %w", err)
		}
	}
	if m.Sigma != "" {
		s, err = strconv.ParseFloat(m.Sigma, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse material sigma. error: %w", err)
		}
	}

	mat = &dev_cooling_model.MaterialData{
		Title:     m.Title,
		Epsilon:   e,
		SigmaAt20: sAt20,
		Sigma:     s,
	}

	return mat, nil
}

func (c *CapData) Parse() (result *dev_cooling_model.CapData, err error) {
	bottomThick, err := strconv.ParseFloat(c.BottomThick, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap bottom thick. error: %w", err)
	}

	var wallThick, l, radius float64
	if c.WallThick != "" {
		wallThick, err = strconv.ParseFloat(c.WallThick, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cap wall thick. error: %w", err)
		}
	}
	flangeThick, err := strconv.ParseFloat(c.FlangeThick, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap flange thick. error: %w", err)
	}
	sideWallThick, err := strconv.ParseFloat(c.SideWallThick, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap side wall thick. error: %w", err)
	}
	innerSize, err := strconv.ParseFloat(c.InnerSize, 64)
	if err != nil {
		return nil, err
	}
	outerSize, err := strconv.ParseFloat(c.OuterSize, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap outer size. error: %w", err)
	}
	depth, err := strconv.ParseFloat(c.Depth, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap depth. error: %w", err)
	}

	if c.L != "" {
		l, err = strconv.ParseFloat(c.L, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cap l. error: %w", err)
		}
	}
	if c.Radius != "" {
		radius, err = strconv.ParseFloat(c.Radius, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cap radius. error: %w", err)
		}
	}

	strength, err := strconv.ParseFloat(c.Strength, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap strength. error: %w", err)
	}
	corrosion, err := strconv.ParseFloat(c.Corrosion, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cap corrosion. error: %w", err)
	}

	var mat *dev_cooling_model.MaterialData
	if c.MarkId == "another" {
		mat, err = c.Material.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse cap material. error: %w", err)
		}
	}

	result = &dev_cooling_model.CapData{
		BottomThick:   bottomThick,
		WallThick:     wallThick,
		FlangeThick:   flangeThick,
		SideWallThick: sideWallThick,
		InnerSize:     innerSize,
		OuterSize:     outerSize,
		Depth:         depth,
		L:             l,
		Strength:      strength,
		Corrosion:     corrosion,
		Radius:        radius,
		MarkId:        c.MarkId,
		Material:      mat,
	}

	return result, nil
}

func (ts *TubeSheetData) Parse() (result *dev_cooling_model.TubeSheetData, err error) {
	ZoneThick, err := strconv.ParseFloat(ts.ZoneThick, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet zone thick. error: %w", err)
	}
	PlaceThick, err := strconv.ParseFloat(ts.PlaceThick, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet place thick. error: %w", err)
	}
	OutZoneThick, err := strconv.ParseFloat(ts.OutZoneThick, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet out zone thick. error: %w", err)
	}
	Width, err := strconv.ParseFloat(ts.Width, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet width. error: %w", err)
	}
	StepLong, err := strconv.ParseFloat(ts.StepLong, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet step long. error: %w", err)
	}
	StepTrans, err := strconv.ParseFloat(ts.StepTrans, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet step trans. error: %w", err)
	}
	Count, err := strconv.Atoi(ts.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet count. error: %w", err)
	}
	Diameter, err := strconv.ParseFloat(ts.Diameter, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet diameter. error: %w", err)
	}
	Corrosion, err := strconv.ParseFloat(ts.Corrosion, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube sheet corrosion. error: %w", err)
	}

	var mat *dev_cooling_model.MaterialData
	if ts.MarkId == "another" {
		mat, err = ts.Material.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse tube sheet materials. error: %w", err)
		}
	}

	result = &dev_cooling_model.TubeSheetData{
		ZoneThick:    ZoneThick,
		PlaceThick:   PlaceThick,
		OutZoneThick: OutZoneThick,
		Width:        Width,
		StepLong:     StepLong,
		StepTrans:    StepTrans,
		Count:        int32(Count),
		Diameter:     Diameter,
		Corrosion:    Corrosion,
		MarkId:       ts.MarkId,
		Material:     mat,
	}

	return result, nil
}

func (t *TubeData) Parse() (result *dev_cooling_model.TubeData, err error) {
	Length, err := strconv.ParseFloat(t.Length, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube length. error: %w", err)
	}
	ReducedLength, err := strconv.ParseFloat(t.ReducedLength, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube reduced length. error: %w", err)
	}
	Diameter, err := strconv.ParseFloat(t.Diameter, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube diameter. error: %w", err)
	}
	Thickness, err := strconv.ParseFloat(t.Thickness, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube thickness. error: %w", err)
	}
	Corrosion, err := strconv.ParseFloat(t.Corrosion, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tube corrosion. error: %w", err)
	}
	var Depth, Size float64
	if t.Depth != "" {
		Depth, err = strconv.ParseFloat(t.Depth, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tube depth. error: %w", err)
		}
	}
	if t.Size != "" {
		Size, err = strconv.ParseFloat(t.Size, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse tube size. error: %w", err)
		}
	}

	var mat *dev_cooling_model.MaterialData
	if t.MarkId == "another" {
		mat, err = t.Material.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse tube materials. error: %w", err)
		}
	}

	result = &dev_cooling_model.TubeData{
		Length:        Length,
		ReducedLength: ReducedLength,
		Diameter:      Diameter,
		Thickness:     Thickness,
		Corrosion:     Corrosion,
		Depth:         Depth,
		Size:          Size,
		MarkId:        t.MarkId,
		Material:      mat,
	}

	return result, nil
}

func (b *BoltData) Parse() (result *dev_cooling_model.BoltData, err error) {
	distance, err := strconv.ParseFloat(b.Distance, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolt distance. error: %w", err)
	}
	length, err := strconv.ParseFloat(b.Lenght, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolt length. error: %w", err)
	}
	count, err := strconv.Atoi(b.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bolt count. error: %w", err)
	}

	var diameter, area float64
	if b.BoltId == "another" {
		diameter, err = strconv.ParseFloat(b.Diameter, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt diameter. error: %w", err)
		}
		area, err = strconv.ParseFloat(b.Area, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt area. error: %w", err)
		}
	}

	var mat *dev_cooling_model.MaterialData
	if b.MarkId == "another" {
		mat, err = b.Material.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse bolt materials. error: %w", err)
		}
	}

	result = &dev_cooling_model.BoltData{
		Distance: distance,
		Lenght:   length,
		Count:    int32(count),
		BoltId:   b.BoltId,
		Diameter: diameter,
		Area:     area,
		MarkId:   b.MarkId,
		Material: mat,
	}

	return result, nil
}

func (g *GasketFullData) Parse() (gasket *dev_cooling_model.GasketData, err error) {
	thickness, err := strconv.ParseFloat(g.Thickness, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket thickness. error: %w", err)
	}
	width, err := strconv.ParseFloat(g.Width, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket width. error: %w", err)
	}
	sizeLong, err := strconv.ParseFloat(g.SizeLong, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket size long. error: %w", err)
	}
	sizeTrans, err := strconv.ParseFloat(g.SizeTrans, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gasket size trans. error: %w", err)
	}
	gasket = &dev_cooling_model.GasketData{
		GasketId:  g.GasketId,
		EnvId:     g.EnvId,
		Thickness: thickness,
		Width:     width,
		SizeLong:  sizeLong,
		SizeTrans: sizeTrans,
	}
	if g.GasketId == "another" {
		data, err := g.Data.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse gasket data. error: %w", err)
		}

		gasket.Data = data
	}

	return gasket, nil
}

func (g *GasketData) Parse() (data *dev_cooling_model.GasketData_Data, err error) {
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

	typeG := dev_cooling_model.GasketData_Type_value[g.TypeG]

	data = &dev_cooling_model.GasketData_Data{
		Title:           g.Title,
		Type:            dev_cooling_model.GasketData_Type(typeG),
		Qo:              q,
		M:               m,
		Compression:     compression,
		Epsilon:         e,
		PermissiblePres: permissiblePres,
	}

	return data, nil
}
