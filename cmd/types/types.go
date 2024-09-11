package types

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vickon16/third-backend-tutorial/cmd/sqlc"
)

type RegisterUserPayload struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=3,max=100"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type JWTClaims struct {
	UserId    string `json:"userId"`
	UserEmail string `json:"userEmail"`
	jwt.RegisteredClaims
}

type Project struct {
	ID           string              `json:"id"`
	UserId       string              `json:"userId"`
	Name         string              `json:"name"`
	Description  string              `json:"description,omitempty"`
	RepoURL      string              `json:"repoURL,omitempty"`
	SiteURL      string              `json:"siteURL,omitempty"`
	Status       sqlc.ProjectsStatus `json:"status"`
	Dependencies string              `json:"dependencies,omitempty"`
}

type CreateProjectPayload struct {
	UserId       string              `json:"userId" validate:"required"`
	Name         string              `json:"name" validate:"required"`
	Description  string              `json:"description" validate:"omitempty"`
	RepoURL      string              `json:"repoURL" validate:"omitempty,url"`
	SiteURL      string              `json:"siteURL" validate:"omitempty,url"`
	Status       sqlc.ProjectsStatus `json:"status" validate:"required,project_status_enum"`
	Dependencies string              `json:"dependencies" validate:"omitempty"`
}

type UpdateProjectPayload struct {
	Name         string              `json:"name" validate:"omitempty"`
	Description  string              `json:"description" validate:"omitempty"`
	RepoURL      string              `json:"repoURL" validate:"omitempty,url"`
	SiteURL      string              `json:"siteURL" validate:"omitempty,url"`
	Status       sqlc.ProjectsStatus `json:"status" validate:"omitempty,project_status_enum"`
	Dependencies string              `json:"dependencies" validate:"omitempty"`
}
