package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/todo/internal/domain/task"
	"github.com/yrss1/todo/internal/service/todo"
	"github.com/yrss1/todo/pkg/server/response"
	"github.com/yrss1/todo/pkg/store"
	"strconv"
)

type TaskHandler struct {
	todoService *todo.Service
}

func NewTaskHandler(s *todo.Service) *TaskHandler {
	return &TaskHandler{todoService: s}
}

func (h *TaskHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/tasks")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)

		api.GET("/:id", h.get)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)
	}
}

// list godoc
// @Summary List tasks
// @Description Get all tasks for the current user with optional filtering, sorting, and pagination
// @Tags tasks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param title query string false "Filter tasks by title"
// @Param status query string false "Filter tasks by status"
// @Param sortBy query string false "Field to sort by (e.g., id, title)" Enums(id, title, status)
// @Param sortOrder query string false "Sort order (asc or desc)" Enums(asc, desc)
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of tasks per page" default(10)
// @Success 200 {array} task.Response "List of tasks"
// @Failure 400 {object} response.Object "Bad Request"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /tasks [get]
func (h *TaskHandler) list(c *gin.Context) {
	userID := c.Value("userID").(string)

	// Extract query parameters
	titleFilter := c.Query("title")
	statusFilter := c.Query("status")
	sortBy := c.DefaultQuery("sortBy", "id")
	sortOrder := c.DefaultQuery("sortOrder", "asc")

	// Pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Validate sortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		response.BadRequest(c, errors.New("invalid sortOrder parameter"), nil)
		return
	}

	res, err := h.todoService.ListTasks(c, userID, titleFilter, statusFilter, sortBy, sortOrder, page, limit)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// add godoc
// @Summary Add a task
// @Description Add a new task for the current user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param task body task.Request true "Task request"
// @Success 200 {object} task.Response "Task created successfully"
// @Failure 400 {object} response.Object "Bad Request"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /tasks [post]
func (h *TaskHandler) add(c *gin.Context) {
	userID := c.Value("userID").(string)

	req := task.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	req.UserID = &userID
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.todoService.CreateTask(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// get godoc
// @Summary Get a task
// @Description Get task by ID for the current user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} task.Response "Task details"
// @Failure 404 {object} response.Object "Task not found"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /tasks/{id} [get]
func (h *TaskHandler) get(c *gin.Context) {
	userID := c.Value("userID").(string)
	taskID := c.Param("id")

	res, err := h.todoService.GetTask(c, userID, taskID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, res)
}

// update godoc
// @Summary Update a task
// @Description Update task by ID for the current user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param task body task.Request true "Task request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Object "Bad Request"
// @Failure 404 {object} response.Object "Task not found"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) update(c *gin.Context) {
	userID := c.Value("userID").(string)
	taskID := c.Param("id")
	req := task.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.IsEmpty("update"); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := h.todoService.UpdateTask(c, userID, taskID, req); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, "ok")
}

// delete godoc
// @Summary Delete a task
// @Description Delete task by ID for the current user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {string} string "Task deleted"
// @Failure 404 {object} response.Object "Task not found"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) delete(c *gin.Context) {
	userID := c.Value("userID").(string)
	taskID := c.Param("id")

	if err := h.todoService.DeleteTask(c, userID, taskID); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, "Task deleted")
}
