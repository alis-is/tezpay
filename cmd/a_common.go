package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"blockwatch.cc/tzgo/tezos"
	"github.com/AlecAivazis/survey/v2"
	"github.com/alis-is/tezpay/clients"
	"github.com/alis-is/tezpay/configuration"
	"github.com/alis-is/tezpay/constants"
	"github.com/alis-is/tezpay/core/common"
	"github.com/alis-is/tezpay/core/reports"
	"github.com/alis-is/tezpay/notifications"
	"github.com/alis-is/tezpay/state"
	"github.com/alis-is/tezpay/utils"
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
)

type configurationAndEngines struct {
	Configuration *configuration.RuntimeConfiguration
	Collector     common.CollectorEngine
	Signer        common.SignerEngine
	Transactor    common.TransactorEngine
}

func (cae *configurationAndEngines) Unwrap() (*configuration.RuntimeConfiguration, common.CollectorEngine, common.SignerEngine, common.TransactorEngine) {
	return cae.Configuration, cae.Collector, cae.Signer, cae.Transactor
}

func loadConfigurationAndEngines() (*configurationAndEngines, error) {
	config, err := configuration.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration - %s", err.Error())
	}

	signerEngine := state.Global.SignerOverride
	if signerEngine == nil {
		signerEngine, err = config.LoadSigner()
		if err != nil {
			return nil, fmt.Errorf("failed to load signer - %s", err.Error())
		}
	}
	// for testing point transactor to testnet
	// transactorEngine, err := clients.InitDefaultTransactor("https://rpc.tzkt.io/ghostnet/", "https://api.ghostnet.tzkt.io/") // (config.Network.RpcUrl, config.Network.TzktUrl)
	transactorEngine, err := clients.InitDefaultTransactor(config.Network.RpcUrl, config.Network.TzktUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to load transactor - %s", err.Error())
	}

	collector, err := clients.InitDefaultRpcAndTzktColletor(config.Network.RpcUrl, config.Network.TzktUrl)
	if err != nil {
		return nil, err
	}

	if utils.IsTty() && state.Global.GetIsInDebugMode() {
		marshaled, _ := json.MarshalIndent(config, "", "\t")
		fmt.Println("Loaded configuration:", string(marshaled))
	}

	return &configurationAndEngines{
		Configuration: config,
		Collector:     collector,
		Signer:        signerEngine,
		Transactor:    transactorEngine,
	}, nil
}

func loadPayoutBlueprintFromFile(fromFile string) (*common.CyclePayoutBlueprint, error) {
	log.Infof("reading payouts from '%s'", fromFile)
	data, err := os.ReadFile(fromFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read payouts from file - %s", err.Error())
	}
	payouts, err := utils.PayoutBlueprintFromJson(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payouts from file - %s", err.Error())
	}
	return payouts, nil
}

func writePayoutBlueprintToFile(toFile string, blueprint *common.CyclePayoutBlueprint) error {
	log.Infof("writing payouts to '%s'", toFile)
	err := os.WriteFile(toFile, utils.PayoutBlueprintToJson(blueprint), 0644)
	if err != nil {
		return fmt.Errorf("failed to write generated payouts to file - %s", err.Error())
	}
	return nil
}

func loadPastPayoutReports(baker tezos.Address, cycle int64) ([]common.PayoutReport, error) {
	reports, err := reports.ReadPayoutReports(cycle)
	if err == nil || os.IsNotExist(err) {
		return utils.FilterReportsByBaker(reports, baker), nil
	}
	return []common.PayoutReport{}, err
}

func requireConfirmation(msg string) error {
	proceed := false
	if utils.IsTty() {
		prompt := &survey.Confirm{
			Message: msg,
		}
		survey.AskOne(prompt, &proceed)
	}
	if !proceed {
		return errors.New("not confirmed")
	}
	return nil
}

func notifyPayoutsProcessed(configuration *configuration.RuntimeConfiguration, summary *common.CyclePayoutSummary, filter string) {
	for _, notificatorConfiguration := range configuration.NotificationConfigurations {
		if filter != "" && string(notificatorConfiguration.Type) != filter {
			continue
		}

		log.Infof("sending notification with %s", notificatorConfiguration.Type)
		notificator, err := notifications.LoadNotificatior(notificatorConfiguration.Type, notificatorConfiguration.Configuration)
		if err != nil {
			log.Warnf("failed to send notification - %s", err.Error())
			continue
		}

		err = notificator.PayoutSummaryNotify(summary)
		if err != nil {
			log.Warnf("failed to send notification - %s", err.Error())
			continue
		}
	}
	log.Info("notifications sent.")
}
func notifyPayoutsProcessedThroughAllNotificators(configuration *configuration.RuntimeConfiguration, summary *common.CyclePayoutSummary) {
	notifyPayoutsProcessed(configuration, summary, "")
}

func notifyAdmin(configuration *configuration.RuntimeConfiguration, msg string) {
	for _, notificatorConfiguration := range configuration.NotificationConfigurations {
		if !notificatorConfiguration.IsAdmin {
			continue
		}

		log.Infof("sending admin notification with %s", notificatorConfiguration.Type)
		notificator, err := notifications.LoadNotificatior(notificatorConfiguration.Type, notificatorConfiguration.Configuration)
		if err != nil {
			log.Warnf("failed to send notification - %s", err.Error())
			continue
		}

		err = notificator.AdminNotify(msg)
		if err != nil {
			log.Warnf("failed to send notification - %s", err.Error())
			continue
		}
	}
	log.Info("admin notifications sent.")
}

func notifyAdminFactory(configuration *configuration.RuntimeConfiguration) func(string) {
	return func(msg string) {
		notifyAdmin(configuration, msg)
	}
}

func printPayoutCycleReport(report *common.PayoutCycleReport) error {
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	fmt.Println("REPORT:", string(data))
	return nil
}

type versionInfo struct {
	Version string `json:"tag_name"`
}

func checkLatestVersion() {
	log.Info("checking for new version")
	// https://api.github.com/repos/tez-capital/tezpay/releases/latest
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", constants.TEZPAY_REPOSITORY))
	if err != nil {
		log.Warnf("Failed to check latest version!")
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Failed to check latest version!")
		return
	}
	var info versionInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Warnf("Failed to check latest version!")
		return
	}
	latestVersion := info.Version
	if latestVersion == "" {
		log.Warnf("failed to check latest version - empty tag!")
		return
	}
	currentVersion := constants.VERSION
	lv, err := version.NewVersion(latestVersion)
	if err != nil {
		log.Warnf("failed to check latest version - invalid version from remote!")
		return
	}
	cv, err := version.NewVersion(currentVersion)
	if err != nil {
		log.Warnf("failed to check latest version - invalid binary version!")
		return
	}

	if cv.GreaterThanOrEqual(lv) {
		log.Info("you are running latest version")
		return
	}
	err = requireConfirmation(fmt.Sprintf("You are not running latest version of tezpay (new version : '%s', current version: '%s').\n Do you want to continue anyway?", latestVersion, currentVersion))
	if err != nil && err.Error() == "not confirmed" {
		log.Infof("You can download new version here:\n\nhttps://github.com/%s/releases\n", constants.TEZPAY_REPOSITORY)
		os.Exit(1)
	}
}
