package vpn

import (
	"context"
)

func (h *Handler) stopVpn(ctx context.Context) error {
	h.logger.Debug("stopping VPN")
	data := h.storage.Data()

	if data.InterfaceName == "" {
		return nil
	}

	h.logger.Debug("clean vpn start")
	if err := h.cleanVpn(ctx, data.InterfaceName); err != nil {
		return err
	}
	h.logger.Debug("clean vpn end")

	h.logger.Debug("clean vpn DNS")
	if err := h.DnsClean(ctx); err != nil {
		return err
	}
	h.logger.Debug("clean DNS start")

	if err := h.storage.Clear(); err != nil {
		return err
	}
	return nil
}
