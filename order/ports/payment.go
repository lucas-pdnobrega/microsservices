package ports

import "github.com/lucas-pdnobrega/microservices/order/internal/application/core/domain"

type PaymentPort interface {
	Charge(order domain.Order) error
}
