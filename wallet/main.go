package main

import (
	"context"
	"flag"
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
	w       *astilectron.Window
	WalletDB *walletdb.WalletDB

	ctx 	context.Context
	client *rpc.Client
	ch = make(chan NodeMsg)
)

type NodeMsg struct {
	t string
	v interface{}
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
	WalletDB ,_ = walletdb.NewWalletDB("/Users/kimmegi/test.ldb")
	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		//Asset:    Asset,
		//AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug: *debuging,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
			//SubMenu: []*astilectron.MenuItemOptions{
			//	{
			//		Label: astilectron.PtrStr("About"),
			//		OnClick: func(e astilectron.Event) (deleteListener bool) {
			//			if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
			//				// Unmarshal payload
			//				var s string
			//				if err := json.Unmarshal(m.Payload, &s); err != nil {
			//					astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
			//					return
			//				}
			//				astilog.Infof("About modal has been displayed and payload is %s!", s)
			//			}); err != nil {
			//				astilog.Error(errors.Wrap(err, "sending about event failed"))
			//			}
			//
			//
			//
			//			return
			//		},
			//	},
			//	{Role: astilectron.MenuItemRoleClose},
			//},
		}},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go func() {
				time.Sleep(time.Second * 2)
				if err := bootstrap.SendMessage(w, "notify_show", ""); err != nil {
					astilog.Error(errors.Wrap(err, "sending check.out.menu event failed"))
				}
				for{
					node := <-ch
					switch node.t {
					case "client":
						client = node.v.(*rpc.Client)
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


