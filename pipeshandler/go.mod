module github.com/Liphium/station/pipeshandler

go 1.23.0

require (
	github.com/Liphium/station/chatserver v0.0.0-20240819105432-c6d7be0605b2
	github.com/Liphium/station/pipes v0.0.0-20240819105432-c6d7be0605b2
	github.com/bytedance/sonic v1.12.2
	github.com/dgraph-io/ristretto v0.1.1
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/gofiber/websocket/v2 v2.2.1
	github.com/golang-jwt/jwt/v5 v5.2.1
)

require (
	github.com/Liphium/station/main v0.0.0-20240819105432-c6d7be0605b2 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/coder/websocket v1.8.12 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fasthttp/websocket v1.5.10 // indirect
	github.com/golang/glog v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/savsgio/gotils v0.0.0-20240704082632-aef3928b8a38 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.55.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/arch v0.9.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
)

replace github.com/Liphium/station/pipes => ../pipes

replace github.com/Liphium/station/pipeshandler => ../pipeshandler

replace github.com/Liphium/station/chatserver => ../chatserver

replace github.com/Liphium/station/spacestation => ../spacestation

replace github.com/Liphium/station/backend => ../backend
