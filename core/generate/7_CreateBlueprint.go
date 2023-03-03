package generate

import (
	"time"

	"blockwatch.cc/tzgo/tezos"
	"github.com/alis-is/tezpay/common"
	"github.com/alis-is/tezpay/constants/enums"
	"github.com/alis-is/tezpay/extension"
	"github.com/samber/lo"
)

func sumValidPayoutsAmount(payouts []common.PayoutRecipe) tezos.Z {
	return lo.Reduce(payouts, func(agg tezos.Z, payout common.PayoutRecipe, _ int) tezos.Z {
		if !payout.IsValid {
			return agg
		}
		return agg.Add(payout.Amount)
	}, tezos.Zero)
}

type AfterPayoutsBlueprintGeneratedHookData = common.CyclePayoutBlueprint

// NOTE: do we want to allow rewriting of blueprint?
func ExecuteAfterPayoutsBlueprintGenerated(data AfterPayoutsBlueprintGeneratedHookData) error {
	return extension.ExecuteHook(enums.EXTENSION_HOOK_AFTER_PAYOUTS_BLUEPRINT_GENERATED, "0.1", &data)
}

func CreateBlueprint(ctx *PayoutGenerationContext, options *common.GeneratePayoutsOptions) (result *PayoutGenerationContext, err error) {
	stageData := ctx.StageData

	blueprint := common.CyclePayoutBlueprint{
		Cycle:   options.Cycle,
		Payouts: stageData.Payouts,
		Summary: common.CyclePayoutSummary{
			Cycle:              options.Cycle,
			Delegators:         len(stageData.CycleData.Delegators),
			PaidDelegators:     stageData.PaidDelegators,
			StakingBalance:     stageData.CycleData.StakingBalance,
			EarnedFees:         stageData.CycleData.BlockFees,
			EarnedRewards:      stageData.CycleData.BlockRewards.Add(stageData.CycleData.EndorsementRewards),
			DistributedRewards: sumValidPayoutsAmount(stageData.Payouts),
			BondIncome:         stageData.BakerBondsAmount,
			FeeIncome:          stageData.BakerFeesAmount,
			IncomeTotal:        stageData.BakerBondsAmount.Add(stageData.BakerFeesAmount),
			DonatedBonds:       stageData.DonateBondsAmount,
			DonatedFees:        stageData.DonateFeesAmount,
			DonatedTotal:       stageData.DonateFeesAmount.Add(stageData.DonateBondsAmount),
			Timestamp:          time.Now(),
		},
	}

	err = ExecuteAfterPayoutsBlueprintGenerated(blueprint)
	if err != nil {
		return ctx, err
	}

	stageData.PayoutBlueprint = &blueprint
	return ctx, nil
}
