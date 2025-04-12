package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Функция для загрузки переменных окружения
func LoadEnv() {
	// Загрузжаем файл .env и помещаем все значения в ENV для запущенного процесса
	err := godotenv.Load("D:/dev/projects/golang/tasks_manager/.env")
	if err != nil {
		fmt.Println("⚠️ No .env file found, using system environment variables.")
	} else {
		fmt.Println("✅ .env file loaded successfully!")
	}
}

// Функция для получения строки Data Source Name (строка, содержащая информацию для подключения к БД)
func GetDatabaseDSN() string {
	LoadEnv()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}
