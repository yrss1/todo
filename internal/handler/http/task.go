package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/todo/internal/domain/task"
	"github.com/yrss1/todo/internal/service/todo"
	"github.com/yrss1/todo/pkg/server/response"
	"github.com/yrss1/todo/pkg/store"
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
// @Description Get all tasks for the current user
// @Tags tasks
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} task.Response "List of tasks"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /tasks [get]
func (h *TaskHandler) list(c *gin.Context) {
	userID := c.Value("userID").(string)

	res, err := h.todoService.ListTasks(c, userID)
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
