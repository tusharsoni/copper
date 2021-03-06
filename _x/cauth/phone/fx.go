package phone

import "go.uber.org/fx"

var Fx = fx.Provide(
	NewSQLRepo,
	NewSvc,

	NewRouter,
	NewLogin,
	NewSignup,
)

var FxMigrations = fx.Invoke(
	RunMigrations,
)

var FxValidators = fx.Invoke(
	AddPhoneNumberValidator,
)
