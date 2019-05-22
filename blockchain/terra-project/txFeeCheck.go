package main

import (
	"fmt"
	"crypto/sha256"
	"encoding/base64"

	"encoding/json"
	"os/exec"
	"strconv"
)

var (
	publicNode string
	chainId    string = "soju-0007"
	txDenom    string = "mluna"
)

type TxStruct struct {
	Height     string   `json":"height"`
	Code       int      `json":"code"`
	Gas_wanted string   `json":"gas_wanted"`
	Gas_used   string   `json":"gas_used"`
	Tags       []TxTags `json":"tags"`
	Tx         Tx       `json":"tx"`
}

type TxTags struct {
	Key   string `json":"key"`
	Value string `json":"value"`
}

type Tx struct {
	Type  string  `json":"type"`
	Value TxValue `json":"value"`
}

type TxValue struct {
	Msg        []TxMsg        `json":"msg"`
	Fee        TxFee          `json":"fee"`
	Signatures []TxSignatures `json":"signatures"`
	Memo       string         `json":"memo"`
}

type TxFee struct {
	Amount []TxAmount `json":"amount"`
	Gas    string     `json":"gas"`
}

type TxSignatures struct {
	Pub_key   Pub_key `json":"pub_key"`
	Signature string  `json":"signature"`
}

type TxMsg struct {
	Type string `json":"type"`
}

type TxAmount struct {
	Denom  string `json":"denom"`
	Amount string `json":"amount"`
}

type Pub_key struct {
	Type  string `json":"type"`
	value string `json":"value"`
}

func TxData(txs string) {

	switch TxCheck(txs) {

	case 0:
		fmt.Println("\n------------------txInfo-------------------")
		fmt.Printf("\ntx: %s\n\n", txs)
		fmt.Println("\n> tx Fee("+txDenom+"): ", float64(TxFeeSearch(TxSearch(txs, 0))))
		fmt.Println("\n")

	case 1:

		errCode, errStatus := TxErrorCheck(TxSearch(TxsConvert(txs), 1))

		if errStatus == true {
			fmt.Println("\n------------------txInfo-------------------")
			fmt.Printf("\ntx: %s\n\n", TxsConvert(txs))
			fmt.Println("\n> errrorStatus: ", errStatus)
			fmt.Println("> errorCode: ", errCode)

		} else {

			fmt.Println("\n------------------txInfo-------------------")
			fmt.Printf("\ntx: %s\n\n", TxsConvert(txs))
			fmt.Println("\n> tx Fee("+txDenom+"): ", float64(TxFeeSearch(TxSearch(TxsConvert(txs), 0))))
			fmt.Println("\n")
		}
	}

}

// tx convert 필요한지 확인
func TxCheck(txs string) int {

	// 변환이 필요없는 tx -> return 0, 변환이 필요한 tx -> return 1
	if len(txs) == 64 {
		return 0
	} else {
		return 1
	}

}

// block에 담긴 txs를 조회 가능하도록 변경
func TxsConvert(txs string) string {

	// txs -> 검색 가능한 tx hash로 변환
	// base64로decode-> sha256 encode -> hax encode
	txsDecode, _ := base64.StdEncoding.DecodeString(txs)

	txsSha256 := sha256.New()
	txsSha256.Write(txsDecode)

	r := txsSha256.Sum(nil)

	return fmt.Sprintf("%x", r)

}

//func TxSearch(tx string ) t.TxStruct {
func TxSearch(tx string, errCheck int) TxStruct {
	//        var txData t.TxStruct
	var txData TxStruct

	cmd := "terracli q tx " + tx + " --chain-id=" + chainId + " --output=json --node=" + publicNode
	out, _ := exec.Command("/bin/bash", "-c", cmd).Output()
	json.Unmarshal(out, &txData)

	if errCheck != 1 {
		fmt.Println(string(out))
		fmt.Println("------------------txInfo-------------------")
	}

	return txData

}

// tx에 대한 Fee 조회
//func TxFeeSearch(txData t.TxStruct) int64 {
func TxFeeSearch(txData TxStruct) int64 {

	var txFee int64

	if len(txData.Tx.Value.Fee.Amount) == 0 {
		txFee = 0
	} else {
		for _, value := range txData.Tx.Value.Fee.Amount {
			fmt.Println("TxFeeSearch: ", value)
			txFeeInt, _ := strconv.ParseInt(value.Amount, 10, 64)
			txFee = txFeeInt
		}

	}
	return txFee
}

//func TxErrorCheck(txData t.TxStruct) (int, bool) {
func TxErrorCheck(txData TxStruct) (int, bool) {

	errCode := txData.Code
	errStatus := false

	if errCode != 0 {
		errStatus = true
	}

	return errCode, errStatus
}

func main() {

	var txs string

	fmt.Println("\n[ Info ]")
        fmt.Println("- Transaction fee check program_Terra-project terrad(\"soju-0007\")")
        fmt.Println("- Program Stop: Enter \"quit\" or \"q\"")
        fmt.Println("- Made by Node A-Team_J ")

	fmt.Println("\n[Initial Setting ]")
	fmt.Printf("- input Peer(192.168.0.1:26657): ")
	fmt.Scanf("%s", &publicNode)


	for {

		fmt.Println("\n[ Start ]")
		fmt.Printf("- input txs: ")
		fmt.Scanf("%s", &txs)

		if txs == "quit" || txs == "q" {
			break
		}
		TxData(txs)
	}
	//	TxData("uQHwYl3uCjIxh8QlCgRtZ2JwEhA4MDI4Mzg2OTIxMTc2NjI3GhTWrYggimP0WInC/NX0I2LrezuAOxITCg0KBW1sdW5hEgQzMDAwEMCaDBpqCibrWumHIQMlenWbBolsL05FGny1CG+Mlj1zx67ziYEpLbll1HMWPhJAGSvdJw5g1oToT0LBhu8xp4LwD8OxtQpH6GSNNwWV5jBodvQfZbrF3X2ch+oXgEJDQytbKtSfaUQaFPR6i/Gmow==")

	//	TxData("ugHwYl3uCkKSHS5OChSbFiDqVqPB6NfhZAyT1bEORIPQJxIUXyCQoZIYfgQpO/O7glYfFLcxQdgaEAoFbWx1bmESBzg5OTkwMDASBBDAmgwaagom61rphyEDQdPRa3AAVs22SbVPsBDrDbEmlWMMpGHJocrgPMlJ9bESQEPBeU3X8E3uParOthAwtlheCbI8Y/NCw90Pvr9+2NW5E3JdQD2ZBd9PJvld/6yD3SuP9AtwBscNwFWxTav6lxk=")

}
