package main

import (
	"encoding/json"
	"github.com/BerithFoundation/berith-chain/rpc"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
)

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "init":
		ch <- NodeMsg{
			t: "init",
			v: nil,
		}
		break
	case "callApi":
		var info map[string]interface{}
		err = json.Unmarshal(m.Payload, &info)
		if err != nil{
			payload = nil
			break
		}

		api := info["api"]
		args := info["args"]
		payload, err = callNodeApi(api, args)
		break
	}
	return
}

func callNodeApi(api interface{}, args ...interface{}) (string, error)  {
	var result json.RawMessage
	err := client.Call(&result, api.(string), args)

	var val string

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
	}

	return val, err
}







