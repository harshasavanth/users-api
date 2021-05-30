package users

import (
	"github.com/harshasavanth/bookstore_users-api/logger"
	"github.com/harshasavanth/users-api/datasources/mysql/users_db"
	"github.com/harshasavanth/utils-go/rest_errors"
)

const (
	queryInsertUser = "INSERT INTO users (id, first_name, last_name, age, email, password, account_used_to_login, acknowledgement, email_verification," +
		"previous_login, previous_password1, previous_password2, previous_password3, date_created,access_token) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	queryGetUser = "SELECT id, first_name, last_name, age, email, password, account_used_to_login, acknowledgement, " +
		"email_verification, previous_login, previous_password1, previous_password2, previous_password3, date_created ,access_token FROM users WHERE id = ?;"
	queryUpdateUser = "UPDATE users SET id=?, first_name=?, last_name=?, age=?, email=?, password=?, account_used_to_login=?, acknowledgement=?, email_verification=?," +
		"previous_login=?, previous_password1=?, previous_password2=?, previous_password3=?, date_created=?, access_token = ? WHERE id=?;"
	queryGetUserByEmail = "SELECT id, first_name, last_name, age, email, password, account_used_to_login, acknowledgement, " +
		"email_verification, previous_login, previous_password1, previous_password2, previous_password3, date_created, access_token FROM users WHERE email = ?;"
	queryUpdateVerifiedEmail = "UPDATE users SET email_verification=? , access_token = ? WHERE email=?"
	queryDeleteUser          = "DELETE from users WHERE id=?;"
)

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	user.Id = GenerateUuid()
	go user.SendVerificationEmail()
	_, saveErr := stmt.Exec(user.Id, user.FirstName, user.LastName, user.Age, user.Email, user.Password, user.AccountUsedToLogin, user.Acknowledgement,
		user.EmailVerification, user.PreviousLogin, user.PreviousPasswords[0], user.PreviousPasswords[1], user.PreviousPasswords[2], user.DateCreated, user.AccessToken)
	if saveErr != nil {
		logger.Error("error when trying to save user ", saveErr)
		return rest_errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Age, &user.Email, &user.Password, &user.AccountUsedToLogin, &user.Acknowledgement,
		&user.EmailVerification, &user.PreviousLogin, &user.PreviousPasswords[0], &user.PreviousPasswords[1], &user.PreviousPasswords[2], &user.DateCreated, &user.AccessToken); getErr != nil {
		logger.Error("error when trying to get user by id ", getErr)
		return rest_errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.FirstName, user.LastName, user.Age, user.Email, user.Password, user.AccountUsedToLogin, user.Acknowledgement,
		user.EmailVerification, user.PreviousLogin, user.PreviousPasswords[0], user.PreviousPasswords[1], user.PreviousPasswords[2], user.DateCreated, user.AccessToken, user.Id)
	if err != nil {
		logger.Error("error when trying to update user ", err)
		return rest_errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement ", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) GetByEmail() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUserByEmail)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Age, &user.Email, &user.Password, &user.AccountUsedToLogin, &user.Acknowledgement,
		&user.EmailVerification, &user.PreviousLogin, &user.PreviousPasswords[0], &user.PreviousPasswords[1], &user.PreviousPasswords[2], &user.DateCreated, &user.AccessToken); getErr != nil {
		logger.Error("error when trying to get user by id ", getErr)
		return rest_errors.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) UpdateEmailVerified() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateVerifiedEmail)
	if err != nil {
		logger.Error("error when trying to prepare update user statement ", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.EmailVerification, user.AccessToken, user.Email)
	if err != nil {
		logger.Error("error when trying to update user ", err)
		return rest_errors.NewInternalServerError("database error")
	}
	return nil
}
