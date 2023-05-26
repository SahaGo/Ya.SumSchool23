package repositories

import (
	cerrors "Ya.SumSchool23/controllers/errors"
	"Ya.SumSchool23/services/model"
	"Ya.SumSchool23/services/service_data"
	"database/sql"
	"github.com/lib/pq"
)

type CourierRepository struct {
	db *sql.DB
}

func NewCourierRepository(db *sql.DB) *CourierRepository {
	return &CourierRepository{
		db: db,
	}
}

func (r *CourierRepository) GetCouriers(limit, offset int64) ([]*model.Courier, error) {

	rows, err := r.db.Query(
		"SELECT id, courier_type, regions, working_hours FROM couriers ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var couriers []*model.Courier
	for rows.Next() {
		courier := &model.Courier{}
		if err = rows.Scan(&courier.CourierId, &courier.CourierType, pq.Array(&courier.Regions), pq.Array(&courier.WorkingHours)); err != nil {
			return nil, err
		}
		couriers = append(couriers, courier)
	}
	return couriers, nil
}

func (r *CourierRepository) GetCourierById(id int64) (*model.Courier, error) {

	row := r.db.QueryRow("SELECT id, courier_type, regions, working_hours FROM couriers WHERE id = $1", id)

	courier := &model.Courier{}
	if err := row.Scan(
		&courier.CourierId,
		&courier.CourierType,
		pq.Array(&courier.Regions),
		pq.Array(&courier.WorkingHours),
	); err == sql.ErrNoRows {
		return nil, cerrors.NotFound.Wrapf(err, "courier with id = '%v' not found", id)
	} else if err != nil {
		return nil, err
	}
	return courier, nil
}

func (r *CourierRepository) CreateCouriers(data []service_data.NewCourierData) ([]int64, error) {

	ids := make([]int64, len(data))

	for i := 0; i < len(data); i++ {
		row := r.db.QueryRow("INSERT INTO couriers(courier_type, regions, working_hours) VALUES ($1,$2,$3) RETURNING id",
			data[i].CourierType, pq.Array(data[i].Regions), pq.StringArray(data[i].WorkingHours))
		if err := row.Scan(&ids[i]); err != nil {
			return nil, err
		}
	}
	return ids, nil
}
