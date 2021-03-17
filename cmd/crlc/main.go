package main

import (
	"flag"
	"fmt"
	"github.com/WineGecko/crlcollector/pkg/tsl"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

type config struct {
	TslUrl         *string
	Filename       *string
	UpdateInterval *time.Duration
	Port           *string
}

func main() {
	logger := log.New(os.Stdout, "crlc: ", log.Lshortfile)
	c := &config{
		TslUrl:         flag.String("tsllink", "https://e-trust.gosuslugi.ru/app/scc/portal/api/v1/portal/ca/getxml", "TSL url"),
		Filename:       flag.String("filename", "tsl.xml", "TSL filename"),
		UpdateInterval: flag.Duration("update", 12*time.Hour, "TSL file update interval"),
		Port:           flag.String("listen", ":8080", "Address:port for API"),
	}
	flag.Parse()

	t, err := tsl.NewTSL(*c.TslUrl, *c.Filename, logger)

	go func() {
		err := t.Update(*c.UpdateInterval)
		if err != nil {
			log.Printf("Unable update TSL. %s", err)
		}
	}()

	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()

	r.GET("/crl/:keyId", ReverseProxy(t, logger))
	r.GET("/debug", func(c *gin.Context) {
		c.JSON(http.StatusOK, t.GetCDPMap())
	})

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method not allowed")
	})

	s := &http.Server{
		Addr:         *c.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	_ = s.ListenAndServe()
}

func ReverseProxy(t *tsl.Tsl, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyId := c.Param("keyId")
		targetUrl := t.GetCDPMap()[strings.ToLower(keyId)]
		u, _ := url.Parse(targetUrl[0])
		logger.Printf("Redirecting to %s for keyId %s", u, keyId)
		director := func(req *http.Request) {
			req.URL = u
			req.Host = u.Host
		}
		proxy := &httputil.ReverseProxy{Director: director}
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.crl", keyId))
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
