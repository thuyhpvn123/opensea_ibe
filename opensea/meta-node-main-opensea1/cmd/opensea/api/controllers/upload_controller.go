package controllers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"gitlab.com/meta-node/meta-node/cmd/opensea/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UploadController struct {
	Collection *mongo.Collection
}

func (c *UploadController)UploadHandler(ctx *gin.Context) {
	
		// single file
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		log.Println(file.Filename)
		dst:="./frontend/public/uploads/"
		// Upload the file to specific dst.
		filename := filepath.Base(file.Filename)
		if err := ctx.SaveUploadedFile(file, dst+filename); err != nil {
			ctx.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}
		imageUrl = dst+filename
		fileName=filename
		// ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
		response := models.Response{
			Code: 200,
			Data: fmt.Sprintf("'%s' uploaded!", file.Filename),
		}
		ctx.JSON(200, response)
		return
}

// type user struct {
// 	// Id    int    `uri:"id"`
// 	// Name  string `form:"name"`
// 	// Email string `form:"email"`
// 	File *multipart.FileHeader `form:"image" binding:"required"`
//    }
// func (c *UploadController) UploadHandler(ctx *gin.Context) {
// 	 var userObj user
// 	 if err := ctx.ShouldBind(&userObj); err != nil {
// 		ctx.String(http.StatusBadRequest, "bad request")
// 	  return
// 	 }
	 
// 	//  if err := ctx.ShouldBindUri(&userObj); err != nil {
// 	// 	ctx.String(http.StatusBadRequest, "bad request")
// 	//   return
// 	//  }
// 	 dst:="./frontend/public/uploads/"
// 	 err := ctx.SaveUploadedFile(userObj.File,  dst+userObj.File.Filename)
// 	 if err != nil {
// 		ctx.String(http.StatusInternalServerError, "unknown error")
// 	  return
// 	 }
   
// 	 ctx.JSON(http.StatusOK, gin.H{
// 	  "status": "ok",
// 	  "data":   userObj,
// 	 })
   
   
// }