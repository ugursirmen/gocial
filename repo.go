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
		sql.Named("Id", model.UserID))

	effectedRows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if effectedRows == 0 {
		return errors.New("An error has occured while trying to update user info")
	}

	return nil
}

func UpdateUserPP(model UpdateUserPPModel) error {
	ctx := context.Background()
	var err error

	tsql := fmt.Sprintf("UPDATE Users SET ProfilePicture = @ProfilePicture WHERE Id = @Id")

	result, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("ProfilePicture", model.ProfilePicture),
		sql.Named("Id", model.UserID))

	effectedRows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if effectedRows == 0 {
		return errors.New("An error has occured while trying to update user pp")
	}

	return nil
}

func UpdateUserPassword(model UpdateUserPasswordModel) error {
	ctx := context.Background()
	var err error

	rows, err := db.QueryContext(ctx, "SELECT Password FROM Users where Id = @Id", sql.Named("Id", model.UserID))

	if err != nil {
		return err
	}

	defer rows.Close()

	var password string

	for rows.Next() {

		err := rows.Scan(&password)
		if err != nil {
			return err
		}
	}

	if password != model.OldPassword {
		return errors.New("Old-password mismatch")
	}

	tsql := fmt.Sprintf("UPDATE Users SET Password = @Password WHERE Id = @Id")

	result, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("Password", model.NewPassword),
		sql.Named("Id", model.UserID))

	effectedRows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if effectedRows == 0 {
		return errors.New("An error has occured while trying to update user password")
	}

	return nil
}

func GetUserInfo(id int) (UserDto, error) {

	ctx := context.Background()
	var err error
	var user UserDto

	rows, err := db.QueryContext(ctx, "SELECT Username,FirstName,LastName,ProfilePicture,CreatedAt FROM Users where Id = @Id", sql.Named("Id", id))

	if err != nil {
		return user, err
	}

	defer rows.Close()

	var username, firstName, lastName, pp string
	var createdAt time.Time

	for rows.Next() {
		err := rows.Scan(&username, &firstName, &lastName, &pp, &createdAt)
		if err != nil {
			return user, err
		}
	}

	user.Username = username
	user.FullName = firstName + " " + lastName
	user.CreatedAt = createdAt
	user.ProfilePicture = pp
	user.Followed = true

	return user, nil
}

func FollowUser(model FollowUserModel) error {
	ctx := context.Background()
	var err error

	tsql := `
      INSERT INTO Followers (UserId, FollowerId) VALUES (@UserId, @FollowerId)
    `

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("UserId", model.UserID),
		sql.Named("FollowerId", model.FollowerID))

	err = row.Err()
	if err != nil {
		return err
	}

	return nil
}

func UnfollowUser(model FollowUserModel) error {
	ctx := context.Background()
	var err error

	tsql := fmt.Sprintf("DELETE FROM Followers WHERE UserId = @UserId AND FollowerId = @FollowerId")

	result, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("UserId", model.UserID),
		sql.Named("FollowerId", model.FollowerID))

	effectedRows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if effectedRows == 0 {
		return errors.New("An error has occured while trying to unfollow user")
	}

	return nil
}
