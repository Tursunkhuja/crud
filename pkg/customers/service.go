package customers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var ErrNotFound = errors.New("item not found")
var ErrInternal = errors.New("internal error")

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

func (s *Service) ByID(ctx context.Context, id int64) (*Customer, error) {
	item := &Customer{}
	err := s.pool.QueryRow(ctx, `
	SELECT id, name ,phone, active, created FROM  customers WHERE  id = $1;
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}
func (s *Service) GetAll(ctx context.Context) ([]*Customer, error) {
	var items []*Customer
	rows, err := s.pool.Query(ctx, `
		SELECT * FROM customers;
	`)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *Service) GetAllActive(ctx context.Context) ([]*Customer, error) {
	items := make([]*Customer, 0)
	rows, err := s.pool.Query(ctx, `
		SELECT * FROM customers where active;
	`)
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		item := &Customer{}
		err := rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			return nil, ErrInternal
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {

		return nil, ErrInternal
	}
	return items, nil
}
func (s *Service) Save(ctx context.Context, item *Customer) (*Customer, error) {
	itemFromDB := &Customer{}
	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `
	INSERT INTO customers (name, phone) VALUES ($1,$2) RETURNING  id,name,phone,active,created;
	`, item.Name, item.Phone).Scan(&itemFromDB.ID, &itemFromDB.Name, &itemFromDB.Phone, &itemFromDB.Active, &itemFromDB.Created)
		if err != nil {
			log.Println(err)
		}
		return itemFromDB, nil
	}
	err := s.pool.QueryRow(ctx, `
	UPDATE  customers SET name=$2, phone=$3 where id=$1 RETURNING  id,name,phone,active,created;
	`, item.ID, item.Name, item.Phone).Scan(&itemFromDB.ID, &itemFromDB.Name, &itemFromDB.Phone, &itemFromDB.Active, &itemFromDB.Created)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return itemFromDB, nil
}
func (s *Service) RemoveByID(ctx context.Context, id int64) (int64, error) {
	_, err := s.ByID(ctx, id)
	if err == ErrNotFound {
		log.Println(ErrNotFound)
		return 0, ErrNotFound
	}
	s.pool.QueryRow(ctx, `
	DELETE FROM customers WHERE id = $1;
	`, id)
	log.Println(id, "id")
	return id, nil
}
func (s *Service) BlockByID(ctx context.Context, id int64) (*Customer, error) {
	itemFromDB := &Customer{}
	_, err := s.ByID(ctx, id)
	if err == ErrNotFound {
		return nil, ErrNotFound
	}
	s.pool.QueryRow(ctx, `
	UPDATE customers SET active=false WHERE id = $1 RETURNING id ,name, phone, active, created;
	`, id).Scan(&itemFromDB.ID, &itemFromDB.Name, &itemFromDB.Phone, &itemFromDB.Active, &itemFromDB.Created)
	return itemFromDB, nil
}
func (s *Service) UnblockByID(ctx context.Context, id int64) (*Customer, error) {
	itemFromDB := &Customer{}
	_, err := s.ByID(ctx, id)
	if err == ErrNotFound {
		return nil, ErrNotFound
	}
	s.pool.QueryRow(ctx, `
		UPDATE customers SET active=true WHERE id = $1 RETURNING id ,name, phone, active, created;
	`, id).Scan(&itemFromDB.ID, &itemFromDB.Name, &itemFromDB.Phone, &itemFromDB.Active, &itemFromDB.Created)
	return itemFromDB, nil
}
