module github.com/bee-travels/bee-travels-go/services/destination-v2

go 1.14

require (
	github.com/elgris/sqrl v0.0.0-20190909141434-5a439265eeec
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/instana/go-sensor v1.29.0
	github.com/iris-contrib/middleware/cors v0.0.0-20210110101738-6d0a4d799b5d
	github.com/jackc/pgconn v1.8.1
	github.com/jackc/pgproto3/v2 v2.0.6
	github.com/jackc/pgx/v4 v4.11.0
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/iris/v12 v12.2.0-alpha8
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.8.1
)

replace github.com/russross/blackfriday/v2 => gopkg.in/russross/blackfriday.v2 v2.1.0
