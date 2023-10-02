package service

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/analytic_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type OrderServiceNew struct {
	repo     repository.OrderNew
	position Position
	zip      Zip
	userApi  user_api.UserServiceClient
	fileApi  file_api.FileServiceClient
}

func NewOrderService_New(repo repository.OrderNew, position Position, zip Zip, userApi user_api.UserServiceClient, fileApi file_api.FileServiceClient,
) *OrderServiceNew {
	return &OrderServiceNew{
		repo:     repo,
		position: position,
		zip:      zip,
		fileApi:  fileApi,
		userApi:  userApi,
	}
}

func (s *OrderServiceNew) Get(ctx context.Context, req *order_api.GetOrder) (*order_model.CurrentOrder, error) {
	o, err := s.repo.Get(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order is not exist")
		}
		return nil, fmt.Errorf("failed to get order. error: %w", err)
	}

	positions, err := s.position.GetFull(ctx, o.Id)
	if err != nil {
		return nil, err
	}

	order := &order_model.CurrentOrder{
		Id:        o.Id,
		Number:    o.Number,
		UserId:    o.UserId,
		Positions: positions,
		Info:      o.Info,
	}

	return order, nil
}

func (s *OrderServiceNew) GetCurrent(ctx context.Context, req *order_api.GetCurrentOrder) (*order_model.CurrentOrder, error) {
	empty := false
	order, err := s.repo.GetCurrent(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			empty = true
		} else {
			return nil, fmt.Errorf("failed to get current order. error: %w", err)
		}
	}
	if empty {
		id, err := s.Create(ctx, &order_api.CreateOrder{UserId: req.UserId, ManagerId: req.ManagerId})
		if err != nil {
			return nil, err
		}

		order = &order_model.CurrentOrder{Id: id}
	}

	positions, err := s.position.GetFull(ctx, order.Id)
	if err != nil {
		return nil, err
	}
	order.Positions = positions

	return order, nil
}

func (s *OrderServiceNew) GetForFile(ctx context.Context, req *order_api.GetOrder) (*models.PositionCount, *order_model.FullOrder, error) {
	order, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get order by id. error: %w", err)
	}

	count, positions, err := s.position.Get(ctx, order.Id)
	if err != nil {
		return nil, nil, err
	}

	order.Positions = positions

	return count, order, nil
}

// TODO больно уж здоровая функция получается. надо бы подумать как уменьшить или разбить на куски
// func (s *OrderServiceNew) GetFile(ctx context.Context, req *order_api.GetOrder) (*bytes.Buffer, string, error) {
// 	if err := s.SetStatus(ctx, &order_api.Status{Status: order_model.OrderStatus_work, OrderId: req.Id}); err != nil {
// 		return nil, "", err
// 	}

// 	count, order, err := s.GetForFile(ctx, req)
// 	if err != nil {
// 		return nil, "", err
// 	}

// 	mainColumn := []interface{}{"№", "Наименование", "Доп. информация", "Количество", "Цена", "Сумма", "Себестоимость", "Шаблон"}

// 	tempColumn := []interface{}{"№", "Наименование", "Доп. информация", "Количество", "Ед. изм.", "Цена", "Сумма", "Шаблон"}

// 	snpColumn := []interface{}{
// 		"№", "Наименование", "Д4", "Д3", "Д2", "Д1", "h", "материал внутр. кольца", "материал каркаса", "материал наполнителя",
// 		"материал нар. кольца", "Перемычка", "Отверстие", "Крепление", "Чертеж", "Себестоимость", "Цена", "Шаблон",
// 	}

// 	putgColumnBase := []interface{}{
// 		"№", "Наименование", "D нар.", "D вн.", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации", "Материал обтюраторов", "Слюда",
// 		"Ингибитор", "Разъемная прокладка", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон",
// 	}
// 	putgColumnRings := []interface{}{
// 		"№", "Наименование", "D нар. огр.", "D нар.", "D вн.", "D вн. огр.", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации",
// 		"Материал нар. огр. кольца", "Материал обтюраторов", "Материал вн. огр. кольца", "Слюда", "Ингибитор", "Разъемная прокладка", "Перемычка",
// 		"Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон",
// 	}
// 	putgColumnForms := []interface{}{
// 		"№", "Наименование", "D нар.", "D вн.", "Поле", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации", "Материал обтюраторов",
// 		"Слюда", "Ингибитор", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон",
// 	}

// 	// ringColumn := []interface{}{
// 	// 	"№", "Наименование", "Тип кольца", "Д нар.", "Д вн.", "Высота", "Плотность", "Вид обтюратора", "Мод. добавка",
// 	// 	"Чертеж", "Себестоимость", "Цена", "Шаблон",
// 	// }

// 	startMain := 1
// 	startAside := 11

// 	snpStart := 1
// 	putgStart := 4 + count.SnpCount

// 	snpCount := snpStart
// 	putgCount := putgStart

// 	units := "шт"

// 	file := excelize.NewFile()
// 	// сделать 2 лист, а первый переименовать в заявку
// 	orderSheet := file.GetSheetName(file.GetActiveSheetIndex())

// 	tempSheetIdx, err := file.NewSheet("для1С")
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to create new sheet. error: %w", err)
// 	}
// 	tempSheet := file.GetSheetName(tempSheetIdx)

// 	headerStyle, err := file.NewStyle(&excelize.Style{
// 		Fill: excelize.Fill{
// 			Type:    "pattern",
// 			Pattern: 1,
// 			Color:   []string{"d9d9d9"},
// 		},
// 		Alignment: &excelize.Alignment{
// 			Horizontal:     "center",
// 			Vertical:       "center",
// 			RelativeIndent: 1,
// 			ShrinkToFit:    true,
// 			Indent:         1,
// 			ReadingOrder:   0,
// 			WrapText:       true,
// 		},
// 		Border: []excelize.Border{
// 			{Type: "left", Color: "000000", Style: 7},
// 			{Type: "top", Color: "000000", Style: 7},
// 			{Type: "bottom", Color: "000000", Style: 7},
// 			{Type: "right", Color: "000000", Style: 7},
// 		},
// 	})
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to create header style. error: %w", err)
// 	}

// 	// стиль для наименования прокладки (и для доп. инфы)
// 	titleStyle, err := file.NewStyle(&excelize.Style{
// 		Border: []excelize.Border{
// 			{Type: "left", Color: "000000", Style: 7},
// 			{Type: "top", Color: "000000", Style: 7},
// 			{Type: "bottom", Color: "000000", Style: 7},
// 			{Type: "right", Color: "000000", Style: 7},
// 		},
// 	})
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to create title style. error: %w", err)
// 	}

// 	cellStyle, err := file.NewStyle(&excelize.Style{
// 		Alignment: &excelize.Alignment{
// 			Horizontal: "center",
// 		},
// 		Border: []excelize.Border{
// 			{Type: "left", Color: "000000", Style: 7},
// 			{Type: "top", Color: "000000", Style: 7},
// 			{Type: "bottom", Color: "000000", Style: 7},
// 			{Type: "right", Color: "000000", Style: 7},
// 		},
// 	})
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to create cell style. error: %w", err)
// 	}

// 	cell, err := excelize.CoordinatesToCellName(startMain, 1)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 	}

// 	// добавление заголовков для основной таблицы
// 	if err = file.SetSheetRow(orderSheet, cell, &mainColumn); err != nil {
// 		return nil, "", fmt.Errorf("failed to create header table. error: %w", err)
// 	}

