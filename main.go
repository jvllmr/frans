package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jvllmr/frans/package/config"
	apiRoutes "github.com/jvllmr/frans/package/routes/api"
	clientRoutes "github.com/jvllmr/frans/package/routes/client"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func setupRouter(configValue config.Config) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()

	r := gin.Default()

	defaultGroup := r.Group(configValue.RootPath)
	clientRoutes.SetupClientRoutes(r, defaultGroup, configValue)
	apiRoutes.SetupAPIRoutes(defaultGroup, configValue)
	return r
}

func main() {
	godotenv.Load()
	configValue, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	config.InitDB(configValue)
	defer config.DBClient.Close()
	config.InitOIDC(configValue)
	if configValue.DevMode {
		log.Println("Info: frans was started in development mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := setupRouter(configValue)
	serveString := fmt.Sprintf("%s:%d", configValue.Host, configValue.Port)
	log.Printf("Info: Serving on %s", serveString)
	r.Run(serveString)

}
