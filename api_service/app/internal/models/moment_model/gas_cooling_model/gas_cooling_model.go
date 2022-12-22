package gas_cooling_model

import (
	"strconv"

	"github.com/Alexander272/sealur_proto/api/moment/calc_api"
	"github.com/Alexander272/sealur_proto/api/moment/calc_api/gas_cooling_model"
	"github.com/Alexander272/sealur_proto/api/moment/models/device_model"
)

type Calc struct {
	Device         Device           `json:"device"`
	Factor         FinningFactor    `json:"factor"`
	Pressure       Pressure         `json:"pressure"`
	Section        SectionExecution `json:"section"`
	TubeCount      TubeCount        `json:"tubeCount"`
	NumberOfMoves  NumberOfMoves    `json:"numberOfMoves"`
	TubeLenght     TubeLenght       `json:"tubeLength"`
	TestPressure   string           `json:"testPressure"`
	TypeBolt       string           `json:"type"`
	Condition      string           `json:"condition"`
	Bolts          BoltsData        `json:"bolts"`
	Gasket         GasketFullData   `json:"gasket"`
	IsNeedFormulas bool             `json:"isNeedFormulas"`
}

type GasketFullData struct {
	GasketId  string     `json:"gasketId"`
	EnvId     string     `json:"envId"`
	Thickness string     `json:"thickness"`
	Name      NameGasket `json:"name"`
}

type BoltsData struct {
	Count    string       `json:"count"`
	BoltId   string       `json:"boltId"`
	MarkId   string       `json:"markId"`
	Material MaterialData `json:"material"`
	Diameter string       `json:"diameter"`
	Area     string       `json:"area"`
	Title    string       `json:"title"`
}

type MaterialData struct {
	Title       string `json:"title"`
	EpsilonAt20 string `json:"epsilonAt20"`
	SigmaAt20   string `json:"sigmaAt20"`
}

type Device struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type Pressure struct {
	Id    string  `json:"id"`
	Value float64 `json:"value"`
}

type TubeCount struct {
	Id    string `json:"id"`
	Value int32  `json:"value"`
}

type FinningFactor struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type SectionExecution struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type TubeLenght struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type NumberOfMoves struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type NameGasket struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func (f *Calc) Parse() (ex *calc_api.GasCoolingRequest, err error) {
	var testPressure float64
	if f.TestPressure != "" {
		testPressure, err = strconv.ParseFloat(f.TestPressure, 64)
		if err != nil {
			return nil, err
		}
	}

	typeBolt := gas_cooling_model.MainData_TypeBolt_value[f.TypeBolt]
	condition := gas_cooling_model.MainData_Condition_value[f.Condition]

	device := &device_model.Device{
		Id:    f.Device.Id,
		Title: f.Device.Title,
	}
	pressure := &device_model.Pressure{
		Id:    f.Pressure.Id,
		Value: f.Pressure.Value,
	}
	factor := &device_model.FinningFactor{
		Id:    f.Factor.Id,
		Value: f.Factor.Value,
	}
	section := &device_model.SectionExecution{
		Id:    f.Section.Id,
		Value: f.Section.Value,
	}
	tubeCount := &device_model.TubeCount{
		Id:    f.TubeCount.Id,
		Value: f.TubeCount.Value,
	}
	number := &device_model.NumberOfMoves{
		Id:    f.NumberOfMoves.Id,
		Value: f.NumberOfMoves.Value,
	}
	tubeLenght := &device_model.TubeLength{
		Id:    f.TubeLenght.Id,
		Value: f.TubeLenght.Value,
	}

	bolts, err := f.Bolts.Parse()
	if err != nil {
		return nil, err
	}

	gasket, err := f.Gasket.Parse()
	if err != nil {
		return nil, err
	}

	data := &gas_cooling_model.MainData{
		Device:        device,
		Factor:        factor,
		Pressure:      pressure,
		Section:       section,
		TubeCount:     tubeCount,
		NumberOfMoves: number,
		TubeLength:    tubeLenght,
		TestPressure:  testPressure,
		TypeBolt:      gas_cooling_model.MainData_TypeBolt(typeBolt),
		Condition:     gas_cooling_model.MainData_Condition(condition),
	}

	ex = &calc_api.GasCoolingRequest{
		IsNeedFormulas: f.IsNeedFormulas,
		Data:           data,
		Bolts:          bolts,
		Gasket:         gasket,
	}
	return ex, nil
}

func (m *MaterialData) Parse() (mat *gas_cooling_model.MaterialData, err error) {
	eAt20, err := strconv.ParseFloat(m.EpsilonAt20, 64)
	if err != nil {
		return nil, err
	}
	sAt20, err := strconv.ParseFloat(m.SigmaAt20, 64)
	if err != nil {
		return nil, err
	}

	mat = &gas_cooling_model.MaterialData{
		Title:       m.Title,
		EpsilonAt20: eAt20,
		SigmaAt20:   sAt20,
	}

	return mat, nil
}

func (b *BoltsData) Parse() (bolts *gas_cooling_model.BoltData, err error) {
	bolts = &gas_cooling_model.BoltData{
		MarkId: b.MarkId,
		BoltId: b.BoltId,
	}

	count, err := strconv.Atoi(b.Count)
	if err != nil {
		return nil, err
	}
	bolts.Count = int32(count)

	if b.BoltId == "another" {
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
		bolts.Material, err = b.Material.Parse()
		if err != nil {
			return nil, err
		}
	}

	return bolts, nil
}

func (g *GasketFullData) Parse() (gasket *gas_cooling_model.GasketData, err error) {
	thickness, err := strconv.ParseFloat(g.Thickness, 64)
	if err != nil {
		return nil, err
	}

	name := &device_model.NameGasket{
		Id:    g.Name.Id,
		Title: g.Name.Title,
	}

	gasket = &gas_cooling_model.GasketData{
		GasketId:   g.GasketId,
		EnvId:      g.EnvId,
		Thickness:  thickness,
		NameGasket: name,
	}

	return gasket, nil
}
