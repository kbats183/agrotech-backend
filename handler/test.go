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

var testIDKey = "testID"

func tests(router chi.Router) {
	router.Route("/{testID}", func(router chi.Router) {
		router.Use(TestContext)
		router.Get("/", getTest)
		router.Route("/answers/{userAuth}", func(router chi.Router) {
			router.Use(UserContext)
			router.Get("/", getTestAnswers)
			router.Get("/count", getTestAnswersCount)
			router.Post("/", addTestAnswer)
		})
		router.Route("/professions/{userAuth}", func(router chi.Router) {
			router.Use(UserContext)
			router.Get("/", getProfessionsByTest)
		})
	})
}

func TestContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testIDStr := chi.URLParam(r, "testID")
		if testIDStr == "" {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("test ID is required")))
			return
		}
		testIDInt, err := strconv.Atoi(testIDStr)
		if err != nil {
			render.Render(w, r, ErrorRenderer(fmt.Errorf("test ID mast be number")))
			return
		}
		ctx := context.WithValue(r.Context(), testIDKey, testIDInt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTest(w http.ResponseWriter, r *http.Request) {
	testID := r.Context().Value(testIDKey).(int)
	test, err := dbInstance.GetTestByID(testID)
	if err != nil {
		if err == db.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &test); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getTestAnswers(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByAuthInContext(w, r)
	if err != nil {
		return
	}
	testID := r.Context().Value(testIDKey).(int)

	answers, err := dbInstance.GetTestAnswers(testID, user.ID)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, &answers); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func addTestAnswer(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByAuthInContext(w, r)
	if err != nil {
		return
	}
	answer := models.TestAnswerUserData{}
	if err := render.Bind(r, &answer); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err = dbInstance.AddTestAnswer(user.ID, answer)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	render.Status(r, http.StatusOK)
}

func getTestAnswersCount(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByAuthInContext(w, r)
	if err != nil {
		return
	}
	testID := r.Context().Value(testIDKey).(int)

	count, err := dbInstance.GetTestAnswersCount(testID, user.ID)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, &count); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}

func getProfessionsByTest(w http.ResponseWriter, r *http.Request) {
	user, err := getUserByAuthInContext(w, r)
	if err != nil {
		return
	}
	testID := r.Context().Value(testIDKey).(int)

	professions, err := dbInstance.GetProfessionByTest(user.ID, testID)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, &professions); err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
}
