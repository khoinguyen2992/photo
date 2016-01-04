package server

import (
	"net/http"
	"photo/api/rethink"
	"photo/api/xhttp"
	"photo/handlers/account"
	"photo/handlers/comment"
	"photo/handlers/follower"
	"photo/handlers/photo"
	"photo/middlewares"
	authService "photo/services/auth"
	"photo/stores"
	"photo/utils/logs"
	"photo/utils/sendmail"

	"github.com/dancannon/gorethink"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

var l = logs.New("photo/server")

type setupStruct struct {
	Config

	Rethink *rethink.Instance
	Handler http.Handler
}

func setup(cfg Config) *setupStruct {
	s := &setupStruct{Config: cfg}
	s.setupRethink()
	s.setupRoutes()

	return s
}

func (s *setupStruct) setupRethink() {
	cfg := s.Config
	re, err := rethink.NewInstance(gorethink.ConnectOpts{
		Address:  cfg.RethinkDB.Addr + ":" + cfg.RethinkDB.Port,
		Database: cfg.RethinkDB.DBName,
	})

	if err != nil {
		l.Fatalln("Could not connect to RethinkDB")
	}

	s.Rethink = re
}

func commonMiddlewares() func(http.Handler) http.Handler {
	logger := middlewares.NewLogger()
	recovery := middlewares.NewRecovery()

	return func(h http.Handler) http.Handler {
		return recovery(logger(h))
	}
}

func authMiddlewares(s *setupStruct) func(http.Handler) http.Handler {
	accountStore := stores.NewAccountStore(s.Rethink)
	accountAuth := authService.NewAccountAuth(s.Config.Auth.SessionName, s.Config.Auth.CookieStore)
	auth := middlewares.NewAuth(accountStore, accountAuth)

	return func(h http.Handler) http.Handler {
		return auth(h)
	}
}

func (s *setupStruct) setupRoutes() {
	commonMids := commonMiddlewares()
	authMids := authMiddlewares(s)

	normal := func(h http.HandlerFunc) httprouter.Handle {
		return xhttp.Adapt(commonMids(h))
	}

	auth := func(h http.HandlerFunc) httprouter.Handle {
		return xhttp.Adapt(commonMids(authMids(h)))
	}

	router := httprouter.New()

	followerStore := stores.NewFollowerStore(s.Rethink)
	photoStore := stores.NewPhotoStore(s.Rethink)
	commentStore := stores.NewCommentStore(s.Rethink)
	accountStore := stores.NewAccountStore(s.Rethink)
	accountAuth := authService.NewAccountAuth(s.Config.Auth.SessionName, s.Config.Auth.CookieStore)
	sendMail := sendmail.NewSendMail(s.Config.Mail.Username, s.Config.Mail.Password, s.Config.Mail.Domain, s.Config.Mail.Subject, s.Config.Mail.Message)

	{
		accountCtrl := account.NewAccountCtrl(accountStore, accountAuth)
		router.GET("/v1/me", auth(accountCtrl.Me))
		router.POST("/v1/login", normal(accountCtrl.Login))
		router.POST("/v1/logout", auth(accountCtrl.Logout))
		router.POST("/v1/register", normal(accountCtrl.Register))
		router.GET("/v1/account", auth(accountCtrl.List))
		router.PUT("/v1/account/change_password", auth(accountCtrl.UpdatePassword))
		router.PUT("/v1/account/change_profile", auth(accountCtrl.UpdateProfile))
	}

	{
		photoCtrl := photo.NewPhotoCtrl(photoStore, s.Config.Upload.UploadDir, s.Config.Upload.PathPrefix)
		router.ServeFiles(s.Config.Upload.PathPrefix+"*filepath", http.Dir(s.Config.Upload.UploadDir))
		router.POST("/v1/photo", auth(photoCtrl.Create))
		router.GET("/v1/photo/me", auth(photoCtrl.ListByOwner))
		router.GET("/v1/photo/follow", auth(photoCtrl.ListByFollowers))
		router.GET("/v1/photo/account/:id", auth(photoCtrl.ListByAccount))
		router.PUT("/v1/photo/private/:id", auth(photoCtrl.UpdatePrivate))
		router.POST("/v1/photo/crop/:savename/:width/:height", auth(photoCtrl.CreateCropImage))
	}

	{
		commentCtrl := comment.NewCommentCtrl(commentStore, photoStore)
		router.POST("/v1/comment", auth(commentCtrl.Create))
		router.GET("/v1/comment/tags", auth(commentCtrl.ListByTags))
		router.GET("/v1/comment/photo/:id", auth(commentCtrl.ListByPhoto))
		router.GET("/v1/comment/notification", auth(commentCtrl.ListByNotification))
		router.PUT("/v1/comment/known/:id", auth(commentCtrl.UpdateNotification))
		router.DELETE("/v1/comment/:id", auth(commentCtrl.Delete))
	}

	{
		followerCtrl := follower.NewFollowerCtrl(followerStore, accountStore, sendMail)
		router.POST("/v1/follow/account/:id", auth(followerCtrl.Create))
		router.DELETE("/v1/follow/:id", auth(followerCtrl.Delete))
	}

	s.Handler = context.ClearHandler(router)
}
