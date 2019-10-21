# berith executables  

> ## index

- <a href="#berit">berith</a>
- <a href="#berithench">berithench</a>
- <a href="#berithkey">berithkey</a>
- <a href="#bootnode">bootnode</a>  
- <a href="#gwizard">gwizard</a>
- <a href="#wallet">wallet</a>

---  

<div id="berith"></div>

---  

<div id="berithench"></div>  

> ## berithench  

berithench is simple stress test tools for berith node.

- execute transfer  
; send only transfer transactions with repeatedly or duration.

> config.toml  
$berithench execute transfer --config config.toml

```bash
ChainID = 36435
# endpoints of nodes
Nodes = [
  "http://localhost:8501",
  "http://localhost:8502"
]

# keystores path
Keystore = "/home/app/workspaces/berith/keystore"

# from accounts for tx.
Addresses = [
  "8676fb254279ef78c53b8a781e228ab439065786",
  "ca7207de79e55c1a69dbc67a4a2e81dfc62c6ac4",
  "d8a25ff31c6174ce7bce74ca4a91c2e816dbf91e"
]

# password to unlock accounts
Password= "/home/app/workspaces/berith/node1/password_temp.txt"

# how long the test will be executed
Duration= "00:00:10"

# how many test runs will be executed
TxCount=3

# interval between transaction requests. default 10ms
TxInterval=20

# initial delay before running. default 0
InitDelay=0

# output dir to write a result file
OutputPath="/home/app/workspaces/berith/berithench-test.txt"
```  

or  

```bash
berithench execute transfer --nodes "http://localhost:8501,http://localhost:8502" \
    --chainid "1234"
    --keystore "/home/app/workspaces/berith/keystore" \
    --addresses "8676fb254279ef78c53b8a781e228ab439065786,ca7207de79e55c1a69dbc67a4a2e81dfc62c6ac4,d8a25ff31c6174ce7bce74ca4a91c2e816dbf91e" \
    --password "/home/app/workspaces/berith/node1/password_temp.txt" \
    --duration "00:00:10" --txcount 0 --txinterval 20 --initdelay 0 \
    --outputpath "/home/app/workspaces/berith/berithench-test.txt"
```

- tps :   



---  

<div id="berithkey"></div>
; TODO

---  

<div id="bootnode"></div>
; TODO  

---  

<div id="gwizard"></div>
; TODO

---  

<div id="wallet"></div>
; TODO
