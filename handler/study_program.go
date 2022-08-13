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

var studyProgramIDKey = "studyProgramID"

func studyProgram(router chi.Router) {
	router.Route("/{studyProgramID}/favourite/{userAuth}", func(router chi.Router) {
		router.Use(StudyProgramContext)
		router.Use(UserContext)
		router.Get("/", getStudyProgramByID)
		router.Post("/", addStudyProgramFavourite)
		router.Delete("/", deleteStudyProgramFavourite)
	})
	router.Route("/profession/{professionID}", func(router chi.Router) {
		router.Use(ProfessionContext)
		router.Get("/", getStudyProgramsByProfessionID)
	})
	router.Route("/favourite/{userAuth}", func(router chi.Router) {
		router.Use(UserContext)
		router.Get("/", getStudyProgramsForFavouriteProfessions)
	})
}

func StudyProgramContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		programIDStr := chi.URLParam(r, "studyProgramID")
		if programIDStr == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("program ID is required")))
			return
		}
		programIDInt, err := strconv.Atoi(programIDStr)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("program ID mast be number")))
			return
		}
		ctx := context.WithValue(r.Context(), studyProgramIDKey, programIDInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getStudyProgramsByProfessionID(w http.ResponseWriter, r *http.Request) {
	professionID := r.Context().Value(professionID).(int)
	programs, err := dbInstance.GetAllProgramsByProfessionID(professionID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &programs); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getStudyProgramsForFavouriteProfessions(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByAuthInContext(w, r)
	if err != nil {
		return
	}
	programs, err := dbInstance.GetAllProgramsForFavouriteProfessions(user.ID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &programs); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getStudyProgramByID(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByAuthInContext(w, r)
	if err != nil {
		return
	}
	programID := r.Context().Value(studyProgramIDKey).(int)
	program, err := dbInstance.GetStudyProgramByID(programID, user.ID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &program); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func addStudyProgramFavourite(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByAuthInContext(w, r)
	if err != nil {
		return
	}

	programID := r.Context().Value(studyProgramIDKey).(int)
	err = dbInstance.AddStudyProgramFavourite(user.ID, programID)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
	} else {
		render.Status(r, http.StatusOK)
	}
}

func deleteStudyProgramFavourite(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByAuthInContext(w, r)
	if err != nil {
		return
	}

	programID := r.Context().Value(studyProgramIDKey).(int)
	err = dbInstance.DeleteStudyProgramFavourite(user.ID, programID)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
	} else {
		render.Status(r, http.StatusOK)
	}
}
