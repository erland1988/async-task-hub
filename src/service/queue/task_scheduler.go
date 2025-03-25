package queue

import (
	"async-task-hub/global"
	"async-task-hub/src/model"
	"context"
	"go.uber.org/zap"
	"time"
)

// 任务调度器
type TaskScheduler struct {
	workerLimit  int // 最大并发数
	QueueService *QueueService
}

func NewTaskScheduler() *TaskScheduler {
	return &TaskScheduler{
		workerLimit:  30,
		QueueService: NewQueueService(),
	}
}

// 启动任务队列监听器
func (s *TaskScheduler) StartTaskQueueListener(ctx context.Context) {
	taskChan := make(chan int, s.workerLimit)

	// 启动 worker
	for i := 0; i < s.workerLimit; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done(): // 优雅退出
					global.Logger.Info("StartTaskQueueListener worker stopped.")
					return
				case taskQueueID, ok := <-taskChan:
					if !ok { // 如果 channel 已关闭，直接退出
						return
					}
					err := s.QueueService.ProcessTaskQueue(ctx, taskQueueID)
					if err != nil {
						var taskQueue model.TaskQueue
						if err := global.DB.Preload("Task").First(&taskQueue, taskQueueID).Error; err != nil {
							global.Logger.Error("StartTaskQueueListener", zap.Error(err))
							continue
						}
						if taskQueue.ID > 0 {
							s.QueueService.PushTaskToQueue(&taskQueue)
						}
					}
				}
			}
		}()
	}

	// 任务派发
	go func() {
		defer func() {
			close(taskChan) // 确保在退出时关闭 channel
			global.Logger.Info("StartTaskQueueListener dispatcher closed after finishing.")
		}()
		for {
			select {
			case <-ctx.Done(): // 优雅退出
				global.Logger.Info("StartTaskQueueListener dispatcher stopped before processing tasks.")
				return
			default:
				taskQueueID, err := s.QueueService.PopTaskFromQueue(ctx)
				if err != nil {
					select {
					case <-ctx.Done(): // 退出时终止sleep
						global.Logger.Info("StartTaskQueueListener dispatcher stopped during sleep.")
						return
					case <-time.After(time.Second):
						global.Logger.Debug("StartTaskQueueListener sleep 1 second")
					}
					continue
				}
				select {
				case taskChan <- taskQueueID:
				case <-ctx.Done(): // 防止在退出时向关闭的 channel 发送数据
					global.Logger.Info("StartTaskQueueListener dispatcher stopped while dispatching task.")
					return
				}
			}
		}
	}()
}

func (s *TaskScheduler) StartTaskQueueMonitor(ctx context.Context) {
	s.QueueService.RecoverLostTasks()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // 优雅退出
			global.Logger.Info("StartTaskQueueMonitor stopped.")
			return
		case <-ticker.C:
			s.QueueService.RecoverLostTasks()
		}
	}
}
