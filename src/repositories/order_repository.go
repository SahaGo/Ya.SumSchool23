package repositories

import (
	cerrors "Ya.SumSchool23/controllers/errors"
	"Ya.SumSchool23/services/model"
	"Ya.SumSchool23/services/service_data"
	"database/sql"
	"github.com/lib/pq"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) GetOrderById(id int64) (*model.Order, error) {

	row := r.db.QueryRow("SELECT order_id, weight, regions, delivery_hours, order_cost, completed_time FROM orders WHERE order_id = $1", id)

	order := &model.Order{}
	if err := row.Scan(
		&order.OrderId,
		&order.Weight,
		&order.Regions,
		pq.Array(&order.DeliveryHours),
		&order.Cost,
		&order.CompletedTime,
	); err == sql.ErrNoRows {
		return nil, cerrors.NotFound.Wrapf(err, "order with id = '%v' not found", id)
	} else if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) GetOrders(limit, offset int64) ([]*model.Order, error) {

	rows, err := r.db.Query("SELECT order_id, weight, regions, delivery_hours, order_cost, completed_time FROM orders ORDER BY order_id LIMIT $1 OFFSET $2",
		limit, offset)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		err = rows.Scan(&order.OrderId, &order.Weight, &order.Regions, pq.Array(&order.DeliveryHours),
			&order.Cost, &order.CompletedTime)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) CreateOrders(data []service_data.NewOrderData) ([]int64, error) {

	ids := make([]int64, len(data))

	for i := 0; i < len(data); i++ {
		row := r.db.QueryRow("INSERT INTO orders(weight, regions, delivery_hours,order_cost) VALUES ($1,$2,$3, $4) RETURNING order_id",
			data[i].Weight, data[i].Regions, pq.StringArray(data[i].DeliveryHours), data[i].Cost)
		if err := row.Scan(&ids[i]); err != nil {
			return nil, err
		}
	}
	return ids, nil
}

func (r *OrderRepository) CreateCompleteOrder(data []service_data.NewCompleteOrderData) ([]int64, error) {

	ids := make([]int64, len(data))
	for i := 0; i < len(data); i++ {
		ids[i] = data[i].OrderId
	}

	var existingIdsCount int
	ordersCountRow := r.db.QueryRow("select count(*) as sum from orders where order_id = ANY($1)", pq.Array(ids))
	err := ordersCountRow.Scan(&existingIdsCount)
	if err != nil {
		return nil, err
	}
	if existingIdsCount != len(ids) {
		return nil, cerrors.BadRequest.Wrapf(err, "some orders not found")
	}

	for i := 0; i < len(data); i++ {
		var completedTime *string
		row := r.db.QueryRow("select completed_time from orders where order_id = $1", data[i].OrderId)
		if err = row.Scan(&completedTime); err != nil {
			return nil, err
		}

		if completedTime == nil {
			_, err = r.db.Exec("update orders set completed_time = $1, completed_courier_id = $2 where order_id = $3",
				data[i].CompleteTime, data[i].CourierId, data[i].OrderId)
			if err != nil {
				return nil, err
			}
		}
	}

	return ids, nil
}
