package restclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RestClient interface {
	BuildUrl(externalApi string, resource string, params ...interface{}) (string, error)
	DoGet(ctx context.Context, url string, result interface{}, additionalHeaders ...Header) error
}

type restClient struct {
	config Config
	client http.Client
}

type Header struct {
	Key   string
	Value string
}

func NewRestClient(config Config) (RestClient, error) {
	client := http.Client{
		Timeout: config.TimeoutMillis * time.Millisecond,
	}
	return &restClient{
		config: config,
		client: client,
	}, nil
}

func (rc restClient) BuildUrl(externalApi string, resource string, params ...interface{}) (string, error) {
	url := ""
	if val, exist := rc.config.ExternalApiCalls[externalApi]; exist {
		url = val.ApiDomain + fmt.Sprintf(val.Resources[resource].RequestUri, params...)
		return url, nil
	}
	return url, errors.New("resource_not_found")
}

func (rc restClient) DoGet(ctx context.Context, url string, result interface{}, additionalHeaders ...Header) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	for _, header := range additionalHeaders {
		req.Header.Add(header.Key, header.Value)
	}
	res, err := rc.client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = res.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	return nil
}
