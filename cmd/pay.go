//go:build !wasm

package cmd

import (
	"fmt"

	"github.com/alis-is/tezpay/common"
	"github.com/alis-is/tezpay/core"
	reporter_engines "github.com/alis-is/tezpay/engines/reporter"
	"github.com/alis-is/tezpay/extension"
	"github.com/alis-is/tezpay/state"
	"github.com/alis-is/tezpay/utils"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var payCmd = &cobra.Command{
	Use:   "pay",
	Short: "manual payout",
	Long:  "runs manual payout",
	Run: func(cmd *cobra.Command, args []string) {
		config, collector, signer, transactor := assertRunWithResult(loadConfigurationEnginesExtensions, common.EXIT_CONFIGURATION_LOAD_FAILURE).Unwrap()
		defer extension.CloseExtensions()

		cycle, _ := cmd.Flags().GetInt64(CYCLE_FLAG)
		skipBalanceCheck, _ := cmd.Flags().GetBool(SKIP_BALANCE_CHECK_FLAG)
		confirmed, _ := cmd.Flags().GetBool(CONFIRM_FLAG)
		mixInContractCalls, _ := cmd.Flags().GetBool(DISABLE_SEPERATE_SC_PAYOUTS_FLAG)
		mixInFATransfers, _ := cmd.Flags().GetBool(DISABLE_SEPERATE_FA_PAYOUTS_FLAG)

		fsReporter := reporter_engines.NewFileSystemReporter(config)
		stdioReporter := reporter_engines.NewStdioReporter(config)

		if !state.Global.IsDonationPromptDisabled() && !config.IsDonatingToTezCapital() {
			assertRequireConfirmation("With your current configuration you are not going to donate to tez.capital. Do you want to proceed?")
		}

		if cycle <= 0 {
			lastCompletedCycle := assertRunWithResultAndErrFmt(collector.GetLastCompletedCycle, common.EXIT_OPERTION_FAILED, "failed to get last completed cycle")
			cycle = lastCompletedCycle + cycle
		}

		var generationResult *common.GeneratePayoutsResult
		fromFile, _ := cmd.Flags().GetString(TO_FILE_FLAG)
		if fromFile != "" {
			generationResult = assertRunWithResult(func() (*common.GeneratePayoutsResult, error) {
				return loadGeneratePayoutsResultFromFile(fromFile)
			}, common.EXIT_PAYOUTS_READ_FAILURE)
		} else {
			generationResult = assertRunWithResult(func() (*common.GeneratePayoutsResult, error) {
				return core.GeneratePayouts(config, common.NewGeneratePayoutsEngines(collector, signer, notifyAdminFactory(config)),
					&common.GeneratePayoutsOptions{
						Cycle:            cycle,
						SkipBalanceCheck: skipBalanceCheck,
					})
			}, common.EXIT_OPERTION_FAILED)
		}
		log.Info("checking past reports")
		preparationResult := assertRunWithResult(func() (*common.PreparePayoutsResult, error) {
			return core.PreparePayouts(generationResult, config, common.NewPreparePayoutsEngineContext(collector, fsReporter, notifyAdminFactory(config)), &common.PreparePayoutsOptions{})
		}, common.EXIT_OPERTION_FAILED)

		if state.Global.GetWantsOutputJson() {
			utils.PrintPayoutsAsJson(preparationResult.ReportsOfPastSuccesfulPayouts)
			utils.PrintPayoutsAsJson(preparationResult.Payouts)
		} else {
			utils.PrintInvalidPayoutRecipes(preparationResult.Payouts, generationResult.Cycle)
			utils.PrintReports(preparationResult.ReportsOfPastSuccesfulPayouts, fmt.Sprintf("Already Successfull - #%d", generationResult.Cycle), true)
			utils.PrintValidPayoutRecipes(preparationResult.Payouts, generationResult.Cycle)
		}

		if len(utils.OnlyValidPayouts(preparationResult.Payouts)) == 0 {
			log.Info("nothing to pay out")
			notificator, _ := cmd.Flags().GetString(NOTIFICATOR_FLAG)
			if notificator != "" { // rerun notification through notificator if specified manually
				notifyPayoutsProcessed(config, &generationResult.Summary, notificator)
			}
			panic(common.PanicStatus{
				ExitCode: common.EXIT_SUCCESS,
				Message:  "nothing to pay out",
			})
		}

		if !confirmed {
			assertRequireConfirmation("Do you want to pay out above VALID payouts?")
		}

		log.Info("executing payout")
		executionResult := assertRunWithResult(func() (common.ExecutePayoutsResult, error) {
			var reporter common.ReporterEngine
			reporter = fsReporter
			if reportToStdout, _ := cmd.Flags().GetBool(REPORT_TO_STDOUT); reportToStdout {
				reporter = stdioReporter
			}
			return core.ExecutePayouts(preparationResult, config, common.NewExecutePayoutsEngineContext(signer, transactor, reporter, notifyAdminFactory(config)), &common.ExecutePayoutsOptions{
				MixInContractCalls: mixInContractCalls,
				MixInFATransfers:   mixInFATransfers,
			})
		}, common.EXIT_OPERTION_FAILED)

		// notify
		failedCount := lo.CountBy(executionResult, func(br common.BatchResult) bool { return !br.IsSuccess })
		if len(executionResult) > 0 && failedCount > 0 {
			log.Errorf("%d of operations failed", failedCount)
			panic(common.PanicStatus{
				ExitCode: common.EXIT_OPERTION_FAILED,
				Error:    fmt.Errorf("%d of operations failed", failedCount),
			})
		}
		if silent, _ := cmd.Flags().GetBool(SILENT_FLAG); !silent {
			notifyPayoutsProcessedThroughAllNotificators(config, &generationResult.Summary)
		}
		utils.PrintBatchResults(executionResult, fmt.Sprintf("Results of #%d", generationResult.Cycle), config.Network.Explorer)
	},
}

func init() {
	payCmd.Flags().Bool(CONFIRM_FLAG, false, "automatically confirms generated payouts")
	payCmd.Flags().Int64P(CYCLE_FLAG, "c", 0, "cycle to generate payouts for")
	payCmd.Flags().Bool(REPORT_TO_STDOUT, false, "prints them to stdout (wont write to file)")
	payCmd.Flags().String(FROM_FILE_FLAG, "", "loads payouts from file instead of generating on the fly")
	payCmd.Flags().Bool(DISABLE_SEPERATE_SC_PAYOUTS_FLAG, false, "disables smart contract separation (mixes txs and smart contract calls within batches)")
	payCmd.Flags().Bool(DISABLE_SEPERATE_FA_PAYOUTS_FLAG, false, "disables fa transfers separation (mixes txs and fa transfers within batches)")
	payCmd.Flags().BoolP(SILENT_FLAG, "s", false, "suppresses notifications")
	payCmd.Flags().String(NOTIFICATOR_FLAG, "", "Notify through specific notificator")
	payCmd.Flags().Bool(SKIP_BALANCE_CHECK_FLAG, false, "skips payout wallet balance check")

	RootCmd.AddCommand(payCmd)
}
