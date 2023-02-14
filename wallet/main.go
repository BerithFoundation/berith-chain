package main

import (
	"context"
	"flag"
	"path/filepath"
	"strings"
	"time"

	"github.com/BerithFoundation/berith-chain/node"
	"github.com/BerithFoundation/berith-chain/rpc"
	walletdb "github.com/BerithFoundation/berith-chain/wallet/database"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Vars
var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
	// === wallet flags
	debuging = flag.Bool("d", false, "enables the debug mode")

	// === node flags
	nodePort       = flag.String("nodeport", "", "node's network listening port")
	verbosity      = flag.Int("verbosity", 3, "logging verbosity: 0=silent, 1=error, 2=warn, 3=info, 4=debug, 5=detail")
	node_testnet   = flag.String("testnet", "", "testnet")
	node_console   = flag.String("console", "", "console")
	node_datadir   = flag.String("datadir", "", "datadir")
	nodeConfig     = flag.String("nodeconfig", "", "config file path")
	nodiscover     = flag.Bool("nodiscover", false, "nodiscover")
	httpFlag       = flag.Bool("http", false, "open http connection")
	httpCorsDomain = flag.String("http.corsdomain", "https://remix.ethereum.org/", "set cors domain")
	httpApi        = flag.String("http.api", "", "set http api")
	w              *astilectron.Window
	WalletDB       *walletdb.WalletDB

	ctx    context.Context
	client *rpc.Client
	stack  *node.Node

	ch  = make(chan NodeMsg)
	ch2 = make(chan bool)
)

type NodeMsg struct {
	t     string
	v     interface{}
	stack interface{}
}

func init() {
	Init()
}
func main() {
	startUi()
}

// db , 일렉트론 초기설정 후 wallet 프로그램 실행 함수
func startUi() {
	// Init
	flag.Parse()
	astilog.FlagInit()
	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/app/images/common/berith_pa.ico",
			AppIconDefaultPath: "resources/app/images/common/berith_pa.ico",
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
			DataDirectoryPath:  filepath.Join(node.DefaultDataDir(), "wallet"),
		},
		Debug: *debuging,
		// MenuOptions: []*astilectron.MenuItemOptions{{
		// 	Label: astilectron.PtrStr("File"),
		// }},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go func() {
				time.Sleep(time.Second * 2)
				if err := bootstrap.SendMessage(w, "notify_show", ""); err != nil {
					astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
				}
				for {
					nodeChannel := <-ch
					switch nodeChannel.t {
					case "client":
						client = nodeChannel.v.(*rpc.Client)
						stack = nodeChannel.stack.(*node.Node)
						ctx = context.TODO()
						if err := bootstrap.SendMessage(w, "notify_hide", ""); err != nil {
							astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
						}
						w.On(astilectron.EventNameWindowEventClosed, func(e astilectron.Event) (deleteListener bool) {
							if stack == nil {
								return false
							}
							stack.Stop()
							return true
						})
						//startPolling()
						break
					}
				}
			}()
			go Start()
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage: "html/login.html",
			//Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(1250),
				Width:           astilectron.PtrInt(1250),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

// 동기화 여부 , 최신블록넘버  반복조회 함수
func startPolling() {
	go func() {
		for {
			//isPolling := <- ch2
			//if !isPolling {
			//	break
			//}
			sync, err := callNodeApi("berith_syncing", nil)
			if err != nil {
				astilog.Error(errors.Wrap(err, "syncing failed"))
			}

			// 동기화 완료시 최신 블록 조회
			if sync == "false" {
				blockNum, err2 := callNodeApi("berith_blockNumber")
				if err2 != nil {
					astilog.Error(errors.Wrap(err, "blockNumber Failed"))
					return
				}
				blockNum = strings.Replace(blockNum, "\"", "", -1)
				blockInfo, err3 := callNodeApi("berith_getBlockByNumber", blockNum, true)
				if err3 != nil {
					astilog.Error(errors.Wrap(err, "getBlockByNumber Failed"))
					return
				}
				if err := bootstrap.SendMessage(w, "getBlockInfo", blockInfo); err != nil {
					astilog.Error(errors.Wrap(err, "getBlockInfo failed"))
					return
				}
			}

			if err := bootstrap.SendMessage(w, "syncing", sync); err != nil {
				astilog.Error(errors.Wrap(err, "syncing failed"))
				return
			}
			//3초간격
			time.Sleep(3 * time.Second)
		}
	}()
}
