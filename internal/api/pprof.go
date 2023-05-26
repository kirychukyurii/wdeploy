package api

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func SetDebug() {
	//debug.SetGCPercent(-1)

	go func() {
		fmt.Println("Start debug server on http://localhost:8090/debug/pprof/")
		err := http.ListenAndServe(":8090", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
}
