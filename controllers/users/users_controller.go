package users

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/harshasavanth/bookstore_users-api/logger"
	"github.com/harshasavanth/users-api/domain/users"
	"github.com/harshasavanth/users-api/services"
	"github.com/harshasavanth/users-api/utils/aws"
	"github.com/harshasavanth/users-api/utils/crypto_utils"
	"github.com/harshasavanth/users-api/utils/rest_errors"
	"os"
	"time"

	"net/http"
)

var (
	UsersController usersControllerInterface = &usersController{}
)

type usersControllerInterface interface {
	Create(*gin.Context)
	Login(ctx *gin.Context)
	Get(c *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetByEmail(*gin.Context)
	VerifyEmail(*gin.Context)
	IsAuthorized(endpoint func(ctx *gin.Context)) gin.HandlerFunc
	ProfilePicUpload(ctx *gin.Context)
}

type usersController struct {
}

func (c *usersController) Create(ctx *gin.Context) {
	var user users.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid JSON body")
		ctx.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		ctx.JSON(saveErr.Status, saveErr)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (c *usersController) Login(ctx *gin.Context) {
	var user users.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid JSON body")
		ctx.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.UsersService.Login(user)
	if saveErr != nil {
		ctx.JSON(saveErr.Status, saveErr)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (c *usersController) Get(ctx *gin.Context) {
	userId := ctx.Param("user_id")

	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		ctx.JSON(getErr.Status, getErr)
		return
	}
	c.Display(ctx, userId, user)
}

func (c *usersController) ProfilePicUpload(ctx *gin.Context) {
	fileHeader, _ := ctx.FormFile("path")
	file, err := fileHeader.Open()
	if err != nil {
		restErr := rest_errors.NewBadRequestError("Unable to open pic")
		ctx.JSON(restErr.Status, restErr)
		return
	}
	userid := ctx.GetHeader("ID")
	logger.Info(userid)
	path, restErr := aws.Upload(file, fileHeader, userid)
	if err != nil {
		ctx.JSON(restErr.Status, restErr)
		return
	}
	user, getErr := services.UsersService.UpdateProfilePic(userid, path)
	if getErr != nil {
		ctx.JSON(getErr.Status, getErr)
		return
	}
	ctx.JSON(http.StatusOK, user)

}

func (c *usersController) GetByEmail(ctx *gin.Context) {
	userId := ctx.Param("email")
	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		ctx.JSON(getErr.Status, getErr)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *usersController) Update(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	var user users.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid JSON body")
		ctx.JSON(restErr.Status, restErr)
		return
	}
	user.Id = userId
	result, err := services.UsersService.UpdateUser(user)
	if err != nil {
		ctx.JSON(err.Status, err)
	}
	c.Display(ctx, userId, result)
}

func (c *usersController) Delete(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	if ctx.GetHeader("ID") == userId {
		if err := services.UsersService.DeleteUser(userId); err != nil {
			ctx.JSON(err.Status, err)
		}
		ctx.JSON(http.StatusOK, map[string]string{"status": "deleted"})
	} else {
		ctx.JSON(http.StatusNotImplemented, "you cannot delete other users")
	}
}

func (c *usersController) VerifyEmail(ctx *gin.Context) {
	id := ctx.Param("id")
	var user users.User
	did, err := crypto_utils.Decrypt(id)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	user.Id = did
	logger.Info(user.Id)
	result, verErr := services.UsersService.VerifyEmail(user)
	if verErr != nil {
		ctx.JSON(verErr.Status, verErr)
		return
	}
	ctx.JSON(http.StatusCreated, result)
}

func (con *usersController) IsAuthorized(endpoint func(*gin.Context)) gin.HandlerFunc {
	var signingKey = []byte(os.Getenv("signedKey"))
	return func(c *gin.Context) {
		if c.GetHeader("Token") != "" {
			token, err := jwt.Parse(c.GetHeader("Token"), func(token *jwt.Token) (interface{}, error) {
				if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
					return nil, fmt.Errorf("there was an error")
				}
				return signingKey, nil
			})
			if err != nil {
				err := rest_errors.NewNotAuthorizedError("not authorized")
				c.JSON(err.Status, err)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				var tm time.Time
				switch iat := claims["expires"].(type) {
				case float64:
					tm = time.Unix(int64(iat), 0)
				case json.Number:
					v, _ := iat.Int64()
					tm = time.Unix(v, 0)
				}
				logger.Info(fmt.Sprintf("%s", tm))
				logger.Info(fmt.Sprintf("%s", time.Now()))
				logger.Info(fmt.Sprintf("%s", tm.Before(time.Now())))

				if tm.Before(time.Now()) {
					err := rest_errors.NewNotAuthorizedError("token expired")
					c.JSON(err.Status, err)
					return
				} else {
					id := fmt.Sprintf("%s", claims["user"])
					c.Request.Header.Set("ID", id)
					endpoint(&*c)
				}
			} else {
				err := rest_errors.NewNotAuthorizedError("invalid token")
				c.JSON(err.Status, err)
				return
			}
		} else {
			err := rest_errors.NewNotAuthorizedError("not authorized")
			c.JSON(err.Status, err)
			return
		}
	}
}

func (c *usersController) Display(ctx *gin.Context, id string, user *users.User) {
	if ctx.GetHeader("ID") == id {
		ctx.JSON(http.StatusOK, user)
	} else {
		ctx.JSON(http.StatusNotImplemented, "not authorized")
	}

}
