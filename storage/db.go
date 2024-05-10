package storage

import (
	"WB_L0_Task/models"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
)

type IDb interface {
	IStorage
	GetAllOrders() (*map[string]*models.Order, error)
}

type OrdersDatabase struct {
	db *sql.DB
}

func (db *OrdersDatabase) GetOrder(key string) *models.Order {
	//TODO implement me
	panic("implement me")
}

func (db *OrdersDatabase) Init() error {
	newDb, err := sql.Open("postgres", "user=postgres password=2917819 dbname=WB_L0_DB sslmode=disable")
	if err != nil {
		return err
	}
	db.db = newDb
	_, err = db.db.Exec(`CREATE TABLE if NOT EXISTS orders (
		order_uid CHARACTER VARYING PRIMARY KEY,
		order_info JSON NOT NULL
	);`)
	if err != nil {
		return err
	}
	return nil
}

func (db *OrdersDatabase) GetAllOrders() (*map[string]*models.Order, error) {
	rows, err := db.db.Query("SELECT order_uid, order_info FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		jsonByte []byte
		index    string
	)
	jsonBuf := new(models.Order)
	orders := new(map[string]*models.Order)
	*orders = map[string]*models.Order{}
	for rows.Next() {
		err := rows.Scan(&index, &jsonByte)
		if err != nil {
			return nil, err
		}
		jsonBuf = nil
		err = json.Unmarshal(jsonByte, &jsonBuf)
		if err != nil {
			return nil, err
		}
		(*orders)[index] = jsonBuf
	}
	return orders, nil
}

func (db *OrdersDatabase) AddOrder(order models.Order, key string) error {
	el, err := json.Marshal(order)
	if err != nil {
		return err
	}
	query := "INSERT INTO orders (order_uid, order_info) VALUES ('" + key + "', '" + string(el) + "')"
	_, err = db.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
