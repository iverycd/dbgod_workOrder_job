package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CreateDBRequest 请求参数结构体
type CreateDBRequest struct {
	WorkOrderId   int    `json:"workorder_id"`   // 工单的id
	DbType        int    `json:"db_type"`        // 数据库类型,oneof枚举类型，限定前端只能传入特定数值
	DbName        string `json:"db_name"`        // 申请创建的数据库名
	ApplicantUser string `json:"applicant_user"` // 申请人的名称
	Department    int    `json:"department"`     // 申请人所在条线,oneof枚举类型，限定前端只能传入特定数值
	Role          int    `json:"role"`           // 权限  1 管理员  2 普通用户 (默认) 3 游客
	Remark        string `json:"remark"`         // 备注信息
	Usage         int    `json:"usage"`          // 使用用途 1 开发 2 测试 3 压测
	TimeLimit     int    `json:"time_limit"`     // 数据库预估使用期限 1 使用数周 2 使用数月 3 长期使用
	InstId        int    `json:"inst_id"`        // 数据库实例信息的id
}

// CreateDBResponse 响应结果结构体
type CreateDBResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// ClientOption 客户端配置选项
type ClientOption func(*DBClient)

// DBClient 数据库客户端
type DBClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// NewDBClient 创建新客户端
func NewDBClient(baseURL string, opts ...ClientOption) *DBClient {
	c := &DBClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // 默认超时时间
		},
		headers: make(map[string]string),
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithTimeout 设置自定义超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *DBClient) {
		c.httpClient.Timeout = timeout
	}
}

// WithHeader 设置自定义请求头
func WithHeader(key, value string) ClientOption {
	return func(c *DBClient) {
		c.headers[key] = value
	}
}

// CreateDatabase 执行创建数据库请求
//func (c *DBClient) CreateDatabase(ctx context.Context, req *CreateDBRequest) (*CreateDBResponse, error) {
//	// 序列化请求体
//	bodyBytes, err := json.Marshal(req)
//	if err != nil {
//		return nil, fmt.Errorf("json marshal error: %w", err)
//	}
//
//	// 创建请求
//	reqHTTP := fasthttp.AcquireRequest()   //获取Request连接池中的连接
//	defer fasthttp.ReleaseRequest(reqHTTP) // 用完需要释放资源
//
//	// 设置请求头
//	reqHTTP.Header.SetContentType("application/json")
//	reqHTTP.Header.SetMethod("POST")
//	for k, v := range c.headers {
//		reqHTTP.Header.Set(k, v)
//	}
//
//	reqHTTP.SetRequestURI(c.baseURL)
//
//	reqHTTP.SetBody(bodyBytes)
//
//	resp := fasthttp.AcquireResponse()             //获取Response连接池中的连接
//	defer fasthttp.ReleaseResponse(resp)           // 用完需要释放资源
//	if err := fasthttp.Do(req, resp); err != nil { //发送请求
//		return err
//	}
//	rspBody := resp.Body()
//	fmt.Println(string(rspBody))
//
//	// 发送请求
//
//
//	// 解析响应
//	var resp CreateDBResponse
//	if err := json.NewDecoder(respHTTP.Body).Decode(&resp); err != nil {
//		return nil, fmt.Errorf("response decode error: %w", err)
//	}
//
//	return &resp, nil
//}

func (c *DBClient) CreateDatabase(ctx context.Context, req *CreateDBRequest) (*CreateDBResponse, error) {
	// 序列化请求体
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}

	// 创建请求
	reqHTTP, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(bodyBytes))
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
	var resp CreateDBResponse
	if err := json.NewDecoder(respHTTP.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("response decode error: %w", err)
	}

	return &resp, nil
}
