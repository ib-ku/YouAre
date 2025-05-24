package mocks

import (
	"product-service/internal/entity"

	"github.com/stretchr/testify/mock"
)

type ProductRepository struct {
	mock.Mock
}

func (m *ProductRepository) Create(product *entity.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *ProductRepository) GetProductByID(ID string) (*entity.Product, error) {
	args := m.Called(ID)
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepository) GetAllProducts() ([]*entity.Product, error) {
	args := m.Called()
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *ProductRepository) Update(product *entity.Product) (*entity.Product, error) {
	args := m.Called(product)
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepository) Decrease(ID string, quantity int) (*entity.Product, error) {
	args := m.Called(ID, quantity)
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepository) Delete(ID string) error {
	args := m.Called(ID)
	return args.Error(0)
}
