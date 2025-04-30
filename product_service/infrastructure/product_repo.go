package infrastructure

import (
	"YouAre/product_service/domain"
	"fmt"
)

type InMemoryProductRepository struct {
	products map[string]*domain.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: make(map[string]*domain.Product),
	}
}
func (repo *InMemoryProductRepository) Save(product *domain.Product) error {
	repo.products[product.ID] = product
	fmt.Printf("Product saved: %+v\n", product)
	return nil
}

func (repo *InMemoryProductRepository) Get(id string) (*domain.Product, error) {
	product, exists := repo.products[id]
	if !exists {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

func (repo *InMemoryProductRepository) GetAll() []*domain.Product {
	var products []*domain.Product
	for _, product := range repo.products {
		products = append(products, product)
	}
	return products
}

func (repo *InMemoryProductRepository) Update(product *domain.Product) error {
	repo.products[product.ID] = product
	fmt.Printf("Product updated: %+v\n", product)
	return nil
}

func (repo *InMemoryProductRepository) Delete(id string) error {
	_, exists := repo.products[id]
	if !exists {
		return fmt.Errorf("product not found")
	}
	delete(repo.products, id)
	fmt.Printf("Product deleted: %s\n", id)
	return nil
}
