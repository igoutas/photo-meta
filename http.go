package main

import (
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
}

func (httpClient *HttpClient) Get(url string) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}

func NewHttpClient() HttpClientInterface {
	return &HttpClient{}
}
