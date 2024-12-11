package middleware

import (
	"dockman/app"
	"dockman/pages"
	"github.com/go-chi/chi/v5"
	"github.com/maddalax/htmgo/framework/h"
	"net/http"
)

func UseLoginRequiredMiddleware(router *chi.Mux) {
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedPaths := []string{
				"/login",
				"/logout",
				h.GetPartialPath(pages.RegisterUser),
				h.GetPartialPath(pages.LoginUser)}

			for _, path := range allowedPaths {
				if r.URL.Path == path {
					handler.ServeHTTP(w, r)
					return
				}
			}

			ctx := h.GetRequestContext(r)
			user, err := app.ValidateSession(ctx)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			ctx.Set("user", user)
			handler.ServeHTTP(w, r)
		})
	})
}
