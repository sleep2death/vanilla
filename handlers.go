package vanilla

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var jwtKey = []byte("vanilla_icecream")
var expire = time.Hour * 100

// login form binding
type login struct {
	Username string `form:"username" json:"username" bson:"username"  binding:"required"`
	Password string `form:"password" json:"password" bson:"password" binding:"required"`
}

// register form binding
type register struct {
	Username string `form:"username" json:"username" bson:"username"  binding:"required"`
	Email    string `form:"email" json:"email" bson:"email"  binding:"required"`
	Password string `form:"password" json:"password" bson:"password" binding:"required"`
}

func getLoginHandler(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// valid username and password pattern first
		if err := json.Validate(); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"reason": "illigal username or password",
			})
			return
		}

		opts := options.FindOne()
		filter := bson.M{"username": json.Username}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		res := db.Collection(UserCollection).FindOne(ctx, filter, opts)

		if err := res.Err(); err != nil {
			if err.Error() == "mongo: no documents in result" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"reason": "username not existed",
				})
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}
		r := &login{}
		if err := res.Decode(r); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		cErr := bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(json.Password))
		if cErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"reason": "password not match",
			})
			return
		}

		claims := &jwt.StandardClaims{
			Id:        r.Username,
			ExpiresAt: time.Now().Add(expire).Unix(),
		}

		// create jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(jwtKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"reason": "failed to generate token",
			})
			return
		}

		// log.Println("token:", tokenStr)
		// c.Header("Authorization", "Bearer "+tokenString)

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"token":  tokenStr,
		})
	}
}

func getRegisterHandler(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json register
		if err := c.ShouldBindJSON(&json); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// valid username and password pattern first
		if err := json.Validate(); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"reason": "illigal username or password",
			})
			return
		}

		// encrypt the password
		hash, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"reason": "password encryption error",
			})
			return
		}

		// log.Printf("User register: %s %s %s", json.Username, json.Email, hash)
		opts := options.FindOneAndUpdate().SetUpsert(true)
		filter := bson.M{"username": json.Username}
		update := bson.M{"$setOnInsert": bson.M{"username": json.Username, "email": json.Email, "password": string(hash)}}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		res := db.Collection(UserCollection).FindOneAndUpdate(ctx, filter, update, opts)

		if err := res.Err(); err != nil {
			if err.Error() == "mongo: no documents in result" {
				claims := &jwt.StandardClaims{
					Id:        json.Username,
					ExpiresAt: time.Now().Add(time.Second * 60).Unix(),
				}

				// create jwt token
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenStr, err := token.SignedString(jwtKey)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"reason": "failed to generate token",
					})
					return
				}

				// log.Println("token:", tokenStr)
				// c.Header("Authorization", "Bearer "+tokenString)

				c.JSON(http.StatusOK, gin.H{
					"status": "ok",
					"token":  tokenStr,
				})
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"reason": "username existed",
		})
	}
}

func getPingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result:": "pong"})
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if len(auth) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"reason": "authorization field empty",
			})
			return
		}

		authHeaderParts := strings.Fields(auth)
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"reason": "authorization format error",
			})
			return
		}

		tokenStr := authHeaderParts[1]
		// log.Println("token", tokenStr)
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			if err.Error() == "Token is expired" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"reason": "token is expired",
				})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"reason": "token is invalid",
				})
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println(claims["jti"], claims["exp"])
			c.Set("username", claims["jti"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"reason": "token is invalid",
			})
			return
		}
	}
}

// CORSMiddleware allow all
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
