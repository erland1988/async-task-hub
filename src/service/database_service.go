package service

import (
	"asynctaskhub/common"
	"asynctaskhub/global"
	"asynctaskhub/src/model"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type DatabaseService struct{}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

func (s *DatabaseService) InitDB() {
	// 初始化数据表
	global.DB.Migrator().DropIndex(&model.Config{}, "key")
	if err := global.DB.AutoMigrate(&model.Admin{}, &model.Application{}, &model.Config{}, &model.Log{}, &model.Login{}, &model.Task{}, &model.TaskLog{}, &model.TaskQueue{}); err != nil {
		global.Logger.Error("failed to migrate database", zap.Error(err))
		return
	}

	// 如果是第一次初始化，插入默认数据
	var count int64
	if global.DB.Migrator().HasTable(&model.Admin{}) {
		global.DB.Model(&model.Admin{}).Count(&count)
		if count > 0 {
			global.Logger.Info("database already initialized")
			return
		}
	}

	var configs []model.Config
	keys := NewConfigService().GetKeys()
	for _, key := range keys {
		configs = append(configs, model.Config{
			Key:   key,
			Value: "",
		})
	}
	if err := global.DB.Create(&configs).Error; err != nil {
		global.Logger.Error("failed to insert config", zap.Error(err))
		return
	}

	admins, err := s.createAdmins()
	if err != nil {
		global.Logger.Error("failed to insert admin", zap.Error(err))
		return
	}

	applications, err := s.createApplications(admins)
	if err != nil {
		global.Logger.Error("failed to insert application", zap.Error(err))
		return
	}

	tasks, err := s.createTasks(applications)
	if err != nil {
		global.Logger.Error("failed to insert task", zap.Error(err))
		return
	}
	err = s.createTaskQueues(tasks)
	if err != nil {
		global.Logger.Error("failed to insert task_queue", zap.Error(err))
		return
	}
	global.Logger.Info("database initialized successfully")
}

func (s *DatabaseService) createAdmins() ([]model.Admin, error) {
	var admins []model.Admin
	admins = append(admins, model.Admin{
		Username:  "root",
		Password:  common.HashMD5("123456"),
		Role:      model.GlobalAdmin,
		ExpiresAt: time.Now().AddDate(10, 0, 0),
	})
	admins = append(admins, model.Admin{
		Username:  "admin",
		Password:  common.HashMD5("123456"),
		Role:      model.GlobalAdmin,
		ExpiresAt: time.Now().AddDate(10, 0, 0),
	})
	admins = append(admins, model.Admin{
		Username:  "test",
		Password:  common.HashMD5("123456"),
		Role:      model.AppAdmin,
		ExpiresAt: time.Now().AddDate(10, 0, 0),
	})
	if err := global.DB.Create(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

func (s *DatabaseService) createApplications(admins []model.Admin) ([]model.Application, error) {
	var applications []model.Application

	for _, admin := range admins {
		for i := 0; i < 3; i++ {
			applications = append(applications, model.Application{
				Name:      admin.Username + "的应用" + strconv.Itoa(i+1),
				AppKey:    "key_" + strconv.Itoa(admin.ID) + "_" + strconv.Itoa(i+1),
				AppSecret: "secret",
				AdminID:   admin.ID,
			})
		}
	}
	if err := global.DB.Create(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (s *DatabaseService) createTasks(applications []model.Application) ([]model.Task, error) {
	var tasks []model.Task
	for _, application := range applications {
		for i := 0; i < 3; i++ {
			tasks = append(tasks, model.Task{
				AppID:       application.ID,
				Name:        application.Name + "的任务" + strconv.Itoa(i+1),
				TaskCode:    "code_" + strconv.Itoa(application.ID) + "_" + strconv.Itoa(i+1),
				ExecutorURL: "https://www.bing.com/",
			})
		}
	}
	if err := global.DB.Create(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *DatabaseService) createTaskQueues(tasks []model.Task) error {
	var taskQueues []model.TaskQueue
	for _, task := range tasks {
		for i := 0; i < 3; i++ {
			taskQueues = append(taskQueues, model.TaskQueue{
				AppID:             task.AppID,
				TaskID:            task.ID,
				RelativeDelayTime: 3,
				ExecutionTime:     int(time.Now().Add(3 * time.Second).Unix()),
				ExecutionStatus:   model.TaskQueuePending,
			})
		}
	}
	if err := global.DB.Create(&taskQueues).Error; err != nil {
		return err
	}
	return nil
}
