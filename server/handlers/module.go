package handlers

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	addUserModule,
	allUsersModule,
)
