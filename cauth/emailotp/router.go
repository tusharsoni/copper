package emailotp

import (
	"net/http"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

type NewRouterParams struct {
	Auth   Svc
	RW     chttp.ReaderWriter
	Logger clogger.Logger
}

func NewRouter(p NewRouterParams) chttp.Router {
	ro := &Router{
		rw:     p.RW,
		logger: p.Logger,

		auth: p.Auth,
	}

	return chttp.NewRouter([]chttp.Route{
		{
			Path:    "/api/auth/email-otp/signup",
			Methods: []string{http.MethodPost},
			Handler: http.HandlerFunc(ro.HandleSignup),
		},
		{
			Path:    "/api/auth/email-otp/login",
			Methods: []string{http.MethodPost},
			Handler: http.HandlerFunc(ro.HandleLogin),
		},
	})
}

type Router struct {
	rw     chttp.ReaderWriter
	logger clogger.Logger

	auth Svc
}

func (ro *Router) HandleSignup(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email" valid:"email,required"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	c, token, err := ro.auth.Signup(r.Context(), body.Email)
	if err != nil {
		ro.logger.Error("Failed to sign up with email", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, map[string]interface{}{
		"user_uuid":     c.UserUUID,
		"session_token": token,
	})
}

func (ro *Router) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email            string `json:"email" valid:"email,required"`
		VerificationCode uint   `json:"verification_code" valid:"required"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	c, sessionToken, err := ro.auth.Login(r.Context(), body.Email, body.VerificationCode)
	if err != nil {
		ro.logger.WithTags(map[string]interface{}{
			"email": body.Email,
		}).Error("Failed to login with email", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, map[string]string{
		"user_uuid":     c.UserUUID,
		"session_token": sessionToken,
	})
}
