package main

import (
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"

	"github.com/zartbot/goflow/lib/metricbeat"

	"net/http"
	_ "net/http/pprof"
)

func main() {

	/*
		pprof hook
	*/
	go func() {
		logrus.Info(http.ListenAndServe("localhost:6666", nil))
	}()
	go metricbeat.StartMetricBeat()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		app := gin.Default()
		app.StaticFS("/static", http.Dir("config"))

		app.RunTLS(":8000", "./config/cert/webhook.cert", "./config/cert/webhook.key")
	}()
	wg.Wait()
}
