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

	stmt, err := db.Prepare("EXEC CreateUser @Username, @Email, @Password, @CreatedAt")
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

	result, err := db.ExecContext(
		ctx,
		"EXEC UpdateUserInfo @Username, @FirstName, @LastName, @Bio, @UserId",
		sql.Named("Username", model.Username),
		sql.Named("FirstName", model.FirstName),
		sql.Named("LastName", model.LastName),
		sql.Named("Bio", model.Bio),
		sql.Named("UserId", model.UserID))

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

	result, err := db.ExecContext(
		ctx,
		"EXEC UpdateUserPP @ProfilePicture, @UserId",
		sql.Named("ProfilePicture", model.ProfilePicture),
		sql.Named("UserId", model.UserID))

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

	result, err := db.ExecContext(
		ctx,
		"EXEC UpdateUserPassword @OldPassword, @NewPassword, @UserId",
		sql.Named("OldPassword", model.OldPassword),
		sql.Named("NewPassword", model.NewPassword),
		sql.Named("UserId", model.UserID))

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

	rows, err := db.QueryContext(ctx, "EXEC GetUserInfo @UserId", sql.Named("UserId", id))

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

	user.ProfilePicture = pp
	user.Followed = true

	return user, nil
}

func FollowUser(model FollowUserModel) error {
	ctx := context.Background()
	var err error

	stmt, err := db.Prepare("EXEC FollowUser @FollowerId, @UserId")
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

	result, err := db.ExecContext(
		ctx,
		"EXEC UnfollowUser @FollowerId, @UserId",
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

	rows, err := db.QueryContext(ctx, "EXEC GetUserFollows @FollowerId", sql.Named("FollowerId", id))

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
		user.CreatedAt = &createdAt
		user.Followed = true

		users = append(users, user)
	}

	return users, nil
}

func GetUserFollowers(id int) ([]UserDto, error) {

	ctx := context.Background()
	var err error

	var users []UserDto

	rows, err := db.QueryContext(ctx, "EXEC GetUserFollowers @UserId", sql.Named("UserId", id))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		var username, email, firstName, lastName, pp string
		var createdAt time.Time
		var followed bool

		err := rows.Scan(&id, &email, &username, &firstName, &lastName, &pp, &createdAt, &followed)
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
		user.CreatedAt = &createdAt
		user.Followed = followed

		users = append(users, user)
	}

	return users, nil
}

func CreatePost(model CreatePostModel) error {
	ctx := context.Background()
	var err error

	stmt, err := db.Prepare("EXEC CreatePost @UserId, @Description, @Image, @CreatedAt")
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

	rows, err := db.QueryContext(ctx, "EXEC GetUserPosts @UserId", sql.Named("UserId", id))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var post PostDto

		var deleted bool

		err := rows.Scan(&post.ID, &post.Description, &post.Image, &deleted, &post.CreatedAt, &post.Liked)
		if err != nil {
			return nil, err
		}

		post.Owner = nil

		posts = append(posts, post)
	}

	return posts, nil
}

func GetPostsArbitrary(model PostsArbitraryModel) ([]PostDto, error) {

	ctx := context.Background()
	var err error

	var posts []PostDto

	uniqueSlice := unique(model.PostIDs)

	postIds := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(uniqueSlice)), ","), "[]")

	rows, err := db.QueryContext(
		ctx, "EXEC GetPostsArbitrary @UserId, @Posts",
		sql.Named("UserId", model.UserID),
		sql.Named("Posts", postIds))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var username, email, firstName, lastName, pp, description, image string
		var createdAt time.Time
		var userID, postID int

		var deleted, followed, liked bool

		err := rows.Scan(&userID, &username, &email, &firstName, &lastName, &pp, &postID, &description, &image, &deleted, &createdAt, &followed, &liked)
		if err != nil {
			return nil, err
		}

		var post PostDto

		post.ID = postID
		post.CreatedAt = createdAt
		post.Description = description
		post.Image = image
		post.Liked = liked

		var owner UserDto

		owner.ID = userID
		owner.Username = username
		if firstName != "" && lastName != "" {
			owner.FullName = firstName + " " + lastName
		}
		owner.ProfilePicture = pp
		owner.Followed = followed

		post.Owner = &owner

		if !deleted {
			posts = append(posts, post)
		} else {
			posts = append(posts, nil...)
		}
	}

	return posts, nil
}

