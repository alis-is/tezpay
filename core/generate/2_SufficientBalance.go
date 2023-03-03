package generate

import (
	"fmt"
	"time"

	"blockwatch.cc/tzgo/tezos"
	"github.com/alis-is/tezpay/common"
	"github.com/alis-is/tezpay/constants"
	"github.com/alis-is/tezpay/constants/enums"
	"github.com/alis-is/tezpay/extension"
	"github.com/alis-is/tezpay/utils"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type CheckBalanceHookData struct {
	SkipTezCheck bool                            `json:"skip_tez_check"`
	IsSufficient bool                            `json:"is_sufficient"`
	Message      string                          `json:"message"`
	Payouts      []PayoutCandidateWithBondAmount `json:"payouts"`
}

func checkBalanceWithHook(data *CheckBalanceHookData) error {
	err := extension.ExecuteHook(enums.EXTENSION_HOOK_CHECK_BALANCE, "0.1", data)
	if err != nil {
		return err
	}
	return nil
}

func checkBalanceWithCollector(data *CheckBalanceHookData, ctx *PayoutGenerationContext) error {
	if data.SkipTezCheck { // skip tez check for cases when pervious hook already checked it
		return nil
	}
	payableBalance, err := ctx.GetCollector().GetBalance(ctx.PayoutKey.Address())
	if err != nil {
		return err
	}

	configuration := ctx.GetConfiguration()

	totalPayouts := len(lo.Filter(data.Payouts, func(candidate PayoutCandidateWithBondAmount, _ int) bool {
		return !candidate.IsInvalid
	}))
	// add all bonds, fees and donations destinations
	totalPayouts = totalPayouts + len(configuration.IncomeRecipients.Bonds) + len(configuration.IncomeRecipients.Fees) + utils.Max(len(configuration.IncomeRecipients.Donations), 1)

	requiredbalance := lo.Reduce(data.Payouts, func(agg tezos.Z, candidate PayoutCandidateWithBondAmount, _ int) tezos.Z {
		if candidate.TxKind == enums.PAYOUT_TX_KIND_TEZ {
			return agg.Add(candidate.BondsAmount)
		}
		return agg
	}, tezos.Zero)

	requiredbalance = ctx.StageData.BakerBondsAmount.Add(requiredbalance)
	requiredbalance = requiredbalance.Add(tezos.NewZ(constants.PAYOUT_FEE_BUFFER).Mul64(int64(totalPayouts)))

	diff := payableBalance.Sub(requiredbalance)
	if diff.IsNeg() || diff.IsZero() {
		data.IsSufficient = false
		data.Message = fmt.Sprintf("required: %s, available: %s", requiredbalance, payableBalance)
	}
	return nil
}

func runBalanceCheck(ctx *PayoutGenerationContext, check func(*CheckBalanceHookData) error, data *CheckBalanceHookData, options *common.GeneratePayoutsOptions) error {
	notificatorTrigger := 0
	for {
		if err := check(data); err != nil {
			if options.WaitForSufficientBalance {
				log.Errorf("failed to check balance - %s, waiting 5 minutes...", err.Error())
				time.Sleep(time.Minute * 5)
				continue
			}
			return err
		}

		if !data.IsSufficient {
			if options.WaitForSufficientBalance {
				log.Warnf("insufficient balance - %s, waiting 5 minutes...", data.Message)
				if notificatorTrigger%12 == 0 { // every hour
					ctx.AdminNotify(fmt.Sprintf("insufficient balance - %s", data.Message))
				}
				time.Sleep(time.Minute * 5)
				notificatorTrigger++
				continue
			}
			return fmt.Errorf("insufficient balance - %s", data.Message)
		}
		break
	}
	return nil
}

/*
Technically we could calculate real required balance by checking all payouts and fees and donations in final stage
but because of potential changes of transaction fees (on-chain state changes) it would not be accurate anyway.
So we just try to estimate with a buffer which should be enough for most cases.
*/

func CheckSufficientBalance(ctx *PayoutGenerationContext, options *common.GeneratePayoutsOptions) (*PayoutGenerationContext, error) {
	if options.SkipBalanceCheck { // skip
		return ctx, nil
	}

	log.Debugf("checking for sufficient balance")
	hookResponse := CheckBalanceHookData{
		IsSufficient: true,
		Payouts:      ctx.StageData.PayoutCandidatesWithBondAmount,
	}

	checks := []func(*CheckBalanceHookData) error{
		func(data *CheckBalanceHookData) error {
			log.Trace("checking balance with hook")
			return checkBalanceWithHook(data)
		},
		func(data *CheckBalanceHookData) error {
			log.Trace("checking tez balance with collector")
			return checkBalanceWithCollector(data, ctx)
		},
	}

	for _, check := range checks {
		err := runBalanceCheck(ctx, check, &hookResponse, options)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}
