package email

import (
	"go.uber.org/fx"
)

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,

	NewRouter,
	NewSignupRoute,
	NewLoginRoute,
	NewVerifyUserRoute,
	NewResendVerificationCodeRoute,
	NewChangePasswordRoute,
	NewResetPasswordRoute,
	NewAddCredentialsRoute,
	NewChangeEmailRoute,
	NewGetCredentialsRoute,
)

var FxMigrations = fx.Invoke(
	RunMigrations,
)