// 	// добавление заголовков для таблицы на листе для 1с
// 	if err = file.SetSheetRow(tempSheet, cell, &tempColumn); err != nil {
// 		return nil, "", fmt.Errorf("failed to create temp header table. error: %w", err)
// 	}

// 	endCell, err := excelize.CoordinatesToCellName(startMain+len(mainColumn)-1, 1)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 	}
// 	// добавление стилей для основной таблицы
// 	err = file.SetCellStyle(orderSheet, cell, endCell, headerStyle)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 	}

// 	endCell, err = excelize.CoordinatesToCellName(startMain+len(tempColumn)-1, 1)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 	}
// 	// добавление стилей для таблицы на листе для 1с
// 	err = file.SetCellStyle(tempSheet, cell, endCell, headerStyle)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 	}

// 	// snp
// 	cell, err = excelize.CoordinatesToCellName(startAside, snpStart)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 	}

// 	// добавление заголовков для таблицы снп
// 	if err = file.SetSheetRow(orderSheet, cell, &snpColumn); err != nil {
// 		return nil, "", fmt.Errorf("failed to create snp header table. error: %w", err)
// 	}

// 	endCell, err = excelize.CoordinatesToCellName(startAside+len(snpColumn)-1, snpStart)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 	}

// 	// добавление стилей
// 	err = file.SetCellStyle(orderSheet, cell, endCell, headerStyle)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 	}

// 	putgRingCount := 0
// 	putgFormCount := 0

// 	for _, fp := range order.Positions {
// 		if fp.Type == position_model.PositionType_Putg {
// 			if fp.PutgData.Main.ConfigurationCode != "round" {
// 				putgFormCount++
// 			}
// 			if utf8.RuneCountInString(fp.PutgData.Material.ConstructionCode) == 3 {
// 				putgRingCount++
// 			}
// 		}
// 	}

// 	// путг
// 	cell, err = excelize.CoordinatesToCellName(startAside, putgStart)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 	}

// 	// добавление заголовков для таблицы путг (основная)
// 	if err = file.SetSheetRow(orderSheet, cell, &putgColumnBase); err != nil {
// 		return nil, "", fmt.Errorf("failed to create putg base header table. error: %w", err)
// 	}

// 	endCell, err = excelize.CoordinatesToCellName(startAside+len(putgColumnBase)-1, putgStart)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 	}

// 	// добавление стилей
// 	err = file.SetCellStyle(orderSheet, cell, endCell, headerStyle)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 	}

// 	if putgRingCount != 0 {
// 		start := putgStart + count.PutgCount - putgRingCount - putgFormCount + 3
// 		cell, err = excelize.CoordinatesToCellName(startAside, start)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 		}
// 		// добавление заголовков для таблицы путг (с кольцами)
// 		if err = file.SetSheetRow(orderSheet, cell, &putgColumnRings); err != nil {
// 			return nil, "", fmt.Errorf("failed to create putg ring header table. error: %w", err)
// 		}

// 		endCell, err = excelize.CoordinatesToCellName(startAside+len(putgColumnRings)-1, start)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 		}

// 		// добавление стилей
// 		err = file.SetCellStyle(orderSheet, cell, endCell, headerStyle)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 		}

// 		count.PutgCount += 3
// 		putgRingCount = start
// 	}
// 	if putgFormCount != 0 {
// 		start := putgStart + count.PutgCount - putgFormCount + 3
// 		cell, err = excelize.CoordinatesToCellName(startAside, start)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 		}
// 		// добавление заголовков для таблицы путг (формы отличные от круглой)
// 		if err = file.SetSheetRow(orderSheet, cell, &putgColumnForms); err != nil {
// 			return nil, "", fmt.Errorf("failed to create putg form header table. error: %w", err)
// 		}

// 		endCell, err = excelize.CoordinatesToCellName(startAside+len(putgColumnForms)-1, start)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 		}

// 		// добавление стилей
// 		err = file.SetCellStyle(orderSheet, cell, endCell, headerStyle)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 		}

// 		count.PutgCount += 3
// 		putgFormCount = start
// 	}

// 	drawings := []string{}

// 	mainTitle, err := excelize.ColumnNumberToName(startMain + 1)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}

// 	asideTitle, err := excelize.ColumnNumberToName(startAside + 1)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}

// 	err = file.SetColWidth(orderSheet, mainTitle, mainTitle, 30)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to set column width. error: %w", err)
// 	}

// 	err = file.SetColWidth(orderSheet, asideTitle, asideTitle, 30)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to set column width. error: %w", err)
// 	}

// 	// получение колонок для вставки формул
// 	countColumn, err := excelize.ColumnNumberToName(startMain + 3)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}
// 	priceColumn, err := excelize.ColumnNumberToName(startMain + 4)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}
// 	costColumn, err := excelize.ColumnNumberToName(startMain + 6)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}
// 	sumColumn, err := excelize.ColumnNumberToName(startMain + 5)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}
// 	templateColumn, err := excelize.ColumnNumberToName(startMain + 7)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}

// 	// получение колонок для вставки формул на листе 1с
// 	tempPriceColumn, err := excelize.ColumnNumberToName(startMain + 5)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}
// 	tempSumColumn, err := excelize.ColumnNumberToName(startMain + 6)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}
// 	tempTemplateColumn, err := excelize.ColumnNumberToName(startMain + 7)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 	}

// 	for i, p := range order.Positions {
// 		mainLine := []interface{}{p.Count, p.Title, p.Info, p.Amount}

// 		cell, err := excelize.CoordinatesToCellName(startMain, 2+i)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 		}
// 		// добавление основных данных
// 		if err = file.SetSheetRow(orderSheet, cell, &mainLine); err != nil {
// 			return nil, "", fmt.Errorf("failed to create main line. error: %w", err)
// 		}

// 		endCell, err := excelize.CoordinatesToCellName(startMain+len(mainColumn)-1, 2+i)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 		}

// 		// добавление стилей
// 		err = file.SetCellStyle(orderSheet, cell, endCell, cellStyle)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 		}
// 		// добавление стилей для наименования
// 		err = file.SetCellStyle(orderSheet, fmt.Sprintf("%s%d", mainTitle, 2+i), fmt.Sprintf("%s%d", mainTitle, 2+i), titleStyle)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 		}

// 		err = file.SetCellFormula(orderSheet, fmt.Sprintf("%s%d", sumColumn, i+2), fmt.Sprintf("=%s%d*%s%d", countColumn, i+2, priceColumn, i+2))
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
// 		}

// 		var cost, price, template string
// 		var line int

// 		// SNP position
// 		if p.Type == position_model.PositionType_Snp {
// 			snpData := p.SnpData
// 			snpThickness := snpData.Size.H
// 			if snpThickness == "" {
// 				snpThickness = snpData.Size.Another
// 			}
// 			jumper := ""
// 			if snpData.Design.HasJumper {
// 				jumper = fmt.Sprintf("%s/%s", snpData.Design.JumperCode, snpData.Design.JumperWidth)
// 			}
// 			hole := ""
// 			if snpData.Design.HasHole {
// 				hole = "есть"
// 			}
// 			mounting := ""
// 			if snpData.Design.HasMounting {
// 				mounting = snpData.Design.MountingCode
// 			}
// 			drawing := ""
// 			if snpData.Design.Drawing != "" {
// 				drawing = "есть"
// 				parts := strings.Split(snpData.Design.Drawing, "/")
// 				drawings = append(drawings, fmt.Sprintf("%d_%s", p.Count, parts[len(parts)-1]))
// 			}

