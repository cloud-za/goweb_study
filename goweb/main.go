package main

import (
	"net/http"
	"sync"

	"github.com/zartbot/goweb/middleware/jwt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zartbot/goflow/lib/metricbeat"
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
		app.StaticFS("/static", http.Dir("static"))

		//token
		j := jwt.NewJWT("test_security_key", 60*60*24*7)

		app.GET("login", func(c *gin.Context) {
			claims := jwt.NewTokenClaims("kevin", "aa", "bb", 2)
			token, _ := j.GenerateToken(claims)
			c.SetCookie("token", token, j.ExpireInSecond, "/", "", true, true)
			c.String(http.StatusOK, "login successfl")
		})

		app.GET("/welcome", j.VerifyToken(), func(c *gin.Context) {
			name1 := c.DefaultQuery("f", "Guest")
			name2 := c.Query("l")
			a, _ := c.Get("claims")

			c.String(http.StatusOK, "hello f:%s l:%s---%+v", name1, name2, a)
		})

		app.RunTLS(":8000", "./config/cert/webhook.cert", "./config/cert/webhook.key")
	}()
	wg.Wait()

}
