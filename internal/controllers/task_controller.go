package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/terlyne/tasks_manager/internal/logger"
	"github.com/terlyne/tasks_manager/internal/models"
	"github.com/terlyne/tasks_manager/internal/repository"
	"gorm.io/gorm"
)

func SetupTaskRouter(db *gorm.DB, r *gin.Engine) {
	log := logger.GetLogger()
	taskGroup := r.Group("/tasks")

	// Ручка для создания models.Task
	taskGroup.POST("/", func(ctx *gin.Context) {
		var task models.Task

		if err := ctx.ShouldBindJSON(&task); err != nil {
			log.Error("Invalid request body", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		createdTask, err := repository.CreateTask(db, &task)
		if err != nil {
			log.Error("Failed to create task", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
			return
		}

		log.Info("Task created successfully", "task_id", createdTask.ID)
		ctx.JSON(http.StatusCreated, createdTask)
	})

	// Ручка для получения всех models.Task
	taskGroup.GET("/", func(ctx *gin.Context) {
		tasks, err := repository.GetAllTasks(db)
		if err != nil {
			log.Error("Failed to retrieve tasks", "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
			return
		}

		log.Info("Retrieved all tasks", "count", len(tasks))
		ctx.JSON(http.StatusOK, tasks)
	})

	// Ручка для получения models.Task по ID
	taskGroup.GET("/:id", func(ctx *gin.Context) {
		id, _ := ctx.Params.Get("id")
		taskID, err := strconv.Atoi(id)

		if err != nil {
			log.Error("Invalid task ID", "id", id, "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		task, err := repository.GetTaskByID(db, uint(taskID))
		if task == nil {
			log.Warn("Task not found", "task_id", taskID)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		} else if err != nil {
			log.Error("Failed to retrieve task", "task_id", taskID, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
			return
		}

		log.Info("Retrieved task", "task_id", taskID)
		ctx.JSON(http.StatusOK, task)
	})

	taskGroup.PUT("/:id", func(ctx *gin.Context) {
		id, _ := ctx.Params.Get("id")

		taskID, err := strconv.Atoi(id)
		if err != nil {
			log.Error("Invalid task ID", "id", id, "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		existingTask, err := repository.GetTaskByID(db, uint(taskID))
		if err != nil || existingTask == nil {
			log.Warn("Task not found", "task_id", taskID)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		var updatedTask models.Task
		if err := ctx.ShouldBindJSON(&updatedTask); err != nil {
			log.Error("Invalid request body", "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if updatedTask.Title != "" {
			existingTask.Title = updatedTask.Title
		}
		if updatedTask.Description != "" {
			existingTask.Description = updatedTask.Description
		}
		existingTask.Done = updatedTask.Done

		result, err := repository.UpdateTask(db, existingTask)
		if err != nil {
			log.Error("Failed to update task", "task_id", taskID, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		log.Info("Task updated successfully", "task_id", taskID)
		ctx.JSON(http.StatusOK, result)
	})

	taskGroup.DELETE(":id", func(ctx *gin.Context) {
		id, _ := ctx.Params.Get("id")
		taskID, err := strconv.Atoi(id)

		if err != nil {
			log.Error("Invalid task ID", "id", id, "error", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		existingTask, err := repository.GetTaskByID(db, uint(taskID))
		if err != nil || existingTask == nil {
			log.Warn("Task not found", "task_id", taskID)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		err = repository.DeleteTask(db, uint(taskID))
		if err != nil {
			log.Error("Failed to delete task", "task_id", taskID, "error", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
			return
		}

		log.Info("Task deleted successfully", "task_id", taskID)
		ctx.JSON(http.StatusOK, "Task deleted successfully")
	})
}
