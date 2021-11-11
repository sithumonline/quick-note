package router

import (
	"github.com/go-chi/chi"

	"github.com/sithumonline/quick-note/api/handler"
	"github.com/sithumonline/quick-note/internal/auth"
	"github.com/sithumonline/quick-note/transact/note"
)

func (o *Router) NoteRouter() chi.Router {
	r := chi.NewRouter()
	noteHandler := handler.NewNoteHandler(note.NewNoteRepo(o.db), *auth.NewAuth(o.db))

	r.Post("/{userId}", noteHandler.CreateNote)
	r.Get("/{userId}", noteHandler.NoteList)
	r.Get("/{userId}/{id}", noteHandler.GetNote)
	r.Put("/{userId}/{id}", noteHandler.UpdateNote)
	r.Delete("/{userId}/{id}", noteHandler.DeleteNote)

	return r
}
