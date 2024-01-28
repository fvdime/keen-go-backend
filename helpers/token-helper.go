package helpers

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fvdime/keen-go-backend/database"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenProps struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

// initializing our db
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var JWT_KEY string = os.Getenv("JWT_KEY")

func GenerateTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &TokenProps{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &TokenProps{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, tokenErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWT_KEY))
	if tokenErr != nil {
		log.Panic(tokenErr)
		return "", "", tokenErr
	}

	refreshToken, refreshErr := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(JWT_KEY))
	if refreshErr != nil {
		log.Panic(refreshErr)
		return "", "", refreshErr
	}

	return token, refreshToken, nil
}


func UpdateTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObject primitive.D

	updateObject = append(updateObject, bson.E{"token", signedToken})
	updateObject = append(updateObject, bson.E{"refresh_token", signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObject = append(updateObject, bson.E{"updated_at", Updated_at})

	//4 mongodb update
	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx, filter, bson.D{{"$set", updateObject}}, &opt,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	// return
}


func ValidateToken(signedToken string) (claims *TokenProps, msg string){
	token, err := jwt.ParseWithClaims(
		signedToken,
		&TokenProps{},
		func(token *jwt.Token)(interface{}, error){
			return []byte(JWT_KEY), nil
		},
	)

	if err != nil {
		msg=err.Error()
		return
	}

	claims, ok:= token.Claims.(*TokenProps)
	if !ok{
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}