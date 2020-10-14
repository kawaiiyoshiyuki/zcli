package params

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Handler struct {
	params map[string]interface{}
	viper  *viper.Viper
}

func New() *Handler {
	return &Handler{
		params: make(map[string]interface{}),
		viper:  viper.New(),
	}
}

type optionConfig struct {
	persistent bool
}

func (h *Handler) getCmdId(cmd *cobra.Command, name string) string {
	return cmd.Use + name
}

func (h *Handler) RegisterString(cmd *cobra.Command, name, defaultValue, description string) {

	var paramValue string

	cmd.Flags().StringVar(&paramValue, name, "", description)

	h.params[h.getCmdId(cmd, name)] = func() string {
		if paramValue != "" {
			return paramValue
		}
		if h.viper.GetString(name) != "" {
			return h.viper.GetString(name)
		}

		return defaultValue
	}
}

func (h *Handler) RegisterPersistentString(cmd *cobra.Command, name, defaultValue, description string) {

	var paramValue string

	cmd.PersistentFlags().StringVar(&paramValue, name, "", description)
	h.viper.BindPFlags(cmd.PersistentFlags())

	h.params[name] = func() string {
		if paramValue != "" {
			return paramValue
		}
		if h.viper.GetString(name) != "" {
			return h.viper.GetString(name)
		}

		return defaultValue
	}
}

func (h *Handler) RegisterUInt32(cmd *cobra.Command, name string, defaultValue uint32, description string) {
	var paramValue uint32

	cmd.Flags().Uint32Var(&paramValue, name, defaultValue, description)

	h.params[name] = func() uint32 {
		if paramValue > 0 {
			return paramValue
		}
		if h.viper.GetInt32(name) != 0 {
			return h.viper.GetUint32(name)
		}

		return defaultValue
	}
}

func (h *Handler) GetPersistentString(name string) string {
	if param, exists := h.params[name]; exists {
		if v, ok := param.(func() string); ok {
			return v()
		}
		return ""
	}
	return ""
}

func (h *Handler) GetString(cmd *cobra.Command, name string) string {
	id := h.getCmdId(cmd, name)
	if param, exists := h.params[id]; exists {
		if v, ok := param.(func() string); ok {
			return v()
		}
		return ""
	}
	return ""
}

func (h *Handler) GetUint32(name string) uint32 {
	if param, exists := h.params[name]; exists {
		if v, ok := param.(func() uint32); ok {
			return v()
		}
		return 0
	}
	return 0
}

func (h *Handler) InitViper() error {

	path, err := os.Getwd()
	if err != nil {
		return err
	}
	h.viper.AddConfigPath(path)
	h.viper.SetConfigName("zcli.config")
	h.viper.AutomaticEnv()

	if err := h.viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", h.viper.ConfigFileUsed())
	}

	return nil
}
