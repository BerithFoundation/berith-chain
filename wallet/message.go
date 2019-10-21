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

// handleMessages handles messages
// js api랑 통신하는 함수
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	var info map[string]interface{}
	err = json.Unmarshal(m.Payload, &info)
	if err != nil{
		payload = nil
		return
	}
	api := info["api"]
	args := info["args"].([]interface{})
	switch m.Name {

	case "init":
		break
	case "polling":
		//ch2 <- 1
		startPolling()
		break
	case  "stopPolling" :
		ch2 <- 0
		break
	case "callApi":

		payload, err = callNodeApi(api, args...)
		break
	case "callDB":
		payload , err = callDB(api , args...)

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

	if err!=nil {
		astilog.Error(err.Error())
	}
	astilog.Debugf("Payload: %s", payload)
	return
}


// cli 관련 api 처리하는 함수
func callNodeApi(api interface{}, args ...interface{}) (string, error)  {
	var apiName = api.(string)

	var result json.RawMessage
	p := make([]interface{}, 0)
	for _, item := range args {
		if item == nil {
			break
		}
		// 트랜잭션시

		if apiName == "berith_sendTransaction" || apiName == "berith_stake" || apiName == "berith_rewardToBalance"  || apiName == "berith_rewardToStake" || apiName == "berith_stopStaking"   {
			temp := reflect.ValueOf(item).Interface()
			itemMap:= temp.(map[string]interface{})
			p = append(p, itemMap)
		} else{
			p = append(p , item)
		}
	}
	err := client.Call(&result, api.(string), p...)

	if (err!= nil) {
		return err.Error(), err
	}
	return string(result), err
}

// db 관련 api 처리하는 함수
func callDB ( api interface{}, args... interface{}) ( interface{}, error){

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

	case "selectContact" :
		contact := make(walletdb.Contact,0)
		err := WalletDB.Select([]byte("c"+acc), &contact)

		if err != nil {
			return nil, err
		}
		return contact, nil

		break
	case "selectMember":
		var member walletdb.Member

		err := WalletDB.Select([]byte("qwer5910"), &member)
		if err != nil {
			return nil, err
		}
		return member, nil
		break
	case "insertContact":
		contact := make(walletdb.Contact, 0)
		WalletDB.Select([]byte("c"+acc), &contact)
		contact[common.HexToAddress(key[0])] = key[1]

		err := WalletDB.Insert([]byte("c"+acc) , contact)
		if err != nil {
			return nil, err
		}
		return  contact , nil
		break
	case "updateContact":
		contact := make(walletdb.Contact, 0)
		WalletDB.Select([]byte("c"+acc), &contact)
		delete(contact, common.HexToAddress(key[0]))
		err := WalletDB.Insert([]byte("c"+acc) , contact)
		if err != nil {
			return nil , err
		}
		return contact, nil
		break
	case "restoreMember":
		var mem walletdb.Member
		err := WalletDB.Select([]byte(key[1]), &mem)
		if err == nil {
			return "err" , err
		}
		member := walletdb.Member{
			Address: common.HexToAddress(key[0]),
			ID : key[1],
			Password: key[2],
			PrivateKey: key[3],
		}
		err = WalletDB.Insert([]byte(key[1]) , member)
		if err != nil {
			return nil , err
		}
		return member, nil
		break
	case "updateMember":
		var mem walletdb.Member
		err := WalletDB.Select([]byte(key[0]), &mem)
		if err != nil {
			return "err" , err
		}
		mem.Password = key[1]
		err = WalletDB.Insert([]byte(key[0]) , mem)
		if err != nil {
			return nil , err
		}
		return mem, nil
		break
	case "selectTxInfo":
		txMaster := make(walletdb.TxHistoryMaster , 0)
		err := WalletDB.Select([]byte("t"+acc) , &txMaster)
		if err != nil {
			return nil , err
		}
		txDetails := make([]walletdb.TxHistory,0)
		for _ , val := range txMaster {
			var txDetail walletdb.TxHistory
			err2 := WalletDB.Select([]byte(val) , &txDetail)
			if err2 != nil {
				return nil , err
			}
			txDetails=append(txDetails,txDetail)
		}
		return txDetails , nil
		break
	case "insertTxInfo":
		var tempTxInfo walletdb.TxHistory
		err := WalletDB.Select([]byte(key[0]), &tempTxInfo)
		if err == nil{
			return "err" ,err
		}
		txMaster := make(walletdb.TxHistoryMaster , 0)
		WalletDB.Select([]byte("t"+acc), &txMaster)
		txMaster[key[0]] = key[0]
		err = WalletDB.Insert([]byte("t"+acc), txMaster)
		if err != nil {
			return nil ,err
		}
		txinfo := walletdb.TxHistory{
			TxAddress: common.HexToAddress(key[1]),
			TxType: key[2],
			TxAmount: key[3],
			Txtime: time.Now().Format("2006-01-02 15:04:05"),
			Hash: common.HexToHash(key[4]),
			GasLimit: key[5],
			GasPrice: key[6],
		}
		err = WalletDB.Insert([]byte(key[0]), txinfo)
		if err != nil {
			return nil , err
		}
		return txinfo, nil

		break
	case "insertMember":
		var mem walletdb.Member
		err := WalletDB.Select([]byte(key[0]), &mem)
		if err == nil {
			return "err" , err
		}
		newAcc ,err := callNodeApi("personal_newAccount", key[1])
		newAcc = strings.ReplaceAll(newAcc , "\"","")
		privateKey , err := callNodeApi("personal_privateKey",newAcc , key[1] )
		privateKey = strings.ReplaceAll(privateKey , "\"","")
		member := walletdb.Member{

			Address: common.HexToAddress(newAcc),
			ID : key[0],

			Password: key[1],
			PrivateKey: privateKey,
		}

		err = WalletDB.Insert([]byte(key[0]) , member)

		if err != nil {
			return nil, err
		}
		return member, nil
		break

	case "checkLogin":
		var member walletdb.Member

		err := WalletDB.Select([]byte(key[0]), &member)
		if err != nil {
			return nil, err
		}
		return member, nil
		break
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
	log.Info("Found keystore dir: ", dir)
	password := args[0].(string)
	targetPath := dir + string(os.PathSeparator) + tempFileName
	er := ZipSecure(dir, targetPath, password)
	if er != nil {
		return nil, er
	}
	log.Info("Successfully created temp file, " + tempFileName + ", at: " + dir)

	zippedFile, err := os.Open(targetPath)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(zippedFile)
	if err != nil {
		fmt.Println(err)
	}

	zippedFile.Close()
	os.Remove(targetPath)
	log.Info("Removed temp file, " + tempFileName + ", from: " + dir)

	return body, nil
}

// 개인키 삽입 함수
func importKeystore(args []interface{}) (error)  {


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
