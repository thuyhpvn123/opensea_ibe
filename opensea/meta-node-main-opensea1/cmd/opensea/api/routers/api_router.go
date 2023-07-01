package routers

import (
	// "fmt"
	"fmt"
	"net/http"
	// "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gitlab.com/meta-node/meta-node/cmd/opensea/api/controllers"
	c_config "gitlab.com/meta-node/meta-node/cmd/opensea/config"
	// "gitlab.com/meta-node/meta-node/cmd/opensea/core"
	"gitlab.com/meta-node/meta-node/cmd/opensea/database"
	"gitlab.com/meta-node/meta-node/pkg/logger"
)

// SetupRouter sets up the API routes and returns the Gin router.
func InitRouter() *gin.Engine {
	server := controllers.Server{}
	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.Config)


	serverapp := server.Init(cConfig)
	client:=serverapp.ConnectionHandler()
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	uploadController := controllers.UploadController{Collection: database.GetDB()}
	createNftController := controllers.CreateNftController{Collection: database.GetDB()}
	
	v1 := router.Group("/api/v1")
	 //http://localhost:3000/api/v1/test/template/
	{
		v1.StaticFS("", http.Dir("frontend/public"))  
		// CreateNFTRoutes(v1)
	}
	router.POST("/upload",uploadController.UploadHandler)
	router.POST("/createNft",createNftController.CreateNft)
	router.GET("/getNft/:tokenid",createNftController.GetNft)
	router.GET("/getAllNft",createNftController.GetAllNft)
	router.POST("/connectWallet",client.Caller.ConnectWallet)
	router.POST("/call",client.Caller.TryCall)
	// fmt.Println("server is running on port 2000")

	return router
}
