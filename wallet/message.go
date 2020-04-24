package main

import (
	"encoding/json"
	"fmt"
	"github.com/BerithFoundation/berith-chain/common"
	"github.com/BerithFoundation/berith-chain/log"
	"github.com/BerithFoundation/berith-chain/wallet/database"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"
)

var (
	WalletDB2 *walletdb.WalletDB
)

// handleMessages handles messages
// js api랑 통신하는 함수
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	var info map[string]interface{}
	err = json.Unmarshal(m.Payload, &info)
	if err != nil {
		payload = nil
		return
	}
	api := info["api"]
	args := info["args"].([]interface{})
	switch m.Name {

	case "init":
		break
	case "polling":
		//ch2 <- true
		startPolling()
		break
	case "stopPolling":
		//ch2 <- false
		fmt.Print("logout!!!!")

		break
	case "callApi":
		payload, err = callNodeApi(api, args...)
		break
	case "callDB":
		payload, err = callDB(api, args...)
		break
	case "exportKeystore":
		args := info["args"].([]interface{})
		payload, err = exportKeystore(args)
		break

	case "importKeystore":
		args := info["args"].([]interface{})
		err = importKeystore(args)
		payload = nil
		break
	}

	if err != nil {
		astilog.Error(err.Error())
	}
	astilog.Debugf("Payload: %s", payload)
	return
}

// cli 관련 api 처리하는 함수
func callNodeApi(api interface{}, args ...interface{}) (string, error) {
	var apiName = api.(string)

	var result json.RawMessage
	p := make([]interface{}, 0)
	for _, item := range args {
		if item == nil {
			break
		}
		// 트랜잭션시

		if apiName == "berith_sendTransaction" || apiName == "berith_stake" || apiName == "berith_rewardToBalance" || apiName == "berith_rewardToStake" || apiName == "berith_stopStaking" {
			temp := reflect.ValueOf(item).Interface()
			itemMap := temp.(map[string]interface{})
			p = append(p, itemMap)
		} else {
			p = append(p, item)
		}
	}
	err := client.Call(&result, api.(string), p...)

	if err != nil {
		return err.Error(), err
	}
	return string(result), err
}

// db 관련 api 처리하는 함수
func callDB(api interface{}, args ...interface{}) (interface{}, error) {

	key := make([]string, 0)
	for _, item := range args {
		key = append(key, item.(string))
	}
	acc, err := callNodeApi("berith_coinbase", nil)
	// acc = strings.ReplaceAll(acc , "\"","")
	acc = strings.Replace(acc, "\"", "", -1)
	if err != nil {
		astilog.Error(errors.Wrap(err, "insert error"))
	}
	switch api.(string) {

	case "selectContact":
		contact := make(walletdb.Contact, 0)
		err := WalletDB.Select([]byte("c"+acc), &contact)

		if err != nil {
			return nil, err
		}
		return contact, nil
	case "selectMember":
		var member walletdb.Member

		err := WalletDB.Select([]byte("qwer5910"), &member)
		if err != nil {
			return nil, err
		}
		return member, nil
	case "insertContact":
		contact := make(walletdb.Contact, 0)
		WalletDB.Select([]byte("c"+acc), &contact)
		contact[common.HexToAddress(key[0])] = key[1]

		err := WalletDB.Insert([]byte("c"+acc), contact)
		if err != nil {
			return nil, err
		}
		return contact, nil
	case "updateContact":
		contact := make(walletdb.Contact, 0)
		WalletDB.Select([]byte("c"+acc), &contact)
		delete(contact, common.HexToAddress(key[0]))
		err := WalletDB.Insert([]byte("c"+acc), contact)
		if err != nil {
			return nil, err
		}
		return contact, nil
	case "restoreMember":
		var mem walletdb.Member
		err := WalletDB.Select([]byte(key[1]), &mem)
		if err == nil {
			mem = walletdb.Member{
				Address:    common.HexToAddress(key[0]),
				ID:         key[1],
				Password:   key[2],
				PrivateKey: key[3],
			}
			return mem, nil
		}
		member := walletdb.Member{
			Address:    common.HexToAddress(key[0]),
			ID:         key[1],
			Password:   key[2],
			PrivateKey: key[3],
		}
		err = WalletDB.Insert([]byte(key[1]), member)
		if err != nil {
			return nil, err
		}
		return member, nil
	case "updateMember":
		var mem walletdb.Member
		err := WalletDB.Select([]byte(key[0]), &mem)
		if err != nil {
			return "err", err
		}
		mem.Password = key[1]
		err = WalletDB.Insert([]byte(key[0]), mem)
		if err != nil {
			return nil, err
		}
		return mem, nil
	case "selectTxInfo":
		txMaster := make(walletdb.TxHistoryMaster, 0)
		err := WalletDB.Select([]byte("t"+acc), &txMaster)
		if err != nil {
			return nil, err
		}
		txDetails := make([]walletdb.TxHistory, 0)
		for _, val := range txMaster {
			var txDetail walletdb.TxHistory
			err2 := WalletDB.Select([]byte(val), &txDetail)
			if err2 != nil {
				return nil, err
			}
			txDetails = append(txDetails, txDetail)
		}
		return txDetails, nil
	case "insertTxInfo":
		var tempTxInfo walletdb.TxHistory
		err := WalletDB.Select([]byte(key[0]), &tempTxInfo)
		if err == nil {
			return "err", err
		}
		txMaster := make(walletdb.TxHistoryMaster, 0)
		WalletDB.Select([]byte("t"+acc), &txMaster)
		txMaster[key[0]] = key[0]
		err = WalletDB.Insert([]byte("t"+acc), txMaster)
		if err != nil {
			return nil, err
		}
		txinfo := walletdb.TxHistory{
			TxBlockNumber: key[0],
			TxAddress:     common.HexToAddress(key[1]),
			TxType:        key[2],
			TxAmount:      key[3],
			Txtime:        time.Now().Format("2006-01-02 15:04:05"),
			Hash:          common.HexToHash(key[4]),
			GasLimit:      key[5],
			GasPrice:      key[6],
			GasUsed:       key[7],
		}
		err = WalletDB.Insert([]byte(key[0]), txinfo)
		if err != nil {
			return nil, err
		}
		return txinfo, nil
	case "insertMember":
		var mem walletdb.Member
		err := WalletDB.Select([]byte(key[0]), &mem)
		if err == nil {
			return "err", err
		}
		newAcc, err := callNodeApi("personal_newAccount", key[1])
		newAcc = strings.Replace(newAcc, "\"", "", -1)
		privateKey, err := callNodeApi("personal_privateKey", newAcc, key[1])
		privateKey = strings.Replace(privateKey, "\"", "", -1)
		member := walletdb.Member{
			Address:    common.HexToAddress(newAcc),
			ID:         key[0],
			Password:   key[1],
			PrivateKey: privateKey,
		}

		err = WalletDB.Insert([]byte(key[0]), member)

		if err != nil {
			return nil, err
		}
		return member, nil
	case "checkLogin":
		var member walletdb.Member

		err := WalletDB.Select([]byte(key[0]), &member)
		if err != nil {
			return nil, err
		}
		return member, nil
	}

	return nil, nil
}

