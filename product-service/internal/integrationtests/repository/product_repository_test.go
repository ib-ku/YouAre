//go:build integration
// +build integration

package repository_test

import (
	"product-service/internal/entity"
	"product-service/internal/integrationtests/testdb"
	"product-service/internal/repository"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	testdb.TestMainWrapper(m)
}

func TestProductRepository_Create(t *testing.T) {
	repo := repository.NewProductRepo(testdb.GetDB())

	product := &entity.Product{
		Name:  "Test Product",
		Price: 19.99,
		Stock: 10,
	}

	err := repo.Create(product)
	require.NoError(t, err)
	require.NotEmpty(t, product.ID)
}
