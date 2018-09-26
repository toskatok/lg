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

package core

import (
	"fmt"

	"github.com/go-resty/resty"
)

// HTTPTransport implements transport interface for http/https protocol
type HTTPTransport struct {
	cli *resty.Client
}

// Init creates a http client
func (ht *HTTPTransport) Init(url string, token string) error {
	// TODO use SetAuthToken instead of setting it manually
	ht.cli = resty.New().SetHeader("Authorization", token).SetHostURL(url)
	return nil
}

// Transmit sends data with given topic that is used as url
func (ht *HTTPTransport) Transmit(topic string, data []byte) error {
	resp, err := ht.cli.R().SetBody(data).Post(topic)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("error response: %s", resp.Body())
	}
	return nil
}