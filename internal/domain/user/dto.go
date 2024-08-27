package user

import (
	"errors"
)

type Request struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (s *Request) Validate() error {
	if s.Name == nil {
		return errors.New("name: cannot be blank")
	}

	if s.Email == nil {
		return errors.New("email: cannot be blank")
	}

	if s.Password == nil {
		return errors.New("password: cannot be blank")
	}

	return nil
}

func (s *Request) IsEmpty(check string) error {
	if check == "update" {
		if s.Name == nil && s.Email == nil && s.Password == nil {
			return errors.New("data cannot be blank")
		}
	}

	if check == "search" {
		if s.Name == nil && s.Email == nil {
			return errors.New("invalid query")
		}
	}

	return nil
}

type Response struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:    data.ID,
		Name:  *data.Name,
		Email: *data.Email,
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
