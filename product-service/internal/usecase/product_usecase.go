package usecase

import (
	"product-service/internal/entity"
	"product-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductUseCase interface {
	CreateProduct(name string, price float64, stock int) (*entity.Product, error)
	GetProduct(ID string) (*entity.Product, error)
	GetAllProducts() ([]*entity.Product, error)
	UpdateProduct(ID string, name string, price float64, stock int) (*entity.Product, error)
	DecreaseStock(ID string, quantity int) (*entity.Product, error)
	DeleteProduct(ID string) error
}

type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) *productUseCase {
	return &productUseCase{
		repo: repo,
	}
}

// realization
func (p *productUseCase) CreateProduct(name string, price float64, stock int) (*entity.Product, error) {
	if name == "" || price <= 0 || stock <= 0 {
		return nil, entity.ErrInvalidCredentials
	}

	product := &entity.Product{
		Name:  name,
		Price: price,
		Stock: stock,
	}

	err := p.repo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productUseCase) GetProduct(ID string) (*entity.Product, error) {
	product, err := p.repo.GetProductByID(ID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productUseCase) GetAllProducts() ([]*entity.Product, error) {
	products, err := p.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *productUseCase) UpdateProduct(ID string, name string, price float64, stock int) (*entity.Product, error) {
	if name == "" || price <= 0 || stock < 0 {
		return nil, entity.ErrInvalidCredentials
	}

	objID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	product := &entity.Product{
		ID:    objID,
		Name:  name,
		Price: price,
		Stock: stock,
	}

	return p.repo.Update(product)
}

func (p *productUseCase) DecreaseStock(ID string, quantity int) (*entity.Product, error) {
	if quantity <= 0 {
		return nil, entity.ErrInvalidCredentials
	}
	product, err := p.repo.Decrease(ID, quantity)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *productUseCase) DeleteProduct(ID string) error {
	return p.repo.Delete(ID)
}
