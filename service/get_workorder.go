package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WorkOrderListResponse 工单列表响应
type WorkOrderListResponse struct {
	Code int `json:"code"`
	Data struct {
		Count int `json:"count"`
		List  []struct {
			Id            int       `json:"id"`
			CreatedAt     time.Time `json:"created_at"`
			BusinessName  string    `json:"business_name"`
			DbName        string    `json:"db_name"`
			ApplicantUser string    `json:"applicant_user"`
			Remark        string    `json:"remark"`
			Department    string    `json:"department"`
			DbType        string    `json:"db_type"`
			Usage         string    `json:"usage"`
			TimeLimit     string    `json:"time_limit"`
			Status        string    `json:"status"`
			Role          string    `json:"role"`
			OperationUser string    `json:"operation_user"`
			FinishedAt    time.Time `json:"finished_at"`
			DbTypeInt     int       `json:"db_type_int"`
			DepartmentInt int       `json:"department_int"`
			UsageInt      int       `json:"usage_int"`
			TimeLimitInt  int       `json:"time_limit_int"`
			StatusInt     int       `json:"status_int"`
			RoleInt       int       `json:"role_int"`
		} `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// GetWorkOrdersList 获取工单列表
func (c *DBClient) GetWorkOrdersList(ctx context.Context) (*WorkOrderListResponse, error) {
	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var result WorkOrderListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
