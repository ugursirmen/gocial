package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

	rows, err := db.QueryContext(ctx, "SELECT Username,Email,COALESCE(FirstName,''),COALESCE(LastName,''),COALESCE(ProfilePicture,''),CreatedAt FROM Users where Id = @Id", sql.Named("Id", id))

	if err != nil {
		return user, err
	}

	defer rows.Close()

	var username, email, firstName, lastName, pp string
	var createdAt time.Time

	for rows.Next() {
		err := rows.Scan(&username, &email, &firstName, &lastName, &pp, &createdAt)
		if err != nil {
			return user, err
		}
	}

	user.ID = id
	user.Username = username
	if firstName != "" && lastName != "" {
		user.FullName = firstName + " " + lastName
	}
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

func GetUserFollows(id int) ([]UserDto, error) {

	ctx := context.Background()
	var err error

	var users []UserDto

	rows, err := db.QueryContext(ctx, "select Id,Username,Email,COALESCE(FirstName,''),COALESCE(LastName,''),COALESCE(ProfilePicture,''),CreatedAt from Followers LEFT JOIN Users U on U.Id = Followers.UserId WHERE FollowerId = @FollowerId ORDER BY CreatedAt DESC", sql.Named("FollowerId", id))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		var username, email, firstName, lastName, pp string
		var createdAt time.Time

		err := rows.Scan(&id, &email, &username, &firstName, &lastName, &pp, &createdAt)
		if err != nil {
			return nil, err
		}

		var user UserDto
		user.ID = id
		user.Username = username
		if firstName != "" && lastName != "" {
			user.FullName = firstName + " " + lastName
		}
		user.ProfilePicture = pp
		user.CreatedAt = createdAt

		users = append(users, user)
	}

	return users, nil
}

func GetUserFollowers(id int) ([]UserDto, error) {

	ctx := context.Background()
	var err error

	var users []UserDto

	rows, err := db.QueryContext(ctx, "select Id,Username,Email,COALESCE(FirstName,''),COALESCE(LastName,''),COALESCE(ProfilePicture,''),CreatedAt from Followers LEFT JOIN Users U on U.Id = Followers.FollowerId WHERE UserId = @UserId ORDER BY CreatedAt DESC", sql.Named("UserId", id))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		var username, email, firstName, lastName, pp string
		var createdAt time.Time

		err := rows.Scan(&id, &email, &username, &firstName, &lastName, &pp, &createdAt)
		if err != nil {
			return nil, err
		}

		var user UserDto
		user.ID = id
		user.Username = username
		if firstName != "" && lastName != "" {
			user.FullName = firstName + " " + lastName
		}
		user.ProfilePicture = pp
		user.CreatedAt = createdAt

		users = append(users, user)
	}

	return users, nil
}

func CreatePost(model CreatePostModel) error {
	ctx := context.Background()
	var err error

	tsql := `
      INSERT INTO Posts (UserId,Description,Image,CreatedAt) VALUES (@UserId, @Description, @Image, @CreatedAt);
      select isNull(SCOPE_IDENTITY(), -1);
    `

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("UserId", model.UserID),
		sql.Named("Description", model.Description),
		sql.Named("Image", model.Image),
		sql.Named("CreatedAt", time.Now()))

	var newID int
	err = row.Scan(&newID)
	if err != nil {
		return err
	}

	if newID == -1 {
		return errors.New("An error has occured while trying to create post")
	}

	return nil
}

