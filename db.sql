create database GocialDB collate SQL_Latin1_General_CP1_CI_AS
go

use GocialDB
go

create table Users
(
	Id int identity
		constraint Users_pk
			primary key nonclustered,
	Username nvarchar(100) not null,
	Email nvarchar(300) not null,
	Password nvarchar(500) not null,
	FirstName nvarchar(500),
	LastName nvarchar(500),
	Bio nvarchar(500),
	ProfilePicture nvarchar(max),
	CreatedAt datetime not null
)
go

create table Followers
(
	UserId int not null
		constraint Followers_Users_Id_fk
			references Users,
	FollowerId int not null
)
go

create unique index Followers_UserId_FollowerId_uindex
	on Followers (UserId, FollowerId)
go

create table Posts
(
	Id int identity
		constraint Posts_pk
			primary key nonclustered,
	Description nvarchar(max),
	Image nvarchar(max),
	IsDeleted bit,
	CreatedAt datetime not null,
	UserId int not null
		constraint Posts_Users_Id_fk
			references Users
)
go

create table Likes
(
	PostId int not null
		constraint Likes_Posts_Id_fk
			references Posts,
	UserId int not null
)
go

create unique index Likes_PostId_UserId_uindex
	on Likes (PostId, UserId)
go

create unique index Posts_Id_uindex
	on Posts (Id)
go

create unique index Users_Id_uindex
	on Users (Id)
go

create unique index Users_Username_uindex
	on Users (Username)
go

create unique index Users_Email_uindex
	on Users (Email)
go

CREATE PROCEDURE CreatePost
    @UserId int,
    @Description nvarchar(max),
    @Image nvarchar(max),
    @CreatedAt datetime
AS
INSERT INTO Posts (UserId, Description, Image, CreatedAt)
VALUES (@UserId, @Description, @Image, @CreatedAt);
select isNull(SCOPE_IDENTITY(), -1);
go

CREATE PROCEDURE CreateUser
    @Username nvarchar(100),
    @Email nvarchar(300),
    @Password nvarchar(500),
    @CreatedAt datetime
AS

INSERT INTO Users (Username, Email, Password, CreatedAt) VALUES (@Username, @Email, @Password, @CreatedAt)
SELECT isNull(SCOPE_IDENTITY(), -1)
go

CREATE PROCEDURE DeletePost @PostId int
AS
UPDATE Posts
SET IsDeleted = 1
WHERE Id = @PostId
go

CREATE PROCEDURE FollowUser
    @FollowerId int,
    @UserId int
AS

INSERT INTO Followers (UserId, FollowerId) VALUES (@UserId, @FollowerId)
go

CREATE PROCEDURE GetPostDetail @PostId int
AS
SELECT Posts.Id,
       COALESCE(Description, ''),
       COALESCE(Image, ''),
       COALESCE(IsDeleted, 0),
       Posts.CreatedAt,
       UserId,
       Username,
       Email,
       COALESCE(FirstName, ''),
       COALESCE(LastName, ''),
       COALESCE(ProfilePicture, '')
FROM Posts
         LEFT JOIN Users U ON U.Id = Posts.UserId
WHERE Posts.Id = @PostId
go

CREATE PROCEDURE GetPostsArbitrary

@UserId nvarchar(max),
@Posts nvarchar(max)

AS
DECLARE @SQL nvarchar(4000);

