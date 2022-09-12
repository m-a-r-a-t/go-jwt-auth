package auth

import "github.com/m-a-r-a-t/go-rest-wrap/pkg/errors"

type ForeignServiceError struct {
	errors.BaseError
}

func NewForeignServiceError(msg string) *ForeignServiceError {

	return &ForeignServiceError{
		BaseError: errors.BaseError{
			Type:     "failed.req.to.foreign.service.error",
			Messsage: "Failed request to a foreign service !",
			Errors:   []string{msg},
		},
	}
}

type CreateNewUserError struct {
	errors.BaseError
}

func NewCreateNewUserError(msg string) *CreateNewUserError {

	return &CreateNewUserError{
		BaseError: errors.BaseError{
			Type:     "failed.create.new.user.error",
			Messsage: "Failed creating new user !",
			Errors:   []string{msg},
		},
	}
}
