package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const requestTimeout = time.Second * 1 // Default timeout.

func doHTTPGetWithTimeout(requestURL string, t time.Duration) (*bytes.Buffer, error) {
	httpClient := &http.Client{
		Timeout: t,
	}

	response, err := httpClient.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", response.StatusCode)
	}

	var b bytes.Buffer
	if _, err := io.Copy(&b, response.Body); err != nil {
		return nil, err
	}
	return &b, nil
}

func doHTTPGet(requestURL string) (*bytes.Buffer, error) {
	return doHTTPGetWithTimeout(requestURL, requestTimeout)
}
