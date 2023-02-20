package enums

type EExtensionRpcKind string

const (
	EXTENSION_STDIO_RPC EExtensionRpcKind = "stdio"
	EXTENSION_TCP_RPC   EExtensionRpcKind = "tcp"
	EXTENSION_HTTP_RPC  EExtensionRpcKind = "http"
	EXTENSION_WS_RPC    EExtensionRpcKind = "ws"
)

var (
	SUPPORTED_EXTENSION_RPC_KINDS = []EExtensionRpcKind{
		EXTENSION_STDIO_RPC,
	}
)

type EExtensionHookMode string

const (
	EXTENSION_HOOK_MODE_READ_ONLY  EExtensionHookMode = "ro"
	EXTENSION_HOOK_MODE_READ_WRITE EExtensionHookMode = "rw"
	EXTENSION_HOOK_MODE_UNKNOWN    EExtensionHookMode = "unknown"
)

var (
	SUPPORTED_EXTENSION_HOOK_MODES = []EExtensionHookMode{
		EXTENSION_HOOK_MODE_READ_ONLY,
		EXTENSION_HOOK_MODE_READ_WRITE,
	}
)

type EExtensionHook string

const (
	EXTENSION_HOOK_ALL          EExtensionHook = "all"
	EXTENSION_HOOK_TEST_NOTIFY  EExtensionHook = "test-notify"
	EXTENSION_HOOK_TEST_REQUEST EExtensionHook = "test-request"
	// can adjust generated candidate list (inject, remove, mutate)
	EXTENSION_HOOK_AFTER_CANDIDATE_GENERATED EExtensionHook = "after_candidate_generated"
	// can adjust generated bond list (inject, remove, mutate)
	EXTENSION_HOOK_AFTER_BONDS_DISTRIBUTED EExtensionHook = "after_bonds_distributed"
	// can provide aditional logic to check balance and optionally skip in tezpay check
	EXTENSION_HOOK_CHECK_BALANCE EExtensionHook = "check_balance"
	// can adjust fees to be collected by baker
	EXTENSION_HOOK_ON_FEES_COLLECTION EExtensionHook = "on_fees_collection"

	// EXTENSION_HOOK_AFTER_PAYOUTS_FINALIZED            EExtensionHook = "after_payouts_finalized"
	// EXTENSION_HOOK_AFTER_PAYOUTS_PREPARE_DISTRIBUTION EExtensionHook = "after_prepare_distribution"
	// EXTENSION_HOOK_AFTER_REWARD_DISTRIBUTED           EExtensionHook = "after_reward_distributed"
	EXTENSION_HOOK_UNKNOWN EExtensionHook = "unknown"
)

var (
	SUPPORTED_EXTENSION_HOOKS = []EExtensionHook{
		EXTENSION_HOOK_AFTER_CANDIDATE_GENERATED,
		EXTENSION_HOOK_AFTER_BONDS_DISTRIBUTED,
		EXTENSION_HOOK_CHECK_BALANCE,
		EXTENSION_HOOK_ON_FEES_COLLECTION,
	}
)

type EExtensionErrorAction string

const (
	EXTENSION_ERROR_ACTION_CONTINUE EExtensionErrorAction = "continue"
	EXTENSION_ERROR_ACTION_STOP     EExtensionErrorAction = "stop"
)

type EExtensionLifespan string

const (
	EXTENSION_LIFESPAN_SCOPED    EExtensionLifespan = "scoped"
	EXTENSION_LIFESPAN_TRANSIENT EExtensionLifespan = "transient"
)
