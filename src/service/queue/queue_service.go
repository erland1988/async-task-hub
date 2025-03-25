package queue

import (
	"async-task-hub/common"
	"async-task-hub/global"
	"async-task-hub/src/model"
	"async-task-hub/src/types"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"time"
)

// 任务队列服务
type QueueService struct {
	ctx                     context.Context
	taskQueueZSetKey        string          // 任务队列zset key
	recoverLostTaskslockKey string          // 恢复丢失任务锁key
	maxBackoff              int             // 最大退避时间
	maxExecutionCount       int             // 最大执行次数
	futureTime              time.Duration   // 未来时间
	recoverTime             time.Duration   // 恢复时间
	ExecutorClient          *ExecutorClient // 执行器客户端
}

func NewQueueService() *QueueService {
	return &QueueService{
		ctx:                     context.Background(),
		taskQueueZSetKey:        global.CacheKey("queue_zset"),
		recoverLostTaskslockKey: global.CacheKey("queue_monitor_lock"),
		maxBackoff:              300,
		maxExecutionCount:       4,
		futureTime:              2 * time.Hour,
		recoverTime:             5 * time.Minute,
		ExecutorClient:          NewExecutorClient(),
	}
}

// 添加任务队列
func (s *QueueService) PushTaskToQueue(taskQueue *model.TaskQueue) error {
	if taskQueue.ID == 0 {
		return errors.New("taskQueue.ID is required")
	}
	if taskQueue.ExecutionTime == 0 {
		return errors.New("taskQueue.ExecutionTime is required")
	}
	futureTime := types.Timestamp(time.Now().Add(s.futureTime).Unix())
	if taskQueue.ExecutionTime > futureTime {
		global.Logger.Info("taskQueue.ExecutionTime is too far in the future", zap.Int("taskQueue.ExecutionTime", int(taskQueue.ExecutionTime)), zap.Int("futureTime", int(futureTime)))

		return nil
	}

	global.Logger.Debug("PushTaskToQueue start", zap.Any("taskQueueID", taskQueue.ID))
	err := global.REDIS.ZAdd(s.ctx, s.taskQueueZSetKey, &redis.Z{
		Score:  float64(taskQueue.ExecutionTime),
		Member: taskQueue.ID,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

var popTaskScript = redis.NewScript(`
	local now = tonumber(ARGV[1])
	local result = redis.call("ZRANGEBYSCORE", KEYS[1], "-inf", now, "LIMIT", 0, 1)
	if #result == 0 then
		return nil
	end
	redis.call("ZREM", KEYS[1], result[1])
	return result[1]
`)

func (s *QueueService) PopTaskFromQueue(ctx context.Context) (int, error) {
	select {
	case <-ctx.Done():
		return 0, context.Canceled
	default:
		now := time.Now().Unix()
		result, err := popTaskScript.Run(s.ctx, global.REDIS, []string{s.taskQueueZSetKey}, now).Result()
		if err != nil || result == nil {
			return 0, errors.New("no task found")
		}
		taskQueueID := common.Str2Int(result.(string))
		global.Logger.Debug("PopTaskFromQueue start", zap.Any("taskQueueID", taskQueueID))
		return taskQueueID, nil
	}
}

func (s *QueueService) ProcessTaskQueue(ctx context.Context, taskQueueID int) error {
	select {
	case <-ctx.Done():
		return context.Canceled
	default:
		var taskQueue model.TaskQueue

		global.Logger.Debug("ProcessTaskQueue start", zap.Any("taskQueueID", taskQueueID))
		if err := global.DB.Preload("Task").First(&taskQueue, taskQueueID).Error; err != nil {
			global.Logger.Error("ProcessTaskQueue", zap.Error(err))
			return err
		}

		var taskLog model.TaskLog
		requestID := common.HashUniqueID()
		if err := s.startTaskQueueConsumer(&taskQueue, taskLog, requestID); err != nil {
			global.Logger.Error("ProcessTaskQueue", zap.Error(err))
			return err
		}
		// 请求执行器
		result, err := s.ExecutorClient.SendRequestToExecutor(ctx, taskQueue, requestID)
		if err != nil {
			// 失败处理
			taskQueue.ExecutionStatus = model.TaskQueueFailed
			taskLog.Message = err.Error()
		} else {
			// 成功处理
			taskQueue.ExecutionStatus = model.TaskQueueCompleted
			taskLog.Message = result
		}

		if err := s.endTaskQueueConsumer(&taskQueue, taskLog, requestID); err != nil {
			global.Logger.Warn(err.Error(), zap.Error(err))
			return err
		}
		return nil
	}
}

// 恢复 Redis 队列中丢失的任务
func (s *QueueService) RecoverLostTasks() {
	var taskQueues []model.TaskQueue

	success, err := global.REDIS.SetNX(context.Background(), s.recoverLostTaskslockKey, "locked", s.recoverTime).Result()
	if err != nil || !success {
		return
	}

	global.Logger.Info("RecoverLostTasks start")
	futureTime := types.Timestamp(time.Now().Add(s.futureTime).Unix())

	query := global.DB
	query = query.Where("execution_status IN (?,?)", model.TaskQueuePending, model.TaskQueueFailed)
	query = query.Where("execution_time <?", futureTime)
	query = query.Where("execution_count <?", s.maxExecutionCount)
	if err := query.Find(&taskQueues).Error; err != nil {
		global.Logger.Error("RecoverLostTasks", zap.Error(err))
		return
	}

	for _, taskQueue := range taskQueues {
		_, err := global.REDIS.ZScore(s.ctx, s.taskQueueZSetKey, strconv.Itoa(taskQueue.ID)).Result()
		if err == redis.Nil {
			if err := s.PushTaskToQueue(&taskQueue); err != nil {
				global.Logger.Error("RecoverLostTasks", zap.Error(err))
			}
		}
	}
	return
}

func (s *QueueService) startTaskQueueConsumer(taskQueue *model.TaskQueue, taskLog model.TaskLog, requestID string) error {
	if taskQueue.ID == 0 {
		return nil
	}

	if taskQueue.ExecutionStatus == model.TaskQueueProcessing {
		return nil
	}

	if taskQueue.ExecutionStatus == model.TaskQueueCompleted {
		return nil
	}

	if taskQueue.ExecutionStatus == model.TaskQueueFailed {
		if taskQueue.ExecutionCount >= s.maxExecutionCount {
			return nil
		}
	}

	if taskQueue.ExecutionTime > types.Timestamp(time.Now().Unix()) {
		return nil
	}

	tx := global.DB.Begin()

	taskLog.AppID = taskQueue.AppID
	taskLog.TaskID = taskQueue.TaskID
	taskLog.TaskQueueID = taskQueue.ID
	taskLog.RequestID = requestID
	taskLog.Action = model.TaskLogActionStart
	taskLog.Message = taskQueue.Parameters
	taskLog.MilliTimestamp = types.MilliTimestamp(time.Now().UnixMilli())
	taskLog.CreatedAt = types.Customtime(time.Now())

	if err := tx.Create(&taskLog).Error; err != nil {
		tx.Rollback()
		global.Logger.Error("startTaskQueueConsumer", zap.Error(err))
		return err
	}

	taskQueue.ExecutionStatus = model.TaskQueueProcessing
	executionCount := taskQueue.ExecutionCount + 1
	taskQueue.ExecutionCount = executionCount
	taskQueue.UpdatedAt = types.Customtime(time.Now())

	if err := tx.Select("execution_status", "execution_count", "updated_at").Save(&taskQueue).Error; err != nil {
		tx.Rollback()
		global.Logger.Error("startTaskQueueConsumer", zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

func (s *QueueService) endTaskQueueConsumer(taskQueue *model.TaskQueue, taskLog model.TaskLog, requestID string) error {
	if taskQueue.ExecutionStatus == model.TaskQueuePending {
		return nil
	}
	if taskQueue.ExecutionStatus == model.TaskQueueProcessing {
		return nil
	}
	tx := global.DB.Begin()
	taskLog.AppID = taskQueue.AppID
	taskLog.TaskID = taskQueue.TaskID
	taskLog.TaskQueueID = taskQueue.ID
	taskLog.RequestID = requestID
	taskLog.Action = model.TaskLogActionEnd
	taskLog.MilliTimestamp = types.MilliTimestamp(time.Now().UnixMilli())
	taskLog.CreatedAt = types.Customtime(time.Now())
	if err := tx.Save(&taskLog).Error; err != nil {
		tx.Rollback()
		global.Logger.Error("endTaskQueueConsumer", zap.Error(err))
		return err
	}

	shouldRetry := taskQueue.ExecutionStatus == model.TaskQueueFailed && taskQueue.ExecutionCount < s.maxExecutionCount
	if shouldRetry {
		delaySeconds := s.getNextRetryDelay(taskQueue.ExecutionCount)
		global.Logger.Debug("getNextRetryDelay", zap.Any("executionCount", taskQueue.ExecutionCount), zap.Any("delaySeconds", delaySeconds))
		taskQueue.ExecutionTime = types.Timestamp(time.Now().Add(time.Duration(delaySeconds) * time.Second).Unix())
	}

	taskQueue.UpdatedAt = types.Customtime(time.Now())
	if err := tx.Select("execution_status", "execution_time", "updated_at").Save(&taskQueue).Error; err != nil {
		tx.Rollback()
		global.Logger.Error("endTaskQueueConsumer", zap.Error(err))
		return err
	}
	tx.Commit()

	if shouldRetry {
		if err := s.PushTaskToQueue(taskQueue); err != nil {
			tx.Rollback()
			global.Logger.Error("endTaskQueueConsumer PushTaskToQueue", zap.Error(err))
			return err
		}
	}
	go func() {
		if err := s.recoverDuration(taskQueue.ID, requestID); err != nil {
			global.Logger.Error("endTaskQueueConsumer recoverDuration", zap.Error(err))
		}
	}()
	return nil
}

func (s *QueueService) recoverDuration(taskQueueID int, requestID string) error {
	var taskLogs []model.TaskLog
	if err := global.DB.Where("request_id = ?", requestID).Order("id desc").Find(&taskLogs).Error; err != nil {
		global.Logger.Error("recoverDuration", zap.Error(err))
		return err
	}
	if len(taskLogs) == 0 {
		return nil
	}

	var taskQueue model.TaskQueue
	if err := global.DB.First(&taskQueue, taskQueueID).Error; err != nil {
		global.Logger.Error("recoverDuration", zap.Error(err))
		return err
	}

	var startTime int64
	var endTime int64
	for _, taskLog := range taskLogs {
		if taskLog.Action == model.TaskLogActionStart {
			startTime = int64(taskLog.MilliTimestamp)
		}
		if taskLog.Action == model.TaskLogActionEnd {
			endTime = int64(taskLog.MilliTimestamp)
		}
	}
	if startTime > 0 {
		startTimeTime := types.Customtime(time.UnixMilli(startTime))
		taskQueue.ExecutionStart = &startTimeTime
	}
	if endTime > 0 {
		endTimeTime := types.Customtime(time.UnixMilli(endTime))
		taskQueue.ExecutionEnd = &endTimeTime
	}
	if startTime > 0 && endTime > 0 && endTime >= startTime {
		taskQueue.ExecutionDuration = endTime - startTime
	}
	taskQueue.UpdatedAt = types.Customtime(time.Now())
	if err := global.DB.Select("execution_start", "execution_end", "execution_time", "execution_duration", "updated_at").Save(&taskQueue).Error; err != nil {
		global.Logger.Error("recoverDuration", zap.Error(err))
		return err
	}
	return nil
}

func (s *QueueService) getNextRetryDelay(retryCount int) int {
	baseDelay := 3 * (1 << retryCount)
	jitter := rand.Intn(baseDelay / 2)
	delay := baseDelay + jitter
	if delay > s.maxBackoff {
		delay = s.maxBackoff
	}
	return delay
}
