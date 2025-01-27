package daemonServer

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/zeropsio/zcli/src/proto/daemon"
	"github.com/zeropsio/zcli/src/vpn"
	"google.golang.org/grpc"
)

type Config struct {
	Socket  string
	Address string
}

type Handler struct {
	daemon.UnimplementedZeropsDaemonProtocolServer

	config Config
	vpn    *vpn.Handler
}

func New(config Config, vpn *vpn.Handler) *Handler {
	return &Handler{
		config: config,
		vpn:    vpn,
	}
}

func (h *Handler) Run(ctx context.Context) error {
	address, err := url.Parse(h.config.Socket)
	if err != nil {
		return err
	}

	err = removeUnusedServerSocket(address)
	if err != nil {
		return err
	}

	var lis net.Listener
	if h.config.Socket != "" {
		lis, err = net.Listen("unix", h.config.Socket)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to listen: %v", err))
		}
		if err := os.Chmod(h.config.Socket, 0666); err != nil {
			return errors.New(fmt.Sprintf("failed to chmod: %v", err))
		}

	} else if h.config.Address != "" {
		lis, err = net.Listen("tcp", h.config.Address)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to listen: %v", err))
		}
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	daemon.RegisterZeropsDaemonProtocolServer(grpcServer, h)

	go func() {
		grpcServer.Serve(lis)
	}()

	<-ctx.Done()

	err = lis.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to close listen: %v", err))
	}

	grpcServer.GracefulStop()

	return nil
}

func removeUnusedServerSocket(address *url.URL) error {

	socketDir := filepath.Dir(address.Path)
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		return fmt.Errorf("unable to create socket directory (%s)", socketDir)
	}

	if err := os.Chmod(socketDir, 0755); err != nil {
		return fmt.Errorf("unable to change socket directory (%s) permissions", socketDir)
	}

	if _, errFound := os.Stat(address.Path); errFound != nil {
		return nil
	}

	conn, err := net.DialTimeout("unix", address.Path, 1*time.Second)
	if serverIsRunning := err == nil; serverIsRunning {
		defer func() { _ = conn.Close() }()
		return fmt.Errorf("socket %s already in use", address.String())
	}

	_ = os.Remove(address.Path)
	if _, errFound := os.Stat(address.Path); errFound == nil {
		return fmt.Errorf("unused socket %s can't be deleted", address.String())
	}
	return nil
}
