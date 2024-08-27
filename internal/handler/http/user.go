package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yrss1/todo/internal/domain/user"
	"github.com/yrss1/todo/internal/service/account"
	"github.com/yrss1/todo/pkg/helpers"
	"github.com/yrss1/todo/pkg/server/response"
	"github.com/yrss1/todo/pkg/store"
)

type UserHandler struct {
	accountService *account.Service
}

func NewUserHandler(s *account.Service) *UserHandler {
	return &UserHandler{accountService: s}
}

func (h *UserHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/users")
	{
		api.GET("/", h.list)
		api.POST("/", h.add)

		api.GET("/:id", h.get)
		api.PUT("/:id", h.update)
		api.DELETE("/:id", h.delete)

		api.GET("/search", h.search)
		api.GET("/email", h.getByEmail)
	}
}

// list godoc
// @Summary List users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200 {array} user.Response "List of users"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /users [get]
func (h *UserHandler) list(c *gin.Context) {
	res, err := h.accountService.ListUsers(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// add godoc
// @Summary Add a user
// @Description Add a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param user body user.Request true "User request"
// @Success 200 {object} user.Response "User created successfully"
// @Failure 400 {object} response.Object "Bad Request"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) add(c *gin.Context) {
	req := user.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.accountService.CreateUser(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// get godoc
// @Summary Get a user
// @Description Get user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} user.Response "User details"
// @Failure 404 {object} response.Object "User not found"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /users/{id} [get]
func (h *UserHandler) get(c *gin.Context) {
	id := c.Param("id")

	res, err := h.accountService.GetUser(c, id)
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
// @Summary Update a user
// @Description Update user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body user.Request true "User request"
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.Object "Bad Request"
// @Failure 404 {object} response.Object "User not found"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /users/{id} [put]
func (h *UserHandler) update(c *gin.Context) {
	id := c.Param("id")
	req := user.Request{}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := req.IsEmpty("update"); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	if err := h.accountService.UpdateUser(c, id, req); err != nil {
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
// @Summary Delete a user
// @Description Delete user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted"
// @Failure 404 {object} response.Object "User not found"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /users/{id} [delete]
func (h *UserHandler) delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.accountService.DeleteUser(c, id); err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			response.NotFound(c, err)
		default:
			response.InternalServerError(c, err)
		}
		return
	}

	response.OK(c, "User deleted")
}

// search godoc
// @Summary Search users
// @Description Search users by name or email
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param name query string false "Name"
// @Param email query string false "Email"
// @Success 200 {array} user.Response "List of users matching the search criteria"
// @Failure 400 {object} response.Object "Bad Request"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /users/search [get]
func (h *UserHandler) search(c *gin.Context) {
	req := user.Request{
		Name:  helpers.GetStringPtr(c.Query("name")),
		Email: helpers.GetStringPtr(c.Query("email")),
	}

	if err := req.IsEmpty("search"); err != nil {
		response.BadRequest(c, errors.New("incorrect query"), nil)
		return
	}

	res, err := h.accountService.SearchUser(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// getByEmail godoc
// @Summary Get user by email
// @Description Get user details by email
// @Tags users
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param email query string true "User Email"
// @Success 200 {object} user.Response "User details"
// @Failure 404 {object} response.Object "User not found"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /users/email [get]
func (h *UserHandler) getByEmail(c *gin.Context) {
	email := c.Query("email")

	res, err := h.accountService.GetUserByEmail(c, email)
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
