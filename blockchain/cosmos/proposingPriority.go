package main

import (
      "time"
       "fmt"
        "os/exec"
        "encoding/json"
        "sort"
//        "sync"
        "strconv"
)

type ProposerStatusStruct struct {
	Block_height    string			`json:"block_height"`
	Validators      []ProposerValidators	`json:"validators"`

}

type ProposerValidators struct {
	Address			string	`json:"address"`
        Pub_key			string	`json:"pub_key"`
        Proposer_priority	string	`json:"proposer_priority"`
        Voting_power		string	`json:"voting_power"`
}



var (

	publicNode string
	valPubAddr string

        proposerStatus ProposerStatusStruct

        // proposer[0] : priority 값
        // proposer[1] : priority 순위
        proposer map[string][]int = make(map[string][]int)

)

func ProposerStatus() int {

//	chainId := "cosmoshub-2"
//        cmd := "gaiacli query tendermint-validator-set  --chain-id=" +chainId +" --output=json --trust-node --node=" +publicNode +":26657"
	cmd := "curl -X GET \"http://" +publicNode +":1317/validatorsets/latest\" -H \"accept: application/json\""
        out, _ := exec.Command("/bin/bash", "-c", cmd).Output()
        json.Unmarshal(out, &proposerStatus)

        for _, value := range proposerStatus.Validators {
                priorityInt, _ := strconv.Atoi(value.Proposer_priority)
                proposer[value.Pub_key] = []int{priorityInt, 0}
        }

        Sort()

//	fmt.Println(proposer)

        return proposer[valPubAddr][1]
}

func Sort() {

        keys := []string{}

        // key(moniker)들을 keys 배열에 추가
        for key := range proposer {
                keys = append(keys, key)
        }

        // keys를 totalStake 기준으로 정렬
        sort.Slice(keys, func(i, j int) bool {
                return proposer[keys[i]][0] > proposer[keys[j]][0]
        })

        for i, key := range keys {

		/*
                if key == valPubAddr {
			// fmt.Println("true")
		}
		*/

                proposer[key][1] = i + 1

                // total rank
//                fmt.Printf("%s, %d, %d\n", key, proposer[key][0], proposer[key][1])
        }
}

func main() {


	fmt.Println("\n[ Info ]")
        fmt.Println("- Proposer Priority Checker_Cosmos(\"cosmoshub-2\")")
	fmt.Println("- Rest-Server: cosmos-main.peer.nodeateam.kr")
        fmt.Println("- Program Stop: Enter \"Ctrl\" + \"C\"")
        fmt.Println("- Made by Node A-Team_J ")

	fmt.Println("\n[Initial Setting ]")
//	fmt.Printf("- input Rest-Server IP(192.168.0.1): ")
//	fmt.Scanf("%s", &publicNode)
	publicNode = "cosmos-main.peer.nodeateam.kr"

	fmt.Println("- input Validator Public Address(gaiad tendermint show-validator -> cosmosvalconspub)")
	fmt.Printf(": ")
        fmt.Scanf("%s", &valPubAddr)

	fmt.Println("\n[ Start ]")
	for {
		fmt.Println(ProposerStatus())
                time.Sleep(1*time.Second)
	}
}