// 			var d4, d1 interface{}

// 			if snpData.Size.D4 != "" {
// 				d4 = snpData.Size.D4
// 			} else {
// 				d4 = nil
// 			}
// 			if snpData.Size.D1 != "" {
// 				d1 = snpData.Size.D1
// 			} else {
// 				d1 = nil
// 			}

// 			snpLine := []interface{}{
// 				p.Count, p.Title,
// 				d4, snpData.Size.D3, snpData.Size.D2, d1, snpThickness,
// 				snpData.Material.InnerRingCode, snpData.Material.FrameCode, snpData.Material.FillerCode, snpData.Material.OuterRingCode,
// 				jumper, hole, mounting, drawing,
// 			}

// 			cell, err = excelize.CoordinatesToCellName(startAside, int(1+snpCount))
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 			}
// 			// добавление данных для снп
// 			if err = file.SetSheetRow(orderSheet, cell, &snpLine); err != nil {
// 				return nil, "", fmt.Errorf("failed to create snp line. error: %w", err)
// 			}

// 			endCell, err := excelize.CoordinatesToCellName(startAside+len(snpLine)+2, int(1+snpCount))
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 			}

// 			// добавление стилей
// 			err = file.SetCellStyle(orderSheet, cell, endCell, cellStyle)
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 			}

// 			template, err = excelize.ColumnNumberToName(startAside + len(snpLine) + 2)
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 			}
// 			price, err = excelize.ColumnNumberToName(startAside + len(snpLine) + 1)
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 			}
// 			cost, err = excelize.ColumnNumberToName(startAside + len(snpLine))
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 			}

// 			line = snpCount + 1
// 			snpCount++
// 		}

// 		// Putg position
// 		if p.Type == position_model.PositionType_Putg {
// 			putgData := p.PutgData
// 			lineCount := putgCount
// 			columnCount := len(putgColumnBase) - 1

// 			jumper := ""
// 			if putgData.Design.HasJumper {
// 				jumper = fmt.Sprintf("%s/%s", putgData.Design.JumperCode, putgData.Design.JumperWidth)
// 			}
// 			hole := ""
// 			if putgData.Design.HasHole {
// 				hole = "есть"
// 			}
// 			drawing := ""
// 			if putgData.Design.Drawing != "" {
// 				drawing = "есть"
// 				parts := strings.Split(putgData.Design.Drawing, "/")
// 				drawings = append(drawings, fmt.Sprintf("%d_%s", p.Count, parts[len(parts)-1]))
// 			}

// 			// TODO возможно стоит еще выводить форму прокладки

// 			// "Слюда", "Ингибитор", "Разъемная прокладка",
// 			mica := 0
// 			inhibitor := 0
// 			if strings.HasSuffix(putgData.Material.TypeCode, "5") {
// 				inhibitor = 1
// 			}
// 			removable := 0
// 			if putgData.Design.HasRemovable {
// 				removable = 1
// 			}

// 			reinforce := 1
// 			if strings.HasSuffix(putgData.Material.TypeCode, "05") || strings.HasSuffix(putgData.Material.TypeCode, "00") {
// 				reinforce = 0
// 			}
// 			construction := strings.TrimLeft(putgData.Material.ConstructionCode, "0")

// 			var field interface{}
// 			if putgData.Main.ConfigurationCode != "round" {
// 				lineCount = putgFormCount
// 				columnCount = len(putgColumnForms) - 1
// 				putgFormCount++

// 				if putgData.Size.UseDimensions {
// 					d4, err := strconv.ParseFloat(putgData.Size.D4, 64)
// 					if err != nil {
// 						return nil, "", fmt.Errorf("failed to parse d4. error: %w", err)
// 					}
// 					d3, err := strconv.ParseFloat(putgData.Size.D3, 64)
// 					if err != nil {
// 						return nil, "", fmt.Errorf("failed to parse d3. error: %w", err)
// 					}
// 					d2, err := strconv.ParseFloat(putgData.Size.D2, 64)
// 					if err != nil {
// 						return nil, "", fmt.Errorf("failed to parse d2. error: %w", err)
// 					}
// 					d1, err := strconv.ParseFloat(putgData.Size.D1, 64)
// 					if err != nil {
// 						return nil, "", fmt.Errorf("failed to parse d1. error: %w", err)
// 					}

// 					field = math.Max(d4-d3, d2-d1)
// 					putgData.Size.D3 = putgData.Size.D4
// 				} else {
// 					field = putgData.Size.D1
// 				}

// 				putgData.Size.D4 = ""
// 				putgData.Size.D1 = ""
// 			} else {
// 				field = nil

// 				if utf8.RuneCountInString(putgData.Material.ConstructionCode) == 3 {
// 					lineCount = putgRingCount
// 					columnCount = len(putgColumnRings) - 1
// 					putgRingCount++
// 				} else {
// 					putgCount++
// 				}
// 			}

// 			var d4, d1 interface{}
// 			if putgData.Size.D4 != "" {
// 				d4 = putgData.Size.D4
// 			} else {
// 				d4 = nil
// 			}
// 			if putgData.Size.D1 != "" {
// 				d1 = putgData.Size.D1
// 			} else {
// 				d1 = nil
// 			}

// 			putgData.Size.H = strings.ReplaceAll(putgData.Size.H, ".", ",")

// 			//"№", "Наименование", "D нар.", "D вн.", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации", "Материал обтюраторов", "Слюда", "Ингибитор", "Разъемная прокладка", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена"

// 			putgLine := []interface{}{
// 				p.Count, p.Title,
// 				putgData.Size.D3, putgData.Size.D2, putgData.Size.H,
// 				construction, reinforce, putgData.Material.RotaryPlugCode,
// 				mica, inhibitor, removable, jumper, hole, drawing,
// 			}

// 			if utf8.RuneCountInString(putgData.Material.ConstructionCode) == 3 {
// 				//"№", "Наименование", "D нар. огр.", "D нар.", "D вн.", "D вн. огр.", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации", "Материал нар. огр. кольца", "Материал обтюраторов", "Материал вн. огр. кольца", "Слюда", "Ингибитор", "Разъемная прокладка", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена",

// 				putgLine = []interface{}{
// 					p.Count, p.Title,
// 					d4, putgData.Size.D3, putgData.Size.D2, d1, putgData.Size.H, construction, reinforce,
// 					putgData.Material.InnerRindCode, putgData.Material.RotaryPlugCode, putgData.Material.OuterRingCode,
// 					mica, inhibitor, removable, jumper, hole, drawing,
// 				}
// 			}

// 			if putgData.Main.ConfigurationCode != "round" {
// 				// "№", "Наименование", "D нар.", "D вн.", "Поле", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации", "Материал обтюраторов", "Слюда", "Ингибитор", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена",

// 				putgLine = []interface{}{
// 					p.Count, p.Title,
// 					putgData.Size.D3, putgData.Size.D2, field, putgData.Size.H, construction,
// 					reinforce, putgData.Material.RotaryPlugCode,
// 					mica, inhibitor, jumper, hole, drawing,
// 				}
// 			}

// 			cell, err = excelize.CoordinatesToCellName(startAside, int(1+lineCount))
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 			}

