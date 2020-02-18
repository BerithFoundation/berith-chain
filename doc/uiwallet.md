## 베리드 UI Wallet

베리드 PC 월렛(이하 월렛)은 베리드 POS 노드(이하 노드) 의 기능을 GUI로 제공하기 위해 만들어진 Electron 기반 프로그램이다. 월렛은 백그라운드로 노드를 실행시키며 Electron 에서 발생한 요청을 RPC로 노드가 처리하여 결과값을 전달하고 전달된 결과값을 통해 월렛이 화면을 전환하는 방식으로 동작한다. 

### 월렛과 노드 사이의 통신

월렛은 javascript 프로그램이고, 노드는 golang 프로그램이다. 베리드는 둘 사이의 통신을 위해 

 [https://github.com/asticode/go-astilectron](https://github.com/asticode/go-astilectron)

라이브러리를 사용한다. 해당 라이브러리는 양쪽 프로그램이 통신할 수 있는 함수를 제공한다.
```
async function sendMessage(methodType, methodName, args) {
    let messagePromise = new Promise(function (resolve) {
        let message = {"name": methodType};
        message.payload = {
            "api": methodName,
            "args": args
        }
        //asticode.loader.show()

        //console.log("Request: ", JSON.stringify(message));syncing
        astilectron.sendMessage(message, function (response) {
            //asticode.loader.hide();
            //console.log("Response: ", JSON.stringify(response));
            resolve(response);
        });
    });
    return messagePromise;
}
```
위의 코드는 월렛 프로그램이 노드와 통신하기 위해 작성된 함수의 내용이다. 코드에서 ```astilection.sendMessage``` 함수를 사용하는 것을 확인할 수 있다.
```
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
    var info map[string]interface{}
    err = json.Unmarshal(m.Payload, &info)
    if err != nil {
        payload = nil
        return
    }
    api := info["api"]
    args := info["args"].([]interface{})
    switch m.Name {

    case "init":
        break
    case "polling":
        //ch2 <- true
        startPolling()
        break
    case  "stopPolling" :
        //ch2 <- false
        fmt.Print("logout!!!!")

        break
    case "callApi":
        payload, err = callNodeApi(api, args...)
        break
    case "callDB":
        payload , err = callDB(api , args...)
        break
    case "exportKeystore":
        args := info["args"].([]interface{})
        payload, err = exportKeystore(args)
        break

    case "importKeystore":
        args := info["args"].([]interface{})
        err = importKeystore(args)
        payload = nil
        break
    }

    if err != nil {
        astilog.Error(err.Error())
    }
    astilog.Debugf("Payload: %s", payload)
    return
}
```
위의 코드는 월렛 프로그램의 요청을 처리하기 위한 노드 프로그램의 핸들러 함수의 내용이다. 

### 월렛과 로그

월렛은 백그라운드에서 동작하는 노드의 상태를 파악하기 위해 노드에서 발생하는 로그를 저장한다. 로그는 두가지 저장소에 저장된다. 베리드 이더리움에서 사용하던 로그 패키지를 그대로 사용한다. 월렛의 추가적인 로그 출력을 위해 추가적인 로그 핸들러를 등록한다.
```
func Init() {

        ...

		app.Before = func(ctx *cli.Context) error {

		//TODO : wallet program shoud export log file without debug flag
		logdir := filepath.Join(node.DefaultDataDir(), "logs")

		batch = log.NewBerithLogBatch(logCh, logdir, time.Hour*24, log.TerminalFormat(false))

		go batch.Loop()

		if err := debug.SetupForWallet(ctx, logCh); err != nil {
			return err
		}
		
		...
		
}

func SetupForWallet(ctx *cli.Context, ch chan *log.Record) error {
	// logging
	log.PrintOrigins(ctx.GlobalBool(debugFlag.Name))

	glogger.SetHandler(log.MultiHandler(ostream, log.ChannelHandler(ch)))
	glogger.Verbosity(log.Lvl(ctx.GlobalInt(verbosityFlag.Name)))
	glogger.Vmodule(ctx.GlobalString(vmoduleFlag.Name))
	glogger.BacktraceAt(ctx.GlobalString(backtraceAtFlag.Name))
	log.Root().SetHandler(glogger)

	// profiling, tracing
	runtime.MemProfileRate = ctx.GlobalInt(memprofilerateFlag.Name)
	Handler.SetBlockProfileRate(ctx.GlobalInt(blockprofilerateFlag.Name))
	if traceFile := ctx.GlobalString(traceFlag.Name); traceFile != "" {
		if err := Handler.StartGoTrace(traceFile); err != nil {
			return err
		}
	}
	if cpuFile := ctx.GlobalString(cpuprofileFlag.Name); cpuFile != "" {
		if err := Handler.StartCPUProfile(cpuFile); err != nil {
			return err
		}
	}

	// pprof server
	if ctx.GlobalBool(pprofFlag.Name) {
		address := fmt.Sprintf("%s:%d", ctx.GlobalString(pprofAddrFlag.Name), ctx.GlobalInt(pprofPortFlag.Name))
		StartPProf(address)
	}
	return nil
}
```
위의 코드는 두가지 함수의 내용중 일부이다.

1. 노드를 실행할 때 가장 먼저 실행되는 함수이다. 위의 코드에서 새롭게 만든 핸들러가 로그를 받는 채널을 로그 패키지에 핸들러를 등록하는 함수의 인자로 넘겨주는 것을 확인할 수 있다.

1. 로그 패키지에 핸들러를 등록하는 함수이다. 위의 코드에서 전달받은 채널을 멀티 핸들러로 기존의 핸들러와 함께 로그 패키지에 등록하는 것을 확인 할 수 있다.

#### 로컬 파일

기본적으로 로그는 노드에서 설정된 ```Datadir``` 경로 하위에 ```logs``` 디렉터리에 파일 형태로 저장된다. 
```
func DefaultDataDir() string {
    // Try to place the data folder in the user's home dir
    home := homeDir()
    if home != "" {
        switch runtime.GOOS {
        case "darwin":
            return filepath.Join(home, "Library", "Berith")
        case "windows":
            // We used to put everything in %HOME%\AppData\Roaming, but this caused
            // problems with non-typical setups. If this fallback location exists and
            // is non-empty, use it, otherwise DTRT and check %LOCALAPPDATA%.
            fallback := filepath.Join(home, "AppData", "Roaming", "Berith")
            appdata := windowsAppData()
            if appdata == "" || isNonEmptyDir(fallback) {
                return fallback
            }
            return filepath.Join(appdata, "Berith")
        default:
            return filepath.Join(home, ".berith")
        }
    }
    // As we cannot guess a stable location, return empty and handle later
    return ""
}
```
위의 코드는 각 os 별로 기본 ```Datadir``` 을 지정하는 함수의 내용이다. ```Datadir``` 은 노드 설정에 따라 달라질 수 있다. 
```
func (b *BerithLogBatch) Loop() {
    for {
        select {
        case record := <-b.ch:
            b.cnt++
            if b.file == nil || time.Now().Sub(b.time) >= b.rotatePeriod {

                if err := os.MkdirAll(b.logdir, 0700); err != nil {
                    continue
                }
                now := time.Now()
                logpath := filepath.Join(b.logdir, strings.Replace(now.Format("060102150405.00"), ".", "", 1)+".log")
                logfile, err := os.Create(logpath)

                if err != nil {
                    continue
                }

                b.file.Close()
                b.file = logfile
                b.time = now

            }
            b.file.Write(b.format.Format(record))
            b.buffer += string(b.format.Format(record))
            
      ...

}
```
위의 코드는 작성된 로그 핸들러에서 파일 로그를 저장하는 부분이다. 현재 시간을 파일명으로 사용하여 파일을 만들고 시간을 저장한다. 로그 입력을 요청하는 신호가 왔을 때, 현재 시간이 저장된 시간보다 하루만큼 크다면 새로운 파일을 생성한다.

#### 로그 수집 서버

베리드 서비스가 안정화 될 때 까지 원활하게 문제를 관리하기 위해 모든 월렛 사용자들의 로그를 사내 서버에 파일 형태로 수집한다. 노드의 로그 출력 요청이 올 때마다 로그 내용을 버퍼에 저장했다가 100번째 요청이 왔을 때 서버로 저장했던  버퍼의 내용을 전송한다. 이 후 요청 갯수와 버퍼를 초기화하여 다음 100개의 요청을 기다린다.
```
func (b *BerithLogBatch) Loop() {
    for {
        select {
        case record := <-b.ch:
            b.cnt++
            
            ...

            if b.cnt == 100 {
                go b.handler(b.buffer)
                b.cnt = 0
                b.buffer = ""
            }

          ...

    }
}

handler := func(buffer string) {
		if stack != nil {
			rpcHandler, err := stack.RPCHandler()

			if err != nil {
				return
			}

			cli := rpc.DialInProc(rpcHandler)

			nodeInfo := p2p.NodeInfo{}
			berithbase := common.Address{}
			if err := cli.CallContext(context.Background(), &nodeInfo, "admin_nodeInfo"); err != nil {
				return
			}

			if err := cli.CallContext(context.Background(), &berithbase, "berith_coinbase"); err != nil {
				return
			}

			jsonByte, err := json.Marshal(LogPost{
				Enode:      nodeInfo.Enode,
				Berithbase: berithbase.Hex(),
				Logs:       buffer,
			})

			if err != nil {
				return
			}

			http.Post("https://baas.berith.co/v1/api/logs/bers", "application/json", bytes.NewReader(jsonByte))
		}
	}
```
위의 코드는 100개의 요청을 버퍼의 저장하는 부분과, 사내 서버로 버퍼의 내용을 전달하는 핸들러의 내용이다. 로그는 사내 aws 계정의 ```Berith-was01(R)``` EC2 인스턴스로 전달된다.

