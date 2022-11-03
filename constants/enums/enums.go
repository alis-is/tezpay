package enums

type WalletMode string

const (
	WALLET_MODE_LOCAL_PRIVATE_KEY WalletMode = "local-private-key"
	WALLET_MODE_REMOTE_SIGNER     WalletMode = "remote-signer"
)

type EPayoutInvalidReason string

const (
	INVALID_DELEGATOR_EMPTIED           EPayoutInvalidReason = "DELEGATOR_EMPTIED"
	INVALID_DELEGATOR_IGNORED           EPayoutInvalidReason = "DELEGATOR_IGNORED"
	INVALID_DELEGATOR_LOW_BAlANCE       EPayoutInvalidReason = "DELEGATOR_LOW_BALANCE"
	INVALID_PAYOUT_BELLOW_MINIMUM       EPayoutInvalidReason = "PAYOUT_BELLOW_MINIMUM"
	INVALID_PAYOUT_ZERO                 EPayoutInvalidReason = "PAYOUT_ZERO"
	INVALID_INVALID_ADDRESS             EPayoutInvalidReason = "PAYOUT_INVALID_RECIPIENT"
	INVALID_KT_IGNORED                  EPayoutInvalidReason = "PAYOUT_KT_IGNORED"
	INVALID_RECIPIENT_TARGETS_PAYOUT    EPayoutInvalidReason = "RECIPIENT_TARGETS_PAYOUT"
	INVALID_FAILED_TO_ESTIMATE_TX_COSTS EPayoutInvalidReason = "FAILED_TO_ESTIMATE_TX_COSTS"
)

type EPayoutKind string

const (
	PAYOUT_KIND_DELEGATOR_REWARD EPayoutKind = "delegator reward"
	PAYOUT_KIND_BAKER_REWARD     EPayoutKind = "baker reward"
	PAYOUT_KIND_DONATION         EPayoutKind = "donation"
	PAYOUT_KIND_FEE_INCOME       EPayoutKind = "fee income"
	PAYOUT_KIND_INVALID          EPayoutKind = "invalid"
)
