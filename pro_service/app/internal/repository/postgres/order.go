package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/sealur/pro_service/internal/models"
	"github.com/Alexander272/sealur_proto/api/pro/models/analytic_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/order_model"
	"github.com/Alexander272/sealur_proto/api/pro/models/position_model"
	"github.com/Alexander272/sealur_proto/api/pro/order_api"
	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Get(ctx context.Context, req *order_api.GetOrder) (order *order_model.FullOrder, err error) {
	var data models.ManagerOrder
	query := fmt.Sprintf("SELECT id, date, count_position, info, number, user_id FROM \"%s\" WHERE id=$1", OrderTable)

	if err := r.db.Get(&data, query, req.Id); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	order = &order_model.FullOrder{
		Id:            data.Id,
		Date:          data.Date,
		CountPosition: data.Count,
		Number:        data.Number,
		UserId:        data.UserId,
		Info:          data.Info,
	}

	return order, nil
}

func (r *OrderRepo) GetCurrent(ctx context.Context, req *order_api.GetCurrentOrder) (order *order_model.CurrentOrder, err error) {
	var data models.OrderNew
	query := fmt.Sprintf("SELECT id, number, info FROM \"%s\" WHERE user_id=$1 AND date=''", OrderTable)

	if err := r.db.Get(&data, query, req.UserId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	order = &order_model.CurrentOrder{
		Id:     data.Id,
		Number: data.Number,
		Info:   data.Info,
	}

	return order, nil
}

func (r *OrderRepo) GetAll(ctx context.Context, req *order_api.GetAllOrders) (orders []*order_model.Order, err error) {
	var data []models.OrderWithPosition
	query := fmt.Sprintf(`SELECT "%s".id, date, "%s".info, count_position, number, %s.id as position_id, title, amount, %s.count as position_count
		FROM "%s" INNER JOIN %s on order_id="%s".id WHERE user_id=$1 AND date != '' ORDER BY number DESC, position_count`,
		OrderTable, OrderTable, PositionTable, PositionTable, OrderTable, PositionTable, OrderTable,
	)

	if err := r.db.Select(&data, query, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i, o := range data {
		if i > 0 && o.Id == orders[len(orders)-1].Id {
			orders[len(orders)-1].Positions = append(orders[len(orders)-1].Positions, &position_model.Position{
				Id:     o.PositionId,
				Count:  o.PositionCount,
				Title:  o.Title,
				Amount: o.Amount,
			})
		} else {
			orders = append(orders, &order_model.Order{
				Id:            o.Id,
				Date:          o.Date,
				CountPosition: o.Count,
				Number:        o.Number,
				Info:          o.Info,
				Positions: []*position_model.Position{{
					Id:     o.PositionId,
					Count:  o.PositionCount,
					Title:  o.Title,
					Amount: o.Amount,
				}},
			})
		}

	}

	return orders, nil
}

func (r *OrderRepo) GetByNumber(ctx context.Context, order *order_api.GetOrderByNumber) (o *analytic_model.FullOrder, err error) {
	var data models.FullOrderAnalytics
	query := fmt.Sprintf(`SELECT id, user_id, manager_id, number, date, status,	(SELECT company FROM "%s" WHERE "%s".id=user_id) as user_company,
		(SELECT name FROM "%s" WHERE "%s".id=user_id) as user_name,
		(SELECT name FROM "%s" WHERE id="%s".manager_id) as manager FROM "%s" WHERE date!='' AND number=$1
		ORDER BY manager_id, user_id LIMIT 1`,
		UserTable, UserTable, UserTable, UserTable, UserTable, OrderTable, OrderTable,
	)

	if err := r.db.Get(&data, query, order.Number); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	o = &analytic_model.FullOrder{
		Id:      data.ManagerId,
		Manager: data.Manager,
		Clients: []*analytic_model.Client{{
			Id:      data.UserId,
			Company: data.Company,
			Name:    data.Name,
			Orders: []*analytic_model.ClientOrder{{
				Id:     data.Id,
				Number: data.Number,
				Date:   data.Date,
				Status: data.Status,
			}},
		}},
	}

	return o, nil
}

func (r *OrderRepo) GetNumber(ctx context.Context, order *order_api.CreateOrder, date string) (int64, error) {
	query := fmt.Sprintf(`UPDATE "%s" SET date=$1, count_position=$2 WHERE id=$3 RETURNING number`, OrderTable)

	row := r.db.QueryRow(query, date, order.Count, order.Id)
	if row.Err() != nil {
		return 0, fmt.Errorf("failed to execute query. error: %w", row.Err())
	}

	var number int64
	if err := row.Scan(&number); err != nil {
		return 0, fmt.Errorf("failed to scan result. error: %w", err)
	}
	return number, nil
}

func (r *OrderRepo) GetOpen(ctx context.Context, managerId string) (orders []*order_model.ManagerOrder, err error) {
	var data []models.ManagerOrder
	// query := fmt.Sprintf(`SELECT "%s".id, user_id, "%s".company, "%s".date, count_position, "number", status FROM "%s"
	// 	INNER JOIN "%s" ON "%s".id=user_id WHERE manager_id=$1 AND status != '%s' AND "%s".date !='' ORDER BY status, "%s".date`,
	// 	OrderTable, UserTable, OrderTable, OrderTable, UserTable, UserTable, order_model.OrderStatus_finish.String(), OrderTable, OrderTable,
	// )

	//
	query := fmt.Sprintf(`SELECT "%s".id, user_id, "%s".company, "%s".date, count_position, "number", status FROM "%s" 
		INNER JOIN "%s" ON "%s".id=user_id WHERE "%s".manager_id=$1 AND status != '%s' AND "%s".date !='' ORDER BY "%s".date`,
		OrderTable, UserTable, OrderTable, OrderTable, UserTable, UserTable, OrderTable, order_model.OrderStatus_finish.String(), OrderTable, OrderTable,
	)

	if err := r.db.Select(&data, query, managerId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, o := range data {
		status := order_model.OrderStatus_value[o.Status]

		orders = append(orders, &order_model.ManagerOrder{
			Id:            o.Id,
			Date:          o.Date,
			UserId:        o.UserId,
			Company:       o.Company,
			CountPosition: o.Count,
			Number:        o.Number,
			Status:        order_model.OrderStatus(status),
		})
	}

	return orders, nil
}

func (r *OrderRepo) GetAnalytics(ctx context.Context, req *order_api.GetOrderAnalytics) (orders []*analytic_model.Order, err error) {
	var data []models.OrderAnalytics
	// query := fmt.Sprintf(`SELECT distinct user_id, manager_id, COUNT(distinct number) as order_count, SUM(amount::integer) as position_count,
	// 	SUM(case when type = 'Snp' then amount::integer end) as position_snp_count
	// 	FROM "%s" INNER JOIN "%s" ON order_id="%s".id
	// 	WHERE date>=$1 AND date<=$2 GROUP BY user_id, manager_id ORDER BY manager_id`,
	// 	OrderTable, PositionTable, OrderTable,
	// )

	// ? Запрос с группировкой по менеджерам и названиям компаний (скорее всего надо добавить поле инн, тк названия могут одинаковые)
	// SELECT distinct manager_id, COUNT(distinct number) as order_count, SUM(amount::integer) as position_count,
	// 	COALESCE(SUM(case when type = 'Snp' then amount::integer end),0) as position_snp_count,
	// 	COALESCE(SUM(case when type = 'Putg' then amount::integer end),0) as position_putg_count,
	// 	COALESCE((SUM(case when type = 'Ring' then amount::integer end)),0) as position_ring_count,
	// 	COALESCE((SUM(case when type = 'RingsKit' then amount::integer end)),0) as position_kit_count,
	// 	(SELECT company FROM "user" as u WHERE u.id=user_id) as user_company,
	// 	/*(SELECT inn FROM "user" as u WHERE u.id=user_id) as user_inn,*/
	// 	(SELECT name FROM "user" as u WHERE u.id=o.manager_id) as manager
	// 	FROM "order" as o INNER JOIN "position" ON order_id=o.id
	// 	WHERE date>=$1 AND date<=$2 GROUP BY manager_id, user_company ORDER BY manager_id

	query := fmt.Sprintf(`SELECT distinct user_id, manager_id, COUNT(distinct number) as order_count, SUM(amount::integer) as position_count,
		COALESCE(SUM(case when type = 'Snp' then amount::integer end),0) as position_snp_count,
		COALESCE(SUM(case when type = 'Putg' then amount::integer end),0) as position_putg_count,
		COALESCE((SUM(case when type = 'Ring' then amount::integer end)),0) as position_ring_count,
		COALESCE((SUM(case when type = 'RingsKit' then amount::integer end)),0) as position_kit_count,
		(SELECT company FROM "%s" AS u WHERE u.id=user_id) as user_company,
		(SELECT name FROM "%s" AS u WHERE u.id=o.manager_id) as manager
		FROM "%s" AS o INNER JOIN "%s" ON order_id=o.id
		WHERE date>=$1 AND date<=$2 GROUP BY user_id, manager_id ORDER BY manager_id`,
		UserTable, UserTable, OrderTable, PositionTable,
	)

	if err := r.db.Select(&data, query, req.PeriodAt, req.PeriodEnd); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	//? можно запихивать все id в map а потом передавать в user_api для получения пользователей
	//? после пробегаться по массиву и подставлять данные в нужные места (не особо производительный вариант)
	//* или похрен сделать inner в таблицу с пользователями (не особо хорошая идея т.к. для пользователей создан отдельный сервис
	//* и может быть будет отдельная бд)

	for i, oa := range data {
		c := &analytic_model.Clients{
			Id:                oa.UserId,
			Name:              oa.Company,
			OrdersCount:       oa.OrderCount,
			PositionCount:     oa.PositionCount,
			SnpPositionCount:  oa.PositionSnpCount,
			PutgPositionCount: oa.PositionPutgCount,
			RingPositionCount: oa.PositionRingCount,
			KitPositionCount:  oa.PositionKitCount,
		}

		if i == 0 || orders[len(orders)-1].Id != oa.ManagerId {
			orders = append(orders, &analytic_model.Order{
				Id:      oa.ManagerId,
				Manager: oa.Manager,
				Clients: []*analytic_model.Clients{c},
			})
		} else {
			orders[len(orders)-1].Clients = append(orders[len(orders)-1].Clients, c)
		}
	}

	return orders, nil
}

func (r *OrderRepo) GetLast(ctx context.Context, order *order_api.GetLastOrders) (orders []*analytic_model.FullOrder, err error) {
	var data []models.FullOrderAnalytics
	query := fmt.Sprintf(`SELECT id, user_id, manager_id, number, date, status,	(SELECT company FROM "%s" WHERE "%s".id=user_id) as user_company,
		(SELECT name FROM "%s" WHERE "%s".id=user_id) as user_name,
		(SELECT name FROM "%s" WHERE id="%s".manager_id) as manager FROM "%s" WHERE date!=''
		ORDER BY date DESC LIMIT 20`,
		UserTable, UserTable, UserTable, UserTable, UserTable, OrderTable, OrderTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i, foa := range data {
		if i == 0 || orders[len(orders)-1].Id != foa.ManagerId {
			orders = append(orders, &analytic_model.FullOrder{
				Id:      foa.ManagerId,
				Manager: foa.Manager,
				Clients: []*analytic_model.Client{{
					Id:      foa.UserId,
					Company: foa.Company,
					Name:    foa.Name,
					Orders: []*analytic_model.ClientOrder{{
						Id:     foa.Id,
						Number: foa.Number,
						Date:   foa.Date,
						Status: foa.Status,
					}},
				}},
			})
		} else {
			if orders[len(orders)-1].Clients[len(orders[len(orders)-1].Clients)-1].Id != foa.UserId {
				orders[len(orders)-1].Clients = append(orders[len(orders)-1].Clients, &analytic_model.Client{
					Id:      foa.UserId,
					Company: foa.Company,
					Name:    foa.Name,
					Orders: []*analytic_model.ClientOrder{{
						Id:     foa.Id,
						Number: foa.Number,
						Date:   foa.Date,
						Status: foa.Status,
					}},
				})
			} else {
				orders[len(orders)-1].Clients[len(orders[len(orders)-1].Clients)-1].Orders =
					append(orders[len(orders)-1].Clients[len(orders[len(orders)-1].Clients)-1].Orders, &analytic_model.ClientOrder{
						Id:     foa.Id,
						Number: foa.Number,
						Date:   foa.Date,
						Status: foa.Status,
					})
			}
		}
	}

	return orders, nil
}

/*
SELECT user_id, name,company, inn, count("number")
	FROM public."order"
	left join "user" on "user".id=user_id
	WHERE "order".date != '' group by user_id, name, company, inn order by count desc
*/
// количество заявок по пользователям

/*
количество заявок по пользователям + количество позиций (всего и снп)
добавил среднее число позиций в заявке (думаю норм идея)
добавил число заявок в которых есть снп

SELECT user_id, company, name, count(DISTINCT o.id) as order_count,

	count(DISTINCT case when type = 'Snp' then o.id end) as order_snp_count,
	count(DISTINCT case when type = 'Putg' then o.id end) as order_putg_count,

	SUM(amount::integer) as position_count,
	COALESCE(SUM(case when type = 'Snp' then amount::integer end),0) as position_snp_count,
	COALESCE(SUM(case when type = 'Putg' then amount::integer end),0) as position_putg_count,

	(SUM(amount::integer)/count(DISTINCT o.id))::real as average_position,
	COALESCE((SUM(case when type = 'Snp' then amount::integer end)/count(DISTINCT case when type = 'Snp' then o.id end))::real,0) as average_snp_position,
	COALESCE((SUM(case when type = 'Putg' then amount::integer end)/count(DISTINCT case when type = 'Putg' then o.id end))::real,0) as average_putg_position

	FROM "order" AS o
	INNER JOIN "position" ON order_id=o.id
	INNER JOIN "user" AS u ON user_id=u.id
	WHERE o.date != '' GROUP BY user_id, company, name ORDER BY position_count DESC
*/
func (r *OrderRepo) GetOrdersCount(ctx context.Context, req *order_api.GetOrderCountAnalytics) (orders []*analytic_model.OrderCount, err error) {
	var data []models.OrderCount
	// query := fmt.Sprintf(`SELECT user_id, name, company, inn, count("number") FROM "%s"
	// 	LEFT JOIN "%s" ON "%s".id=user_id WHERE "%s".date != '' GROUP BY user_id, name, company, inn ORDER BY count DESC`,
	// 	OrderTable, UserTable, UserTable, OrderTable,
	// )

	query := fmt.Sprintf(`SELECT user_id, company, name, count(DISTINCT o.id) as order_count, 
		count(DISTINCT case when type = 'Snp' then o.id end) as order_snp_count,
		count(DISTINCT case when type = 'Putg' then o.id end) as order_putg_count,
		count(DISTINCT case when type = 'Ring' then o.id end) as order_ring_count,
		count(DISTINCT case when type = 'RingsKit' then o.id end) as order_kit_count,

		SUM(amount::integer) as position_count,	
		COALESCE(SUM(case when type = 'Snp' then amount::integer end),0) as position_snp_count,
		COALESCE(SUM(case when type = 'Putg' then amount::integer end),0) as position_putg_count,
		COALESCE(SUM(case when type = 'Ring' then amount::integer end),0) as position_ring_count,
		COALESCE(SUM(case when type = 'RingsKit' then amount::integer end),0) as position_kit_count,
		
		(SUM(amount::integer)/count(DISTINCT o.id))::real as average_position,
		COALESCE((SUM(case when type = 'Snp' then amount::integer end)/count(DISTINCT case when type = 'Snp' then o.id end))::real,0) as average_snp_position,
		COALESCE((SUM(case when type = 'Putg' then amount::integer end)/count(DISTINCT case when type = 'Putg' then o.id end))::real,0) as average_putg_position,
		COALESCE((SUM(case when type = 'Ring' then amount::integer end)/count(DISTINCT case when type = 'Ring' then o.id end))::real,0) as average_ring_position,
		COALESCE((SUM(case when type = 'RingsKit' then amount::integer end)/count(DISTINCT case when type = 'RingsKit' then o.id end))::real,0) as average_kit_position

		FROM "%s" AS o
		INNER JOIN "%s" ON order_id=o.id
		INNER JOIN "%s" AS u ON user_id=u.id
		WHERE o.date != '' GROUP BY user_id, company, name ORDER BY position_count DESC`,
		OrderTable, PositionTable, UserTable,
	)

	if err := r.db.Select(&data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for _, oc := range data {
		orders = append(orders, &analytic_model.OrderCount{
			Id:                  oc.UserId,
			Name:                oc.Name,
			Company:             oc.Company,
			OrderCount:          oc.OrderCount,
			SnpOrderCount:       oc.SnpOrderCount,
			PutgOrderCount:      oc.PutgOrderCount,
			RingOrderCount:      oc.RingOrderCount,
			KitOrderCount:       oc.KitOrderCount,
			PositionCount:       oc.PositionCount,
			SnpPositionCount:    oc.PositionSnpCount,
			PutgPositionCount:   oc.PositionPutgCount,
			RingPositionCount:   oc.PositionRingCount,
			KitPositionCount:    oc.PositionKitCount,
			AveragePosition:     oc.AveragePosition,
			AverageSnpPosition:  oc.AverageSnpPosition,
			AveragePutgPosition: oc.AveragePutgPosition,
			AverageRingPosition: oc.AverageRingPosition,
			AverageKitPosition:  oc.AverageKitPosition,
		})
	}

	return orders, nil
}

func (r *OrderRepo) GetFullAnalytics(ctx context.Context, req *order_api.GetFullOrderAnalytics) (orders []*analytic_model.FullOrder, err error) {
	var condition string
	var params []interface{}

	if req.UserId != "" {
		params = append(params, req.UserId)
		condition += fmt.Sprintf(" AND user_id=$%d", len(params))
	}
	if req.PeriodAt != "" {
		params = append(params, req.PeriodAt, req.PeriodEnd)
		condition += fmt.Sprintf(" AND date>=$%d AND date<=$%d", len(params)-1, len(params))
	}

	var data []models.FullOrderAnalytics
	query := fmt.Sprintf(`SELECT id, user_id, manager_id, number, date, status,	(SELECT company FROM "%s" WHERE "%s".id=user_id) as user_company,
		(SELECT name FROM "%s" WHERE "%s".id=user_id) as user_name,
		(SELECT name FROM "%s" WHERE id="%s".manager_id) as manager FROM "%s" WHERE date!='' %s
		ORDER BY manager_id, user_id, number`,
		UserTable, UserTable, UserTable, UserTable, UserTable, OrderTable, OrderTable, condition,
	)

	if err := r.db.Select(&data, query, params...); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	for i, foa := range data {
		if i == 0 || orders[len(orders)-1].Id != foa.ManagerId {
			orders = append(orders, &analytic_model.FullOrder{
				Id:      foa.ManagerId,
				Manager: foa.Manager,
				Clients: []*analytic_model.Client{{
					Id:      foa.UserId,
					Company: foa.Company,
					Name:    foa.Name,
					Orders: []*analytic_model.ClientOrder{{
						Id:     foa.Id,
						Number: foa.Number,
						Date:   foa.Date,
						Status: foa.Status,
					}},
				}},
			})
		} else {
			if orders[len(orders)-1].Clients[len(orders[len(orders)-1].Clients)-1].Id != foa.UserId {
				orders[len(orders)-1].Clients = append(orders[len(orders)-1].Clients, &analytic_model.Client{
					Id:      foa.UserId,
					Company: foa.Company,
					Name:    foa.Name,
					Orders: []*analytic_model.ClientOrder{{
						Id:     foa.Id,
						Number: foa.Number,
						Date:   foa.Date,
						Status: foa.Status,
					}},
				})
			} else {
				orders[len(orders)-1].Clients[len(orders[len(orders)-1].Clients)-1].Orders =
					append(orders[len(orders)-1].Clients[len(orders[len(orders)-1].Clients)-1].Orders, &analytic_model.ClientOrder{
						Id:     foa.Id,
						Number: foa.Number,
						Date:   foa.Date,
						Status: foa.Status,
					})
			}
		}
	}

	return orders, nil
}

func (r *OrderRepo) Create(ctx context.Context, order *order_api.CreateOrder, date string) error {
	query := fmt.Sprintf(`INSERT INTO "%s" (id, user_id, date, count_position, manager_id, info) VALUES ($1, $2, $3, $4, $5, $6)`, OrderTable)

	_, err := r.db.Exec(query, order.Id, order.UserId, date, order.Count, order.ManagerId, order.Info)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) SetInfo(ctx context.Context, order *order_api.Info) error {
	query := fmt.Sprintf(`UPDATE "%s" SET info=$1 WHERE id=$2`, OrderTable)

	_, err := r.db.Exec(query, order.Info, order.OrderId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) SetStatus(ctx context.Context, status *order_api.Status) error {
	query := fmt.Sprintf(`UPDATE "%s" SET status=$1, %s_date=$2 WHERE id=$3 AND status!=$1`, OrderTable, status.Status.String())

	_, err := r.db.Exec(query, status.Status.String(), status.Date, status.OrderId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *OrderRepo) SetManager(ctx context.Context, manager *order_api.Manager) error {
	query := fmt.Sprintf(`UPDATE "%s" SET manager_id=$1 WHERE id=$2`, OrderTable)

	_, err := r.db.Exec(query, manager.ManagerId, manager.OrderId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
