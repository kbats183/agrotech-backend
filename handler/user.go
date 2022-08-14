package handler

import (
	"context"
	"fmt"
	"githab.com/kbats183/argotech/backend/cv"
	"githab.com/kbats183/argotech/backend/db"
	"githab.com/kbats183/argotech/backend/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

var userAuthKey = "userAuth"

func users(router chi.Router) {
	router.Get("/", getAllUsers)
	router.Post("/", createUser)
	router.Route("/{userAuth}", func(router chi.Router) {
		router.Use(UserContext)
		router.Get("/", getUser)
		router.Put("/", updateUser)
		router.Delete("/", deleteUser)
		router.Put("/profile", updateUserProfile)
		router.Route("/cv", func(router chi.Router) {
			router.Get("/", getCVByLogin)
			router.Post("/", updateCVByLogin)
			router.Get("/document.rtf", getCVByLoginRTF)
		})
	})
}

func UserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAuth := chi.URLParam(r, "userAuth")
		if userAuth == "" {
			_ = render.Render(w, r, ErrorRenderer(fmt.Errorf("user Auth is required")))
			return
		}
		ctx := context.WithValue(r.Context(), userAuthKey, userAuth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := dbInstance.GetAllUsers()
	if err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, users); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	userData := models.UserData{}
	if err := render.Bind(r, &userData); err != nil {
		_ = render.Render(w, r, ErrBadRequest)
		return
	}
	user := &models.User{UserData: userData}
	if err := dbInstance.AddUser(user); err != nil {
		_ = render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, user); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userAuth := r.Context().Value(userAuthKey).(string)
	user, err := dbInstance.GetUserByAuth(userAuth)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &user); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userAuthKey).(int)
	userData := models.UserData{}
	if err := render.Bind(r, &userData); err != nil {
		_ = render.Render(w, r, ErrBadRequest)
		return
	}
	user, err := dbInstance.UpdateUser(userID, userData)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &user); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(userAuthKey).(int)
	err := dbInstance.DeleteUser(userId)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
}

func updateUserProfile(w http.ResponseWriter, r *http.Request) {
	userAuth := r.Context().Value(userAuthKey).(string)
	userProfile := models.UserProfile{}
	if err := render.Bind(r, &userProfile); err != nil {
		_ = render.Render(w, r, ErrBadRequest)
		return
	}
	user, err := dbInstance.UpdateUserProfile(userAuth, userProfile)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &user); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getUserByAuthInContext(w http.ResponseWriter, r *http.Request) (models.User, error) {
	userAuth := r.Context().Value(userAuthKey).(string)
	user, err := dbInstance.GetUserByAuth(userAuth)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ErrorRenderer(err))
		}
	}
	return user, err
}

func getCVByLogin(w http.ResponseWriter, r *http.Request) {
	userAuth := r.Context().Value(userAuthKey).(string)
	user, err := dbInstance.GetCVByLogin(userAuth)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &user); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func updateCVByLogin(w http.ResponseWriter, r *http.Request) {
	userAuth := r.Context().Value(userAuthKey).(string)
	userCV := models.UserCV{}
	if err := render.Bind(r, &userCV); err != nil {
		_ = render.Render(w, r, ErrBadRequest)
		return
	}
	err := dbInstance.UpdateCV(userAuth, userCV)
	if err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	user, err := dbInstance.GetCVByLogin(userAuth)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &user); err != nil {
		_ = render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getCVByLoginRTF(w http.ResponseWriter, r *http.Request) {
	userAuth := r.Context().Value(userAuthKey).(string)
	cvDetails, err := dbInstance.GetCVByLogin(userAuth)
	bytes := cv.GenerateResume(cvDetails)
	if err != nil {
		if err == db.ErrNoMatch {
			_ = render.Render(w, r, ErrNotFound)
		} else {
			_ = render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	w.Header().Set("Content-Type", "application/rtf")
	_, _ = w.Write(bytes)
}
