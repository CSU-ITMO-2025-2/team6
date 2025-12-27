package app

import (
	"context"

	"github.com/go-chi/chi/v5"
)

func (a *App) initRouter(ctx context.Context) error {
	a.router = chi.NewRouter()

	//user := a.serviceProvider.UserImpl(ctx)

	// users
	{
		//a.router.Post("/users", user.Create)
		//a.router.Post("/users/setIsActive", user.Update)
		//a.router.Get("/users/getReview", user.GetUserReviews)
	}

	return nil
}
