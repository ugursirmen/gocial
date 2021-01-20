package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func CreateUser(model CreateUserModel) error {
	ctx := context.Background()
	var err error

	tsql := `
      INSERT INTO Users (Username, Email, Password, CreatedAt) VALUES (@Username, @Email, @Password, @CreatedAt);
      select isNull(SCOPE_IDENTITY(), -1);
    `

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Username", model.Username),
		sql.Named("Email", model.Email),
		sql.Named("Password", model.Password),
		sql.Named("CreatedAt", time.Now()))

	var newID int
	err = row.Scan(&newID)
	if err != nil {
		return err
	}

	if newID == -1 {
		return errors.New("An error has occured while trying to create user")
	}

	return nil
}

func UpdateUserInfo(model UpdateUserInfoModel) error {
	ctx := context.Background()
	var err error

	tsql := fmt.Sprintf("UPDATE Users SET Username = @Username, FirstName = @FirstName, LastName = @LastName, Bio = @Bio WHERE Id = @Id")

	result, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("Username", model.Username),
		sql.Named("FirstName", model.FirstName),
		sql.Named("LastName", model.LastName),
		sql.Named("Bio", model.Bio),
		sql.Named("Id", model.UserId))

	effectedRows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if effectedRows == 0 {
		return errors.New("An error has occured while trying to update user info")
	}

	return nil
}
