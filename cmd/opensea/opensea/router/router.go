package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.com/meta-node/meta-node/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	c_config "gitlab.com/meta-node/meta-node/cmd/client/config"
)

// Person represents a log document in MongoDB
type LogModel struct {
	ID        primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string    `json:"name,omitempty" bson:"name,omitempty"`
	Description int `json:"description,omitempty" bson:"description,omitempty"`
	Image string `json:"image,omitempty" bson:"image,omitempty"`
	Attributes []AttributeModel `json:"attributes,omitempty" bson:"attributes,omitempty"`
	TokenId int `json:"tokenid,omitempty" bson:"tokenid,omitempty"`
}

type AttributeModel struct {
	TraitType string `json:"trait_type,omitempty" bson:"trait_type,omitempty"`
	Value int `json:"value,omitempty" bson:"value,omitempty"`
	MaxValue int `json:"max_value,omitempty" bson:"max_value,omitempty"`
}
var collection *mongo.Collection
const (
	CONFIG_FILE_PATH = "config/conf.json"
)

func InitDatabase(){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://minigame:chinhtoi@cluster0.qpzfe.mongodb.net/?retryWrites=true&w=majority")
	client, _ := mongo.Connect(ctx, clientOptions)
	collection = client.Database("opensea").Collection("marketplace")

}
func InitRouter()*mux.Router{
	config, err := c_config.LoadConfig(CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)

	router := mux.NewRouter().StrictSlash(true)
	server := Server{}
	server.Init(cConfig)
	router.HandleFunc("/nft", CreateNft).Methods("POST")
	router.HandleFunc("/upload", UploadHandler).Methods("POST")

	router.HandleFunc("/nft/{tokenid}", GetNft).Methods("GET")
	router.PathPrefix("/public").Handler(http.FileServer(http.Dir("frontend")))
	router.HandleFunc("/ws", server.WebsocketHandler)
	fmt.Println("server is running on port 3000")
	return router
}

func CreateNft(response http.ResponseWriter, request *http.Request,) {
	response.Header().Set("content-type", "application/json")
	var log *LogModel
	json.NewDecoder(request.Body).Decode(&log)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, &log)
	id := result.InsertedID
	log.ID = id.(primitive.ObjectID)
	json.NewEncoder(response).Encode(log)
}

func GetNft(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var log LogModel
	vars :=mux.Vars(request)
	tokenid,err:=strconv.Atoi(vars["tokenid"])
	fmt.Println("token:",vars["tokenid"])
	filter:=bson.D{{"tokenid",tokenid}}
	fmt.Println("filter:",filter)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor := collection.FindOne(ctx,filter)
	fmt.Println("cursor:",cursor)
	err=cursor.Decode(&log)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(log)
}
const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB


func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {

		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)

		return

	}
	// The argument to FormFile must match the name attribute
	// of the file input on the frontend
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./frontend/public/uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("./frontend/public/uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Upload successful")

}
