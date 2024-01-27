package middleware

import(
	"fmt"
	"net/http"
	"github.com/fvdime/keen-go-backend/helpers"
	"github.com/gin-gonic/gin"
)


func Authenticate() gin.HandlerFunc{
	return func(c *gin.Context){
		userToken := c.Request.Header.Get("token")
		if userToken == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error":fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(userToken)
		if err !="" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid",claims.Uid)
		c.Set("user_type", claims.User_type)
		c.Next()
	}
}