// 			// добавление данных для путг
// 			if err = file.SetSheetRow(orderSheet, cell, &putgLine); err != nil {
// 				return nil, "", fmt.Errorf("failed to create putg line. error: %w", err)
// 			}

// 			endCell, err := excelize.CoordinatesToCellName(startAside+columnCount, int(1+lineCount))
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 			}

// 			// добавление стилей
// 			err = file.SetCellStyle(orderSheet, cell, endCell, cellStyle)
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 			}

// 			template, err = excelize.ColumnNumberToName(startAside + columnCount)
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 			}
// 			price, err = excelize.ColumnNumberToName(startAside + columnCount - 1)
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 			}
// 			cost, err = excelize.ColumnNumberToName(startAside + columnCount - 2)
// 			if err != nil {
// 				return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
// 			}

// 			line = lineCount + 1
// 		}

// 		// добавление стилей для наименований
// 		err = file.SetCellStyle(orderSheet, fmt.Sprintf("%s%d", asideTitle, line), fmt.Sprintf("%s%d", asideTitle, line), titleStyle)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 		}

// 		err = file.SetCellFormula(orderSheet, fmt.Sprintf("%s%d", costColumn, i+2), fmt.Sprintf("=%s%d", cost, line))
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
// 		}
// 		err = file.SetCellFormula(orderSheet, fmt.Sprintf("%s%d", priceColumn, i+2), fmt.Sprintf("=%s%d", price, line))
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
// 		}
// 		err = file.SetCellFormula(orderSheet, fmt.Sprintf("%s%d", templateColumn, i+2), fmt.Sprintf("=%s%d", template, line))
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
// 		}

// 		// строка для 1с
// 		tempLine := []interface{}{p.Count, p.Title, p.Info, p.Amount, units}

// 		cell, err = excelize.CoordinatesToCellName(startMain, 2+i)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
// 		}
// 		// добавление данных для таблицы на листе для 1с
// 		if err = file.SetSheetRow(tempSheet, cell, &tempLine); err != nil {
// 			return nil, "", fmt.Errorf("failed to create main line. error: %w", err)
// 		}

// 		endCell, err = excelize.CoordinatesToCellName(startMain+len(tempColumn)-1, 2+i)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to get end cell. error: %w", err)
// 		}

// 		// добавление стилей для таблицы на листе для 1с
// 		err = file.SetCellStyle(tempSheet, cell, endCell, cellStyle)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell style. error: %w", err)
// 		}

// 		// добавление формул на листе 1с
// 		err = file.SetCellFormula(tempSheet, fmt.Sprintf("%s%d", tempPriceColumn, i+2), fmt.Sprintf("=%s!%s%d", orderSheet, priceColumn, i+2))
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
// 		}
// 		err = file.SetCellFormula(tempSheet, fmt.Sprintf("%s%d", tempSumColumn, i+2), fmt.Sprintf("=%s!%s%d", orderSheet, sumColumn, i+2))
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
// 		}
// 		err = file.SetCellFormula(tempSheet, fmt.Sprintf("%s%d", tempTemplateColumn, i+2), fmt.Sprintf("=%s!%s%d", orderSheet, templateColumn, i+2))
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
// 		}
// 	}

// 	fileName := fmt.Sprintf("Заявка %d", order.Number)
// 	var buffer *bytes.Buffer

// 	if len(drawings) > 0 {
// 		stream, err := s.fileApi.GroupDownload(ctx, &file_api.GroupDownloadRequest{
// 			Bucket: "pro",
// 			Group:  order.Id,
// 		})
// 		if err != nil {
// 			return nil, "", fmt.Errorf("failed to download drawing. err :%w", err)
// 		}

// 		res, err := stream.Recv()
// 		if err != nil && !strings.Contains(err.Error(), "file not found") {
// 			return nil, "", fmt.Errorf("failed to get data. err: %w", err)
// 		}
// 		meta := res.GetMetadata()
// 		fileData := bytes.Buffer{}

// 		if meta != nil {
// 			for {
// 				logger.Debug("waiting to receive more data")

// 				req, err := stream.Recv()
// 				if err == io.EOF {
// 					logger.Debug("no more data")
// 					break
// 				}
// 				if err != nil {
// 					return nil, "", fmt.Errorf("failed to get chunk. err %w", err)
// 				}

// 				chunk := req.GetFile().Content
// 				_, err = fileData.Write(chunk)
// 				if err != nil {
// 					return nil, "", fmt.Errorf("failed to write chunk. err %w", err)
// 				}
// 			}
// 		}

// 		buffer, err = s.zip.CreateWithDrawings(fileName+".xlsx", file, fileData, drawings)
// 		if err != nil {
// 			return nil, "", err
// 		}
// 	} else {
// 		buffer, err = s.zip.Create(fileName+".xlsx", file)
// 		if err != nil {
// 			return nil, "", err
// 		}
// 	}

