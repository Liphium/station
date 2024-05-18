module github.com/Liphium/station/pipes

go 1.22.1

require (
	github.com/bytedance/sonic v1.11.3
	github.com/dgraph-io/ristretto v0.1.1
	nhooyr.io/websocket v1.8.10
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20230717121745-296ad89f973d // indirect
	github.com/chenzhuoyu/iasm v0.9.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/glog v1.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	golang.org/x/arch v0.7.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)

replace github.com/Liphium/station/pipes => ../pipes

replace github.com/Liphium/station/pipeshandler => ../pipeshandler

replace github.com/Liphium/station/chatserver => ../chatserver

replace github.com/Liphium/station/spacestation => ../spacestation

replace github.com/Liphium/station/backend => ../backend
