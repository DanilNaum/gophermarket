package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	// "time"

	"github.com/DanilNaum/gophermarket/internal/repository/order"
	"github.com/google/uuid"
)

type orderRepository interface {
	CreateOrder(ctx context.Context, order *order.Order) error
	UpdateOrder(ctx context.Context, order *order.Order) error
	GetOrdersByUserId(ctx context.Context, userID uuid.UUID) ([]*order.Order, error)
}

type orderService struct {
	orders chan string

	orderRepository orderRepository

	url string
}

type Order struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *int   `json:"accrual,omitempty"`
}

func NewOrderService(orderRepository orderRepository, bufferSize int, url string) *orderService {
	return &orderService{
		orderRepository: orderRepository,
		orders:          make(chan string, bufferSize),
		url:             url,
	}
}

func (o *orderService) NewOrder(orderID string) {
	o.orders <- orderID
}

func (o *orderService) Start(appCtx context.Context) {
	for {
		select {
		case orderID, ok := <-o.orders:
			if !ok {
				return
			}
			err := o.processOrder(appCtx, orderID)
			if err != nil {
				o.orders <- orderID
			}
		case <-appCtx.Done():
			return
		}
	}
}

var errToManyRequests = errors.New("too many requests")

const defaultTimeout = 60

func (o *orderService) processOrder(appCtx context.Context, orderID string) error {
	// ctx, cancel := context.WithTimeout(appCtx, 10*time.Second)
	// defer cancel()

	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", o.url, orderID))
	if err != nil {
		//log error
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if resp.Header.Get("Content-Type") != "application/json" {
			return errors.New("invalid content type")
		}
		body, err := io.ReadAll(resp.Body)

		resp.Body.Close()

		if err != nil {
			return err
		}

		ord := &Order{}

		err = json.Unmarshal(body, ord)

		if err != nil {
			return err
		}

		switch ord.Status {

		case "PROCESSED":
			err := o.orderRepository.UpdateOrder(appCtx, &order.Order{
				ID:      ord.Order,
				Status:  ord.Status,
				Accrual: ord.Accrual,
			})
			if err != nil {
				// TODO log
			}
			if ord.Accrual != nil {

			}

		case "INVALID":
			err := o.orderRepository.UpdateOrder(appCtx, &order.Order{
				ID:     ord.Order,
				Status: ord.Status,
			})
			if err != nil {
				// TODO log
			}

		case "PROCESSING":
			err := o.orderRepository.UpdateOrder(appCtx, &order.Order{
				ID:      ord.Order,
				Status:  ord.Status,
				Accrual: ord.Accrual,
			})
			if err != nil {
				// TODO log
			}
			return errors.New("not processed yet")

		case "REGISTERED":
			return errors.New("not processed yet")
		}

	case http.StatusTooManyRequests:
		timeout, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			<-time.After(time.Duration(defaultTimeout) * time.Second)
			return errToManyRequests
		}
		<-time.After(time.Duration(timeout) * time.Second)
		return errToManyRequests

	case http.StatusNoContent:
		// todo: log error
		return nil

	default:
		return errors.New("unexpected status code")
	}

	return nil
}
