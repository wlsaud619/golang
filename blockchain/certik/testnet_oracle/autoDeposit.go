package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
//	"strconv"
	"go.uber.org/zap"

	utils "github.com/node-a-team/Cosmos-IE/utils"
)

var (
	keyName				string = "ATEAM"
	operator_addr			string  = "certik1hgf3anv20wg7n9g60tprfw4lfgyk8k05f8t86r"
	check_balances, check_rewards	float64 = 1000.0, 1000.0
)

type account struct {
	Value struct {
		Address string
		Coins   []coins
	}
}

type operator struct {
	Address             string
	Proposer            string
	Collateral          []coins
	Accumulated_rewards []coins
	Name                string
}

type coins struct {
	Denom  string
	Amount string
}

func main() {

	log,_ := zap.NewDevelopment()
        defer log.Sync()

	fmt.Printf("\033[31m*****************************************************************\033[0m\n")
	fmt.Printf("\033[31mDeposit condition:\033[0m \033[32mbalances > %.2f \033[0m or \033[33mrewards > %.2f \033[0m\n", check_balances, check_rewards)
	fmt.Printf("\033[31m*****************************************************************\033[0m\n")
	for {
		balances := queryBalances()
		rewards := queryRewards()

		if balances > check_balances || rewards > check_rewards {
			log.Info("Balance Check", zap.Bool("Condition", true), zap.String("Notice", "The condition is satisfie"),)
			log.Info("Balance Check", zap.String("Current balances/rewards", fmt.Sprintf("%.2f/%.2f CTK", balances, rewards)),)
			deposit(balances, log)
		} else {
			log.Info("Balance Check", zap.Bool("Condition", false), zap.String("Notice", "The condition is not met"),)
                        log.Info("Balance Check", zap.String("Current balances/rewards", fmt.Sprintf("%.2f/%.2f CTK", balances, rewards)),)
		}

		time.Sleep(5 * time.Second)
	}

}

func queryBalances() float64 {
	var a account

	res, _ := SHELL("certikcli q account " + operator_addr + " -o json ")
        json.Unmarshal(res, &a)

	if len(a.Value.Coins) == 0 {
                return 0.0
        } else {
                return utils.StringToFloat64(a.Value.Coins[0].Amount) / 1000000
	}
}

func queryRewards() float64 {
        var o operator

	res, _ := SHELL("certikcli q oracle operator " + operator_addr + " -o json ")
        json.Unmarshal(res, &o)


	if len(o.Accumulated_rewards) == 0 {
		return 0.0
	} else {
		return utils.StringToFloat64(o.Accumulated_rewards[0].Amount) / 1000000
	}
}

func deposit(old_balances float64, log *zap.Logger) {

	// Claim
	log.Info("Deposit", zap.String("Notice", "Claim rewards"),)
	SHELL("certikcli tx oracle claim-reward " + operator_addr + " --from " +keyName +" -y")

	time.Sleep(10 * time.Second)
	new_balances := queryBalances()
/*
	for {
		if new_balances > old_balances {
			break
		} else {

			time.Sleep(1 * time.Second)
			new_balances = queryBalances()

			fmt.Println("old, new: ", old_balances, new_balances)
		}
	}
*/
	// Deposit
	new_balances = (new_balances - 0.050000) * 1000000
//	new_balances = (110 - 0.050000) * 1000000
	log.Info("Deposit", zap.String("Deposit amount", fmt.Sprintf("%.2f CTK", new_balances/1000000)),)

	SHELL("certikcli tx oracle deposit-collateral " +operator_addr +" " +fmt.Sprint(int(new_balances)) +"uctk" +" --from " +keyName +" --fees 10000uctk -y")
//	fmt.Println("certikcli tx oracle deposit-collateral " +operator_addr +" " +fmt.Sprint(int(new_balances)) +"uctk" +" --from " +keyName +" --fees 10000uctk -y")
	time.Sleep(10 * time.Second)
}

func SHELL(cmd string) ([]uint8, error) {
	out, err := exec.Command("/bin/bash", "-c", cmd).Output()
	fmt.Println(cmd)
	return out, err

}
