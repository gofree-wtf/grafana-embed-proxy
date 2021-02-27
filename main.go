package main

import (
	"flag"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	var grafanaUrl string
	var grafanaAPIKey string

	flag.StringVar(&grafanaUrl, "grafana-url", "", "Grafana URL")
	flag.StringVar(&grafanaAPIKey, "grafana-api-key", "", "Grafana API Key")
	flag.Parse()

	log.Info().Str("grafanaUrl", grafanaUrl).Msg("")
	log.Debug().Str("grafanaAPIKey", grafanaAPIKey).Msg("")

	grafanaTargetUrl, err := url.Parse(grafanaUrl)
	if err != nil {
		log.Panic().Err(err).Msg("failed to parse grafana url")
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = grafanaTargetUrl.Scheme
			req.URL.Host = grafanaTargetUrl.Host
			req.Host = grafanaTargetUrl.Host // for ingress
			req.Header.Set("Authorization", "Bearer "+grafanaAPIKey)

			if req.Method == http.MethodPost {
				// change POST to GET
				req.Method = http.MethodGet

				// remove body
				req.Header.Del("Content-Type")
				req.Header.Del("Content-Length")
				req.ContentLength = 0
				req.Body = nil
			}
		},
	}

	r := gin.New()
	r.Use(logger.SetLogger(logger.Config{Logger: &log.Logger}))
	r.Use(gin.Recovery())

	r.POST("/*proxyPath",
		setCookie("localhost"), // TODO configuration
		func(c *gin.Context) {
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	r.GET("/*proxyPath",
		checkCookie,
		func(c *gin.Context) {
			proxy.ServeHTTP(c.Writer, c.Request)
		})

	err = r.Run()
	if err != nil {
		log.Error().Err(err).Msg("")
	}
}

type ProxyFormData struct {
	ApiKey string `form:"Proxy-API-Key"`
}

const cookieKey = "Proxy-API-Key"

func setCookie(proxyHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		formData := &ProxyFormData{}
		if err := c.ShouldBind(formData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid form data",
			})
			return
		}

		if formData.ApiKey != "testtest" { // TODO: implement auth
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		cookieValue, err := c.Cookie(cookieKey)
		if err != nil || cookieValue != formData.ApiKey {
			c.SetCookie(cookieKey, formData.ApiKey, 3600, "/", proxyHost, false, true)
		}

		c.Next()
	}
}

func checkCookie(c *gin.Context) {
	cookieValue, err := c.Cookie(cookieKey)
	if err != nil || cookieValue != "testtest" { // TODO: implement auth
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	c.Next()
}
