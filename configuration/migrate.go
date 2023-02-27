package configuration

import (
	"bytes"
	"fmt"
	"os"

	"github.com/alis-is/tezpay/common"
	tezpay_configuration "github.com/alis-is/tezpay/configuration/v"
	"github.com/alis-is/tezpay/constants"
	"github.com/alis-is/tezpay/state"

	"github.com/hjson/hjson-go/v4"
)

type LatestConfigurationType = tezpay_configuration.ConfigurationV0

func WriteMigratedConfiguration(source string, configuration LatestConfigurationType) error {
	var marshaled []byte
	var err error
	encoderOptions := hjson.DefaultOptions()
	encoderOptions.IndentBy = "\t"
	marshaled, err = hjson.MarshalWithOptions(configuration, encoderOptions)
	if err != nil {
		return err
	}
	err = os.WriteFile(source, marshaled, 0644)
	return err
}

func migrate(sourceBytes []byte, versionInfo *common.ConfigurationVersionInfo) ([]byte, error) {
	switch versionInfo.TPVersion {
	/* here go future migrations */
	}

	return sourceBytes, nil
}

func Migrate(sourceBytes []byte, versionInfo *common.ConfigurationVersionInfo, persist bool) (*LatestConfigurationType, error) {
	originalSourceBytes := sourceBytes

	sourceBytes, err := migrate(sourceBytes, versionInfo)
	if err != nil {
		return nil, err
	}

	// load final config
	configuration := tezpay_configuration.GetDefaultV0()
	err = hjson.Unmarshal(sourceBytes, &configuration)
	if err != nil {
		return nil, err
	}
	configuration.SourceBytes = sourceBytes // inject bytes for processing in future

	// persist migrated config
	isMigrated := !bytes.Equal(originalSourceBytes, sourceBytes)
	if isMigrated && persist {
		source := state.Global.GetConfigurationFilePath()
		os.Rename(source, source+constants.CONFIG_FILE_BACKUP_SUFFIX)
		err := WriteMigratedConfiguration(source, configuration)
		if err != nil {
			return nil, fmt.Errorf("failed to write migrated configuration - %s", err.Error())
		}
	}

	return &configuration, nil
}
