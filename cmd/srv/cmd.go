package srv

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/weedge/wedis/internal/srv"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "srv",
		Short: "start wedis server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			server, err := srv.NewServer(ctx)
			if err != nil {
				return err
			}
			return server.Run(ctx)
		},
	}
}
