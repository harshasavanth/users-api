package app

import (
	"github.com/gin-gonic/gin"
	"os"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	//logger.Info("about to start the application...")
	router.Run(":8080" + os.Getenv("PORT"))
}
