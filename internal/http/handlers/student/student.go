package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/0xNaim/students-api/internal/storage"
	"github.com/0xNaim/students-api/internal/types"
	"github.com/0xNaim/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(errors.New("request body cannot be empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate the student data
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidateError(validateErrs))
			return
		}

		// Create the student in the storage
		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]any{
			"message": "Student created successfully",
			"student": struct {
				ID    int64  `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
				Age   int    `json:"age"`
			}{
				ID:    lastId,
				Name:  student.Name,
				Email: student.Email,
				Age:   student.Age,
			},
		})

	}
}
