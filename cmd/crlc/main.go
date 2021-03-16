package main

import (
	"github.com/WineGecko/crlcollector/pkg/tsl"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

//Ссылка на TSL https://e-trust.gosuslugi.ru/app/scc/portal/api/v1/portal/ca/getxml

func main() {
	logger := log.New(os.Stdout, "crlc: ", log.Lshortfile)

	t, err := tsl.NewTSL("https://e-trust.gosuslugi.ru/app/scc/portal/api/v1/portal/ca/getxml", "tsl.xml", logger)

	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()

	r.GET("/crl/:keyId", ReverseProxy(t))
	r.GET("/debug", func(c *gin.Context) {
		c.JSON(http.StatusOK, t.GetCDPMap())
	})

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method not allowed")
	})

	r.Run(":8080")
}

func ReverseProxy(t *tsl.Tsl) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyId := c.Param("keyId")
		targetUrl := t.GetCDPMap()[keyId]
		u, _ := url.Parse(targetUrl)
		director := func(req *http.Request) {
			req = c.Request
			req.URL = u
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
