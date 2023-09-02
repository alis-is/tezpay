package tezpay_configuration

import (
	"encoding/json"

	"blockwatch.cc/tzgo/tezos"
	"github.com/alis-is/tezpay/common"
	"github.com/alis-is/tezpay/constants"
	"github.com/alis-is/tezpay/constants/enums"
)

type IncomeRecipientsV0 struct {
	Bonds       map[string]float64 `json:"bonds,omitempty"`
	Fees        map[string]float64 `json:"fees,omitempty"`
	Donate      *float64           `json:"donate,omitempty"`
	DonateFees  *float64           `json:"donate_fees,omitempty"`
	DonateBonds *float64           `json:"donate_bonds,omitempty"`
	Donations   map[string]float64 `json:"donations,omitempty"`
}

type DelegatorRequirementsV0 struct {
	MinimumBalance                        float64                   `json:"minimum_balance,omitempty"`
	BellowMinimumBalanceRewardDestination *enums.ERewardDestination `json:"below_minimum_reward_destination,omitempty"`
}

type DelegatorOverrideV0 struct {
	Recipient                    tezos.Address `json:"recipient,omitempty"`
	Fee                          *float64      `json:"fee,omitempty"`
	MinimumBalance               float64       `json:"minimum_balance,omitempty"`
	IsBakerPayingTxFee           *bool         `json:"baker_pays_transaction_fee,omitempty"`
	IsBakerPayingAllocationTxFee *bool         `json:"baker_pays_allocation_fee,omitempty"`
	MaximumBalance               *float64      `json:"maximum_balance,omitempty"`
}

type DelegatorsConfigurationV0 struct {
	Requirements DelegatorRequirementsV0        `json:"requirements,omitempty"`
	Ignore       []tezos.Address                `json:"ignore,omitempty"`
	Overrides    map[string]DelegatorOverrideV0 `json:"overrides,omitempty"`
	FeeOverrides map[string][]tezos.Address     `json:"fee_overrides,omitempty"`
}

type TezosNetworkConfigurationV0 struct {
	RpcUrl                 string `json:"rpc_url,omitempty"`
	TzktUrl                string `json:"tzkt_url,omitempty"`
	Explorer               string `json:"explorer,omitempty"`
	DoNotPaySmartContracts bool   `json:"ignore_kt,omitempty"`
	IgnoreProtocolChanges  bool   `json:"ignore_protocol_changes,omitempty"`
}

type OverdelegationConfigurationV0 struct {
	IsProtectionEnabled bool `json:"protect,omitempty"`
}

type PayoutConfigurationV0 struct {
	WalletMode                 enums.EWalletMode `json:"wallet_mode"`
	PayoutMode                 enums.EPayoutMode `json:"payout_mode"`
	Fee                        float64           `json:"fee,omitempty"`
	IsPayingTxFee              bool              `json:"baker_pays_transaction_fee,omitempty"`
	IsPayingAllocationTxFee    bool              `json:"baker_pays_allocation_fee,omitempty"`
	MinimumAmount              float64           `json:"minimum_payout_amount,omitempty"`
	IgnoreEmptyAccounts        bool              `json:"ignore_empty_accounts,omitempty"`
	TxGasLimitBuffer           *int64            `json:"transaction_gas_limit_buffer,omitempty"`
	TxDeserializationGasBuffer *int64            `json:"transaction_deserialization_gas_buffer,omitempty"`
	TxFeeBuffer                *int64            `json:"transaction_fee_buffer,omitempty"`
	KtTxFeeBuffer              *int64            `json:"kt_transaction_fee_buffer,omitempty"`
}

type ExtensionConfigurationV0 = common.ExtensionDefinition

type ConfigurationV0 struct {
	Version                    uint                          `json:"tezpay_config_version"`
	BakerPKH                   tezos.Address                 `json:"baker"`
	PayoutConfiguration        PayoutConfigurationV0         `json:"payouts"`
	Delegators                 DelegatorsConfigurationV0     `json:"delegators,omitempty"`
	IncomeRecipients           IncomeRecipientsV0            `json:"income_recipients,omitempty"`
	Network                    TezosNetworkConfigurationV0   `json:"network,omitempty"`
	Overdelegation             OverdelegationConfigurationV0 `json:"overdelegation,omitempty"`
	NotificationConfigurations []json.RawMessage             `json:"notifications,omitempty"`
	Extensions                 []ExtensionConfigurationV0    `json:"extensions,omitempty"`
	SourceBytes                []byte                        `json:"-"`
	DisableAnalytics           bool                          `json:"disable_analytics,omitempty"`
}

type NotificatorConfigurationBase struct {
	Type  string `json:"type"`
	Admin bool   `json:"admin"`
}

func GetDefaultV0() ConfigurationV0 {
	gasLimitBuffer := int64(constants.DEFAULT_TX_GAS_LIMIT_BUFFER)
	deserializaGasBuffer := int64(constants.DEFAULT_TX_DESERIALIZATION_GAS_BUFFER)
	feeBuffer := int64(constants.DEFAULT_TX_FEE_BUFFER)
	ktFeeBUffer := int64(constants.DEFAULT_KT_TX_FEE_BUFFER)

	delegatorBellowMinimumBalanceRewardDestination := enums.REWARD_DESTINATION_NONE

	return ConfigurationV0{
		Version:  0,
		BakerPKH: tezos.InvalidKey.Address(),
		Delegators: DelegatorsConfigurationV0{
			Requirements: DelegatorRequirementsV0{
				MinimumBalance:                        constants.DEFAULT_DELEGATOR_MINIMUM_BALANCE,
				BellowMinimumBalanceRewardDestination: &delegatorBellowMinimumBalanceRewardDestination,
			},
			Overrides: make(map[string]DelegatorOverrideV0),
			Ignore:    make([]tezos.Address, 0),
		},
		Network: TezosNetworkConfigurationV0{
			RpcUrl:                 constants.DEFAULT_RPC_URL,
			TzktUrl:                constants.DEFAULT_TZKT_URL,
			Explorer:               constants.DEFAULT_EXPLORER_URL,
			DoNotPaySmartContracts: false,
			IgnoreProtocolChanges:  false,
		},
		Overdelegation: OverdelegationConfigurationV0{
			IsProtectionEnabled: true,
		},
		PayoutConfiguration: PayoutConfigurationV0{
			WalletMode:                 enums.WALLET_MODE_LOCAL_PRIVATE_KEY,
			PayoutMode:                 enums.PAYOUT_MODE_ACTUAL,
			Fee:                        constants.DEFAULT_BAKER_FEE,
			IsPayingTxFee:              false,
			IsPayingAllocationTxFee:    false,
			MinimumAmount:              constants.DEFAULT_PAYOUT_MINIMUM_AMOUNT,
			TxGasLimitBuffer:           &gasLimitBuffer,
			TxDeserializationGasBuffer: &deserializaGasBuffer,
			TxFeeBuffer:                &feeBuffer,
			KtTxFeeBuffer:              &ktFeeBUffer,
		},
		IncomeRecipients:           IncomeRecipientsV0{},
		NotificationConfigurations: make([]json.RawMessage, 0),
		SourceBytes:                []byte{},
		DisableAnalytics:           false,
	}
}
