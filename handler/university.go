package handler

import (
	"context"
	"fmt"
	"githab.com/kbats183/argotech/backend/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

var universityIDKey = "universityID"

func university(router chi.Router) {
	router.Get("/", getAllUniversity)
	router.Route("/{universityID}", func(router chi.Router) {
		router.Use(UniversityContext)
		router.Get("/", getUniversityByID)
	})
}

func UniversityContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		universityIDSStr := chi.URLParam(r, "universityID")
		if universityIDSStr == "" {
			_ = render.Render(w, r, ErrorRenderer(fmt.Errorf("university ID is required")))
			return
		}
		universityIDIInt, err := strconv.Atoi(universityIDSStr)
		if err != nil {
			_ = render.Render(w, r, ErrorRenderer(fmt.Errorf("university ID mast be number")))
			return
		}
		ctx := context.WithValue(r.Context(), universityIDKey, universityIDIInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllUniversity(w http.ResponseWriter, r *http.Request) {
	university, err := dbInstance.GetAllUniversity()
	if err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, &university); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
	}
}

func getUniversityByID(w http.ResponseWriter, r *http.Request) {
	universityID := r.Context().Value(universityIDKey).(int)
	profession, err := dbInstance.GetUniversityByID(universityID)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &profession); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
