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
	"bytes"
	"fmt"
	"net/http"
)

// HTTPTransport implements transport interface for http/https protocol
type HTTPTransport struct {
	url string
}

// Init creates a http client
func (ht *HTTPTransport) Init(url string) error {
	ht.url = url
	return nil
}

// Transmit sends data with given topic that is used as url
func (ht *HTTPTransport) Transmit(topic string, data []byte) error {
	if _, err := http.Post(fmt.Sprintf("%s/%s", ht.url, topic), "application/json", bytes.NewReader(data)); err != nil {
		return err
	}
	return nil
}
