package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MdZunaed/students-api/internal/types"
	"github.com/MdZunaed/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Students api"))
	}
}

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		// For empty request error
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		// For general error
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		// For validation error
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		response.WriteJson(w, http.StatusAccepted, map[string]bool{"success": true})
	}
}
