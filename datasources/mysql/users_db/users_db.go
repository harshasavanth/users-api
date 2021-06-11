package users_db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

var (
	Client   *sql.DB
	username = os.Getenv("mysql_users_username")
	password = os.Getenv("mysql_users_password")
	host     = os.Getenv("mysql_users_host")
	schema   = os.Getenv("mysql_users_schema")
)

func init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		username, password, host, schema)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	query := "CREATE TABLE IF NOT EXISTS `users_db`.`users` ( `id` VARCHAR(255) NOT NULL,  `first_name` VARCHAR(45) NULL,  `last_name` VARCHAR(45) NULL," +
		"  `over_eighteen` CHAR(1) NOT NULL,  `email` VARCHAR(100) NOT NULL, `password` VARCHAR(45) NOT NULL, `profile_image` VARCHAR(250) NULL , `account_used_to_login` VARCHAR(45) NULL," +
		"  `acknowledgement` CHAR(1) NOT NULL,  `email_verification` CHAR(1) NULL,  `previous_login` VARCHAR(45) NULL,  `previous_password1` VARCHAR(45) NULL," +
		" `previous_password2` VARCHAR(45) NULL,  `previous_password3` VARCHAR(45) NULL,  `date_created` VARCHAR(45) NULL, `access_token` VARCHAR(250) NULL, PRIMARY KEY (`id`),  UNIQUE INDEX `id_UNIQUE` (`id` ASC) VISIBLE);"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelfunc()
	_, err = Client.ExecContext(ctx, query)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("connection success")
}
