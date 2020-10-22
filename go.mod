module github.com/10000ms/neko_server_go

go 1.14

require (
	neko_server_go v0.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/satori/go.uuid v1.2.0
)

replace (
	neko_server_go v0.0.0 => ./
)
