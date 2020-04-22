/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 26-09-2018
 * |
 * | File Name:     http.go
 * +===============================================
 */

package transport

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// HTTP implements transport interface for http/https protocol
type HTTP struct {
	cli *resty.Client
}

// Init creates a http client
func (ht *HTTP) Init(url string, token string) error {
	ht.cli = resty.New().SetHeader(
		"Authorization",
		token,
	).SetHeader(
		"Content-Type",
		"application/json",
	).SetHostURL(fmt.Sprintf("http://%s", url))

	return nil
}

// Transmit sends data with given topic that is used as url
func (ht *HTTP) Transmit(topic string, data []byte) error {
	resp, err := ht.cli.R().SetBody(data).Post(topic)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("error response: %s", resp.Body())
	}

	return nil
}
