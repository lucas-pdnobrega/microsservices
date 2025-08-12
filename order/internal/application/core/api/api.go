package api

import (
	"fmt"

	"github.com/lucas-pdnobrega/microservices/order/internal/application/core/domain"
	"github.com/lucas-pdnobrega/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	var totalQuantity int32
	for _, orderItem := range order.OrderItems {
		totalQuantity += orderItem.Quantity
	}

	if totalQuantity > 50 {
		return domain.Order{}, status.Errorf(codes.InvalidArgument, 
			fmt.Sprintf("Order's total item quantity is %d, it cannot exceed 50", totalQuantity)).Err()
	}
	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}
	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}
	return order, nil
}
