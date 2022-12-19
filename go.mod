module berith-chain

go 1.14

require (
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.3.0
	github.com/BerithFoundation/berith-chain v0.0.0-00010101000000-000000000000
	github.com/alexmullins/zip v0.0.0-20180717182244-4affb64b04d0
	github.com/allegro/bigcache v1.2.1
	github.com/aristanetworks/goarista v0.0.0-20200429182514-19402535e24e
	github.com/asticode/go-astilectron v0.9.1
	github.com/asticode/go-astilectron-bootstrap v0.1.0
	github.com/asticode/go-astilog v1.2.0
	github.com/cespare/cp v1.1.1
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1
	github.com/docker/docker v1.13.1
	github.com/elastic/gosigar v0.14.2
	github.com/fatih/color v1.9.0
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5
	github.com/gizak/termui v2.3.0+incompatible
	github.com/go-stack/stack v1.8.0
	github.com/golang/protobuf v1.4.2
	github.com/golang/snappy v0.0.4
	github.com/gookit/color v1.2.5
	github.com/hashicorp/golang-lru v0.5.4
	github.com/huin/goupnp v1.0.0
	github.com/influxdata/influxdb v1.8.0
	github.com/jackpal/go-nat-pmp v1.0.2
	github.com/julienschmidt/httprouter v1.3.0
	github.com/karalabe/hid v1.0.0
	github.com/mattn/go-colorable v0.1.13
	github.com/mattn/go-isatty v0.0.16
	github.com/naoina/toml v0.1.1
	github.com/pborman/uuid v1.2.0
	github.com/peterh/liner v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/prometheus v1.8.2
	github.com/rjeczalik/notify v0.9.2
	github.com/robertkrimen/otto v0.0.0-20211024170158-b87d35c0b86f
	github.com/rs/cors v1.7.0
	github.com/stretchr/testify v1.7.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	golang.org/x/crypto v0.4.0
	golang.org/x/net v0.4.0
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	golang.org/x/sys v0.3.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200316214253-d7b0ff38cac9
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace github.com/BerithFoundation/berith-chain => ./
