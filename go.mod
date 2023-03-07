module github.com/kostua16/go_simple_logger

go 1.16

require (
	github.com/glebarez/sqlite v1.4.3
	github.com/go-chi/chi v1.5.4
	github.com/hudl/fargo v1.4.0
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	golang.org/x/net v0.7.0 // indirect
	gorm.io/gorm v1.23.5
)

// +heroku goVersion go1.16
// +heroku install ./cmd/...
