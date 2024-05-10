package storage

import (
	"WB_L0_Task/models"
	"errors"
	"sync"
)

type ICache interface {
	IStorage
	FillCache(cache *map[string]*models.Order)
	IndexExist(index string) bool
}

type OrdersCache struct {
	cache *map[string]*models.Order
	mutex sync.Mutex
}

func (c *OrdersCache) Init() error {
	c.cache = new(map[string]*models.Order)
	c.mutex = sync.Mutex{}
	return nil
}

func (c *OrdersCache) AddOrder(order models.Order, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := (*c.cache)[key]; !ok {
		(*c.cache)[key] = &order
		return nil
	} else {
		return errors.New("element is already exists")
	}
}

func (c *OrdersCache) GetOrder(key string) *models.Order {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if val, ok := (*c.cache)[key]; ok {
		return val
	}
	return nil
}

func (c *OrdersCache) FillCache(cache *map[string]*models.Order) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = cache
}

func (c *OrdersCache) IndexExist(index string) bool {
	if _, ok := (*c.cache)[index]; !ok {
		return false
	}
	return true
}
