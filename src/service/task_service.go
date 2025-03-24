package service

import (
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"errors"
	"go.uber.org/zap"
)

type TaskService struct{}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s *TaskService) DeleteQueue(queueID int) error {
	if queueID == 0 {
		return errors.New("参数异常")
	}
	var queue model.TaskQueue
	if err := global.DB.First(&queue, queueID).Error; err != nil {
		global.Logger.Warn("DeleteQueue", zap.Error(err))
		return err
	}
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Logger.Warn("DeleteQueue", zap.Error(tx.Error))
		return tx.Error
	}
	if err := tx.Delete(&queue).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteQueue", zap.Error(err))
		return err
	}
	if err := tx.Where("queue_id =?", queue.ID).Delete(&model.TaskLog{}).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteQueue", zap.Error(err))
		return err
	}
	if err := tx.Commit().Error; err != nil {
		global.Logger.Warn("DeleteQueue", zap.Error(err))
		return err
	}
	return nil
}
func (s *TaskService) DeleteTask(taskID int) error {
	if taskID == 0 {
		return errors.New("参数异常")
	}
	var task model.Task
	if err := global.DB.First(&task, taskID).Error; err != nil {
		global.Logger.Warn("DeleteTask", zap.Error(err))
		return err
	}
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Logger.Warn("DeleteTask", zap.Error(tx.Error))
		return tx.Error
	}
	if err := tx.Delete(&task).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteTask", zap.Error(err))
		return err
	}
	if err := tx.Where("task_id =?", task.ID).Delete(&model.TaskQueue{}).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteTask", zap.Error(err))
		return err
	}
	if err := tx.Where("task_id =?", task.ID).Delete(&model.TaskLog{}).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteTask", zap.Error(err))
		return err
	}
	if err := tx.Commit().Error; err != nil {
		global.Logger.Warn("DeleteTask", zap.Error(err))
		return err
	}
	return nil
}
func (s *TaskService) DeleteApp(appID int) error {
	if appID == 0 {
		return errors.New("参数异常")
	}
	var application model.Application
	if err := global.DB.First(&application, appID).Error; err != nil {
		global.Logger.Warn("DeleteApp", zap.Error(err))
		return err
	}
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Logger.Warn("DeleteApp", zap.Error(tx.Error))
		return tx.Error
	}
	if err := tx.Delete(&application).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteApp", zap.Error(err))
		return err
	}
	if err := tx.Where("app_id =?", application.ID).Delete(&model.Task{}).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteApp", zap.Error(err))
		return err
	}
	if err := tx.Where("app_id =?", application.ID).Delete(&model.TaskQueue{}).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteApp", zap.Error(err))
		return err
	}
	if err := tx.Where("app_id =?", application.ID).Delete(&model.TaskLog{}).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteApp", zap.Error(err))
		return err
	}
	if err := tx.Commit().Error; err != nil {
		global.Logger.Warn("DeleteApp", zap.Error(err))
		return err
	}
	return nil
}
func (s *TaskService) DeleteAdmin(adminID int) error {
	if adminID == 0 {
		return errors.New("参数异常")
	}
	var admin model.Admin
	if err := global.DB.First(&admin, adminID).Error; err != nil {
		return err
	}
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Logger.Warn("DeleteAdmin", zap.Error(tx.Error))
		return tx.Error
	}
	if err := tx.Delete(&admin).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteAdmin", zap.Error(err))
		return err
	}
	var applications []model.Application
	if err := tx.Where("admin_id =?", admin.ID).Find(&applications).Error; err != nil {
		tx.Rollback()
		global.Logger.Warn("DeleteAdmin", zap.Error(err))
		return err
	}
	for _, app := range applications {
		err := s.DeleteApp(app.ID)
		if err != nil {
			tx.Rollback()
			global.Logger.Warn("DeleteAdmin", zap.Error(err))
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		global.Logger.Warn("DeleteAdmin", zap.Error(err))
		return err
	}
	return nil
}
