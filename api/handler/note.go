package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sithumonline/quick-note/internal/auth"

	"github.com/sithumonline/quick-note/models"
	"github.com/sithumonline/quick-note/transact/note"
)

type NoteHandler struct {
	repo note.NoteRepo
	auth auth.Auth
}

func NewNoteHandler(repo note.NoteRepo, auth auth.Auth) *NoteHandler {
	return &NoteHandler{
		repo: repo,
		auth: auth,
	}
}

// CreateNote
// @Summary Create a new note
// @Tags Note
// @Accept json
// @Produce json
// @Param   id	path	string	true	"ID"
// @Param data body	models.Note	true	"data"
// @Success 200 {string} string	"successfully note created"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router	/note/{id}	[post]
func (p *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	v, httpStatus, err := p.auth.TokenValidation(r, true)
	if !v {
		RespondWithError(w, httpStatus, err.Error())
		return
	}

	t := models.Note{}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := p.repo.Save(&t); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "successfully note created")
}

// NoteList
// @Summary Get note list
// @Tags Note
// @Accept json
// @Produce json
// @Param   id	path	string	true	"ID"
// @Success 200 {object} []models.Note
// @Failure 404 {string} string
// @Router	/note/{id}	[get]
func (p *NoteHandler) NoteList(w http.ResponseWriter, r *http.Request) {
	v, httpStatus, err := p.auth.TokenValidation(r, true)
	if !v {
		RespondWithError(w, httpStatus, err.Error())
		return
	}

	userID := chi.URLParam(r, "userId")
	list, err := p.repo.GetList(userID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, list)
}

// GetNote
// @Summary Get note
// @Tags Note
// @Accept json
// @Produce json
// @Param   id	path	string	true	"ID"
// @Param   id	path	string	true	"ID"
// @Success 200 {object} models.Note
// @Failure 404 {string} string
// @Router /note/{id}/{id} [get]
func (p *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")
	o, err := p.repo.Get(ID, userID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	v, httpStatus, err := p.auth.TokenValidation(r, true)
	if !v && !o.Public {
		RespondWithError(w, httpStatus, err.Error())
		return
	}
	if err != nil && !o.Public {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, o)
}

// DeleteNote
// @Summary Delete note
// @Tags Note
// @Accept json
// @Produce json
// @Param   id	path	string	true	"ID"
// @Param   id	path	string	true	"ID"
// @Success 200 {nil}	nil
// @Failure 404 {string}	string
// @Router /note/{id}/{id} [delete]
func (p *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	v, httpStatus, err := p.auth.TokenValidation(r, true)
	if !v {
		RespondWithError(w, httpStatus, err.Error())
		return
	}

	ID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")
	if err := p.repo.Delete(ID, userID); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusNoContent, nil)
}

// UpdateNote
// @Summary Update note
// @Tags Note
// @Description Update note
// @Accept  json
// @Produce  json
// @Param   id	path	string	true	"ID"
// @Param   id	path	string	true	"ID"
// @Success 200 {string} string	"successfully updated"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /note/{id}/{id} [put]
func (p *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	v, httpStatus, err := p.auth.TokenValidation(r, true)
	if !v {
		RespondWithError(w, httpStatus, err.Error())
		return
	}

	t := models.Note{}

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "userId")
	if err := p.repo.Update(&t, ID, userID); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "successfully updated")
}
