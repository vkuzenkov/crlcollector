package main

import (
	"github.com/WineGecko/crlcollector/pkg/tsl"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

//Ссылка на TSL https://e-trust.gosuslugi.ru/app/scc/portal/api/v1/portal/ca/getxml

func main()  {
	logger := log.New(os.Stdout, "crlc: ", log.Lshortfile)

	t, err := tsl.NewTSL("https://e-trust.gosuslugi.ru/app/scc/portal/api/v1/portal/ca/getxml", "tsl.xml", logger)

	if err != nil {
		log.Fatal(err)
	}

	r := gin.New()

	r.GET("/debug", func(c *gin.Context) {
		c.JSON(http.StatusOK, t)
	})

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method not allowed")
	})

	r.Run(":8080")
}