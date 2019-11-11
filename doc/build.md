## Building the source

Building ```berith``` requires both a Go (1.10 or later) and a C compiler.  

> berith node  

```bash
make berith
```

> full suite of utilities  

```bash
make all
```  

> berith wallet  

`TODO : will support cross compile`

Building `wallet` also requires astilectron-bundler  

> install astilectron-bundler  

don't forget to add `$GOPATH/bin` to your `$PATH`.  

```bash
$ go get -u github.com/asticode/go-astilectron-bundler/...
$ go install github.com/asticode/go-astilectron-bundler/astilectron-bundler
$ astilectron-bundler -v
```  

> build wallet  

```bash
$ cd $GOPATH/src/github/BerithFoundation/berith-chain/wallet
$ astilectron-bundler -o /outputpath
```  