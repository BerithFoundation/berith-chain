module berith-chain

go 1.14

require (
	github.com/Azure/azure-storage-blob-go v0.7.0
	github.com/BerithFoundation/berith-chain v0.0.0-00010101000000-000000000000
	github.com/allegro/bigcache v1.2.1
	github.com/aristanetworks/goarista v0.0.0-20200429182514-19402535e24e
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/cespare/cp v1.1.1
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1
	github.com/docker/docker v1.13.1
	github.com/elastic/gosigar v0.10.5
	github.com/fatih/color v1.9.0
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08
	github.com/go-stack/stack v1.8.1
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.3
	github.com/gookit/color v1.2.5
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/holiman/uint256 v1.2.0
	github.com/huin/goupnp v1.0.1-0.20210310174557-0ca763054c88
	github.com/influxdata/influxdb v1.8.3
	github.com/jackpal/go-nat-pmp v1.0.2
	github.com/julienschmidt/httprouter v1.3.0
	github.com/karalabe/hid v1.0.0
	github.com/mattn/go-colorable v0.1.7
	github.com/mattn/go-isatty v0.0.12
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/pborman/uuid v1.2.1
	github.com/peterh/liner v1.2.2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/prometheus v1.8.2
	github.com/rjeczalik/notify v0.9.2
	github.com/robertkrimen/otto v0.0.0-20211024170158-b87d35c0b86f
	github.com/rs/cors v1.7.0
	github.com/status-im/keycard-go v0.0.0-20220906070205-e43cb0f06ae9
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2
	golang.org/x/sys v0.0.0-20211117180635-dee7805ff2e1
	golang.org/x/text v0.3.7
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20200619000410-60c24ae608a6
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace github.com/BerithFoundation/berith-chain => ./
