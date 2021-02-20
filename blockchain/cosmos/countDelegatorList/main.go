package main

import (
	"fmt"
	"os/exec"
	"os"
	"bufio"
//	"strings"
	"strconv"
	"encoding/json"
	"encoding/csv"
	utils "github.com/node-a-team/Cosmos-IE/utils"
)

var (
        ValidatorsList []string
)


type validators struct {
//	Height string `json:"height"`
	Validators []validator
}
type validator struct {
	OperAddr        string `json:"operator_address"`
	ConsPubKey      string `json:"consensus_pubkey"`
	Jailed          bool   `json:"jailed"`
	Status          int    `json:"status"`
	Tokens          string `json:"tokens"`
	DelegatorShares string `json:"delegator_shares"`
	Description     struct {
		Moniker  string `json:"moniker"`
		Identity string `json:"identity"`
		Website  string `json:"website"`
		Details  string `json:"details"`
	}
	UnbondingHeight string `json:"unbonding_height"`
	UnbondingTime   string `json:"unbonding_time"`
	Commission      struct {
		Commission_rates struct {
			Rate          string `json:"rate"`
			Max_rate       string `json:"max_rate"`
			Max_change_rate string `json:"max_change_rate"`
		}
		UpdateTime string `json:"update_time"`
	}
	MinSelfDelegation string `json:"min_self_delegation"`
}

type delegators struct {
	Delegation_responses []struct {
		Delegation struct {
			Delegator_address string
		}
		Balance struct {
			Denom string
			Amount string
		}
	}
}

func main() {


	var v validators
	var d delegators

	delegatorsList := make(map[string]float64)


	logDir := "/data/cosmos/github.com/test/"
	csvDelegator, _ := os.OpenFile(logDir +"delegatorList.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	csvWriter := csv.NewWriter(bufio.NewWriter(csvDelegator))

	defer csvWriter.Flush()

	res, _ := runRESTCommand("/cosmos/staking/v1beta1/validators?pagination.limit=300")
//	res, _ := runRESTCommand("/cosmos/staking/v1beta1/validators?pagination.limit=300&status=BOND_STATUS_BONDED")
	json.Unmarshal(res, &v)

	fmt.Printf("Validator Check:")
	count := 0
	for _, value := range  v.Validators {

		count++
//		fmt.Println("i, value] ", i+1, value.OperAddr)
//		ValidatorsList = append(ValidatorsList, value.OperAddr)

		res, _ = runRESTCommand("/cosmos/staking/v1beta1/validators/" +value.OperAddr +"/delegations?pagination.limit=10000")
	        json.Unmarshal(res, &d)

		fmt.Printf("\n %d] %s, %d\n", count, value.OperAddr, len(d.Delegation_responses))

		for _, dValue := range d.Delegation_responses {

			fmt.Printf("- %s: %f ", dValue.Delegation.Delegator_address, delegatorsList[dValue.Delegation.Delegator_address])
			delegatorsList[dValue.Delegation.Delegator_address] = delegatorsList[dValue.Delegation.Delegator_address] +(utils.StringToFloat64(dValue.Balance.Amount)/1000000)
			fmt.Printf("-> %f\n", delegatorsList[dValue.Delegation.Delegator_address])

		}
	}

	fmt.Printf("ValidatorsList count: %d\n", count)
	fmt.Printf("DelegatorsList count: %d\n", len(delegatorsList)+1)


	csvWriter.Write([]string{"Delegator", "Amount" })

	fmt.Printf("Delegator Check: ")

	count = 0
	for key, val := range delegatorsList {

		count++
		fmt.Println("", count, key, val)

		csvWriter.Write([]string{key, strconv.FormatFloat(val, 'f', -1, 64)})
	}


}


func runRESTCommand(str string) ([]uint8, error) {
	cmd := "curl -s -XGET localhost:1317"  +str +" -H \"accept:application/json\""
        out, err := exec.Command("/bin/bash", "-c", cmd).Output()

        return out, err
}