//		return buffer, fileName, nil
//	}
func (s *OrderServiceNew) GetFile(ctx context.Context, req *order_api.GetOrder) (*bytes.Buffer, string, error) {
	if err := s.SetStatus(ctx, &order_api.Status{Status: order_model.OrderStatus_work, OrderId: req.Id}); err != nil {
		return nil, "", err
	}

	_, order, err := s.GetForFile(ctx, req)
	if err != nil {
		return nil, "", err
	}

	mainColumns := []interface{}{"№", "Наименование", "Доп. информация", "Количество", "Цена", "Сумма", "Себестоимость", "Шаблон"}
	templateColumns := []interface{}{"№", "Наименование", "Доп. информация", "Количество", "Ед. изм.", "Цена", "Сумма", "Шаблон"}

	snpColumns := []interface{}{
		"№", "Наименование", "Д4", "Д3", "Д2", "Д1", "h", "материал внутр. кольца", "материал каркаса", "материал наполнителя",
		"материал нар. кольца", "Перемычка", "Отверстие", "Крепление", "Чертеж", "Себестоимость", "Цена", "Шаблон",
	}

	putgColumns := []interface{}{
		"№", "Наименование", "D нар.", "D вн.", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации", "Материал обтюраторов", "Слюда",
		"Ингибитор", "Разъемная прокладка", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон",
	}
	putgWithRingsColumns := []interface{}{
		"№", "Наименование", "D нар. огр.", "D нар.", "D вн.", "D вн. огр.", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации",
		"Материал нар. огр. кольца", "Материал обтюраторов", "Материал вн. огр. кольца", "Слюда", "Ингибитор", "Разъемная прокладка", "Перемычка",
		"Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон",
	}
	putgNotRoundColumns := []interface{}{
		"№", "Наименование", "D нар.", "D вн.", "Поле", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации", "Материал обтюраторов",
		"Слюда", "Ингибитор", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон",
	}

	ringTwistedColumns := []interface{}{
		"№", "Наименование", "Тип кольца", "Д нар.", "Д вн.", "Высота", "Плотность", "Вид обтюратора", "Мод. добавка",
		"Чертеж", "Себестоимость", "Цена", "Шаблон",
	}
	ringPuffColumns := []interface{}{
		"№", "Наименование", "Тип кольца", "Д нар.", "Д вн.", "Высота", "Мод. добавка", "Чертеж", "Себестоимость", "Цена", "Шаблон",
	}
	ringWickerColumns := []interface{}{
		"№", "Наименование", "Тип кольца", "Д нар.", "Д вн.", "Высота", "Подпрессовка", "Наименование набивки", "Чертеж", "Себестоимость", "Цена", "Шаблон",
	}
	ringCompositeColumns := []interface{}{"№", "Наименование", "Тип кольца"}

	ringsKitColumns := []interface{}{"№", "Наименование", "Д нар.", "Д вн.", "Высота комплекта", "Чертеж", "Себестоимость", "Цена", "Шаблон"}

	file := excelize.NewFile()
	// сделать 2 лист, а первый переименовать в заявку
	orderSheet := file.GetSheetName(file.GetActiveSheetIndex())

	tempSheetIdx, err := file.NewSheet("для1С")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create new sheet. error: %w", err)
	}
	templateSheet := file.GetSheetName(tempSheetIdx)

	headerStyle, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"d9d9d9"},
		},
		Alignment: &excelize.Alignment{
			Horizontal:     "center",
			Vertical:       "center",
			RelativeIndent: 1,
			ShrinkToFit:    true,
			Indent:         1,
			ReadingOrder:   0,
			WrapText:       true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 7},
			{Type: "top", Color: "000000", Style: 7},
			{Type: "bottom", Color: "000000", Style: 7},
			{Type: "right", Color: "000000", Style: 7},
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create header style. error: %w", err)
	}

	// стиль для наименования прокладки (и для доп. инфы)
	titleStyle, err := file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 7},
			{Type: "top", Color: "000000", Style: 7},
			{Type: "bottom", Color: "000000", Style: 7},
			{Type: "right", Color: "000000", Style: 7},
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create title style. error: %w", err)
	}

	cellStyle, err := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 7},
			{Type: "top", Color: "000000", Style: 7},
			{Type: "bottom", Color: "000000", Style: 7},
			{Type: "right", Color: "000000", Style: 7},
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create cell style. error: %w", err)
	}

	row := 1
	col := 1
	asideCol := col + len(mainColumns) + 2
	asideRow := 1

	drawings := []string{}

	snp := make([]string, 0, 16)

	putg := make([]string, 0, 16)
	putgWithRings := make([]string, 0, 16)
	putgNotRound := make([]string, 0, 16)

	ringTwisted := make([]string, 0, 16)
	ringPuff := make([]string, 0, 16)
	ringWicker := make([]string, 0, 16)
	ringComposite := make([]string, 0, 16)

	ringsKit := make([]string, 0, 16)

	type Line struct {
		Index       int
		Row         int
		Line        []interface{}
		CostCol     string
		PriceCol    string
		TemplateCol string
	}
	lines := make(map[string]Line, 0)

	for i, p := range order.Positions {
		if p.Type == position_model.PositionType_Snp {
			snp = append(snp, p.Id)

			snpData := p.SnpData
			snpThickness := snpData.Size.H
			if snpThickness == "" {
				snpThickness = snpData.Size.Another
			}
			jumper := ""
			if snpData.Design.HasJumper {
				jumper = fmt.Sprintf("%s/%s", snpData.Design.JumperCode, snpData.Design.JumperWidth)
			}
			hole := ""
			if snpData.Design.HasHole {
				hole = "есть"
			}
			mounting := ""
			if snpData.Design.HasMounting {
				mounting = snpData.Design.MountingCode
			}
			drawing := ""
			if snpData.Design.Drawing != "" {
				drawing = "есть"
				parts := strings.Split(snpData.Design.Drawing, "/")
				drawings = append(drawings, fmt.Sprintf("%d_%s", p.Count, parts[len(parts)-1]))
			}

			var d4, d1 interface{}

			if snpData.Size.D4 != "" {
				d4 = snpData.Size.D4
			} else {
				d4 = nil
			}
			if snpData.Size.D1 != "" {
				d1 = snpData.Size.D1
			} else {
				d1 = nil
			}

			snpLine := []interface{}{
				p.Count, p.Title,
				d4, snpData.Size.D3, snpData.Size.D2, d1, snpThickness,
				snpData.Material.InnerRingCode, snpData.Material.FrameCode, snpData.Material.FillerCode, snpData.Material.OuterRingCode,
				jumper, hole, mounting, drawing, 0, 0, "",
			}

			lines[p.Id] = Line{Index: i, Line: snpLine}
		}
		if p.Type == position_model.PositionType_Putg {
			putgData := p.PutgData
			var line []interface{}

			jumper := ""
			if putgData.Design.HasJumper {
				jumper = fmt.Sprintf("%s/%s", putgData.Design.JumperCode, putgData.Design.JumperWidth)
			}
			hole := ""
			if putgData.Design.HasHole {
				hole = "есть"
			}
			drawing := ""
			if putgData.Design.Drawing != "" {
				drawing = "есть"
				parts := strings.Split(putgData.Design.Drawing, "/")
				drawings = append(drawings, fmt.Sprintf("%d_%s", p.Count, parts[len(parts)-1]))
			}

			// "Слюда", "Ингибитор", "Разъемная прокладка",
			mica := 0
			inhibitor := 0
			if strings.HasSuffix(putgData.Material.TypeCode, "5") {
				inhibitor = 1
			}
			removable := 0
			if putgData.Design.HasRemovable {
				removable = 1
			}

			reinforce := 1
			if strings.HasSuffix(putgData.Material.TypeCode, "05") || strings.HasSuffix(putgData.Material.TypeCode, "00") {
				reinforce = 0
			}
			construction := strings.TrimLeft(putgData.Material.ConstructionCode, "0")

			var field interface{}
			if putgData.Main.ConfigurationCode != "round" {
				if putgData.Size.UseDimensions {
					d4, err := strconv.ParseFloat(putgData.Size.D4, 64)
					if err != nil {
						return nil, "", fmt.Errorf("failed to parse d4. error: %w", err)
					}
					d3, err := strconv.ParseFloat(putgData.Size.D3, 64)
					if err != nil {
						return nil, "", fmt.Errorf("failed to parse d3. error: %w", err)
					}
					d2, err := strconv.ParseFloat(putgData.Size.D2, 64)
					if err != nil {
						return nil, "", fmt.Errorf("failed to parse d2. error: %w", err)
					}
					d1, err := strconv.ParseFloat(putgData.Size.D1, 64)
					if err != nil {
						return nil, "", fmt.Errorf("failed to parse d1. error: %w", err)
					}

					field = math.Max(d4-d3, d2-d1)
					putgData.Size.D3 = putgData.Size.D4
				} else {
					field = putgData.Size.D1
				}

				putgData.Size.D4 = ""
				putgData.Size.D1 = ""
			} else {
				field = nil
			}

			var d4, d1 interface{}
			if putgData.Size.D4 != "" {
				d4 = putgData.Size.D4
			} else {
				d4 = nil
			}
			if putgData.Size.D1 != "" {
				d1 = putgData.Size.D1
			} else {
				d1 = nil
			}

			putgData.Size.H = strings.ReplaceAll(putgData.Size.H, ".", ",")

			if putgData.Main.ConfigurationCode != "round" {
				// "№", "Наименование", "D нар.", "D вн.", "Поле", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации",
				//"Материал обтюраторов", "Слюда", "Ингибитор", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон"

				line = []interface{}{
					p.Count, p.Title,
					putgData.Size.D3, putgData.Size.D2, field, putgData.Size.H, construction,
					reinforce, putgData.Material.RotaryPlugCode,
					mica, inhibitor, jumper, hole, drawing,
					0, 0, "",
				}
				putgNotRound = append(putgNotRound, p.Id)
			} else if utf8.RuneCountInString(putgData.Material.ConstructionCode) == 3 {
				//"№", "Наименование", "D нар. огр.", "D нар.", "D вн.", "D вн. огр.", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации",
				//"Материал нар. огр. кольца", "Материал обтюраторов", "Материал вн. огр. кольца", "Слюда", "Ингибитор", "Разъемная прокладка",
				//"Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон"

				line = []interface{}{
					p.Count, p.Title,
					d4, putgData.Size.D3, putgData.Size.D2, d1, putgData.Size.H, construction, reinforce,
					putgData.Material.InnerRindCode, putgData.Material.RotaryPlugCode, putgData.Material.OuterRingCode,
					mica, inhibitor, removable, jumper, hole, drawing,
					0, 0, "",
				}
				putgWithRings = append(putgWithRings, p.Id)
			} else {
				//"№", "Наименование", "D нар.", "D вн.", "Толщина прокладки, мм", "Обтюраторы", "Материал перфорации", "Материал обтюраторов",
				//"Слюда", "Ингибитор", "Разъемная прокладка", "Перемычка", "Отверстие", "Чертеж", "Себестоимость", "Цена", "Шаблон"

				line = []interface{}{
					p.Count, p.Title,
					putgData.Size.D3, putgData.Size.D2, putgData.Size.H,
					construction, reinforce, putgData.Material.RotaryPlugCode,
					mica, inhibitor, removable, jumper, hole, drawing,
					0, 0, "",
				}
				putg = append(putg, p.Id)
			}

			lines[p.Id] = Line{Index: i, Line: line}
		}
		if p.Type == position_model.PositionType_Ring {
			ring := p.RingData
			var line []interface{}

			parts := strings.Split(ring.Size, "×")

			drawing := ""
			if ring.Drawing != "" {
				drawing = "есть"
				parts := strings.Split(ring.Drawing, "/")
				drawings = append(drawings, fmt.Sprintf("%d_%s", p.Count, parts[len(parts)-1]))
			}
			modifying := ring.Modifying
			if modifying == "" {
				modifying = "0"
			}

			thickness, mod, density, construction := "", "", "", ""
			if ring.Thickness != "" {
				thickness = "х" + ring.Thickness
			}
			if ring.Modifying != "" {
				mod = "-" + ring.Modifying
			}
			if ring.DensityCode != "" {
				density = "-" + ring.DensityCode
			}
			if ring.ConstructionCode != "" {
				construction = "-" + ring.ConstructionCode
			}

			designation := fmt.Sprintf("%sх%s%s %s%s%s-%s%s", parts[0], parts[1], thickness, ring.TypeCode, density, construction, ring.Material, mod)
			if ring.Drawing != "" {
				designation += " (черт.)"
			}

			if ring.TypeCode == "С" {
				line = []interface{}{p.Count, p.Title, ring.TypeCode, parts[0], parts[1], ring.Thickness, modifying, drawing, 0, 0, designation}
				ringPuff = append(ringPuff, p.Id)
			} else if ring.TypeCode == "П" {
				line = []interface{}{p.Count, p.Title, ring.TypeCode, parts[0], parts[1], ring.Thickness, 1, ring.Material, drawing, 0, 0, designation}
				ringWicker = append(ringWicker, p.Id)
			} else if ring.TypeCode == "К" {
				line = []interface{}{p.Count, p.Title, ring.TypeCode}
				ringComposite = append(ringComposite, p.Id)
			} else {
				density := ring.DensityCode
				if ring.TypeCode == "ВА" {
					density = "ВА-ГН"
				}
				line = []interface{}{
					p.Count, p.Title, ring.TypeCode, parts[0], parts[1], ring.Thickness, density, ring.ConstructionBaseCode, modifying,
					drawing, 0, 0, designation,
				}
				ringTwisted = append(ringTwisted, p.Id)
			}

			lines[p.Id] = Line{Index: i, Line: line}
		}
		if p.Type == position_model.PositionType_RingsKit {
			kit := p.KitData

			ringsKit = append(ringsKit, p.Id)

			parts := strings.Split(kit.Size, "×")

			drawing := ""
			if kit.Drawing != "" {
				drawing = "есть"
				parts := strings.Split(kit.Drawing, "/")
				drawings = append(drawings, fmt.Sprintf("%d_%s", p.Count, parts[len(parts)-1]))
			}

			thickness, modifying := "", ""
			if kit.Thickness != "" {
				thickness = "х" + kit.Thickness
			}
			if kit.Modifying != "" {
				modifying = "-" + kit.Modifying
			}

			designation := fmt.Sprintf("%sх%s%s %s-%s-%s-%s%s", parts[0], parts[1], thickness, kit.Type, kit.ConstructionCode,
				strings.ReplaceAll(kit.Count, "×", "*"), kit.Material, modifying,
			)
			if kit.Drawing != "" {
				designation += " (черт.)"
			}

			line := []interface{}{p.Count, p.Title, parts[0], parts[1], kit.Thickness, drawing, 0, 0, designation}
			lines[p.Id] = Line{Index: i, Line: line}
		}
	}

	var cost, price, template string
	if len(snp) > 0 {
		// вставка заголовков для снп
		if err := s.insertHeader(ctx, file, orderSheet, snpColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(snpColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(snpColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(snpColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных снп
	for _, v := range snp {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	if len(putg) > 0 {
		// отступ
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		asideRow += 2

		// вставка заголовков для круглых путг без колец
		if err := s.insertHeader(ctx, file, orderSheet, putgColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(putgColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(putgColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(putgColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных круглых путг без колец
	for _, v := range putg {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	if len(putgWithRings) > 0 {
		// отступ
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		asideRow += 2

		// вставка заголовков для круглых путг с кольцами
		if err := s.insertHeader(ctx, file, orderSheet, putgWithRingsColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(putgWithRingsColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(putgWithRingsColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(putgWithRingsColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных круглых путг с кольцами
	for _, v := range putgWithRings {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	if len(putgNotRound) > 0 {
		// отступ
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		asideRow += 2

		// вставка заголовков для не круглых путг
		if err := s.insertHeader(ctx, file, orderSheet, putgNotRoundColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(putgNotRoundColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(putgNotRoundColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(putgNotRoundColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных не круглых путг
	for _, v := range putgNotRound {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	if len(ringTwisted) > 0 {
		// отступ
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		asideRow += 2

		// вставка заголовков для витых колец
		if err := s.insertHeader(ctx, file, orderSheet, ringTwistedColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(ringTwistedColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(ringTwistedColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(ringTwistedColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных витых колец
	for _, v := range ringTwisted {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	if len(ringPuff) > 0 {
		// отступ
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		asideRow += 2

		// вставка заголовков для слоенных колец
		if err := s.insertHeader(ctx, file, orderSheet, ringPuffColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(ringPuffColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(ringPuffColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(ringPuffColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных слоенных колец
	for _, v := range ringPuff {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	if len(ringWicker) > 0 {
		// отступ
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		asideRow += 2

		// вставка заголовков для плетенных колец
		if err := s.insertHeader(ctx, file, orderSheet, ringWickerColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(ringWickerColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(ringWickerColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(ringWickerColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных плетенных колец
	for _, v := range ringWicker {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	if len(ringComposite) > 0 {
		// отступ
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		asideRow += 2

		// вставка заголовков для композиционных колец
		if err := s.insertHeader(ctx, file, orderSheet, ringCompositeColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(ringCompositeColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(ringCompositeColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(ringCompositeColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных композиционных колец
	for _, v := range ringComposite {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	if len(ringsKit) > 0 {
		// отступ
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		if err := s.insertData(ctx, file, orderSheet, []interface{}{}, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		asideRow += 2

		// вставка заголовков для комплектов колец
		if err := s.insertHeader(ctx, file, orderSheet, ringsKitColumns, asideCol, asideRow, headerStyle); err != nil {
			return nil, "", err
		}
		asideRow++

		template, err = excelize.ColumnNumberToName(asideCol + len(ringsKitColumns) - 1)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		price, err = excelize.ColumnNumberToName(asideCol + len(ringsKitColumns) - 2)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
		cost, err = excelize.ColumnNumberToName(asideCol + len(ringsKitColumns) - 3)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
		}
	}
	// вставка данных комплектов колец
	for _, v := range ringsKit {
		l := lines[v]
		if err := s.insertData(ctx, file, orderSheet, l.Line, asideCol, asideRow, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}
		lines[v] = Line{Index: l.Index, Row: asideRow, Line: l.Line, CostCol: cost, PriceCol: price, TemplateCol: template}
		asideRow++
	}

	// добавление заголовков для таблицы
	if err := s.insertHeader(ctx, file, orderSheet, mainColumns, col, row, headerStyle); err != nil {
		return nil, "", err
	}
	// добавление заголовков для таблицы на листе для 1с
	if err := s.insertHeader(ctx, file, templateSheet, templateColumns, col, row, headerStyle); err != nil {
		return nil, "", err
	}
	row++

	// получение колонок для вставки формул
	countColumn, err := excelize.ColumnNumberToName(col + 3)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
	}
	priceColumn, err := excelize.ColumnNumberToName(col + 4)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
	}
	costColumn, err := excelize.ColumnNumberToName(col + 6)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
	}
	sumColumn, err := excelize.ColumnNumberToName(col + 5)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
	}
	templateColumn, err := excelize.ColumnNumberToName(col + 7)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
	}

	// получение колонок для вставки формул на листе 1с
	tempPriceColumn, err := excelize.ColumnNumberToName(col + 5)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
	}
	tempSumColumn, err := excelize.ColumnNumberToName(col + 6)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
	}
	tempTemplateColumn, err := excelize.ColumnNumberToName(col + 7)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get column name. error: %w", err)
	}

	units := "шт"

	for _, p := range order.Positions {
		mainLine := []interface{}{p.Count, p.Title, p.Info, p.Amount, "", "", "", ""}

		// добавление данных для основной таблицы
		if err := s.insertData(ctx, file, orderSheet, mainLine, col, row, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}

		l := lines[p.Id]

		err = file.SetCellFormula(orderSheet, fmt.Sprintf("%s%d", sumColumn, row), fmt.Sprintf("=%s%d*%s%d", countColumn, row, priceColumn, row))
		if err != nil {
			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
		}
		err = file.SetCellFormula(orderSheet, fmt.Sprintf("%s%d", costColumn, row), fmt.Sprintf("=%s%d", l.CostCol, l.Row))
		if err != nil {
			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
		}
		err = file.SetCellFormula(orderSheet, fmt.Sprintf("%s%d", priceColumn, row), fmt.Sprintf("=%s%d", l.PriceCol, l.Row))
		if err != nil {
			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
		}
		err = file.SetCellFormula(orderSheet, fmt.Sprintf("%s%d", templateColumn, row), fmt.Sprintf("=%s%d", l.TemplateCol, l.Row))
		if err != nil {
			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
		}

		// строка для 1с
		tempLine := []interface{}{p.Count, p.Title, p.Info, p.Amount, units}
		// добавление данных для таблицы на листе для 1с
		if err := s.insertData(ctx, file, templateSheet, tempLine, col, row, cellStyle, titleStyle); err != nil {
			return nil, "", err
		}

		// добавление формул на листе 1с
		err = file.SetCellFormula(templateSheet, fmt.Sprintf("%s%d", tempPriceColumn, row), fmt.Sprintf("=%s!%s%d", orderSheet, priceColumn, row))
		if err != nil {
			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
		}
		err = file.SetCellFormula(templateSheet, fmt.Sprintf("%s%d", tempSumColumn, row), fmt.Sprintf("=%s!%s%d", orderSheet, sumColumn, row))
		if err != nil {
			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
		}
		err = file.SetCellFormula(templateSheet, fmt.Sprintf("%s%d", tempTemplateColumn, row), fmt.Sprintf("=%s!%s%d", orderSheet, templateColumn, row))
		if err != nil {
			return nil, "", fmt.Errorf("failed to set cell formula. error: %w", err)
		}

		row++
	}

	fileName := fmt.Sprintf("Заявка %d", order.Number)
	var buffer *bytes.Buffer

	if len(drawings) > 0 {
		stream, err := s.fileApi.GroupDownload(ctx, &file_api.GroupDownloadRequest{
			Bucket: "pro",
			Group:  order.Id,
		})
		if err != nil {
			return nil, "", fmt.Errorf("failed to download drawing. err :%w", err)
		}

		res, err := stream.Recv()
		if err != nil && !strings.Contains(err.Error(), "file not found") {
			return nil, "", fmt.Errorf("failed to get data. err: %w", err)
		}
		meta := res.GetMetadata()
		fileData := bytes.Buffer{}

		if meta != nil {
			for {
				logger.Debug("waiting to receive more data")

				req, err := stream.Recv()
				if err == io.EOF {
					logger.Debug("no more data")
					break
				}
				if err != nil {
					return nil, "", fmt.Errorf("failed to get chunk. err %w", err)
				}

				chunk := req.GetFile().Content
				_, err = fileData.Write(chunk)
				if err != nil {
					return nil, "", fmt.Errorf("failed to write chunk. err %w", err)
				}
			}
		}

		buffer, err = s.zip.CreateWithDrawings(fileName+".xlsx", file, fileData, drawings)
		if err != nil {
			return nil, "", err
		}
	} else {
		buffer, err = s.zip.Create(fileName+".xlsx", file)
		if err != nil {
			return nil, "", err
		}
	}

	return buffer, fileName, nil
}
func (s *OrderServiceNew) insertHeader(ctx context.Context, file *excelize.File, sheet string, columns []interface{}, col, row, style int) error {
	// TODO

	// получение координат ячейки
	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return fmt.Errorf("failed to get cell. error: %w", err)
	}

	// добавление заголовков таблицы
	if err = file.SetSheetRow(sheet, cell, &columns); err != nil {
		return fmt.Errorf("failed to create header table. error: %w", err)
	}

	// получение координат ячейки
	endCell, err := excelize.CoordinatesToCellName(col+len(columns)-1, row)
	if err != nil {
		return fmt.Errorf("failed to get end cell. error: %w", err)
	}

	// добавление стилей для таблицы
	err = file.SetCellStyle(sheet, cell, endCell, style)
	if err != nil {
		return fmt.Errorf("failed to set cell style. error: %w", err)
	}

	return nil
}
func (s *OrderServiceNew) insertData(ctx context.Context, file *excelize.File, sheet string, line []interface{}, col, row, style, titleStyle int) error {
	// получение координат ячейки
	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return fmt.Errorf("failed to get cell. error: %w", err)
	}

	// добавление данных
	if err = file.SetSheetRow(sheet, cell, &line); err != nil {
		return fmt.Errorf("failed to create main line. error: %w", err)
	}

	if len(line) > 0 {
		// получение координат ячейки
		endCell, err := excelize.CoordinatesToCellName(col+len(line)-1, row)
		if err != nil {
			return fmt.Errorf("failed to get end cell. error: %w", err)
		}

		// добавление стилей
		err = file.SetCellStyle(sheet, cell, endCell, style)
		if err != nil {
			return fmt.Errorf("failed to set cell style. error: %w", err)
		}

		// получение буквы ячейки
		mainTitle, err := excelize.ColumnNumberToName(col + 1)
		if err != nil {
			return fmt.Errorf("failed to get column name. error: %w", err)
		}

		// добавление стилей для наименования
		err = file.SetCellStyle(sheet, fmt.Sprintf("%s%d", mainTitle, row), fmt.Sprintf("%s%d", mainTitle, row), titleStyle)
		if err != nil {
			return fmt.Errorf("failed to set cell style. error: %w", err)
		}
	}

	return nil
}

func (s *OrderServiceNew) GetAll(ctx context.Context, req *order_api.GetAllOrders) ([]*order_model.Order, error) {
	orders, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all orders. error: %w", err)
	}
	return orders, nil
}

func (s *OrderServiceNew) GetOpen(ctx context.Context, req *order_api.GetManagerOrders) ([]*order_model.ManagerOrder, error) {
	orders, err := s.repo.GetOpen(ctx, req.ManagerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get open orders. error: %w", err)
	}
	return orders, nil
}

func (s *OrderServiceNew) GetByNumber(ctx context.Context, req *order_api.GetOrderByNumber) (*analytic_model.FullOrder, error) {
	order, err := s.repo.GetByNumber(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by number. error: %w", err)
	}
	return order, nil
}

func (s *OrderServiceNew) GetLast(ctx context.Context, req *order_api.GetLastOrders) ([]*analytic_model.FullOrder, error) {
	orders, err := s.repo.GetLast(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get last orders. error: %w", err)
	}
	return orders, nil
}

func (s *OrderServiceNew) GetAnalytics(ctx context.Context, req *order_api.GetOrderAnalytics) (*order_api.Analytics, error) {
	data, err := s.repo.GetAnalytics(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics order. error: %w", err)
	}

	a, err := s.position.GetAnalytics(ctx, req)
	if err != nil {
		return nil, err
	}

	userAnalytics, err := s.userApi.GetAnalytics(ctx, &user_api.GetUserAnalytics{PeriodAt: req.PeriodAt, PeriodEnd: req.PeriodEnd})
	if err != nil {
		return nil, err
	}

	//TODO добавить путг и кольца
	analytics := &order_api.Analytics{
		OrdersCount:        a.OrdersCount,
		UsersCountRegister: userAnalytics.UsersCountRegister,
		UserCountLink:      userAnalytics.UserCountLink,
		UserCount:          a.UserCount,
		PositionCount:      a.PositionCount,
		SnpPositionCount:   a.SnpPositionCount,
		NewUserCount:       userAnalytics.NewUserCount,
		NewUserCountLink:   userAnalytics.NewUserCountLink,
		Orders:             data,
	}

	return analytics, nil
}

func (s *OrderServiceNew) GetOrdersCount(ctx context.Context, req *order_api.GetOrderCountAnalytics) ([]*analytic_model.OrderCount, error) {
	orders, err := s.repo.GetOrdersCount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders count. error: %w", err)
	}
	return orders, nil
}

func (s *OrderServiceNew) GetFullAnalytics(ctx context.Context, req *order_api.GetFullOrderAnalytics) ([]*analytic_model.FullOrder, error) {
	data, err := s.repo.GetFullAnalytics(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order analytics. error: %w", err)
	}
	return data, nil
}

func (s *OrderServiceNew) Save(ctx context.Context, order *order_api.CreateOrder) (*order_api.OrderNumber, error) {
	var orderId = order.Id
	var number int64
	if orderId == "" {
		orderId = uuid.New().String()

		order.Id = orderId
		if err := s.repo.Create(ctx, order, fmt.Sprintf("%d", time.Now().UnixMilli())); err != nil {
			return nil, fmt.Errorf("failed to create order. error: %w", err)
		}
	} else {
		var err error
		number, err = s.repo.GetNumber(ctx, order, fmt.Sprintf("%d", time.Now().UnixMilli()))
		if err != nil {
			return nil, fmt.Errorf("failed to get order number. error: %w", err)
		}
	}

	if len(order.Positions) > 0 {
		if err := s.position.CreateSeveral(ctx, order.Positions, orderId); err != nil {
			return nil, err
		}
	}

	return &order_api.OrderNumber{Number: number}, nil
}

func (s *OrderServiceNew) Copy(ctx context.Context, order *order_api.CopyOrder) error {
	positions, err := s.position.GetAll(ctx, order.FromId)
	if err != nil {
		return fmt.Errorf("failed to get positions. error: %w", err)
	}

	// _, err = s.fileApi.CopyGroup(context.Background(), &file_api.CopyGroupRequest{
	// 	Bucket:   "pro",
	// 	Group:    order.FromId,
	// 	NewGroup: order.TargetId,
	// })
	// if err != nil {
	// 	return fmt.Errorf("failed to copy group files. error: %w", err)
	// }

	for i, fp := range positions {
		_, err := s.position.Copy(ctx, &position_api.CopyPosition{Id: fp.Id, Count: order.Count + int64(i), OrderId: order.TargetId, FromOrderId: order.FromId})
		if err != nil {
			return err
		}
	}
	return nil
}

// можно для заказа запоминать id менеджера, для более точной статистики и для того, чтобы можно было передать только один заказ, а не все заказы от данного клиента
func (s *OrderServiceNew) Create(ctx context.Context, order *order_api.CreateOrder) (string, error) {
	var orderId = order.Id
	if orderId == "" {
		orderId = uuid.New().String()

		order.Id = orderId
		if err := s.repo.Create(ctx, order, ""); err != nil {
			return "", fmt.Errorf("failed to create order. error: %w", err)
		}
	}

	if len(order.Positions) > 0 {
		if err := s.position.CreateSeveral(ctx, order.Positions, orderId); err != nil {
			return "", err
		}
	}

	return orderId, nil
}

func (s *OrderServiceNew) SetInfo(ctx context.Context, info *order_api.Info) error {
	if err := s.repo.SetInfo(ctx, info); err != nil {
		return fmt.Errorf("failed to set info. error: %w", err)
	}
	return nil
}

func (s *OrderServiceNew) SetStatus(ctx context.Context, status *order_api.Status) error {
	status.Date = fmt.Sprintf("%d", time.Now().UnixMilli())

	if err := s.repo.SetStatus(ctx, status); err != nil {
		return fmt.Errorf("failed to set status order. error: %w", err)
	}
	return nil
}

func (s *OrderServiceNew) SetManager(ctx context.Context, manager *order_api.Manager) error {
	if err := s.repo.SetManager(ctx, manager); err != nil {
		return fmt.Errorf("failed to set manager. error: %w", err)
	}
	return nil
}
