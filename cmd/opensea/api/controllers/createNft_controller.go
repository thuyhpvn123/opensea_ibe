package controllers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	// "mime/multipart"
	"net/http"

	// "strconv"
	"gitlab.com/meta-node/meta-node/cmd/opensea/utils"

	"github.com/gin-gonic/gin"
	"gitlab.com/meta-node/meta-node/cmd/opensea/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateNftController struct {
	Collection *mongo.Collection
}
var tokenid int
var fileName string
var imageUrl string
// func (c *CreateNftController) CreateNft(ctx *gin.Context) {
// 	name := ctx.PostForm("name")
// 	fmt.Println("name:",name)
// 	description := ctx.PostForm("description")
// 	fmt.Println("description:",description)


// 	var attributeRequest struct {
// 		TraitType string `json:"trait_type,omitempty" bson:"trait_type,omitempty"`
// 		Value int `json:"value,omitempty" bson:"value,omitempty"`
// 		MaxValue int `json:"max_value,omitempty" bson:"max_value,omitempty"`
// 	}

// 	attributes1 :=attributeRequest
// 	attributes1.TraitType =ctx.PostForm("trait_type")
// 	attributes1.Value,_ = strconv.Atoi(ctx.PostForm("value"))
// 	attributes1.MaxValue,_ = strconv.Atoi(ctx.PostForm("max_value"))
// 	imagePath:= "localhost:3000/public/uploads/"+ fileName 
// 	// create nft
// 	tokenid++
// 	attributeArr :=[]models.AttributeModel{}
// 	attributeArr = append(attributeArr,attributes1)

// 	nftCreated := models.LogModel{
// 		Name:         name,
// 		Description:  description,
// 		Image: imagePath,
// 		Attributes:          attributeArr,
// 		TokenId:      tokenid,
// 	}
// 	result, err := c.Collection.InsertOne(ctx, &nftCreated)
// 	id := result.InsertedID
// 	nftCreated.ID = id.(primitive.ObjectID)
// 	if err != nil {
// 		// Handle error saving the deposit order
// 		response := models.Response{
// 			Code: http.StatusInternalServerError,
// 			Data: gin.H{"error": "Failed to create deposit order"},
// 		}
// 		ctx.JSON(http.StatusInternalServerError, response)
// 		return
// 	}

// 	response := models.Response{
// 		Code: 200,
// 		Data: nftCreated,
// 	}
// 	ctx.JSON(200, response)
// 	return
// }
func rename(oldname string ,newname string, path string ) {
	dst:= path
	 // Rename and Remove a file
	 // Using Rename() function
	 Original_Path := dst+oldname
	 New_Path := dst+newname
	 e := os.Rename(Original_Path, New_Path)
	 if e != nil {
		 log.Fatal(e)
	 }
	   
 }
func (c *CreateNftController) CreateNft(ctx *gin.Context) {
	var request struct {
		Name     string    `json:"name,omitempty" bson:"name,omitempty"`
		Description string `json:"description,omitempty" bson:"description,omitempty"`
		// File  *multipart.FileHeader `form:"information" binding:"required" json:"file,omitempty" bson:"file,omitempty"`
		Attributes []models.AttributeModel `json:"attributes,omitempty" bson:"attributes,omitempty"`	
	}
	// Parse the request body into the request struct
	if err := ctx.ShouldBind(&request); err != nil {
		// Handle error parsing request
		response := models.Response{
			Code: http.StatusBadRequest,
			Data: gin.H{"error": "Invalid request"},
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	// tokenid++
	uid := utils.GenerateRandomMD5String(16)
	ext := filepath.Ext(fileName)
	fmt.Println("ext:", ext)
	nfileName:=uid+ext
	dst:="./frontend/public/uploads/"
	path:="localhost:2000/public/uploads/"
	imagePath:= path+ nfileName 
	tokenid ,err:=strconv.Atoi(uid)
	// generate new id with size as md5 string
	rename(fileName,nfileName,dst)
	// create nft
	nftCreated := models.LogModel{
		Name:         request.Name,
		Description:  request.Description,
		Image: imagePath,
		Attributes:          request.Attributes,
		TokenId:     tokenid ,
	}
	result, err := c.Collection.InsertOne(ctx, &nftCreated)
	id := result.InsertedID
	nftCreated.ID = id.(primitive.ObjectID)
	if err != nil {
		// Handle error saving the deposit order
		response := models.Response{
			Code: http.StatusInternalServerError,
			Data: gin.H{"error": "Failed to create deposit order"},
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := models.Response{
		Code: 200,
		Data: nftCreated,
	}
	ctx.JSON(200, response)
	return
}

func (c *CreateNftController) GetNft(ctx *gin.Context) {
	var log models.LogModel
	// Get tokenid from the request parameters
	tokenid,err := strconv.Atoi(ctx.Param("tokenid"))
	if err != nil {
		// Handle the case when wallet is not found in the database
		// Return nonce as 0
		response := models.Response{
			Code: http.StatusInternalServerError,
			Data: gin.H{
				"error":"`" + err.Error() + "`",
			},
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("token:",tokenid)

	filter:=bson.D{{"tokenid",tokenid}}
	fmt.Println("filter:",filter)
	cursor := c.Collection.FindOne(ctx,filter)
	fmt.Println("cursor:",cursor)
	err=cursor.Decode(&log)
	if err != nil {
		// Handle the case when wallet is not found in the database
		// Return nonce as 0
		response := models.Response{
			Code: 200,
			Data: gin.H{
				"nft-quantity": "0",
			},
		}
		ctx.JSON(200, response)
		return
	}
	// Wallet found in the database, return the Nft info
	
	response := models.Response{
		Code: 200,
		Data: gin.H{
			"name": log.Name,
			"description":log.Description,
			"image":log.Image,
			"attributes":log.Attributes,
			"tokenid":log.TokenId,
		},
	}
	ctx.JSON(200, response)
}

func (c *CreateNftController) GetAllNft(ctx *gin.Context) {
	
	// Get tokenid from the request parameters
	// cursor := c.Collection.FindOne(ctx,filter)
	var logs []models.LogModel

	cursor, err := c.Collection.Find(ctx, bson.M{})
	if err != nil {
		// Handle error saving the deposit order
		response := models.Response{
			Code: http.StatusInternalServerError,
			Data: gin.H{"error":"`" + err.Error() + "`"},
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	fmt.Println("cursor:",cursor)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var log models.LogModel
		// var onepool Pool
		cursor.Decode(&log)
		logs = append(logs, log)
	}
	if err := cursor.Err(); err != nil {
		// Handle the case when wallet is not found in the database
		// Return nonce as 0
		response := models.Response{
			Code: 200,
			Data: gin.H{"error":"`" + err.Error() + "`"},
		}
		ctx.JSON(200, response)
		return
	}
	// Wallet found in the database, return the Nft info
	
	response := models.Response{
		Code: 200,
		Data: gin.H{
			"data":logs ,		
		},
	}
	ctx.JSON(200, response)
	
}



