package main

import (
	"github.com/gin-gonic/gin"
	"github.com/terlyne/tasks_manager/internal/config"
	"github.com/terlyne/tasks_manager/internal/controllers"
	"github.com/terlyne/tasks_manager/internal/logger"
	"github.com/terlyne/tasks_manager/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Функция для создания мрашрутизатора
func SetupRouter(db *gorm.DB) *gin.Engine {

	// Создаем маршрутизатор без предустановленных middleware
	r := gin.New()

	// Добавляем Recovery middleware для обработки паники
	r.Use(gin.Recovery())

	// Добавляем собственный middleware для логирования HTTP-запросов
	r.Use(logger.GinLogger())

	// Middleware для добавления базы данных в контекст
	r.Use(func(ctx *gin.Context) {
		// Устанавливаем объект БД в контекст, обработчики получают доступ к БД через контекст gin.Context
		ctx.Set("db", db)
		// Вызываем следующий middleware или обработчик маршрута в цепочке
		ctx.Next()
	})

	// Создаем маршрутизатор для задач
	controllers.SetupTaskRouter(db, r)

	// Возвращаем объект маршрутизатора
	return r
}

func main() {
	// Инициализируем логгер
	logger.InitLogger("local")
	log := logger.GetLogger()
	log.Info("Starting application")

	// Инициализируем БД
	dsn := config.GetDatabaseDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		panic("Ошибка подключения к БД")
	}
	log.Info("Successfully connected to database")

	// Миграция схемы (создание таблиц)
	db.AutoMigrate(&models.Task{})
	log.Info("Database migration completed")

	r := SetupRouter(db)
	log.Info("Server is starting on :8080")
	r.Run(":8080")
}
