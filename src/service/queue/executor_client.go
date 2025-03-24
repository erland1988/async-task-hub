package queue

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"bytes"
	"context"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

// 执行器客户端
type ExecutorClient struct {
	timeout time.Duration
	client  *http.Client
}

func NewExecutorClient() *ExecutorClient {
	return &ExecutorClient{
		timeout: time.Hour * 2, // 超时时间
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        200,             // 最大空闲连接数
				MaxIdleConnsPerHost: 100,             // 每个主机的最大空闲连接数
				MaxConnsPerHost:     150,             // 每个主机的最大连接数
				IdleConnTimeout:     time.Minute * 5, // 空闲连接超时时间
			},
		},
	}
}

func (s *ExecutorClient) generateSignature(appKey, appSecret, timestamp string) string {
	signature := common.HashMD5(appKey + appSecret + timestamp)
	return signature
}

func (s *ExecutorClient) getHeader(taskQueue model.TaskQueue, requestID string) (http.Header, error) {
	header := http.Header{}

	var application model.Application
	if err := global.DB.Where("id = ?", taskQueue.Task.AppID).First(&application).Error; err != nil {
		return nil, err
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := s.generateSignature(application.AppKey, application.AppSecret, timestamp)
	header.Set("Content-Type", "application/json")
	header.Set("X-App-Key", application.AppKey)
	header.Set("X-Task-Code", taskQueue.Task.TaskCode)
	header.Set("X-Queue-ID", strconv.Itoa(taskQueue.ID))
	header.Set("X-Execution-Count", strconv.Itoa(taskQueue.ExecutionCount))
	header.Set("X-Request-ID", requestID)
	header.Set("X-Timestamp", timestamp)
	header.Set("X-Signature", signature)
	return header, nil
}

func (s *ExecutorClient) SendRequestToExecutor(ctx context.Context, taskQueue model.TaskQueue, requestID string) (string, error) {
	global.Logger.Debug("SendRequestToExecutor start", zap.Any("taskQueue", taskQueue), zap.String("requestID", requestID))

	// 创建带超时控制的子 context，确保传入的 ctx 控制退出
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	reqBody := bytes.NewBuffer([]byte(taskQueue.Parameters))
	req, err := http.NewRequestWithContext(ctx, "POST", taskQueue.Task.ExecutorURL, reqBody)
	if err != nil {
		return "", err
	}

	header, err := s.getHeader(taskQueue, requestID)
	if err != nil {
		return "", err
	}
	req.Header = header

	resp, err := s.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", errors.New("请求超时，可能对方执行时间过长")
		}
		if errors.Is(err, context.Canceled) {
			return "", errors.New("请求已取消，系统正在优雅退出")
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("请求失败，状态码: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if len(body) > 512 {
		body = append(body[:512], []byte("...")...)
	}
	return string(body), nil
}
