package task

import (
	"errors"
	"github.com/yrss1/todo/pkg/helpers"
)

type Request struct {
	ID          string  `json:"id"`
	UserID      *string `json:"user_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	DueDate     *string `json:"due_date"`
}

func (s *Request) Validate() error {
	if s.UserID == nil {
		return errors.New("user_id: cannot be blank")
	}

	if s.Title == nil {
		return errors.New("title: cannot be blank")
	}

	if s.Status == nil {
		s.Status = helpers.GetStringPtr("active")
	}

	if s.Status != nil && (*s.Status != "active" && *s.Status != "done") {
		return errors.New("status must be either 'active' or 'done'")
	}

	return nil
}

func (s *Request) IsEmpty(check string) error {
	if check == "update" {
		if s.UserID == nil && s.Title == nil && s.Description == nil && s.DueDate == nil && s.Status == nil {
			return errors.New("data cannot be blank")
		}
		if s.Status != nil && (*s.Status != "active" && *s.Status != "done") {
			return errors.New("status must be either 'active' or 'done'")
		}
	}

	if check == "search" {
		if s.Title == nil && s.Status == nil {
			return errors.New("invalid query")
		}
	}

	return nil
}

type Response struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:     data.ID,
		UserID: *data.UserID,
		Title:  *data.Title,
	}
	if data.Description != nil {
		res.Description = *data.Description
	}
	if data.Status != nil {
		res.Status = *data.Status
	}
	if data.DueDate != nil {
		res.DueDate = *data.DueDate
	}
	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
