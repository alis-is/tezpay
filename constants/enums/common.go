package enums

type EWalletMode string

const (
	WALLET_MODE_LOCAL_PRIVATE_KEY  EWalletMode = "local-private-key"
	WALLET_MODE_LOCAL_PRIVATE_KEY2 EWalletMode = "local_private_key"
	WALLET_MODE_REMOTE_SIGNER      EWalletMode = "remote-signer"
	WALLET_MODE_REMOTE_SIGNER2     EWalletMode = "remote_signer"
)

var (
	SUPPORTED_WALLET_MODES = []EWalletMode{
		WALLET_MODE_LOCAL_PRIVATE_KEY,
		WALLET_MODE_LOCAL_PRIVATE_KEY2,
		WALLET_MODE_REMOTE_SIGNER,
		WALLET_MODE_REMOTE_SIGNER2,
	}
)

type EPayoutMode string

const (
	PAYOUT_MODE_ACTUAL EPayoutMode = "actual"
	PAYOUT_MODE_IDEAL  EPayoutMode = "ideal"
)

var (
	SUPPORTED_PAYOUT_MODES = []EPayoutMode{
		PAYOUT_MODE_ACTUAL,
		PAYOUT_MODE_IDEAL,
	}
)

type EPayoutInvalidReason string

const (
	INVALID_DELEGATOR_EMPTIED            EPayoutInvalidReason = "DELEGATOR_EMPTIED"
	INVALID_DELEGATOR_IGNORED            EPayoutInvalidReason = "DELEGATOR_IGNORED"
	INVALID_DELEGATOR_LOW_BAlANCE        EPayoutInvalidReason = "DELEGATOR_LOW_BALANCE"
	INVALID_PAYOUT_BELLOW_MINIMUM        EPayoutInvalidReason = "PAYOUT_BELLOW_MINIMUM"
	INVALID_PAYOUT_ZERO                  EPayoutInvalidReason = "PAYOUT_ZERO"
	INVALID_INVALID_ADDRESS              EPayoutInvalidReason = "PAYOUT_INVALID_RECIPIENT"
	INVALID_KT_IGNORED                   EPayoutInvalidReason = "PAYOUT_KT_IGNORED"
	INVALID_RECIPIENT_TARGETS_PAYOUT     EPayoutInvalidReason = "RECIPIENT_TARGETS_PAYOUT"
	INVALID_FAILED_TO_ESTIMATE_TX_COSTS  EPayoutInvalidReason = "FAILED_TO_ESTIMATE_TX_COSTS"
	ITERMEDIATE_FAILED_TO_ESTIMATE_BATCH EPayoutInvalidReason = "FAILED_TO_ESTIMATE_BATCH"
)

type EPayoutKind string

const (
	PAYOUT_KIND_DELEGATOR_REWARD EPayoutKind = "delegator reward"
	PAYOUT_KIND_BAKER_REWARD     EPayoutKind = "baker reward"
	PAYOUT_KIND_DONATION         EPayoutKind = "donation"
	PAYOUT_KIND_FEE_INCOME       EPayoutKind = "fee income"
	PAYOUT_KIND_INVALID          EPayoutKind = "invalid"
)

type EPayoutTransactionKind string

const (
	PAYOUT_TX_KIND_TEZ   EPayoutTransactionKind = "tez"
	PAYOUT_TX_KIND_FA1_2 EPayoutTransactionKind = "fa12"
	PAYOUT_TX_KIND_FA2   EPayoutTransactionKind = "fa2"
)

var (
	FA_OPERATION_KINDS = []EPayoutTransactionKind{
		PAYOUT_TX_KIND_FA1_2,
		PAYOUT_TX_KIND_FA2,
	}
)
