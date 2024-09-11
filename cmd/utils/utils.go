package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vickon16/third-backend-tutorial/cmd/validator"
)

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type AppError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func ParseJsonAndValidate(r *http.Request, payload interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	// validate the payloads
	if err := validator.NewVal.Struct(payload); err != nil {
		return fmt.Errorf("invalid payload: \n %v", err)
	}

	return nil
}

func WriteError(w http.ResponseWriter, status int, err error) {
	if status > 499 {
		WriteJSON(w, status, AppError{
			Status:  status,
			Message: err.Error(),
			Error:   "Internal server error",
		})
	}

	WriteJSON(w, status, AppError{
		Status:  status,
		Message: err.Error(),
		Error:   "Something Went Wrong",
	})
}

func ParseUUID(id string) (uuid.UUID, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}
	return parsedId, nil
}

func AssignNullString(payloadValue string, currentValue sql.NullString) sql.NullString {
	if payloadValue != "" {
		return sql.NullString{String: payloadValue, Valid: true}
	}
	return currentValue
}
