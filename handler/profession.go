package handler

import (
	"context"
	"fmt"
	"githab.com/kbats183/argotech/backend/db"
	"githab.com/kbats183/argotech/backend/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

var professionID = "professionID"

func professions(router chi.Router) {
	router.Get("/", getAllProfessions)
	router.Route("/favourite/{userAuth}", func(router chi.Router) {
		router.Use(UserContext)
		router.Get("/", getAllProfessionsWithRating)
	})
	router.Route("/{professionID}", func(router chi.Router) {
		router.Use(ProfessionContext)
		router.Get("/", getProfession)
		router.Route("/favourite/{userAuth}", func(router chi.Router) {
			router.Use(UserContext)
			router.Get("/", getProfessionWithRating)
			router.Post("/", addProfessionFavourite)
			router.Delete("/", deleteProfessionFavourite)
		})
	})
}

func ProfessionContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		professionIDStr := chi.URLParam(r, "professionID")
		if professionIDStr == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("profession ID is required")))
			return
		}
		professionIDInt, err := strconv.Atoi(professionIDStr)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("profession ID mast be number")))
			return
		}
		ctx := context.WithValue(r.Context(), professionID, professionIDInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllProfessions(w http.ResponseWriter, r *http.Request) {
	professions, err := dbInstance.GetAllProfession()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, &professions); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getProfession(w http.ResponseWriter, r *http.Request) {
	professionID := r.Context().Value(professionID).(int)
	profession, err := dbInstance.GetProfessionByID(professionID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &profession); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getAllProfessionsWithRating(w http.ResponseWriter, r *http.Request) {
	userAuth := r.Context().Value(userAuthKey).(string)
	professions, err := dbInstance.GetAllProfessionWithRating(models.UserAuth(userAuth))
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, &professions); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getProfessionWithRating(w http.ResponseWriter, r *http.Request) {
	professionID := r.Context().Value(professionID).(int)
	userAuth := r.Context().Value(userAuthKey).(string)
	profession, err := dbInstance.GetProfessionWithRatingByID(models.UserAuth(userAuth), professionID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &profession); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func addProfessionFavourite(w http.ResponseWriter, r *http.Request) {
	userAuth := r.Context().Value(userAuthKey).(string)
	user, err := dbInstance.GetUserByAuth(userAuth)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}

	professionID := r.Context().Value(professionID).(int)
	err = dbInstance.AddProfessionFavourite(user.ID, professionID)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
	} else {
		render.Status(r, http.StatusOK)
	}
}

func deleteProfessionFavourite(w http.ResponseWriter, r *http.Request) {
	userAuth := r.Context().Value(userAuthKey).(string)
	user, err := dbInstance.GetUserByAuth(userAuth)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}

	professionID := r.Context().Value(professionID).(int)
	err = dbInstance.DeleteProfessionFavourite(user.ID, professionID)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
	} else {
		render.Status(r, http.StatusOK)
	}
}
