package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/renanmedina/xgh-life/configs"
	"github.com/renanmedina/xgh-life/handlers"
	"github.com/renanmedina/xgh-life/integrations"
)

var (
	//go:embed frontend/static/*
	staticEmbeded embed.FS
	//go:embed frontend/templates/*
	templatesEmbeded embed.FS
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	appConfigs := configs.NewApplicationConfigs()

	router := gin.Default()
	router.Use(cors.Default())

	configureNewRelic(router, appConfigs)
	configureStatic(router)
	configureHandlers(router)

	err := bootServer(router)

	if err != nil {
		panic("[XGH-LIFE:WEBSERVER-LOG] failed to start server due to: " + err.Error())
	}
}

func configureNewRelic(router *gin.Engine, configs *configs.ApplicationConfigs) {
	if configs.NewRelicEnabled {
		newRelicApp, err := integrations.NewRelicApp()
		if err != nil {
			panic(fmt.Sprintf("[XGH-LIFE:NEWRELIC-INTEGRATION] %s", err))
		}

		router.Use(nrgin.Middleware(newRelicApp))
	}
}

func detectEnvPort() string {
	port := os.Getenv("PORT")

	if port != "" {
		return port
	}

	return os.Getenv("XGH_HTTP_PORT")
}

func configureStatic(router *gin.Engine) {
	assetsFS, _ := fs.Sub(staticEmbeded, "frontend")
	router.StaticFS("/assets", http.FS(assetsFS))

	templ := template.Must(template.New("").ParseFS(templatesEmbeded, "frontend/templates/*.tmpl"))
	router.SetHTMLTemplate(templ)
}

func configureHandlers(router *gin.Engine) {
	router.GET("/infra/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"is-alive": true})
	})

	router.GET("/", func(c *gin.Context) {
		c.AddParam("id", "roulette")
		handlers.AxiomDetailsHandlerHtml(staticEmbeded, c)
	})

	router.StaticFileFS("/robots.txt", "frontend/static/robots.txt", http.FS(staticEmbeded))

	router.GET("/:id", func(ctx *gin.Context) {
		handlers.AxiomDetailsHandlerHtml(staticEmbeded, ctx)
	})

	router.GET("/axioms/:id", func(ctx *gin.Context) {
		handlers.AxiomDetailsHandlerHtml(staticEmbeded, ctx)
	})

	api := router.Group("/api")
	{
		api.GET("/axioms", handlers.AxiomsListHandler)
		api.GET("/axioms/:id", handlers.AxiomDetailsHandlerJson)
	}

	slack := router.Group("/slack")
	{
		slack.POST("/axioms", handlers.SlackBotHandler)
	}

	github := router.Group("/github/pull-requests")
	{
		github.POST("/auto-approve", handlers.GithubAutoApprovePullRequestHandler)
	}
}

func bootServer(router *gin.Engine) error {
	port := detectEnvPort()
	logger := integrations.NewApplicationLogger()

	if port == "" {
		logger.Println("[XGH-LIFE:WEBSERVER-LOG] No provided port by environment, using default 8080")
		port = "8080"
	}

	host := fmt.Sprintf(":%s", port)
	logger.Printf("[XGH-LIFE:WEBSERVER-LOG] Listening and serving HTTP on %s", host)

	return router.Run(host)
}
