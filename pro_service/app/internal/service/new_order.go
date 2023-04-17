package service

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Alexander272/sealur/pro_service/internal/repository"
	"github.com/Alexander272/sealur/pro_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/file_api"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/Alexander272/sealur_proto/api/pro/position_api"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type OrderServiceNew struct {
	repo     repository.OrderNew
	position Position
	zip      Zip
	fileApi  file_api.FileServiceClient
}

func NewOrderService_New(repo repository.OrderNew, position Position, zip Zip, fileApi file_api.FileServiceClient) *OrderServiceNew {
	return &OrderServiceNew{
		repo:     repo,
		position: position,
		zip:      zip,
		fileApi:  fileApi,
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

func (s *OrderServiceNew) GetForFile(ctx context.Context, req *order_api.GetOrder) (*order_model.FullOrder, error) {
	order, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id. error: %w", err)
	}

	positions, err := s.position.Get(ctx, order.Id)
	if err != nil {
		return nil, err
	}

	order.Positions = positions

	return order, nil
}
func (s *OrderServiceNew) GetFile(ctx context.Context, req *order_api.GetOrder) (*bytes.Buffer, string, error) {
	if err := s.SetStatus(ctx, &order_api.Status{Status: order_model.OrderStatus_work, OrderId: req.Id}); err != nil {
		return nil, "", err
	}

	order, err := s.GetForFile(ctx, req)
	if err != nil {
		return nil, "", err
	}

	mainColumn := []interface{}{"№", "Наименование", "Количество"}
	snpColumn := []interface{}{"№", "Наименование", "Д4", "Д3", "Д2", "Д1", "h", "материал внутр. кольца", "материал каркаса", "материал наполнителя", "материал нар. кольца", "Перемычка", "Отверстие", "Крепление", "Чертеж"}

	file := excelize.NewFile()
	sheetName := file.GetSheetName(file.GetActiveSheetIndex())

	cell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
	}
	if err = file.SetSheetRow(sheetName, cell, &mainColumn); err != nil {
		return nil, "", fmt.Errorf("failed to create header table. error: %w", err)
	}

	cell, err = excelize.CoordinatesToCellName(6, 1)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
	}
	if err = file.SetSheetRow(sheetName, cell, &snpColumn); err != nil {
		return nil, "", fmt.Errorf("failed to create snp header table. error: %w", err)
	}

	drawings := []string{}

	snpCount := 1

	for i, p := range order.Positions {
		mainLine := []interface{}{1 + i, p.Title, p.Amount}

		cell, err := excelize.CoordinatesToCellName(1, 2+i)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
		}
		if err = file.SetSheetRow(sheetName, cell, &mainLine); err != nil {
			return nil, "", fmt.Errorf("failed to create main line. error: %w", err)
		}

		if p.Type == position_model.PositionType_Snp {
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
				drawings = append(drawings, fmt.Sprintf("%d_%s", snpCount, parts[len(parts)-1]))
			}

			snpLine := []interface{}{
				snpCount, p.Title,
				snpData.Size.D4, snpData.Size.D3, snpData.Size.D2, snpData.Size.D1, snpThickness,
				snpData.Material.InnerRingCode, snpData.Material.FrameCode, snpData.Material.FillerCode, snpData.Material.OuterRingCode,
				jumper, hole, mounting, drawing,
			}

			cell, err = excelize.CoordinatesToCellName(6, int(1+snpCount))
			if err != nil {
				return nil, "", fmt.Errorf("failed to get cell. error: %w", err)
			}
			if err = file.SetSheetRow(sheetName, cell, &snpLine); err != nil {
				return nil, "", fmt.Errorf("failed to create snp line. error: %w", err)
			}
			snpCount++
		}

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
