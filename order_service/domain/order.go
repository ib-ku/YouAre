package domain

type Order struct {
    ID        string
    ProductID string
    Quantity  int
}

type OrderRepository interface {
    CreateOrder(order *Order) error
}
