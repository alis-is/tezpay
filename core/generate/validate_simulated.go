package generate

import (
	"blockwatch.cc/tzgo/tezos"
	"github.com/alis-is/tezpay/configuration"
	"github.com/alis-is/tezpay/constants/enums"
	log "github.com/sirupsen/logrus"
)

type PayoutSimulatedCandidateValidation func(candidate *PayoutCandidateSimulated, configuration *configuration.RuntimeConfiguration, overrides *configuration.RuntimeDelegatorOverride)
type PayoutSimulatedCandidateValidator struct {
	Id       string
	Validate PayoutSimulatedCandidateValidation
}

type PayoutSimulatedValidationContext struct {
	Configuration *configuration.RuntimeConfiguration
	Overrides     *configuration.RuntimeDelegatorOverride
	Payout        *PayoutCandidateSimulated
}

func (validationContext *PayoutSimulatedValidationContext) ToPayoutCandidateSimulated() *PayoutCandidateSimulated {
	return validationContext.Payout
}

func (validationContext *PayoutSimulatedValidationContext) Validate(validators ...PayoutSimulatedCandidateValidator) *PayoutSimulatedValidationContext {
	if validationContext.Payout.IsInvalid || len(validators) == 0 {
		return validationContext
	}
	for _, validator := range validators {
		log.Tracef("validating payout to %s with %s", validationContext.Payout.Recipient, validator.Id)
		validator.Validate(validationContext.Payout, validationContext.Configuration, validationContext.Overrides)
		log.Tracef("payout to %s validation result: %t", validationContext.Payout.Recipient, validationContext.Payout.IsInvalid)
		if validationContext.Payout.IsInvalid {
			break
		}
	}
	return validationContext
}

// validation

func ValidateSimulatedMinumumAmount(candidate *PayoutCandidateSimulated, configuration *configuration.RuntimeConfiguration, _ *configuration.RuntimeDelegatorOverride) {
	treshhold := configuration.PayoutConfiguration.MinimumAmount
	if treshhold.IsNeg() || candidate.TxKind != enums.PAYOUT_TX_KIND_TEZ { // if payout is not tezos we respect anything above 0
		treshhold = tezos.Zero
	}
	diff := candidate.BondsAmount.Sub(treshhold)
	if diff.IsNeg() || diff.IsZero() {
		candidate.IsInvalid = true
		candidate.InvalidBecause = enums.INVALID_PAYOUT_BELLOW_MINIMUM
	}
}

// Validators
var (
	MinumumAmountSimulatedValidator = PayoutSimulatedCandidateValidator{Id: "MinumumAmountValidator", Validate: ValidateSimulatedMinumumAmount}
)
