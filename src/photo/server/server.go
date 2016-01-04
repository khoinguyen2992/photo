package server

import "net/http"

type Config struct {
	Server struct {
		Port string `json:"API_PORT"`
		Addr string `json:"API_ADDR"`
	} `json:"server"`

	RethinkDB struct {
		Port   string `json:"RETHINKDB_PORT"`
		Addr   string `json:"RETHINKDB_ADDR"`
		DBName string `json:"RETHINKDB_DBNAME"`
	} `json:"rethinkdb"`

	Auth struct {
		SessionName string `json:"SESSION_NAME"`
		CookieStore string `json:"COOKIE_STORE"`
		UserId      string `json:"USER_ID"`
	} `json:"auth"`

	Upload struct {
		UploadDir  string `json:"UPLOAD_DIR"`
		PathPrefix string `json:"PATH_PREFIX"`
	} `json:"upload"`

	Mail struct {
		Username string `json:"USERNAME"`
		Password string `json:"PASSWORD"`
		Port     string `json:"PORT"`
		Domain   string `json:"DOMAIN"`
		Subject  string `json:"SUBJECT"`
		Message  string `json:"MESSAGE"`
	} `json:"mail"`
}

func Start(cfg Config) {
	s := setup(cfg)

	listenAddr := cfg.Server.Addr + ":" + cfg.Server.Port

	l.Println("photo-server is listening on", listenAddr)
	http.ListenAndServe(listenAddr, s.Handler)
}
