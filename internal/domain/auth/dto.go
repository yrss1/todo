package auth

import "errors"

type Request struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (s *Request) Validate() error {
	if s.Email == nil {
		return errors.New("email: cannot be blank")
	}

	if s.Password == nil {
		return errors.New("password: cannot be blank")
	}

	return nil
}
