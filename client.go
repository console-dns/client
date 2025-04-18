package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/console-dns/spec/models"
	"github.com/pkg/errors"
)

type ConsoleDnsClient struct {
	Server string `json:"server" yaml:"server" toml:"server"`
	Token  string `json:"token" yaml:"token" toml:"token"`
}

func NewConsoleDnsClient(server string, token string) *ConsoleDnsClient {
	return &ConsoleDnsClient{
		Server: server,
		Token:  token,
	}
}

func (c *ConsoleDnsClient) newRequest(method, path string, data string) (string, *http.Response, error) {
	var b io.Reader
	if data != "" {
		b = bytes.NewBuffer([]byte(data))
	}
	req, _ := http.NewRequest(method, fmt.Sprintf("%s%s", c.Server, path), b)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("User-Agent", "CoreDNS-plugin-console")
	if data != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", resp, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", resp, errors.New(string(body))
	}
	return string(body), resp, nil
}

func (c *ConsoleDnsClient) ListZones() (*models.Zones, *http.Response, error) {
	body, resp, err := c.newRequest("GET", "/api/v1/zones", "")
	if err != nil {
		return nil, resp, err
	}
	nextZone := &models.Zones{}
	err = json.Unmarshal([]byte(body), nextZone)
	if err != nil {
		return nil, resp, err
	}
	return nextZone, resp, nil
}

func (c *ConsoleDnsClient) ListZone(zone string) (*models.Zone, *http.Response, error) {
	body, resp, err := c.newRequest("GET", "/api/v1/zones/"+zone, "")
	if err != nil {
		return nil, resp, err
	}
	nextZone := &models.Zone{}
	err = json.Unmarshal([]byte(body), nextZone)
	if err != nil {
		return nil, resp, err
	}
	return nextZone, resp, nil
}

func (c *ConsoleDnsClient) CreateRecord(zone, rName, rType string, record any) (*http.Response, error) {
	data, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.newRequest("POST", "/api/v1/zones/"+zone+"/"+rName+"/"+rType, string(data))
	return resp, err
}

func (c *ConsoleDnsClient) UpdateRecord(zone, rName, rType string, old, next any) (*http.Response, error) {
	data, err := json.Marshal(map[string]any{
		"src": old,
		"dst": next,
	})
	if err != nil {
		return nil, err
	}
	_, resp, err := c.newRequest("POST", "/api/v1/zones/"+zone+"/"+rName+"/"+rType+"/edit", string(data))
	return resp, err
}

func (c *ConsoleDnsClient) DeleteRecord(zone, rName, rType string, record any) (*http.Response, error) {
	data, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.newRequest("POST", "/api/v1/zones/"+zone+"/"+rName+"/"+rType+"/delete", string(data))
	return resp, err
}
