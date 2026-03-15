package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gophermart/internal/entities"
)

var ErrRetryLater = errors.New("retry later")

type AccrualClient interface {
	GetOrder(ctx context.Context, number string) (*entities.Order, error)
}

type accrualClient struct {
	baseURL string
	client  *http.Client
}

func NewAccrualClient(addr string) AccrualClient {

	return &accrualClient{
		baseURL: addr,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type accrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (c *accrualClient) GetOrder(
	ctx context.Context,
	number string,
) (*entities.Order, error) {

	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, number)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK:

		var r accrualResponse

		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return nil, err
		}

		order := &entities.Order{
			Number:  r.Order,
			Status:  entities.OrderStatus(r.Status),
			Accrual: r.Accrual,
		}

		return order, nil

	case http.StatusNoContent:
		return nil, nil

	case http.StatusTooManyRequests:

		retryAfter := resp.Header.Get("Retry-After")

		seconds, err := strconv.Atoi(retryAfter)
		if err != nil {
			seconds = 60
		}

		select {
		case <-time.After(time.Duration(seconds) * time.Second):
			return nil, ErrRetryLater
		case <-ctx.Done():
			return nil, ctx.Err()
		}

	case http.StatusInternalServerError:
		return nil, ErrRetryLater

	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
