package main

import (
	"flag"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/cobra"
	"github.com/weedge/wedis/cmd/srv"
	"github.com/weedge/wedis/pkg/configparser"
	"github.com/weedge/wedis/pkg/version"
)

var (
	moduleName = version.Get().Module
	rootCmd    = &cobra.Command{
		Use:   moduleName,
		Short: fmt.Sprintf("%s module", moduleName),
	}
)

func main() {
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	rootCmd.AddCommand(
		srv.NewCommand(),
	)

	configparser.Flags(rootCmd.PersistentFlags())
	if err := rootCmd.Execute(); err != nil {
		klog.Fatal(err.Error())
	}
}
