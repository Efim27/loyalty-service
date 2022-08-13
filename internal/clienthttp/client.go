package clienthttp

import (
	"github.com/go-resty/resty/v2"
	"loyalty-service/internal/config"
)

func NewClientHTTP(config config.HTTPClientConfig) (client *resty.Client) {
	client = &resty.Client{}

	client.
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(config.RetryWaitTime).
		SetRetryMaxWaitTime(config.RetryMaxWaitTime)

	return
}
