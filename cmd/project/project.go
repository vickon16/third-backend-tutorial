package project

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	mw "github.com/vickon16/third-backend-tutorial/cmd/middleware"
	"github.com/vickon16/third-backend-tutorial/cmd/sqlc"
	"github.com/vickon16/third-backend-tutorial/cmd/types"
	"github.com/vickon16/third-backend-tutorial/cmd/utils"
)

type Handler struct {
	db *sqlc.Queries
}

func NewHandler(db *sqlc.Queries) *Handler {
	return &Handler{db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	authChain := alice.New(mw.Log, mw.JWT(h.db))

	router.Handle("/projects", authChain.ThenFunc(h.createProject)).Methods(http.MethodPost)
	router.Handle("/projects", authChain.ThenFunc(h.getProjects)).Methods(http.MethodGet)
	router.Handle("/projects/{id}", authChain.ThenFunc(h.getProject)).Methods(http.MethodGet)
	router.Handle("/projects/{id}", authChain.ThenFunc(h.updateProject)).Methods(http.MethodPut)
	router.Handle("/projects/{id}", authChain.ThenFunc(h.deleteProject)).Methods(http.MethodDelete)
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProjectPayload
	if err := utils.ParseJsonAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	claims := r.Context().Value("claims").(*types.JWTClaims)
	userId := claims.UserId

	description, repoURL, siteURL, dependencies := sql.NullString{}, sql.NullString{}, sql.NullString{}, sql.NullString{}
	if payload.Description != "" {
		description = sql.NullString{String: payload.Description, Valid: true}
	}

	if payload.RepoURL != "" {
		repoURL = sql.NullString{String: payload.RepoURL, Valid: true}
	}

	if payload.SiteURL != "" {
		siteURL = sql.NullString{String: payload.SiteURL, Valid: true}
	}

	if payload.Dependencies != "" {
		dependencies = sql.NullString{String: payload.Dependencies, Valid: true}
	}

	err := h.db.CreateProject(r.Context(), sqlc.CreateProjectParams{
		ID:           uuid.New().String(),
		Userid:       userId,
		Name:         payload.Name,
		Description:  description,
		Repourl:      repoURL,
		Siteurl:      siteURL,
		Status:       sqlc.ProjectsStatus(payload.Status),
		Dependencies: dependencies,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "Project created successfully",
	})
}

func (h *Handler) getProjects(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*types.JWTClaims)
	userId := claims.UserId

	projects, err := h.db.GetProjectsByUserId(r.Context(), userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	newProjects := []types.Project{}

	for _, newProject := range projects {
		newProjects = append(newProjects, types.Project{
			ID:           newProject.ID,
			UserId:       newProject.Userid,
			Name:         newProject.Name,
			Description:  newProject.Dependencies.String,
			RepoURL:      newProject.Repourl.String,
			SiteURL:      newProject.Siteurl.String,
			Status:       sqlc.ProjectsStatus(newProject.Status),
			Dependencies: newProject.Dependencies.String,
		})
	}

	utils.WriteJSON(w, http.StatusOK, newProjects)
}

func (h *Handler) customGetProject(r *http.Request, id string) (*types.Project, error) {
	project, err := h.db.GetProjectById(r.Context(), id)
	if err != nil {
		return nil, err
	}

	return &types.Project{
		ID:           project.ID,
		UserId:       project.Userid,
		Name:         project.Name,
		Description:  project.Dependencies.String,
		RepoURL:      project.Repourl.String,
		SiteURL:      project.Siteurl.String,
		Status:       sqlc.ProjectsStatus(project.Status),
		Dependencies: project.Dependencies.String,
	}, nil
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid params"))
		return
	}

	projects, err := h.customGetProject(r, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, projects)
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid params"))
		return
	}

	var payload types.UpdateProjectPayload
	if err := utils.ParseJsonAndValidate(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	claims := r.Context().Value("claims").(*types.JWTClaims)
	userId := claims.UserId

	// check if project exist first
	currentProject, err := h.db.GetProjectById(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("project does not exist"))
		return
	}

	if payload.Name == "" {
		payload.Name = currentProject.Name
	}
	if payload.Status == "" {
		payload.Status = currentProject.Status
	}

	description := utils.AssignNullString(payload.Description, currentProject.Description)
	repoURL := utils.AssignNullString(payload.RepoURL, currentProject.Repourl)
	siteURL := utils.AssignNullString(payload.SiteURL, currentProject.Siteurl)
	dependencies := utils.AssignNullString(payload.Dependencies, currentProject.Dependencies)

	if err := h.db.UpdateProject(r.Context(), sqlc.UpdateProjectParams{
		ID:           id,
		Userid:       userId,
		Name:         payload.Name,
		Description:  description,
		Repourl:      repoURL,
		Siteurl:      siteURL,
		Status:       payload.Status,
		Dependencies: dependencies,
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	project, err := h.customGetProject(r, id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, project)
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid params"))
		return
	}

	claims := r.Context().Value("claims").(*types.JWTClaims)
	userId := claims.UserId

	// check if project exist first
	if _, err := h.db.GetProjectById(r.Context(), id); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("project does not exist"))
		return
	}

	if err := h.db.DeleteProject(r.Context(), sqlc.DeleteProjectParams{
		ID:     id,
		Userid: userId,
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete user: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Project Deleted successfully",
	})
}