func GetUserPosts(id int) ([]PostDto, error) {

	ctx := context.Background()
	var err error

	var posts []PostDto

	rows, err := db.QueryContext(ctx, "SELECT Id,COALESCE(Description,''),COALESCE(Image,''),COALESCE(IsDeleted, 0),CreatedAt FROM Posts WHERE UserId = @UserId ORDER BY CreatedAt DESC", sql.Named("UserId", id))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var post PostDto

		var deleted bool

		err := rows.Scan(&post.ID, &post.Description, &post.Image, &deleted, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.Owner = nil

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPosts(userId int, postIds []int) ([]PostDto, error) {

	ctx := context.Background()
	var err error

	var posts []PostDto

	var query = "select U.Id,Username,Email,COALESCE(FirstName,''),COALESCE(LastName,''),COALESCE(ProfilePicture,''),P.Id ,COALESCE(Description,''),COALESCE(Image,''),COALESCE(IsDeleted, 0),P.CreatedAt from Followers LEFT JOIN Users U on U.Id = Followers.UserId LEFT JOIN Posts P on U.Id = P.UserId WHERE P.Id is not null and FollowerId = @FollowerId ORDER BY CreatedAt DESC"

	if len(postIds) > 0 {

		postIdsString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(postIds)), ","), "[]")

		query = "select U.Id,Username,Email,COALESCE(FirstName,''),COALESCE(LastName,''),COALESCE(ProfilePicture,''),P.Id ,COALESCE(Description,''),COALESCE(Image,''),COALESCE(IsDeleted, 0),P.CreatedAt from Followers LEFT JOIN Users U on U.Id = Followers.UserId LEFT JOIN Posts P on U.Id = P.UserId WHERE P.Id is not null and P.Id IN (" + postIdsString + ") AND FollowerId = @FollowerId"
	}

	rows, err := db.QueryContext(ctx, query, sql.Named("FollowerId", userId))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var username, email, firstName, lastName, pp, description, image string
		var createdAt time.Time
		var userID, postID int

		var deleted bool

		err := rows.Scan(&userID, &username, &email, &firstName, &lastName, &pp, &postID, &description, &image, &deleted, &createdAt)
		if err != nil {
			return nil, err
		}

		var post PostDto

		post.ID = postID
		post.CreatedAt = createdAt
		post.Description = description
		post.Image = image
		post.Liked = false

		var owner UserDto

		owner.ID = userID
		owner.Username = username
		if firstName != "" && lastName != "" {
			owner.FullName = firstName + " " + lastName
		}
		owner.ProfilePicture = pp

		post.Owner = &owner

		if !deleted {
			posts = append(posts, post)
		} else {
			posts = append(posts, PostDto{})
		}
	}

	return posts, nil
}

func GetPostDetail(id int) (*PostDto, error) {

	ctx := context.Background()
	var err error

	var post PostDto

	rows, err := db.QueryContext(ctx, "select Posts.Id,COALESCE(Description,''),COALESCE(Image,''),COALESCE(IsDeleted,0),Posts.CreatedAt,UserId,Username,Email,COALESCE(FirstName,''),COALESCE(LastName,''),COALESCE(ProfilePicture,'') from Posts LEFT JOIN Users U on U.Id = Posts.UserId where Posts.Id = @PostId", sql.Named("PostId", id))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var username, email, firstName, lastName, pp string
	var userID int

	var deleted bool

	for rows.Next() {
		err := rows.Scan(&post.ID, &post.Description, &post.Image, &deleted, &post.CreatedAt, &userID, &username, &email, &firstName, &lastName, &pp)
		if err != nil {
			return nil, err
		}
	}

	if deleted {
		return nil, nil
	}

	var owner UserDto

	owner.ID = userID
	owner.Username = username
	if firstName != "" && lastName != "" {
		owner.FullName = firstName + " " + lastName
	}
	owner.ProfilePicture = pp

	post.Owner = &owner

	return &post, nil
}

func DeletePost(id int) error {
	ctx := context.Background()
	var err error

	tsql := fmt.Sprintf("UPDATE Posts SET IsDeleted = @IsDeleted WHERE Id = @Id")

	result, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("IsDeleted", 1),
		sql.Named("Id", id))

	effectedRows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if effectedRows == 0 {
		return errors.New("An error has occured while trying to delete the post")
	}

	return nil
}

func LikePost(model LikePostModel) error {
	ctx := context.Background()
	var err error

	tsql := `
      INSERT INTO Likes (PostId, UserId) VALUES (@PostId, @UserId)
    `

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("PostId", model.PostID),
		sql.Named("UserId", model.UserID))

	err = row.Err()
	if err != nil {
		return err
	}

	return nil
}

func UnlikePost(model LikePostModel) error {
	ctx := context.Background()
	var err error

	tsql := fmt.Sprintf("DELETE FROM Likes WHERE PostId = @PostId AND UserId = @UserId")

	result, err := db.ExecContext(
		ctx,
		tsql,
		sql.Named("PostId", model.PostID),
		sql.Named("UserId", model.UserID))

	effectedRows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if effectedRows == 0 {
		return errors.New("An error has occured while trying to unlike the post")
	}

	return nil
}
