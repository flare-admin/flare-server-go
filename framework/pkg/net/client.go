package net

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ResponseModel 通用响应模型
type ResponseModel[T any] struct {
	Code   int    `json:"code"`
	Data   T      `json:"data"`
	Msg    string `json:"msg"`
	Reason string `json:"reason"`
}

// HttpClient struct
type HttpClient struct {
	client *http.Client
}

// NewHttpClient 创建一个新的HttpClient实例
func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{Timeout: timeout},
	}
}

// Get 发送GET请求并解析响应
func (c *HttpClient) Get(url string, responseModel interface{}) error {
	resp, err := c.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("unexpected status code: %d ，status:%s  ,err :%s", resp.StatusCode, resp.Status, err.Error())
			}
			fmt.Printf("Received 400 response:\n%s\n", string(body))
		} else {
			return fmt.Errorf("unexpected status code: %d ，status:%s", resp.StatusCode, resp.Status)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, responseModel)
}

// Post 发送POST请求并解析响应
func (c *HttpClient) Post(url string, requestModel interface{}, responseModel interface{}) error {
	requestBody, err := json.Marshal(requestModel)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, responseModel)
}

// PostReturnBody 发送POST请求并解析响应
func (c *HttpClient) PostReturnBody(url string, requestModel interface{}, responseModel interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(requestModel)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, json.Unmarshal(body, responseModel)
}

// PostPar 发送POST请求并解析响应
func (c *HttpClient) PostPar(url string, pas string, responseModel interface{}) error {
	resp, err := c.client.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(pas))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, responseModel)
}

// PostForm 发送POST请求并解析响应
func (c *HttpClient) PostForm(url string, dataMap map[string]interface{}, responseModel interface{}) error {
	// 将map转换为url.Values
	data := MapToURLValues(dataMap)
	resp, err := c.client.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, responseModel)
}

// MapToURLValues 将map转换为url.Values
func MapToURLValues(dataMap map[string]interface{}) url.Values {
	data := url.Values{}
	for key, value := range dataMap {
		data.Set(key, fmt.Sprintf("%v", value)) // 将 interface{} 转换为字符串
	}
	return data
}
