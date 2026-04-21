package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Build-time variables injected via ldflags:
//
//	go build -ldflags "-X cmd.Version=1.2.3 -X cmd.Commit=abc1234 -X cmd.BuildDate=2024-01-15"
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

type versionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

var versionJSON bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the current version, commit hash, build date, and runtime information for vaultpull.`,
	RunE:  runVersion,
}

func init() {
	versionCmd.Flags().BoolVar(&versionJSON, "json", false, "Output version info as JSON")
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	info := versionInfo{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	if versionJSON {
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "vaultpull %s\n", info.Version)
	fmt.Fprintf(cmd.OutOrStdout(), "  commit:     %s\n", info.Commit)
	fmt.Fprintf(cmd.OutOrStdout(), "  built:      %s\n", info.BuildDate)
	fmt.Fprintf(cmd.OutOrStdout(), "  go version: %s\n", info.GoVersion)
	fmt.Fprintf(cmd.OutOrStdout(), "  os/arch:    %s/%s\n", info.OS, info.Arch)
	return nil
}
