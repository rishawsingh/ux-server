package utils

import (
	"sync"

	"github.com/remotestate/golang/models"
)

type InMemScore struct {
	lock           sync.RWMutex
	data           []models.ProductAttributeList
	segregatedData map[string]map[string][]models.ProductAttributeList
}

func NewInMemScore() *InMemScore {
	return &InMemScore{
		data:           make([]models.ProductAttributeList, 0),
		segregatedData: make(map[string]map[string][]models.ProductAttributeList),
	}
}

func (m *InMemScore) Get() []models.ProductAttributeList {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.data
}

func (m *InMemScore) Set(values []models.ProductAttributeList) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data = values
	m.segregate()
}

func (m *InMemScore) segregate() {
	for i := range m.data {
		_, ok := m.segregatedData[m.data[i].SubcategoryID]
		brandID := m.data[i].BrandID
		if !ok {
			m.segregatedData[m.data[i].SubcategoryID] = make(map[string][]models.ProductAttributeList)
		}
		m.segregatedData[m.data[i].SubcategoryID][brandID] = append(m.segregatedData[m.data[i].SubcategoryID][brandID], m.data[i])
	}
}

func (m *InMemScore) GetSegregatedData() map[string]map[string][]models.ProductAttributeList {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.segregatedData
}

func (m *InMemScore) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.data)
}
