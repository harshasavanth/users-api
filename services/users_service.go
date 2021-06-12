package services

import (
	"github.com/harshasavanth/bookstore_users-api/logger"
	"github.com/harshasavanth/users-api/domain/users"
	"github.com/harshasavanth/users-api/utils/crypto_utils"
	"github.com/harshasavanth/users-api/utils/date_utils"
	"github.com/harshasavanth/users-api/utils/rest_errors"
	"net/http"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {
}
type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *rest_errors.RestErr)
	Login(users.User) (*users.User, *rest_errors.RestErr)
	GetUser(string) (*users.User, *rest_errors.RestErr)
	GetUserByEmail(string) (*users.User, *rest_errors.RestErr)
	UpdateUser(users.User) (*users.User, *rest_errors.RestErr)
	UpdateProfilePic(string, string) (*users.User, *rest_errors.RestErr)
	DeleteUser(string) *rest_errors.RestErr
	VerifyEmail(users.User) (*users.User, *rest_errors.RestErr)
}

func (s *usersService) CreateUser(user users.User) (*users.User, *rest_errors.RestErr) {
	exists := users.User{Email: user.Email}

	if exists.GetByEmail() == nil {
		return nil, rest_errors.NewInvalidInputError("user already exists")
	}
	if err := user.RegisterValidate(); err != nil {
		return nil, err
	}
	user.ProfileImage = ""
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)

	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}
func (s *usersService) Login(user users.User) (*users.User, *rest_errors.RestErr) {
	current := &users.User{Email: user.Email}
	if err := current.GetByEmail(); err != nil {
		if err.Status == http.StatusInternalServerError {
			return nil, err
		}
		return nil, rest_errors.NewInvalidInputError("enter valid email")
	}
	if err := current.LoginAuthentication(user.Password); err != nil {
		return nil, err
	}
	if current.EmailVerification {
		token, err := user.GenerateJWT()
		if err != nil {
			return nil, rest_errors.NewInternalServerError("error while generating token")
		}
		logger.Info("email verified")
		current.AccessToken = token
	}

	logger.Info(current.AccessToken)
	return current, nil
}

func (s *usersService) GetUser(userId string) (*users.User, *rest_errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) GetUserByEmail(email string) (*users.User, *rest_errors.RestErr) {
	result := &users.User{Email: email}
	if err := result.GetByEmail(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) UpdateUser(user users.User) (*users.User, *rest_errors.RestErr) {
	current := &users.User{Id: user.Id}
	if err := current.Get(); err != nil {
		return nil, err
	}

	current.FirstName = user.FirstName
	current.LastName = user.LastName
	current.OverEighteen = user.OverEighteen
	current.Email = user.Email
	current.Password = user.Password
	if err := current.RegisterValidate(); err != nil {
		return nil, err
	}
	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) UpdateProfilePic(userId string, path string) (*users.User, *rest_errors.RestErr) {
	current := &users.User{Id: userId}
	if err := current.Get(); err != nil {
		return nil, err
	}

	current.ProfileImage = path
	if err := current.UpdateProfilePic(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) DeleteUser(userId string) *rest_errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *usersService) VerifyEmail(user users.User) (*users.User, *rest_errors.RestErr) {
	current := &users.User{Id: user.Id}
	err := current.Get()
	if err != nil {
		err = rest_errors.NewInternalServerError("cannot find user")
		return nil, err
	}
	current.EmailVerification = true
	accessToken, err := current.GenerateJWT()
	if err != nil {
		return nil, err
	}
	current.AccessToken = accessToken
	if err := current.UpdateEmailVerified(); err != nil {
		return nil, err
	}
	return current, nil
}
