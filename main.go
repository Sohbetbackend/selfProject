package main

import (
	"log"

	"github.com/Sohbetbackend/selfProject/config"
	"github.com/Sohbetbackend/selfProject/internal/api"
	"github.com/Sohbetbackend/selfProject/internal/store"
	"github.com/Sohbetbackend/selfProject/internal/store/pgx"
	"github.com/Sohbetbackend/selfProject/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	defer utils.InitLogs().Close()
	config.LoadConfig()
	defer store.Init().(*pgx.PgxStore).Close()

	routes := gin.Default()
	api.Routes(routes)
	if err := routes.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
