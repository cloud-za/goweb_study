package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/sirupsen/logrus"

	"github.com/googollee/go-socket.io"
	"github.com/zartbot/goflow/lib/metricbeat"
	"github.com/zartbot/goweb/example/graphql"

	"net/http"
	_ "net/http/pprof"
)

func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("before middleware")
		c.Set("request", "clinet_request")
		c.Next()
		fmt.Println("after middleware")
	}
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("session_id"); err == nil {
			value := cookie.Value
			fmt.Println(value)
			if value == "123" {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("CORS middleware loaded...")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, WS, WSS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

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

		//basic usage
		//ref: https://www.jianshu.com/p/a31e4ee25305/

		app.GET("/", func(context *gin.Context) {
			val, _ := context.Cookie("name")
			context.SetCookie("name", "hahhaha", 60*60*24, "/", "", true, true)
			context.String(200, "Cookie:%s", val)
		})

		app.GET("/user/:name", func(c *gin.Context) {
			name := c.Param("name")
			c.String(http.StatusOK, "Hello %s", name)
		})

		app.GET("/use/:name/*action", func(c *gin.Context) {
			name := c.Param("name")
			action := c.Param("action") //match all of the reset as a string
			c.String(http.StatusOK, "Hello %s is %s", name, action)
		})

		app.GET("/welecome", func(c *gin.Context) {
			name1 := c.DefaultQuery("f", "Guest")
			name2 := c.Query("l")
			c.String(http.StatusOK, "hello f:%s l:%s", name1, name2)

		})

		app.POST("/form_post", func(c *gin.Context) {
			message := c.PostForm("message")
			nick := c.DefaultPostForm("nick", "anonymous")

			c.JSON(http.StatusOK, gin.H{
				"status": gin.H{
					"status_code": http.StatusOK,
					"status":      "ok",
				},
				"message":  message,
				"nickname": nick,
			})
		})

		//combined post with parameters
		app.PUT("/post", func(c *gin.Context) {
			id := c.Query("id")
			page := c.DefaultQuery("page", "0")
			name := c.PostForm("name")
			message := c.PostForm("message")
			fmt.Printf("id: %s; page: %s; name: %s; message: %s \n", id, page, name, message)
			c.JSON(http.StatusOK, gin.H{
				"status_code": http.StatusOK,
			})
		})

		//single file upload
		app.POST("/upload", func(c *gin.Context) {
			name := c.PostForm("name")
			fmt.Println(name)
			file, header, err := c.Request.FormFile("filename")
			if err != nil {
				c.String(http.StatusBadRequest, "Bad request")
				return
			}
			filename := header.Filename

			fmt.Println(file, err, filename)

			out, err := os.Create("upload/" + filename)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()
			_, err = io.Copy(out, file)
			if err != nil {
				log.Fatal(err)
			}
			c.String(http.StatusCreated, "upload successful")
		})

		//multifile upload
		app.POST("/multi/upload", func(c *gin.Context) {
			err := c.Request.ParseMultipartForm(200000)
			if err != nil {
				log.Fatal(err)
			}

			formdata := c.Request.MultipartForm

			files := formdata.File["upload"]
			for i, _ := range files {
				file, err := files[i].Open()
				defer file.Close()
				if err != nil {
					log.Fatal(err)
				}

				out, err := os.Create("upload/" + files[i].Filename)

				defer out.Close()

				if err != nil {
					log.Fatal(err)
				}

				_, err = io.Copy(out, file)

				if err != nil {
					log.Fatal(err)
				}

				c.String(http.StatusCreated, "upload successful")

			}

		})

		type User struct {
			Username string `form:"username" json:"username" binding:"required"`
			Passwd   string `form:"passwd" json:"passwd" bdinding:"required"`
			Age      int    `form:"age" json:"age"`
		}

		app.POST("/login", func(c *gin.Context) {
			var user User
			var err error
			contentType := c.Request.Header.Get("Content-Type")

			switch contentType {
			case "application/json":
				err = c.BindJSON(&user)
			case "application/x-www-form-urlencoded":
				err = c.BindWith(&user, binding.Form)
			}

			if err != nil {
				fmt.Println(err)
				log.Fatal(err)
			}

			c.JSON(http.StatusOK, gin.H{
				"user":   user.Username,
				"passwd": user.Passwd,
				"age":    user.Age,
			})

		})
		app.GET("/redirect/163", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "https://www.163.com")
		})

		//render on different type
		app.GET("/render", func(c *gin.Context) {
			contentType := c.DefaultQuery("content_type", "json")
			if contentType == "json" {
				c.JSON(http.StatusOK, gin.H{
					"user":   "rsj217",
					"passwd": "123",
				})
			} else if contentType == "xml" {
				c.XML(http.StatusOK, gin.H{
					"user":   "rsj217",
					"passwd": "123",
				})
			}

		})

		v1 := app.Group("/v1")

		v1.GET("/login", func(c *gin.Context) {
			c.String(http.StatusOK, "v1 login")
		})

		v2 := app.Group("/v2")

		v2.GET("/login", func(c *gin.Context) {
			c.String(http.StatusOK, "v2 login")
		})

		//middle ware

		//app decorator
		app.Use(MiddleWare())
		{
			app.GET("/middleware", func(c *gin.Context) {
				request := c.MustGet("request").(string)
				req, _ := c.Get("request")
				logrus.Warn("Middleware")
				c.JSON(http.StatusOK, gin.H{
					"middile_request": request,
					"request":         req,
				})
			})
		}

		//middleware in route
		app.GET("/before", MiddleWare(), func(c *gin.Context) {
			request := c.MustGet("request").(string)
			c.JSON(http.StatusOK, gin.H{
				"middile_request": request,
			})
		})

		app.GET("/auth/signin", func(c *gin.Context) {
			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    "123",
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(c.Writer, cookie)
			c.String(http.StatusOK, "Login successful")
		})

		app.GET("/home", AuthMiddleWare(), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": "home"})
		})

		soIO, err := socketio.NewServer(nil)
		if err != nil {
			panic(err)
		}
		soIO.On("connection", func(so socketio.Socket) {
			fmt.Println("on connection")
			so.Join("chat")
			so.On("chat message", func(msg string) {
				fmt.Println("emit:", so.Emit("chat message", msg))
				so.BroadcastTo("chat", "chat message", msg)
			})

			so.On("disconnection", func() {
				fmt.Println("on disconnect")
			})
		})
		soIO.On("error", func(so socketio.Socket, err error) {
			fmt.Printf("[ WebSocket ] Error : %v", err.Error())
		})

		app.GET("/socket.io/", gin.WrapH(soIO))

		app.GET("/graphiql", gin.WrapH(handler.Playground("GraphQL playground", "/query")))
		gql := gin.WrapH(handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}})))

		app.GET("/query", gql)
		app.POST("/query", gql)

		app.RunTLS(":8000", "./config/cert/webhook.cert", "./config/cert/webhook.key")
	}()
	wg.Wait()
}