func GetPostsByIds(postIds []int) ([]PostEntity, error) {
	ctx := context.Background()
	var err error

	var posts []PostEntity

	uniqueSlice := unique(postIds)

	postIdsString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(uniqueSlice)), ","), "[]")

	rows, err := db.QueryContext(
		ctx, "SELECT Id,COALESCE(Description,''),COALESCE(Image,''),COALESCE(IsDeleted,0),CreatedAt,UserId FROM Posts WHERE Id IN ("+postIdsString+")")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var postEntity PostEntity

		err := rows.Scan(&postEntity.ID, &postEntity.Description, &postEntity.Image, &postEntity.Deleted, &postEntity.CreatedAt, &postEntity.UserID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, postEntity)
	}

	return posts, nil
}

func GetUserById(userId int) (UserEntity, error) {
	ctx := context.Background()
	var err error

	var user UserEntity

	rows, err := db.QueryContext(
		ctx, "SELECT Id,Username,FirstName,LastName,Bio,ProfilePicture FROM Users Where Id = @UserId",
		sql.Named("UserId", userId))

	if err != nil {
		return user, err
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Bio, &user.ProfilePicture)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func Newsfeed(userId int) ([]PostDto, error) {

	ctx := context.Background()
	var err error

	var posts []PostDto

	rows, err := db.QueryContext(ctx, "EXEC Newsfeed @UserId", sql.Named("UserId", userId))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var username, email, firstName, lastName, pp, description, image string
		var createdAt time.Time
		var userID, postID int

		var deleted, followed, liked bool

		err := rows.Scan(&userID, &username, &email, &firstName, &lastName, &pp, &postID, &description, &image, &deleted, &createdAt, &followed, &liked)
		if err != nil {
			return nil, err
		}

		var post PostDto

		post.ID = postID
		post.CreatedAt = createdAt
		post.Description = description
		post.Image = image
		post.Liked = liked

		var owner UserDto

		owner.ID = userID
		owner.Username = username
		if firstName != "" && lastName != "" {
			owner.FullName = firstName + " " + lastName
		}
		owner.ProfilePicture = pp
		owner.Followed = followed

		post.Owner = &owner

		if !deleted {
			posts = append(posts, post)
		} else {
			posts = append(posts, nil...)
		}
	}

	return posts, nil
}

func GetPostDetail(id int) (*PostDto, error) {

	ctx := context.Background()
	var err error

	var post PostDto

	rows, err := db.QueryContext(ctx, "EXEC GetPostDetail @PostId", sql.Named("PostId", id))

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

	result, err := db.ExecContext(
		ctx,
		"EXEC DeletePost @PostId",
		sql.Named("PostId", id))

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

	stmt, err := db.Prepare("EXEC LikePost @PostId, @UserId")
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

	result, err := db.ExecContext(
		ctx,
		"EXEC UnlikePost @PostId, @UserId",
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

func IsFollowed(userID int, followerID int) (bool, error) {

	ctx := context.Background()
	var err error
	var followed bool

	tsql := `
	SELECT CASE WHEN EXISTS
	(
		SELECT * FROM Followers WHERE FollowerId = @FollowerId AND UserId = @UserId
	)
	THEN CAST(1 AS BIT)
		ELSE CAST(0 AS BIT)
	END
    `

	rows, err := db.QueryContext(ctx, tsql, sql.Named("UserId", userID), sql.Named("FollowerId", followerID))

	if err != nil {
		return followed, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&followed)
		if err != nil {
			return followed, err
		}
	}

	return followed, nil
}

func IsLiked(userID int, postID int) (bool, error) {

	ctx := context.Background()
	var err error
	var liked bool

	tsql := `
	SELECT CASE
	WHEN EXISTS
		(
			SELECT * FROM Likes WHERE UserId = @UserId AND PostId = @PostId
		)
		THEN CAST(1 AS BIT)
	ELSE CAST(0 AS BIT)
	END
    `

	rows, err := db.QueryContext(ctx, tsql, sql.Named("UserId", userID), sql.Named("PostId", postID))

	if err != nil {
		return liked, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&liked)
		if err != nil {
			return liked, err
		}
	}

	return liked, nil
}
