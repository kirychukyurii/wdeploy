package lib

import (
	"github.com/kirychukyurii/wdeploy/internal/lib/ansible"
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"go.uber.org/fx"
)

// Module exports dependency
var Module = fx.Options(
	fx.Provide(logger.NewLogger),
	fx.Provide(ansible.NewExecutor),
)
