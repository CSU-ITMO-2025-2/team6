package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) initRouter(ctx context.Context) error {
	a.router = chi.NewRouter()

	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.Logger)

	//user := a.serviceProvider.UserImpl(ctx)
	study := a.serviceProvider.StudyImpl(ctx)

	// users
	{
		//a.router.Post("/users", user.Create)
		//a.router.Post("/users/setIsActive", user.Update)
		//a.router.Get("/users/getReview", user.GetUserReviews)
	}

	// study
	{
		a.router.Post("/study", study.Create)
	}

	return nil
}
