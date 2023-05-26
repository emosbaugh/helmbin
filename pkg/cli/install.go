package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/emosbaugh/helmbin/pkg/config"
	"github.com/k0sproject/k0s/cmd/install"
	"github.com/k0sproject/k0s/cmd/start"
	"github.com/spf13/cobra"
)

// NewCmdInstall returns a cobra command for installing the server as a systemd service
func NewCmdInstall(cli *CLI) *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Installs and starts the server as a systemd service",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: options
			config := config.Default()
			// Hack so you can re-run this command
			_ = os.RemoveAll("/etc/systemd/system/k0scontroller.service")
			k0scmd := install.NewInstallCmd()
			k0scmd.SetArgs([]string{
				"controller",
				"--enable-worker",
				"--no-taints",
				fmt.Sprintf("--data-dir=%s", filepath.Join(config.DataDir, "k0s")),
				fmt.Sprintf("--config=%s", config.K0sConfigFile),
			})
			if err := k0scmd.ExecuteContext(cmd.Context()); err != nil {
				return fmt.Errorf("failed to install k0s: %w", err)
			}
			k0scmd = start.NewStartCmd()
			k0scmd.SetArgs([]string{})
			if err := cmd.ExecuteContext(cmd.Context()); err != nil {
				return fmt.Errorf("failed to start k0s: %w", err)
			}
			return nil
		},
	}
}
