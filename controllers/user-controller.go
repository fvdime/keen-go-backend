package controllers

import (
	"github.com/fvdime/keen-go-backend/database"
	"github.com/fvdime/keen-go-backend/helpers"
	"github.com/fvdime/keen-go-backend/models"

	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	checkPassword := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect!")
		checkPassword = false
	}

	return checkPassword, msg
}

func SignUp() gin.HandlerFunc {
	// return func(ctx *gin.Context) {
	// 	var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// 	var user models.User

	// 	if err := ctx.BindJSON(&user); err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	validationError := validate.Struct(user)
	// 	if validationError != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
	// 		return
	// 	}

	// 	count, err := userCollection.CountDocuments(c, bson.M{"email": user.Email})
	// 	if err != nil {
	//   log.Panic(err)
	//   ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
	//   return
	// 	}
	// 	cancel()

	// 	password := HashPassword(*user.Password)
	// 	user.Password = &password

	// 	count, err = userCollection.CountDocuments(c, bson.M{"phone": user.Phone})
	// 	defer cancel()
	// 	if err != nil {
	// 		log.Panic(err)
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the phone number"})
	// 	}

	// 	if count > 0 {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
	// 	}

	// 	user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// 	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// 	user.ID = primitive.NewObjectID()
	// 	user.User_Id = user.ID.Hex()
	// 	token, refreshToken, _ := helpers.GenerateTokens(*user.Email, *user.First_Name, *user.Last_Name, *user.User_Type, *&user.User_Id)
	// 	user.Token = &token
	// 	user.Refresh_Token = &refreshToken

	// 	// insert to the db
	// 	resultInsertNumber, insertError := userCollection.InsertOne(c, user)
	// 	if insertError != nil {
	// 			msg := fmt.Sprintf("user not created!")
	// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
	// 			return
	// 	}

	// 	// Use resultInsertNumber only if there is no error
	// 	defer cancel()
	// 	ctx.JSON(http.StatusOK, resultInsertNumber)

	// }
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			return
		}

		phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the phone number"})
			return
		}

		if emailCount > 0 || phoneCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email or phone number already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At = time.Now()
		user.Updated_At = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_Id = user.ID.Hex()

		token, refreshToken, err := helpers.GenerateTokens(*user.Email, *user.First_Name, *user.Last_Name, *user.User_Type, *&user.User_Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating tokens"})
			return
		}

		user.Token = &token
		user.Refresh_Token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, resultInsertionNumber)
	}

}

func SignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		var existingUser models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&existingUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect!"})
			return
		}

		isValidPassword, msg := VerifyPassword(*user.Password, *existingUser.Password)
		defer cancel()
		if isValidPassword != true {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if existingUser.Email == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found."})
		}

		token, refreshToken, _ := helpers.GenerateTokens(*existingUser.First_Name, *existingUser.Last_Name, *existingUser.Email, *existingUser.User_Type, *&existingUser.User_Id)
		helpers.UpdateTokens(token, refreshToken, existingUser.User_Id)

		err = userCollection.FindOne(c, bson.M{"user_id": existingUser.User_Id}).Decode(&existingUser)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, existingUser)
	}
}

// func GetUsers() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
// 			return
// 		}

// 		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		itemPerPage, err := strconv.Atoi(ctx.Query("itemPerPage"))
// 		if err != nil || itemPerPage <1{
// 			itemPerPage = 10
// 		}
// 		page, errP := strconv.Atoi(ctx.Query("page"))
// 		if errP !=nil || page<1{
// 			page = 1
// 		}

// 		startIndex := (page - 1) * itemPerPage
// 		startIndex, err = strconv.Atoi(ctx.Query("startIndex"))

// 		// I don't have any idea bout it
// 		matchStage := bson.D{{"$match", bson.D{{}}}}
// 		groupStage := bson.D{{"$group", bson.D{
// 			{"_id", bson.D{{"_id", "null"}}},
// 			{"total_count", bson.D{{"$sum", 1}}},
// 			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
// 		projectStage := bson.D{
// 			{"$project", bson.D{
// 				{"_id", 0},
// 				{"total_count", 1},
// 				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, itemPerPage}}}},}}}

// 			result,err := userCollection.Aggregate(ctx, mongo.Pipeline{
// 			matchStage, groupStage, projectStage})
// 			defer cancel()
// 			if err!=nil{
// 				c.JSON(http.StatusInternalServerError, gin.H{"error":"error occurred while listing user items"})
// 			}
// 			var allusers []bson.M
// 			if err = result.All(ctx, &allusers); err!=nil{
// 				log.Fatal(err)
// 			}
// 			c.JSON(http.StatusOK, allusers[0])}}

// 	}
// }

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		if err := helpers.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		// Uses the FindOne method on the userCollection to find a document with the specified user ID. It decodes the result into the user variable.
		err := userCollection.FindOne(c, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}
