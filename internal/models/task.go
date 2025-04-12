package models

import "gorm.io/gorm"

// Структура для нашей модели, данную модель мы будем использовать как сущность для gorm,
// а также будем принимать на вход такой json
type Task struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}
