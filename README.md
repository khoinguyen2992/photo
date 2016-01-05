# Photo

This is a server of the Golang test app. It makes sure appropriate application logic and provides necessary API for client side development.

# How to run

*  recommend linux or vagrant
*  install go1.4 - [gvm](https://github.com/moovweb/gvm)
*  install node - [nvm](https://github.com/creationix/nvm)
*  install rethinkdb - [rethinkdb.com](http://rethinkdb.com/)
*  git clone git@github.com:khoinguyen2992/photo.git
*  cd photo
* export GOPATH=$PWD
* go install photo/cmd/photo-server
* node res/load.js res/createDatabase.js

#References
* Architecture - [grox](https://github.com/golang-vietnam/grox)
*  Gorethink - [gorethink](https://github.com/dancannon/gorethink) - for communicating with rethinkdb
*  Imaging - [imaging](https://github.com/disintegration/imaging) - for croping images
*  Gomail - [gomail](https://github.com/go-gomail/gomail/) - for sending email
*  Govalidator - [govalidator](https://github.com/asaskevich/govalidator) - for validation

# API Docs

```
	[<field>: <type>]
```

fields inside the square bracket are hidden from client 

###Account

```
	{
		"id": string,
		"username": string,
		"avatar": string,
		"first_name": string,
		"last_name": string,
		"followers": [],
		"created_time": time,
		"updated_time": time,
		["salt": string,]
		["secret": string,]
	}
```

```
	GET /v1/me
```

to check if user has logined or not

```
	GET /v1/account
```

to list user profiles

```
	POST /v1/login

	{
		"username": string
		"password": string
	}
```

to login account

```
	POST /v1/logout
```

to logout

```
	POST /v1/register

	{
		"username": string,
		"password": string
	}
```

to register new account

```
	PUT /v1/account/change_password

	{
		"old_password": string
		"new_password": string
	}
```
to change password

```
	PUT /v1/account/change_profile

	{
		"avatar": string,
		"first_name" : string,
		"last_name": string
	}
```

###Photo

```
	{
		"id": string,
		"account_id": string,
		"save_name:" string, 
		"uri": string,
		"created_time": time,
		"updated_time": time,
		["is_private": bool,]
	}
```

```
	GET /v1/photo/me
```

to list your own photos

```
	GET /v1/photo/follow
```

to list the photos of accounts you follow

```
	GET /v1/photo/account/:id
```

to list the photos of an account with id

```
	POST /v1/photo

	form-data
	{
		file: file
	}
```

to create new photo

```
	POST /v1/photo/crop/:savename/:width/:height
```

to crop the photo (center)

```
	PUT /v1/photo/private/:id
```

to set photo to private mode

###Comment

```
	{
		"id": string,
		"photo_id": string,
		"text": string,
		"account_id": string,
		["notification_id": string,]
		["tags": []string,]
		["is_known": bool]
	}
```

```
	GET /v1/comment/tags?tags=string
```

to list photos by comment tags. not really make sense :( 

```
	GET /v1/comment/photo/:id
```

to list comments by photo id

```
	GET /v1/comment/notification
```

to list new comments (has not been known yet)

```
	POST /v1/comment

	{
		"photo_id": string,
		"text": string
	}
```

to create new comment

```
	PUT /v1/comment/known/:id
```

to change the comment to be known

```
	DELETE /v1/comment/:id
```

to delete comment

###Followers

```
	{
		"id": string,
		"account_id": string, //account to be followed
		"follower_id": string, //account follow
		"created_time": time,
		"updated_time": time
	}
```

```
	POST /v1/follow/account/:id
```

to follow account with the id. include send mail to the account followed

```
	DELETE /v1/follow/:follow_id
```

to delete follow connection, both follow and be-followed sides have the permission to delete 
