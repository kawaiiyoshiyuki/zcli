package dns

import (
	"context"
	"os/exec"

	"github.com/zerops-io/zcli/src/constants"
	"github.com/zerops-io/zcli/src/daemonStorage"
	"github.com/zerops-io/zcli/src/dnsServer"
	"github.com/zerops-io/zcli/src/utils"
	"github.com/zerops-io/zcli/src/utils/cmdRunner"
	"github.com/zerops-io/zcli/src/utils/logger"
)

func CleanDns(_ context.Context, _ logger.Logger, data daemonStorage.Data, dnsServer *dnsServer.Handler) error {

	switch data.DnsManagement {
	case daemonStorage.LocalDnsManagementUnknown:
		return nil
	case daemonStorage.LocalDnsManagementSystemdResolve:
		return nil
	case daemonStorage.LocalDnsManagementResolveConf:
		cmd := exec.Command("resolvconf", "-d", data.InterfaceName)
		_, err := cmdRunner.Run(cmd)
		if err != nil {
			return err
		}
	case daemonStorage.LocalDnsManagementFile:
		err := utils.RemoveFirstLine(constants.ResolvFilePath, "nameserver "+data.DnsIp.String())
		if err != nil {
			return err
		}
	case
		daemonStorage.LocalDnsManagementNetworkSetup,
		daemonStorage.LocalDnsManagementScutil:
		if err := setDnsByNetworksetup(data, dnsServer, false); err != nil {
			return err
		}

	case daemonStorage.LocalDnsManagementWindows:
		return nil
	default:
		return nil
	}
	return nil
}
