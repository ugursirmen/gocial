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

