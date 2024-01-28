package controllers

import (
	"github.com/fvdime/keen-go-backend/database"
	"github.com/fvdime/keen-go-backend/models"

	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var postCollection *mongo.Collection = database.OpenCollection(database.Client, "post")

func CreatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var post models.Post

		if err := ctx.BindJSON(&post); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		post.Created_At = time.Now()
		post.Updated_At = time.Now()

		insertResult, insertErr := postCollection.InsertOne(c, post)
		if insertErr != nil {
			msg := fmt.Sprintf("Post was not created: %v", insertErr)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "message": msg})
			return
		}

		defer cancel()

		ctx.JSON(http.StatusOK, gin.H{"success": true, "data": insertResult.InsertedID, "message": "post creating is successful"})
	}
}

func DeletePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		postId := ctx.Param("post_id")

		objectId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "message": "Invalid post ID format"})
		}

		deleteResult, err := postCollection.DeleteOne(c, bson.M{"_id": objectId})
		if err != nil {
			msg := fmt.Sprintf("error occurred while deleting the post: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "message": msg})
			return
		}

		defer cancel()

		if deleteResult.DeletedCount == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"success": false, "data": nil, "message": "Post not found"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"success": true, "data": nil, "message": "Post deleted successfully"})
	}
}

func UpdatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		postId := ctx.Param("post_id")

		var updatedPost models.Post
		if err := ctx.BindJSON(&updatedPost); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		objectId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		//defining the filter to identify the post by it is id
		filter := bson.M{"_id": objectId}

		update := bson.M{"$set": bson.M{
			"title": updatedPost.Title,
			"body":  updatedPost.Body,
			"image": updatedPost.Image,
		}}

		updateResult, err := postCollection.UpdateOne(c, filter, update)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "message": fmt.Sprintf("Error updating post: %v", err)})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"success": true, "data": updateResult.ModifiedCount, "message": "Post updated successfully"})
	}
}

func GetPost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		postId := ctx.Param("post_id")
		objectId, err := primitive.ObjectIDFromHex(postId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "data": nil, "message": "Invalid post ID format"})
		}

		var post models.Post

		err = postCollection.FindOne(c, bson.M{"_id": objectId}).Decode(&post)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "data": nil, "message": "Post not found"})
				return
			}
			msg := fmt.Sprintf("Error fetching post: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "message": msg})
			return
		}

		defer cancel()

		ctx.JSON(http.StatusOK, gin.H{"success": true, "data": post, "message": "Post fetched successfully"})
	}
}

func GetPosts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		cursor, err := postCollection.Find(c, bson.D{})
		if err != nil {
			msg := fmt.Sprintf("Error fetching posts: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "message": msg})
			return
		}

		defer cursor.Close(c)

		var posts []models.Post
		if err := cursor.All(c, &posts); err != nil {
			msg := fmt.Sprintf("Error fetching posts: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "message": msg})
			return
		}

		defer cancel()

		ctx.JSON(http.StatusOK, gin.H{"success": true, "data": posts, "message": "Posts fetched successfully"})
	}
}
