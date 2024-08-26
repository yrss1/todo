package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/todo/internal/domain/task"
	"github.com/yrss1/todo/internal/service/todo"
	"github.com/yrss1/todo/pkg/helpers"
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

		api.GET("/search", h.search)
	}
}

// list godoc
// @Summary List tasks
// @Description Get all tasks
// @Tags tasks
// @Accept  json
// @Produce  json
// @Success 200 {array} task.Response
// @Failure 500 {object} response.Object
// @Router /tasks [get]
func (h *TaskHandler) list(c *gin.Context) {
	res, err := h.todoService.ListTasks(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// add godoc
// @Summary Add a task
// @Description Add a new task
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param task body task.Request true "Task request"
// @Success 200 {object} task.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /tasks [post]
func (h *TaskHandler) add(c *gin.Context) {
	req := task.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
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
// @Description Get task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path string true "Task ID"
// @Success 200 {object} task.Response
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /tasks/{id} [get]
func (h *TaskHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.todoService.GetTask(c, id)
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
// @Description Update task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path string true "Task ID"
// @Param task body task.Request true "Task request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Object
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /tasks/{id} [put]
func (h *TaskHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := task.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.IsEmpty("update"); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := h.todoService.UpdateTask(c, id, req); err != nil {
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
// @Description Delete task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path string true "Task ID"
// @Success 200 {string} string "Task deleted"
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /tasks/{id} [delete]
func (h *TaskHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.todoService.DeleteTask(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, id)
}

// search godoc
// @Summary Search tasks
// @Description Search tasks by name or email
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param title query string false "Title"
// @Param status query string false "Status"
// @Success 200 {array} task.Response
// @Failure 400 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /tasks/search [get]
func (h *TaskHandler) search(c *gin.Context) {
	req := task.Request{
		Title:  helpers.GetStringPtr(c.Query("title")),
		Status: helpers.GetStringPtr(c.Query("status")),
	}

	if err := req.IsEmpty("search"); err != nil {
		response.BadRequest(c, errors.New("incorrect query"), nil)
		return
	}

	res, err := h.todoService.SearchTask(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}
