module berith-chain

go 1.19

replace github.com/BerithFoundation/berith-chain => ./

require (
	github.com/Azure/azure-storage-blob-go v0.0.0-20180712005634-eaae161d9d5e
	github.com/BerithFoundation/berith-chain v1.1.0
	github.com/alexmullins/zip v0.0.0-20180717182244-4affb64b04d0
	github.com/allegro/bigcache v1.1.1-0.20181022200625-bff00e20c68d
	github.com/aristanetworks/goarista v0.0.0-20170210015632-ea17b1a17847
	github.com/asticode/go-astilectron v0.8.1-0.20190813121736-df875a09e6cc
	github.com/asticode/go-astilectron-bootstrap v0.0.0-20190816065004-25b857285999
	github.com/asticode/go-astilog v1.0.1-0.20190608125316-952ff13d3f86
	github.com/btcsuite/btcd v0.0.0-20171128150713-2e60448ffcc6
	github.com/cespare/cp v0.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea
	github.com/docker/docker v17.12.0-ce-rc1.0.20180625184442-8e610b2b55bf+incompatible
	github.com/elastic/gosigar v0.8.1-0.20180330100440-37f05ff46ffa
	github.com/fatih/color v1.3.0
	github.com/fjl/memsize v0.0.0-20180418122429-ca190fb6ffbc
	github.com/gizak/termui v2.2.1-0.20170117222342-991cd3d38091+incompatible
	github.com/go-stack/stack v1.5.4
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.0-20170215233205-553a64147049
	github.com/hashicorp/golang-lru v0.0.0-20160813221303-0a025b7e63ad
	github.com/huin/goupnp v0.0.0-20161224104101-679507af18f3
	github.com/influxdata/influxdb v1.2.3-0.20180221223340-01288bdb0883
	github.com/jackpal/go-nat-pmp v1.0.2-0.20160603034137-1fa385a6f458
	github.com/julienschmidt/httprouter v1.2.0
	github.com/karalabe/hid v0.0.0-20181128192157-d815e0c1a2e2
	github.com/mattn/go-colorable v0.1.0
	github.com/mattn/go-isatty v0.0.5-0.20180830101745-3fb116b82035
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/opentracing/opentracing-go v1.0.3-0.20180606204148-bd9c31933947
	github.com/pborman/uuid v0.0.0-20170112150404-1b00554d8222
	github.com/peterh/liner v1.0.1-0.20170902204657-a37ad3984311
	github.com/pkg/errors v0.8.1-0.20171216070316-e881fd58d78e
	github.com/prometheus/prometheus v0.0.0-20170814170113-3101606756c5
	github.com/rjeczalik/notify v0.9.1
	github.com/robertkrimen/otto v0.2.0
	github.com/rs/cors v0.0.0-20160617231935-a62a804a8a00
	github.com/stretchr/testify v1.8.1
	github.com/syndtr/goleveldb v0.0.0-20181128100959-b001fa50d6b2
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.2.0
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	golang.org/x/sys v0.2.0
	gopkg.in/check.v1 v1.0.0-20161208181325-20d25e280405
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20180302121509-abf0ba0be5d5
	gopkg.in/urfave/cli.v1 v1.20.0
)

require (
	github.com/Azure/azure-pipeline-go v0.0.0-20180607212504-7571e8eb0876 // indirect
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6 // indirect
	github.com/akavel/rsrc v0.8.0 // indirect
	github.com/asticode/go-astilectron-bundler v0.0.0-20190426172205-155c2a10bbb1 // indirect
	github.com/asticode/go-astitools v1.2.1-0.20190929114647-d157a994ecbd // indirect
	github.com/asticode/go-bindata v0.0.0-20151023091102-a0ff2567cfb7 // indirect
	github.com/go-ole/go-ole v1.2.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/maruel/panicparse v0.0.0-20160720141634-ad661195ed0e // indirect
	github.com/mattn/go-runewidth v0.0.3 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
	github.com/nsf/termbox-go v0.0.0-20170211012700-3540b76b9c77 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/xhandler v0.0.0-20160618193221-ed27b6fd6521 // indirect
	github.com/sam-kamerer/go-plister v1.1.2 // indirect
	github.com/sirupsen/logrus v1.4.3-0.20190807103436-de736cf91b92 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
