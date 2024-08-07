module github.com/Liphium/station/chatserver

go 1.22.1

require (
	github.com/Liphium/station/main v0.0.0-20240516155328-010091dd4965
	github.com/Liphium/station/pipes v0.0.0-20240516155328-010091dd4965
	github.com/Liphium/station/pipeshandler v0.0.0-20240516155328-010091dd4965
	github.com/gofiber/contrib/jwt v1.0.9
	github.com/gofiber/fiber/v2 v2.52.4
	github.com/valyala/fasthttp v1.54.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgraph-io/ristretto v0.1.1
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/glog v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.26.0 // indirect
	nhooyr.io/websocket v1.8.11 // indirect
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gorm.io/driver/postgres v1.5.7
	gorm.io/gorm v1.25.10
)

require (
	github.com/MicahParks/keyfunc/v2 v2.1.0 // indirect
	github.com/google/uuid v1.6.0
)

require (
	github.com/bytedance/sonic v1.11.8
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.8.0 // indirect
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/fasthttp/websocket v1.5.9 // indirect
	github.com/gofiber/websocket/v2 v2.2.1 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/savsgio/gotils v0.0.0-20240303185622-093b76447511 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
)

replace github.com/Liphium/station/main => ../main

replace github.com/Liphium/station/pipes => ../pipes

replace github.com/Liphium/station/pipeshandler => ../pipeshandler

replace github.com/Liphium/station/chatserver => ../chatserver

replace github.com/Liphium/station/spacestation => ../spacestation

replace github.com/Liphium/station/backend => ../backend
