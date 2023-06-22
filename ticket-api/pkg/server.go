package pkg

import (
	"ticket-api/config"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// HTTPServer Router -> Gin Router
type HTTPServer struct {
	*gin.Engine
}

// NewHTTPServer NewRouter : all the routes are defined here
func NewHTTPServer(env config.Config, logger Logger) HTTPServer {

	if env.Environment != "local" && env.SentryDSN != "" {
		if err := sentry.Init(
			sentry.ClientOptions{
				Dsn:         env.SentryDSN,
				Environment: `ticket-api` + env.Environment,
			},
		); err != nil {
			logger.Infof("Sentry initialization failed: %v\n", err)
		}
	}

	appEnv := env.Environment
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	httpRouter := gin.Default()

	httpRouter.Use(gzip.Gzip(gzip.DefaultCompression))
	httpRouter.MaxMultipartMemory = env.MaxMultipartMemory

	httpRouter.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
				AllowHeaders:     []string{"*"},
				AllowCredentials: true,
			},
		),
	)

	// Attach sentry middleware
	httpRouter.Use(
		sentrygin.New(
			sentrygin.Options{
				Repanic: true,
			},
		),
	)

	httpRouter.GET(
		"/health-check", func(c *gin.Context) {
			// utils.SendSentryMsg(c, "Error")
			c.JSON(http.StatusOK, gin.H{"data": "ticket API Up and Running"})
		},
	)

	return HTTPServer{httpRouter}
}

func (sh *HTTPServer) Start(env config.Config, logger Logger) {
	if env.ServerPort == "" {
		if err := sh.Run(); err != nil {
			logger.Fatal(err)
			return
		}
	}

	if err := sh.Run(":" + env.ServerPort); err != nil {
		logger.Fatal(err)
		return
	}
}
