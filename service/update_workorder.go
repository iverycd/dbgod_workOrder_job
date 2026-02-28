package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UpdateWorkorderRequest 请求参数结构体
type UpdateWorkorderRequest struct {
	Status        int    `json:"status"`
	DbName        string `json:"db_name"` // 申请创建的数据库名
	ID            int    `json:"id"`      // 工单的id
	UsageInt      int    `json:"usage_int"`
	TimeLimitInt  int    `json:"time_limit_int"`
	DepartmentInt int    `json:"department_int"`
	DbTypeInt     int    `json:"db_type_int"`
}

// UpdateWorkorderResponse 响应结果结构体
type UpdateWorkorderResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func (c *DBClient) UpdateWorkorder(ctx context.Context, req *UpdateWorkorderRequest) (*UpdateWorkorderResponse, error) {
	// 序列化请求体
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	// 创建请求
	reqHTTP, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	// 设置请求头
	reqHTTP.Header.Set("Content-Type", "application/json")
	for k, v := range c.headers {
		reqHTTP.Header.Set(k, v)
	}

	// 发送请求
	respHTTP, err := c.httpClient.Do(reqHTTP)
	if err != nil {
		return nil, fmt.Errorf("http do error: %w", err)
	}
	defer respHTTP.Body.Close()

	// 解析响应
	var resp UpdateWorkorderResponse
	if err := json.NewDecoder(respHTTP.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("response decode error: %w", err)
	}

	return &resp, nil
}
