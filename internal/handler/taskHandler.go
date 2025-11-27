package handler

import (
	"golang-task-manager/internal/dtos"
	"golang-task-manager/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
    service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
    return &TaskHandler{service: s}
}

func (h *TaskHandler) Create(c *gin.Context) {
    userIDStr, _ := c.Get("user_id")
    userID, _ := uuid.Parse(userIDStr.(string))

    var req dtos.CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.service.Create(c, userID, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, resp)
}

func (h *TaskHandler) Get(c *gin.Context) {
    userIDStr, _ := c.Get("user_id")
    userID, _ := uuid.Parse(userIDStr.(string))

    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    resp, err := h.service.GetByID(c, id, userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *TaskHandler) List(c *gin.Context) {
    userIDStr, _ := c.Get("user_id")
    userID, _ := uuid.Parse(userIDStr.(string))

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    status := c.Query("status")
    keyword := c.Query("keyword")
	timeStr := c.Query("due_date")
	var dueDate time.Time
	if timeStr != "" {
		parsed, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"malformed datetime"})
			return
		}
		dueDate = time.Date(parsed.Year(), parsed.Month(), parsed.Day(),0,0,0,0,parsed.Location(),)
	}
    filter := dtos.TaskFilter{
        Status:  status,
        Keyword: keyword,
		DueDate: dueDate,
        Page:    page,
        Limit:   limit,
    }

    resp, err := h.service.List(c, userID, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *TaskHandler) Update(c *gin.Context) {
    userIDStr, _ := c.Get("user_id")
    userID, _ := uuid.Parse(userIDStr.(string))

    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req dtos.UpdateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.service.Update(c, id, userID, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *TaskHandler) Delete(c *gin.Context) {
    userIDStr, _ := c.Get("user_id")
    userID, _ := uuid.Parse(userIDStr.(string))

    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    if err := h.service.Delete(c, id, userID); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusNoContent, nil)
}

func (h *TaskHandler) GetPublic(c *gin.Context) {
    slug := c.Param("slug")

    resp, err := h.service.GetPublic(c, slug)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}
