package integration_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"os"
	"testing"
)

const (
	requestURL = "http://localhost:8080/api/times"
)

var headerContentTypeJson = []byte("application/json")

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestAskTimeCorrectlyHttp(t *testing.T) {
	client := &fasthttp.Client{}

	reqEntity := map[string]string{
		"ask": "What time is it?",
	}
	reqEntityBytes, _ := json.Marshal(reqEntity)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestURL)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(reqEntityBytes)
	req.Body()
	resp := fasthttp.AcquireResponse()

	err := client.Do(req, resp)
	assert.Nil(t, err)
	status := resp.StatusCode()

	assert.Equal(t, status, fasthttp.StatusOK)
}
func TestAskTimeInCorrectlyHttp(t *testing.T) {
	client := &fasthttp.Client{}

	reqEntity := map[string]string{
		"ask": " time is it?",
	}
	reqEntityBytes, _ := json.Marshal(reqEntity)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestURL)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(reqEntityBytes)
	req.Body()
	resp := fasthttp.AcquireResponse()

	err := client.Do(req, resp)
	assert.Nil(t, err)
	status := resp.StatusCode()

	assert.Equal(t, status, fasthttp.StatusBadRequest)
}
