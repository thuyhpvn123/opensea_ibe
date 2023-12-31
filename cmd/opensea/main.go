package main

import (
	"fmt"
	c_config "gitlab.com/meta-node/meta-node/cmd/opensea/config"
	"gitlab.com/meta-node/meta-node/cmd/opensea/api/routers"
	"gitlab.com/meta-node/meta-node/cmd/opensea/database"
	"gitlab.com/meta-node/meta-node/pkg/logger"
)


func main() {
	// load config
	config, err := c_config.LoadConfig(c_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.Config)
	// Initialize the database connection
	database.InitDatabase()
	// Initialize the Gin router
	router:=routers.InitRouter()
	// Run the server
	if err := router.Run(cConfig.ServerAddress); err != nil {
		panic(err)
	}
}
