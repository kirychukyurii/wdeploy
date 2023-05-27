package api

import (
	"github.com/kirychukyurii/wdeploy/internal/lib/logger"
	"net/http"
	_ "net/http/pprof"
)

func SetDebug(logger logger.Logger) {
	//debug.SetGCPercent(-1)

	go func() {
		logger.Zap.Infof("Start debug server on http://localhost:8090/debug/pprof/")
		err := http.ListenAndServe(":8090", nil)
		if err != nil {
			logger.Zap.Error(err)
		}
	}()
}
