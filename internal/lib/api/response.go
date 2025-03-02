package api

import "github.com/go-playground/validator/v10"

type Ans struct {
	Status string `json:"status"`
	Msg    string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "ERROR"
	LengthAlias = 6
)

func Ok() *Ans {
	return &Ans{
		Status: StatusOk,
	}
}

func Error(msg string) *Ans {
	return &Ans{
		Status: StatusError,
		Msg:    msg,
	}
}

func Validation(errs validator.ValidationErrors) Ans {
	var errList string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errList += "required field in not valid: " + err.Field() + "\n"
		case "url":
			errList += "url is not valid: " + err.Field() + "\n"
		default:
			errList += "field is not valid: " + err.Field() + "\n"
		}
	}
	return Ans{
		Status: StatusError,
		Msg:    errList,
	}
}
