package repository

import (
	"github.com/terlyne/tasks_manager/internal/config"
	"github.com/terlyne/tasks_manager/internal/logger"
	"github.com/terlyne/tasks_manager/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Создаем псевдоним типа модели задачи
type Task = models.Task

func SetupDatabase() *gorm.DB {
	log := logger.GetLogger()

	dsn := config.GetDatabaseDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic("Failed to connect to database")
	}

	// Здесь мы автоматически создаем ИЛИ обновляем таблицу в БД, которая соответствует модели Task
	db.AutoMigrate(&Task{})
	log.Info("Database migration completed")
	return db
}

func CreateTask(db *gorm.DB, task *models.Task) (*models.Task, error) {
	log := logger.GetLogger()

	// Создаем сущность models.Task в БД
	if err := db.Create(task).Error; err != nil {
		log.Error("Failed to create task in database", "error", err)
		return nil, err
	}

	// Переменная в которой будем хранить созданную models.Task
	var createdTask models.Task
	// Кладем в эту переменную созданную models.Task
	if err := db.First(&createdTask, task.ID).Error; err != nil {
		log.Error("Failed to retrieve created task", "task_id", task.ID, "error", err)
		return nil, err
	}

	log.Info("Task created successfully", "task_id", createdTask.ID)
	return &createdTask, nil
}

func GetAllTasks(db *gorm.DB) ([]models.Task, error) {
	log := logger.GetLogger()

	// Переменная в которой будем хранить все models.Task
	var tasks []models.Task
	// Кладем в эту переменную все models.Task
	if err := db.Find(&tasks).Error; err != nil {
		log.Error("Failed to retrieve all tasks", "error", err)
		return nil, err
	}

	return tasks, nil
}

func GetTaskByID(db *gorm.DB, id uint) (*Task, error) {
	log := logger.GetLogger()

	// Переменная в которой будем хранить найденный models.Task по ID
	var task Task
	// Кладем в эту переменную найденную по ID models.Task
	if err := db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("Task not found", "task_id", id)
			return nil, nil
		}
		log.Error("Failed to retrieve task", "task_id", id, "error", err)
		return nil, err
	}

	log.Info("Retrieved task", "task_id", id)
	return &task, nil
}

func UpdateTask(db *gorm.DB, task *Task) (*Task, error) {
	log := logger.GetLogger()

	if err := db.Save(task).Error; err != nil {
		log.Error("Failed to update task", "task_id", task.ID, "error", err)
		return nil, err
	}

	log.Info("Task updated successfully", "task_id", task.ID)
	return task, nil
}

func DeleteTask(db *gorm.DB, id uint) error {
	log := logger.GetLogger()

	if err := db.Delete(&Task{}, id).Error; err != nil {
		log.Error("Failed to delete task", "task_id", id, "error", err)
		return err
	}

	log.Info("Task deleted successfully", "task_id", id)
	return nil
}