// 개인키 내보내기 함수

func exportKeystore(args []interface{}) (interface{}, error) {
	tempFileName := "keystore.zip"
	dir, err := stack.FetchKeystoreDir()
	if err != nil {
		return nil, err
	}
	tmmp, _ := ioutil.TempDir(dir, "tmp")

	// export 하는 계정에 대한 keysotre 파일만 exportTemp 폴더로 이동하는 부분
	keyAccount := strings.Replace(args[0].(string), "Bx", "", -1)
	tempFile, _ := ioutil.ReadDir(dir)
	for i, value := range tempFile {
		value.Name()
		if strings.Contains(value.Name(), keyAccount) {
			dir3 := tmmp + "/" + value.Name()
			os.Create(dir3)
			bytes, err := ioutil.ReadFile(dir + "/" + value.Name())
			if err != nil {
				panic(err)
			}
			//파일 쓰기
			err = ioutil.WriteFile(dir3, bytes, 0)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("name[", i, "]  :", value.Name())
	}
	// 끝
	// export 하는 계정에 관련된 db 정보만 따로 추출하는 부분
	WalletDB2, _ = walletdb.NewWalletDB(tmmp + "/test.ldb")
	var mem walletdb.Member
	var txInfo walletdb.TxHistory
	contact := make(walletdb.Contact, 0)
	txMaster := make(walletdb.TxHistoryMaster, 0)
	WalletDB.Select([]byte("c"+args[0].(string)), &contact)
	WalletDB.Select([]byte("t"+args[0].(string)), &txMaster)
	for key, _ := range txMaster {
		WalletDB.Select([]byte(key), &txInfo)
		WalletDB2.Insert([]byte(key), txInfo)
	}
	err = WalletDB.Select([]byte(args[2].(string)), &mem)
	if err != nil {
		return nil, err
	}
	err = WalletDB2.Insert([]byte(args[2].(string)), mem)
	if err != nil {
		return nil, err
	}
	err = WalletDB2.Insert([]byte("c"+args[0].(string)), contact)
	if err != nil {
		return nil, err
	}
	err = WalletDB2.Insert([]byte("t"+args[0].(string)), txMaster)
	if err != nil {
		return nil, err
	} // 끝
	password := args[1].(string)
	targetPath := dir + string(os.PathSeparator) + tempFileName
	er := ZipSecure(tmmp, targetPath, password)
	if er != nil {
		return nil, er
	}
	zippedFile, err := os.Open(targetPath)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(zippedFile)
	if err != nil {
		fmt.Println(err)
	}

	zippedFile.Close()
	fmt.Println("targetPath :: ", targetPath)
	fmt.Println("dir2 :: ", dir+string(os.PathSeparator)+"exportTemp")
	os.Remove(targetPath)
	WalletDB2.CloseDB()
	result := os.RemoveAll(tmmp)
	fmt.Println("result :: ", result)
	return body, nil
}

// 개인키 삽입 함수
func importKeystore(args []interface{}) error {

	dir, err := stack.FetchKeystoreDir()
	if err != nil {
		return err
	}
	log.Info("Found keystore dir: ", dir)

	inputFilePath := args[0].(string)
	password := args[1].(string)
	log.Debug("Input keystore file path: ", dir)

	er := UnzipSecure(inputFilePath, dir, password)
	if er != nil {
		return er
	}
	log.Info("Successfully imported keystore folder")
	return err
}
