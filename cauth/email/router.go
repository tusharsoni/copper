package email

import (
	"net/http"

	"github.com/tusharsoni/copper/cerror"
	"gorm.io/gorm"

	"github.com/tusharsoni/copper/cauth"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
)

type Router struct {
	rw     chttp.ReaderWriter
	logger clogger.Logger

	auth   Svc
	config Config
}

type NewRouterParams struct {
	RW     chttp.ReaderWriter
	Logger clogger.Logger

	Auth      Svc
	SessionMW chttp.MiddlewareFunc
	Config    Config
}

func NewRouter(p NewRouterParams) chttp.Router {
	ro := &Router{
		rw:     p.RW,
		logger: p.Logger,
		auth:   p.Auth,
		config: p.Config,
	}

	return chttp.NewRouter([]chttp.Route{
		{
			Path:    "/api/auth/email/signup",
			Methods: []string{http.MethodPost},
			Handler: http.HandlerFunc(ro.Signup),
		},
		{
			Path:    "/api/auth/email/login",
			Methods: []string{http.MethodPost},
			Handler: http.HandlerFunc(ro.Login),
		},
		{
			MiddlewareFuncs: []chttp.MiddlewareFunc{p.SessionMW},
			Path:            "/api/auth/email/verify",
			Methods:         []string{http.MethodPost},
			Handler:         http.HandlerFunc(ro.VerifyUser),
		},
		{
			MiddlewareFuncs: []chttp.MiddlewareFunc{p.SessionMW},
			Path:            "/api/auth/email/resend-verification-code",
			Methods:         []string{http.MethodPost},
			Handler:         http.HandlerFunc(ro.ResendVerificationCode),
		},
		{
			Path:    "/api/auth/email/change-password",
			Methods: []string{http.MethodPost},
			Handler: http.HandlerFunc(ro.ChangePassword),
		},
		{
			Path:    "/api/auth/email/reset-password",
			Methods: []string{http.MethodPost},
			Handler: http.HandlerFunc(ro.ResetPassword),
		},
		{
			Path:            "/api/auth/email/credentials",
			MiddlewareFuncs: []chttp.MiddlewareFunc{p.SessionMW},
			Methods:         []string{http.MethodPost},
			Handler:         http.HandlerFunc(ro.HandleAddCredentials),
		},
		{
			Path:            "/api/auth/email/change-email",
			MiddlewareFuncs: []chttp.MiddlewareFunc{p.SessionMW},
			Methods:         []string{http.MethodPost},
			Handler:         http.HandlerFunc(ro.HandleChangeEmail),
		},
		{
			Path:            "/api/auth/email/credentials",
			MiddlewareFuncs: []chttp.MiddlewareFunc{p.SessionMW},
			Methods:         []string{http.MethodGet},
			Handler:         http.HandlerFunc(ro.HandleGetCredentials),
		},
	})
}

func (ro *Router) Signup(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"runelength(4|32)"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	c, sessionToken, err := ro.auth.Signup(r.Context(), body.Email, body.Password)
	if err != nil && err != cauth.ErrUserAlreadyExists {
		ro.logger.Error("Failed to signup user with email and password", err)
		ro.rw.InternalErr(w)
		return
	} else if err == cauth.ErrUserAlreadyExists {
		ro.rw.BadRequest(w, cauth.ErrUserAlreadyExists)
		return
	}

	ro.rw.Created(w, map[string]string{
		"user_uuid":     c.UserUUID,
		"session_token": sessionToken,
	})
}

func (ro *Router) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"runelength(4|32)"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	c, sessionToken, err := ro.auth.Login(r.Context(), body.Email, body.Password)
	if err != nil && err != cauth.ErrInvalidCredentials {
		ro.logger.Error("Failed to login user with email and password", err)
		ro.rw.InternalErr(w)
		return
	} else if err == cauth.ErrInvalidCredentials {
		ro.rw.Unauthorized(w)
		return
	}

	ro.rw.OK(w, map[string]string{
		"user_uuid":     c.UserUUID,
		"session_token": sessionToken,
	})
}

func (ro *Router) VerifyUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		VerificationCode string `json:"verification_code" valid:"printableascii"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	userUUID := cauth.GetCurrentUserUUID(r.Context())

	err := ro.auth.VerifyUser(r.Context(), userUUID, body.VerificationCode)
	if err != nil && err != cauth.ErrInvalidCredentials {
		ro.logger.Error("Failed to verify user", err)
		ro.rw.InternalErr(w)
		return
	} else if err == cauth.ErrInvalidCredentials {
		ro.rw.BadRequest(w, err)
		return
	}

	ro.rw.OK(w, nil)
}

func (ro *Router) ResendVerificationCode(w http.ResponseWriter, r *http.Request) {
	userUUID := cauth.GetCurrentUserUUID(r.Context())

	err := ro.auth.ResendVerificationCode(r.Context(), userUUID)
	if err != nil {
		ro.logger.Error("Failed to resend verification code", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, nil)
}

func (ro *Router) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email       string `json:"email" valid:"email"`
		OldPassword string `json:"old_password" valid:"printableascii"`
		NewPassword string `json:"new_password" valid:"printableascii"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	err := ro.auth.ChangePassword(r.Context(), body.Email, body.OldPassword, body.NewPassword)
	if err != nil {
		ro.logger.Error("Failed to change password", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, nil)
}

func (ro *Router) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email" valid:"email"`
	}

	if !ro.rw.Read(w, r, &body) {
		return
	}

	err := ro.auth.ResetPassword(r.Context(), body.Email)
	if err != nil {
		ro.logger.Error("Failed to reset password", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, nil)
}

func (ro *Router) HandleAddCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
		body     struct {
			Email    string `json:"email" valid:"required,email"`
			Password string `json:"password" valid:"runelength(4|32)"`
		}
	)

	if !ro.rw.Read(w, r, &body) {
		return
	}

	err := ro.auth.AddCredentials(r.Context(), userUUID, body.Email, body.Password)
	if err != nil {
		ro.logger.Error("Failed to add credentials", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, nil)
}

func (ro *Router) HandleChangeEmail(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
		body     struct {
			NewEmail string `json:"new_email" valid:"required,email"`
		}
	)

	err := ro.auth.ChangeEmail(ctx, userUUID, body.NewEmail)
	if err != nil {
		ro.logger.Error("Failed to change email", err)
		ro.rw.InternalErr(w)
		return
	}

	ro.rw.OK(w, nil)
}

func (ro *Router) HandleGetCredentials(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
	)

	c, err := ro.auth.GetCredentials(ctx, userUUID)
	if err != nil && !cerror.HasCause(err, gorm.ErrRecordNotFound) {
		ro.logger.Error("Failed to get credentials", err)
		ro.rw.InternalErr(w)
		return
	}

	if c == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ro.rw.OK(w, c)
}