SET @SQL = 'SELECT U.Id,
       Username,
       Email,
       COALESCE(FirstName, ''''),
       COALESCE(LastName, ''''),
       COALESCE(ProfilePicture, ''''),
       P.Id,
       COALESCE(Description, ''''),
       COALESCE(Image, ''''),
       COALESCE(IsDeleted, 0),
       P.CreatedAt,
(
           SELECT CASE
                      WHEN EXISTS
                          (
                              SELECT * FROM Followers WHERE FollowerId = '+@UserId+' AND UserId = U.Id
                          )
                          THEN CAST(1 AS BIT)
                      ELSE CAST(0 AS BIT)
                      END) AS Followed,
       (
           SELECT CASE
                      WHEN EXISTS
                          (
                              SELECT * FROM Likes WHERE PostId = P.Id AND UserId = '+@UserId+'
                          )
                          THEN CAST(1 AS BIT)
                      ELSE CAST(0 AS BIT)
                      END) AS Liked
FROM Followers
         LEFT JOIN Users U on U.Id = Followers.UserId
         LEFT JOIN Posts P on U.Id = P.UserId
WHERE P.Id is not null
  AND P.Id IN ('+@Posts+')';

EXECUTE (@SQL);
go

CREATE PROCEDURE GetUserFollowers @UserId int
AS

SELECT Id,
       Username,
       Email,
       COALESCE(FirstName, ''),
       COALESCE(LastName, ''),
       COALESCE(ProfilePicture, ''),
       CreatedAt,
       (
           SELECT CASE
                      WHEN EXISTS
                          (
                              SELECT * FROM Followers WHERE FollowerId = @UserId AND UserId = U.Id
                          )
                          THEN CAST(1 AS BIT)
                      ELSE CAST(0 AS BIT)
                      END) AS Followed
FROM Followers
         LEFT JOIN Users U on U.Id = Followers.FollowerId
WHERE UserId = @UserId
ORDER BY CreatedAt DESC
go

CREATE PROCEDURE GetUserFollows @FollowerId int
AS

SELECT Id,
       Username,
       Email,
       COALESCE(FirstName, ''),
       COALESCE(LastName, ''),
       COALESCE(ProfilePicture, ''),
       CreatedAt
FROM Followers
         LEFT JOIN Users U ON U.Id = Followers.UserId
WHERE FollowerId = @FollowerId
ORDER BY CreatedAt DESC
go

CREATE PROCEDURE GetUserInfo
    @UserId int
AS

SELECT Username,Email,COALESCE(FirstName,''),COALESCE(LastName,''),COALESCE(ProfilePicture,''),CreatedAt FROM Users where Id = @UserId
go

CREATE PROCEDURE GetUserPosts @UserId int
AS

SELECT Id, COALESCE(Description, ''), COALESCE(Image, ''), COALESCE(IsDeleted, 0), CreatedAt,
       (
           SELECT CASE
                      WHEN EXISTS
                          (
                              SELECT * FROM Likes WHERE PostId = Posts.Id AND UserId = @UserId
                          )
                          THEN CAST(1 AS BIT)
                      ELSE CAST(0 AS BIT)
                      END) AS Liked
FROM Posts
WHERE UserId = @UserId
ORDER BY CreatedAt DESC
go

CREATE PROCEDURE LikePost @PostId int,
@UserId int
AS
INSERT INTO Likes (PostId, UserId) VALUES (@PostId, @UserId)
go

CREATE PROCEDURE Newsfeed @UserId int
AS
SELECT U.Id,
       Username,
       Email,
       COALESCE(FirstName, ''),
       COALESCE(LastName, ''),
       COALESCE(ProfilePicture, ''),
       P.Id,
       COALESCE(Description, ''),
       COALESCE(Image, ''),
       COALESCE(IsDeleted, 0),
       P.CreatedAt,
       (
           SELECT CASE
                      WHEN EXISTS
                          (
                              SELECT * FROM Followers WHERE FollowerId = @UserId AND UserId = U.Id
                          )
                          THEN CAST(1 AS BIT)
                      ELSE CAST(0 AS BIT)
                      END) AS Followed,
       (
           SELECT CASE
                      WHEN EXISTS
                          (
                              SELECT * FROM Likes WHERE PostId = P.Id AND UserId = @UserId
                          )
                          THEN CAST(1 AS BIT)
                      ELSE CAST(0 AS BIT)
                      END) AS Liked
FROM Followers
         LEFT JOIN Users U ON U.Id = Followers.UserId
         LEFT JOIN Posts P ON U.Id = P.UserId
WHERE P.Id IS NOT NULL
  and FollowerId = @UserId
ORDER BY CreatedAt DESC
go

CREATE PROCEDURE UnfollowUser
    @FollowerId int,
    @UserId int
AS

DELETE FROM Followers WHERE UserId = @UserId AND FollowerId = @FollowerId
go

CREATE PROCEDURE UnlikePost @PostId int,
@UserId int
AS
DELETE FROM Likes WHERE PostId = @PostId AND UserId = @UserId
go

CREATE PROCEDURE UpdateUserInfo
    @Username nvarchar(100),
    @FirstName nvarchar(500),
    @LastName nvarchar(500),
    @Bio nvarchar(500),
    @UserId int
AS

UPDATE Users SET Username = @Username, FirstName = @FirstName, LastName = @LastName, Bio = @Bio WHERE Id = @UserId
go

CREATE PROCEDURE UpdateUserPP
    @ProfilePicture nvarchar(max),
    @UserId int
AS

UPDATE Users SET ProfilePicture = @ProfilePicture WHERE Id = @UserId
go

CREATE PROCEDURE UpdateUserPassword
    @OldPassword nvarchar(500),
    @NewPassword nvarchar(500),
    @UserId int
AS

UPDATE Users SET Password = @NewPassword WHERE Id = @UserId AND Password = @OldPassword
go

