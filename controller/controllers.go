// controller/controllers.go
package controller

import (
	"PROTOBUF/config"
	"PROTOBUF/gen_proto" // Import the correct package
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

func CreateUserHandle(c *gin.Context) {

	//taking data from user
	var user gen_proto.UserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : cant take user input")
		return
	}

	//hashing the user given password
	hashpash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error occures when hashing the pass")
		return
	}
	user.Password = string(hashpash)

	//generating object id for ID field
	objectID := primitive.NewObjectID()

	//inserting that id into id field
	user.Id = objectID.Hex()

	//Encoding the json formated data into protobuf
	protoData, err := proto.Marshal(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error while marshaling the data")
		return
	}

	//creating a time out context

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Entering data into mongoDB
	coll := config.Getcollection()
	_, err = coll.InsertOne(ctx, bson.M{"protobuf_data": protoData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error occures while inserting data into database")
		return
	}
	c.JSON(http.StatusOK, "successfully inserted data")

}

func GetUser(c *gin.Context) {

	//this variable store the data
	var ss struct {
		Protobufata []byte `bson:"protobuf_data"` // Change the field name based on how it's stored in MongoDB
	}

	var UnmarshalData gen_proto.UserInfo

	//taking the id from param
	id := c.Query("id")

	//coverting the user given id to object id type

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//find that data followed by that param
	coll := config.Getcollection()
	err = coll.FindOne(ctx, bson.M{"_id": objID}).Decode(&ss)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, "cant find the data")
			return
		} else {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

	}

	err = proto.Unmarshal(ss.Protobufata, &UnmarshalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "cant unmarshal the data")
		return
	}
	c.JSON(http.StatusOK, &UnmarshalData)

}
