package main

import (
	"log"
	"net/http"

	"gitlab.com/meta-node/meta-node/cmd/opensea/opensea/router"

)


func main() {
	router.InitDatabase()
	router:=router.InitRouter()
	log.Fatal(http.ListenAndServe(":3000", router))
	
}
