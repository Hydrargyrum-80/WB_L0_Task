package storage

import (
	"WB_L0_Task/models"
)

type IStorage interface {
	Init() error
	AddOrder(order models.Order, key string) error
	GetOrder(key string) *models.Order
}

type OrdersStorage struct {
	OrderCache ICache
	Db         IDb
}

func (c *OrdersStorage) Init() error {
	c.Db = &OrdersDatabase{}
	err := c.Db.Init()
	if err != nil {
		return err
	}
	c.OrderCache = &OrdersCache{}
	err = c.OrderCache.Init()
	if err != nil {
		return err
	}
	var orders *map[string]*models.Order
	orders, err = c.Db.GetAllOrders()
	if err != nil {
		return err
	}
	c.OrderCache.FillCache(orders)
	return nil
}

func (c *OrdersStorage) AddOrder(order models.Order, key string) error {
	err := c.OrderCache.AddOrder(order, key)
	if err != nil {
		return err
	}
	err = c.Db.AddOrder(order, key)
	if err != nil {
		return err
	}
	return nil
}

func (c *OrdersStorage) GetOrder(key string) *models.Order {
	return c.OrderCache.GetOrder(key)
}
