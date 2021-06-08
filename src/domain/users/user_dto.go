package users

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/harshasavanth/users-api/src/crypto_utils"
	"github.com/harshasavanth/users-api/src/rest_errors"

	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	signedKey   = "secret"
	senderEmail = "harshasavanth123@gmail.com"
	password    = "htcdesire"
)

type User struct {
	Id                 string    `json:"id"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	OverEighteen       bool      `json:"over_eighteen"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	ProfileImage       string    `json:"profile_image"`
	AccountUsedToLogin string    `json:"account_used_to_login"`
	EmailVerification  bool      `json:"verified"`
	Acknowledgement    bool      `json:"acknowledgement"`
	DateCreated        string    `json:"date_created"`
	PreviousPasswords  [3]string `json:"previous_passwords"`
	PreviousLogin      string    `json:"previous_login"`
	AccessToken        string    `json:"access_token"`
}

func (user *User) IsValidEmail() *rest_errors.RestErr {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if len(user.Email) < 3 || len(user.Email) > 254 {
		return rest_errors.NewInvalidInputError(fmt.Sprintf("%s length is too small", user.Email))
	}
	if emailRegex.MatchString(user.Email) {
		return nil
	}
	return rest_errors.NewInvalidInputError(fmt.Sprintf("%s is not a valid emai address", user.Email))
}

func (user *User) IsValidPassword() *rest_errors.RestErr {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	const minPassLength = 8
	const maxPassLength = 64
	var passLen int
	var errorString string

	for _, ch := range user.Password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("at least one numeric character required")
	}
	if !specialCharPresent {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		return rest_errors.NewInvalidInputError(fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength))
	}

	if len(errorString) != 0 {
		return rest_errors.NewInvalidInputError(errorString)
	}
	return nil
}

func (user *User) RegisterValidate() *rest_errors.RestErr {
	if err := user.IsValidEmail(); err != nil {
		return err
	}
	if err := user.IsValidPassword(); err != nil {
		return err
	}
	if user.Acknowledgement == false {
		return rest_errors.NewInvalidInputError("Please acknowledge terms")
	}
	if user.OverEighteen == false {
		return rest_errors.NewInvalidInputError("Must be over eighteen years")
	}
	return nil
}

func GenerateUuid() string {
	u := uuid.New()
	return u.String()
}

func (user *User) SendVerificationEmail() *rest_errors.RestErr {
	hostUrl := "smtp.gmail.com"
	hostPort := "587"
	emailSender := senderEmail
	password := password
	emailReceiever := user.Email

	emailAuth := smtp.PlainAuth(
		"",
		emailSender,
		password,
		hostUrl,
	)
	token, enEerr := crypto_utils.Encrypt(user.Id)
	if enEerr != nil {
		return enEerr
	}
	msg := []byte("To: " + emailReceiever + "\r\n" +
		"Subject: " + "Hello" + "\r\n" + "Please click below link to verify\nhttp://localhost:8080/users/verifyemail/" + token)

	err := smtp.SendMail(
		hostUrl+":"+hostPort,
		emailAuth,
		emailSender,
		[]string{emailReceiever},
		msg)

	if err != nil {
		return rest_errors.NewInvalidInputError("invalid email")
	}
	return nil

}

func (user *User) GenerateJWT() (string, *rest_errors.RestErr) {
	var signedKey = []byte(os.Getenv(signedKey))
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user.Id
	claims["expires"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(signedKey)
	if err != nil {
		return "", rest_errors.NewBadRequestError("something happened while creating token")

	}
	return tokenString, nil
}
