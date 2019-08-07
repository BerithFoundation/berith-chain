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
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	var info map[string]interface{}
	err = json.Unmarshal(m.Payload, &info)
	if err != nil{
		payload = nil
		return
	}
	api := info["api"]
	args := info["args"].([]interface{})

	astilog.Debugf("Message type: %s, Method: %s", m.Name, api.(string))

	switch m.Name {
	case "init":


		/*ch <- NodeMsg{
			t: "init",
			v: nil,
			stack: nil,
		}*/
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

func callNodeApi(api interface{}, args ...interface{}) (string, error)  {
	var apiName = api.(string)
	var result json.RawMessage
	p := make([]interface{}, 0)
	for _, item := range args{
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

	/*var val string
	switch err := err.(type) {
	case nil:
		if result == nil {

		} else {
			val = string(result)
			return val, err
		}
	case rpc.Error:
		return val, err
	default:
		return val, err
	}*/

	//return val, err
}

func callDB ( api interface{}, args... interface{}) ( interface{}, error){
	key := make([]string, 0)
	for _, item := range args{
		 key  = append(key, item.(string))
	}
	acc ,err := callNodeApi("berith_coinbase", nil)
	acc = strings.ReplaceAll(acc , "\"","")
	if err != nil {
		astilog.Error(errors.Wrap(err, "insert error"))
	}
	switch api.(string) {
	case "selectContact" :
		contact := make(walletdb.Contact,0)
		err := WalletDB.Select([]byte(acc), &contact)
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
		WalletDB.Select([]byte(acc), &contact)
		contact[common.HexToAddress(key[0])] = key[1]
		err := WalletDB.Insert([]byte(acc) , contact)
		if err != nil {
			return nil, err
		}
		return  contact , nil
		break
	case "updateContact":
		contact := make(walletdb.Contact, 0)
		WalletDB.Select([]byte(acc), &contact)
		delete(contact, common.HexToAddress(key[0]))
		err := WalletDB.Insert([]byte(acc) , contact)
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
			return nil , err
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

	return nil ,nil
}

func exportKeystore(args []interface{}) (interface{}, error) {
	tempFileName:= "keystore.zip"

	dir, err := stack.FetchKeystoreDir()
	if (err!=nil) {
		return nil, err
	}
	log.Info("Found keystore dir: ", dir)
	password:= args[0].(string)
	targetPath := dir + string(os.PathSeparator) + tempFileName
	er := ZipSecure(dir,targetPath,password)
	if er != nil {
		return nil,er
	}
	log.Info("Successfully created temp file, "+tempFileName+", at: " +dir)

	zippedFile, err := os.Open(targetPath)
	if err != nil {
		return nil,err
	}

	body, err := ioutil.ReadAll(zippedFile)
	if err != nil {
		fmt.Println(err)
	}

	zippedFile.Close()
	os.Remove(targetPath)
	log.Info("Removed temp file, "+tempFileName+", from: " +dir)

	return body, nil
}

func importKeystore(args []interface{}) (error)  {

	dir, err := stack.FetchKeystoreDir()
	if (err!=nil) {
		return err
	}
	log.Info("Found keystore dir: ", dir)

	inputFilePath:= args[0].(string)
	password:= args[1].(string)
	log.Debug("Input keystore file path: ", dir)

	er := UnzipSecure(inputFilePath,dir,password)
	if er != nil {
		return er
	}
	log.Info("Successfully imported keystore folder")
	return err
}

