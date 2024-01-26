package middleware

import (
	"context"
	"github.com/hansengotama/disbursement/internal/lib/httphelper"
	"github.com/hansengotama/disbursement/internal/lib/postgres"
	accesstokenrepo "github.com/hansengotama/disbursement/internal/repository/accesstoken"
	"net/http"
	"strings"
	"time"
)

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" || !strings.HasPrefix(authorizationHeader, "Bearer ") {
			httphelper.Response(w, httphelper.HTTPResponse{
				Code:       http.StatusUnauthorized,
				ErrMessage: httphelper.ErrUnauthorized.Error(),
			})
			return
		}

		token := strings.TrimPrefix(authorizationHeader, "Bearer ")
		repo := accesstokenrepo.AccessTokenDB{}
		res, err := repo.GetAccessToken(accesstokenrepo.GetAccessTokenParam{
			Context:  r.Context(),
			Executor: postgres.GetConnection(),
			Token:    token,
		})
		if err != nil || res == nil || res.UserID == 0 {
			httphelper.Response(w, httphelper.HTTPResponse{
				Code:       http.StatusUnauthorized,
				ErrMessage: httphelper.ErrUnauthorized.Error(),
			})
			return
		}

		if res.ExpirationTime.Before(time.Now()) {
			httphelper.Response(w, httphelper.HTTPResponse{
				Code:       http.StatusUnauthorized,
				ErrMessage: "Token has expired",
			})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", res.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
