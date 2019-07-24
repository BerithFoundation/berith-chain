package main

import (
	"context"
	"flag"
	"github.com/BerithFoundation/berith-chain/node"
	"github.com/BerithFoundation/berith-chain/rpc"
	"github.com/BerithFoundation/berith-chain/wallet/database"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"time"
)


// Vars
var (
	AppName string
	BuiltAt string
	debuging   = flag.Bool("d", true, "enables the debug mode")
	node_testnet = flag.String("testnet", "", "testnet")
	node_console = flag.String("console", "", "console")
	node_datadir = flag.String("datadir", "", "datadir")
	//node_berithbase = flag.String("miner.berithbase", "", "berithbase")
	w       *astilectron.Window
	WalletDB *walletdb.WalletDB

	ctx 	context.Context
	client *rpc.Client
	stack *node.Node

	ch = make(chan NodeMsg)
)

type NodeMsg struct {
	t string
	v interface{}
	stack interface{}
}

func init(){
	Init()
}
func main() {
	start_ui()
}
func start_ui(){
	// Init
	flag.Parse()
	astilog.FlagInit()
	WalletDB ,_ = walletdb.NewWalletDB("/Users/usman/test.ldb")
	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			DataDirectoryPath:"C:\\Users\\Usman\\go\\src\\github.com\\BerithFoundation\\berith-chain\\wallet",
		},
		Debug: *debuging,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
		}},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go func() {
				time.Sleep(time.Second * 2)
				if err := bootstrap.SendMessage(w, "notify_show", ""); err != nil {
					astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
				}
				for{
					nodeChannel := <-ch
					switch nodeChannel.t {
					case "client":
						client = nodeChannel.v.(*rpc.Client)
						stack = nodeChannel.stack.(*node.Node)
						ctx = context.TODO()
						if err := bootstrap.SendMessage(w, "notify_hide", ""); err != nil {
							astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
						}
						//startPolling()

						break

					}
				}
			}()
			go Start()
			return nil
		},
		//RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "/html/login.html",
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

func startPolling(){
	go func() {
		for {
			// polling 로직
			val, err := callNodeApi("berith_syncing", nil)
			if err != nil {
				astilog.Error(errors.Wrap(err, "polling failed"))
			}

			if err := bootstrap.SendMessage(w, "polling", val); err != nil {
				astilog.Error(errors.Wrap(err, "polling failed"))
			}
			//coinbase 조회
			val2, err2 := callNodeApi("berith_coinbase" , nil)
			if err2 != nil {
				astilog.Error(errors.Wrap(err2, "coinbase null"))
			}
			if err2 := bootstrap.SendMessage(w, "coinbase", val2); err2 != nil {
				astilog.Error(errors.Wrap(err2, "coinbase null"))
			}

			time.Sleep(3 * time.Second)
		}
	}()
}


