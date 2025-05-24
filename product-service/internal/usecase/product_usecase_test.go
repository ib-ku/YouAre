package usecase

import (
	"product-service/internal/entity"
	"product-service/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductUseCase_CreateProduct(t *testing.T) {
	// Инициализируем мок репозитория
	mockRepo := new(mocks.ProductRepository)

	// Создаем экземпляр usecase с моком
	uc := NewProductUseCase(mockRepo)

	t.Run("successful product creation", func(t *testing.T) {
		// Настраиваем ожидания для мока
		mockRepo.On("Create", mock.AnythingOfType("*entity.Product")).
			Return(nil).Once()

		// Вызываем тестируемый метод
		product, err := uc.CreateProduct("Test Product", 10.99, 100)

		// Проверяем результаты
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, 10.99, product.Price)
		assert.Equal(t, 100, product.Stock)

		// Проверяем, что мок был вызван как ожидалось
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty name should return error", func(t *testing.T) {
		product, err := uc.CreateProduct("", 10.99, 100)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInvalidCredentials, err)
		assert.Nil(t, product)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("zero price should return error", func(t *testing.T) {
		product, err := uc.CreateProduct("Test", 0, 100)

		assert.Error(t, err)
		assert.Nil(t, product)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("negative stock should return error", func(t *testing.T) {
		product, err := uc.CreateProduct("Test", 10.99, -1)

		assert.Error(t, err)
		assert.Nil(t, product)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repository error should propagate", func(t *testing.T) {
		expectedError := entity.ErrInternalServerError
		mockRepo.On("Create", mock.AnythingOfType("*entity.Product")).
			Return(expectedError).Once()

		product, err := uc.CreateProduct("Test", 10.99, 100)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, product)
		mockRepo.AssertExpectations(t)
	})
}
