package middlewares

import (
	"context"
	"fmt"
	"golang-web-core/src/domain"
	sessions_user_int "golang-web-core/src/interface/sessions/user"
	"golang-web-core/src/srv/routes"
	"golang-web-core/src/srv/srverr"
	"net/http"

	"github.com/gorilla/mux"
)

func AuthMiddleware(srvRoutes map[string]routes.Route, router *mux.Router, sessionsRepo domain.SessionRepository) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			// return if method is option because the request isn't going to have any auth headers
			if req.Method == http.MethodOptions {
				next.ServeHTTP(rw, req)
				return
			}

			// figure out the original route template so that we can figure out
			// what the route for something like /media/user/123456 originally was
			muxroute, err := mux.CurrentRoute(req).GetPathTemplate()
			if err != nil {
				srverr.Raise(rw, req, err, http.StatusNotFound)
				return
			}
			// see if the route is in the manually registered routes (if it is not then it is a frontend resource)
			route, ok := srvRoutes[muxroute]
			if !ok {
				next.ServeHTTP(rw, req)
				return
			}

			// attempt to get the current session from the request
			session, err := sessions_user_int.FetchCurrentSession(req, sessionsRepo)
			if err != nil {
				srverr.Raise(rw, req, err, http.StatusInternalServerError)
				return
			}
			// inject the current session into the request context if the session exists
			// we do this so that we don't have to retrieve the session from the database every time we want it
			var ctx context.Context
			if session != nil {
				ctx = context.WithValue(req.Context(), domain.CurrentSessionKey, *session)
			}

			// if the session isn't nil, serve the context that contains the session even if the route
			// doesn't require authentication just in case we want to use it
			if !route.RequiresAuth && session != nil {
				next.ServeHTTP(rw, req.WithContext(ctx))
				return
			}
			if !route.RequiresAuth && session == nil {
				next.ServeHTTP(rw, req)
				return
			}

			// this code will only run if the route requires auth, so if the session is nil, the user can't
			// proceed any further
			if session == nil {
				srverr.Raise(rw, req, fmt.Errorf("unauthorized - you must sign in first"), http.StatusUnauthorized)
				return
			}
			// if route requires admin privs and the user doesn't have it, return no permission
			if route.RequiresAdmin && !session.User.IsAdmin {
				srverr.Raise(rw, req, fmt.Errorf("unauthorized - no permission"), http.StatusForbidden)
				return
			}

			// if all of the checks pass, continue with the request with the session injected into the context
			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}
