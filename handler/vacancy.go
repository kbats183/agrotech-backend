package handler

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

var vacancyIDKey = "vacancyID"

func vacancies(router chi.Router) {
	router.Get("/", getAllVacancies)
}

func VacancyContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vacancyIDStr := chi.URLParam(r, "vacancyID")
		if vacancyIDStr == "" {
			_ = render.Render(w, r, ErrorRenderer(fmt.Errorf("vacancy ID is required")))
			return
		}
		vacancyIDInt, err := strconv.Atoi(vacancyIDStr)
		if err != nil {
			_ = render.Render(w, r, ErrorRenderer(fmt.Errorf("vacancy ID mast be number")))
			return
		}
		ctx := context.WithValue(r.Context(), vacancyIDKey, vacancyIDInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllVacancies(w http.ResponseWriter, r *http.Request) {
	vacancies, err := dbInstance.GetAllVacancies()
	if err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, &vacancies); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
	}
}
