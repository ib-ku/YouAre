//go:build integration
// +build integration

package usecase_test

import (
	"product-service/internal/integrationtests/testdb"
	"product-service/internal/repository"
	"product-service/internal/usecase"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	testdb.TestMainWrapper(m)
}

func TestProductUseCase_Integration(t *testing.T) {
	repo := repository.NewProductRepo(testdb.GetDB())
	uc := usecase.NewProductUseCase(repo)

	t.Run("full flow", func(t *testing.T) {
		product, err := uc.CreateProduct("Test", 10.99, 5)
		require.NoError(t, err)
		require.NotEmpty(t, product.ID)

		found, err := uc.GetProduct(product.ID.Hex())
		require.NoError(t, err)
		require.Equal(t, product.Name, found.Name)
	})
}
