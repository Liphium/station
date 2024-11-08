module github.com/Liphium/station/pipes

go 1.23.0

require (
	github.com/bytedance/sonic v1.12.4
	github.com/coder/websocket v1.8.12
	github.com/dgraph-io/ristretto v0.2.0
)

require (
	github.com/bytedance/sonic/loader v0.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.12.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
)

replace github.com/Liphium/station/pipes => ../pipes

replace github.com/Liphium/station/pipeshandler => ../pipeshandler

replace github.com/Liphium/station/chatserver => ../chatserver

replace github.com/Liphium/station/spacestation => ../spacestation

replace github.com/Liphium/station/backend => ../backend
