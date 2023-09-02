//go:build !wasm

package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/alis-is/tezpay/common"
	"github.com/alis-is/tezpay/constants"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer <destination> <amount tez>",
	Short: "transfers tez to specified address",
	Long:  "transfers tez to specified address from payout wallet",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, signer, transactor := assertRunWithResult(loadConfigurationEnginesExtensions, EXIT_CONFIGURATION_LOAD_FAILURE).Unwrap()
		mutez, _ := cmd.Flags().GetBool(MUTEZ_FLAG)

		if len(args)%2 != 0 {
			log.Error("invalid number of arguments (expects pairs of destination and amount)")
			panic(PanicStatus{
				ExitCode: EXIT_IVNALID_ARGS,
				Error:    fmt.Errorf("invalid number of arguments (expects pairs of destination and amount)"),
			})
		}
		total := int64(0)

		destinations := make([]string, 0)

		op := codec.NewOp().WithSource(signer.GetPKH())
		op.WithTTL(constants.MAX_OPERATION_TTL)
		for i := 0; i < len(args); i += 2 {
			destination, err := tezos.ParseAddress(args[i])
			if err != nil {
				log.Errorf("invalid destination address '%s' - '%s'", args[i], err.Error())
				panic(PanicStatus{
					ExitCode: EXIT_IVNALID_ARGS,
					Error:    fmt.Errorf("invalid destination address '%s' - '%s'", args[i], err.Error()),
				})
			}

			amount, err := strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				log.Errorf("invalid amount '%s' - '%s'", args[i+1], err.Error())
				panic(PanicStatus{
					ExitCode: EXIT_IVNALID_ARGS,
					Error:    fmt.Errorf("invalid amount '%s' - '%s'", args[i+1], err.Error()),
				})
			}
			if !mutez {
				amount *= constants.MUTEZ_FACTOR
			}

			mutez := int64(math.Floor(amount))
			total += mutez
			destinations = append(destinations, destination.String())
			op.WithTransfer(destination, mutez)
		}

		if err := requireConfirmation(fmt.Sprintf("do you really want to transfer %s to %s", common.MutezToTezS(total), strings.Join(destinations, ", "))); err != nil {
			panic(PanicStatus{
				ExitCode: EXIT_OPERTION_CANCELED,
				Error:    fmt.Errorf("operation canceled"),
			})
		}
		err := transactor.Complete(op, signer.GetKey())
		if err != nil {
			log.Errorf("failed to complete tx - %s", err.Error())
			panic(PanicStatus{
				ExitCode: EXIT_OPERTION_FAILED,
				Error:    fmt.Errorf("failed to complete tx - %s", err.Error()),
			})
		}

		err = signer.Sign(op)
		if err != nil {
			log.Errorf("failed to sign tx - %s", err.Error())
			panic(PanicStatus{
				ExitCode: EXIT_OPERTION_FAILED,
				Error:    fmt.Errorf("failed to sign tx - %s", err.Error()),
			})
		}

		log.Infof("transfering tez... waiting for %d confirmations", constants.DEFAULT_REQUIRED_CONFIRMATIONS)

		dispatchResult, err := transactor.Dispatch(op, &common.DispatchOptions{
			Confirmations: constants.DEFAULT_REQUIRED_CONFIRMATIONS,
			TTL:           rpc.DefaultOptions.TTL,
		})
		if err != nil {
			log.Errorf("failed to dispatch tx - %s", err.Error())
			panic(PanicStatus{
				ExitCode: EXIT_OPERTION_FAILED,
				Error:    fmt.Errorf("failed to confirm tx - %s", err.Error()),
			})
		}

		err = dispatchResult.WaitForApply()
		if err != nil {
			log.Errorf("failed tx - %s", err.Error())
			panic(PanicStatus{
				ExitCode: EXIT_OPERTION_FAILED,
				Error:    fmt.Errorf("failed tx - %s", err.Error()),
			})
		}

		log.Info("transfer successful")
	},
}

func init() {
	transferCmd.Flags().Bool(MUTEZ_FLAG, false, "amount in mutez")
	transferCmd.Flags().Bool(CONFIRM_FLAG, false, "automatically confirms transfer")
	RootCmd.AddCommand(transferCmd)
}
