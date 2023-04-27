package app

import (
	"github.com/kirychukyurii/wdeploy/internal/app/ansible"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(ansible.NewExecutor),
)
