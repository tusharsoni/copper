package cauth

import (
	"context"
	"errors"
	"net/http"

	"github.com/tusharsoni/copper/cacl"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

type SessionMiddleware chttp.MiddlewareFunc

type NewSessionMiddlewareParams struct {
	RW     chttp.ReaderWriter
	Svc    Svc
	Logger clogger.Logger
	ACL    cacl.Svc
}

func NewSessionMiddleware(p NewSessionMiddlewareParams) SessionMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userUUID, sessionToken, ok := r.BasicAuth()
			if !ok || userUUID == "" || sessionToken == "" {
				p.RW.Unauthorized(w)
				return
			}

			ok, err := p.Svc.VerifySessionToken(r.Context(), userUUID, sessionToken)
			if err != nil {
				p.Logger.WithTags(map[string]interface{}{
					"userUUID": userUUID,
				}).Error("Failed to verify session token", err)
				p.RW.InternalErr(w)
				return
			}
			if !ok {
				p.RW.Unauthorized(w)
				return
			}

			impersonatedUserUUID := r.Header.Get("x-user-uuid")
			if impersonatedUserUUID != "" {
				if p.ACL == nil {
					p.Logger.Error("Failed to impersonate user", errors.New("acl is not configured"))
					p.RW.InternalErr(w)
					return
				}

				ok, err := p.ACL.UserHasPermission(r.Context(), userUUID, "cauth/session", "impersonate")
				if err != nil {
					p.Logger.WithTags(map[string]interface{}{
						"userUUID": userUUID,
					}).Error("Failed to impersonate user", err)
					p.RW.InternalErr(w)
					return
				}

				if !ok {
					p.Logger.WithTags(map[string]interface{}{
						"userUUID": userUUID,
					}).Error("Failed to impersonate user", errors.New("user does not have permission to impersonate"))
					p.RW.InternalErr(w)
					return
				}

				userUUID = impersonatedUserUUID
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeySession, userUUID)))
		})
	}
}
