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

	"github.com/DanilNaum/gophermarket/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type orderRepository interface {
	CreateOrder(ctx context.Context, order *repository.Order) error
	UpdateOrder(ctx context.Context, order *repository.Order) error
	GetOrdersByUserId(ctx context.Context, userID uuid.UUID) ([]*repository.Order, error)
}

type orderService struct {
	orders chan string

	orderRepository orderRepository

	url string

	workerNum int
}

type Order struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *int   `json:"accrual,omitempty"`
}

const workerNum = 10

func NewOrderService(orderRepository orderRepository, bufferSize int, url string) *orderService {
	return &orderService{
		orderRepository: orderRepository,
		orders:          make(chan string, bufferSize),
		url:             url,
		workerNum:       workerNum,
	}
}

func (o *orderService) NewOrder(orderID string) {
	o.orders <- orderID
}

func (o *orderService) Start(ctx context.Context) {
	group, ctx := errgroup.WithContext(ctx)
	for _ = range o.workerNum {
		group.Go(func() error {
			return o.startWorker(ctx)
		})
	}
	_ = group.Wait()
	// todo log err
	return
}
func (o *orderService) startWorker(ctx context.Context) error {
	for {
		select {
		case orderID, ok := <-o.orders:
			if !ok {
				return errors.New("chan close")
			}
			err := o.processOrder(ctx, orderID)
			if err != nil {
				o.orders <- orderID
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

var errToManyRequests = errors.New("too many requests")

const (
	defaultTimeout = 60

	statusPROCESSED  = "PROCESSED"
	statusINVALID    = "INVALID"
	statusPROCESSING = "PROCESSING"
	statusREGISTERED = "REGISTERED"
)

func (o *orderService) processOrder(ctx context.Context, orderID string) error {

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

		case statusPROCESSED:
			err := o.orderRepository.UpdateOrder(ctx, &repository.Order{
				ID:      ord.Order,
				Status:  ord.Status,
				Accrual: ord.Accrual,
			})
			if err != nil {
				// TODO log
			}

		case statusINVALID:
			err := o.orderRepository.UpdateOrder(ctx, &repository.Order{
				ID:     ord.Order,
				Status: ord.Status,
			})
			if err != nil {
				// TODO log
			}

		case statusPROCESSING:
			err := o.orderRepository.UpdateOrder(ctx, &repository.Order{
				ID:      ord.Order,
				Status:  ord.Status,
				Accrual: ord.Accrual,
			})
			if err != nil {
				// TODO log
			}
			return errors.New("not processed yet")

		case statusREGISTERED:
			return errors.New("not processed yet")
		}

	case http.StatusTooManyRequests:
		timeout, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			select {
			case <-time.After(time.Duration(defaultTimeout) * time.Second):
				return errToManyRequests
			case <-ctx.Done():
				return ctx.Err()
			}
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
