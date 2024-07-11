package cli_controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/controllers/cli/commands"
	"github.com/spf13/cobra"
)

type CliHandler struct {
	RootCmd *cobra.Command
}

func NewCliHandler(cfg *config.Config) *CliHandler {
	rootCmd := &cobra.Command{
		Use:   "khachapuri",
		Short: "Create and manage kachapuri deployments",
	}

	allowanceCommand := &cobra.Command{
		Use:   "allowance [provider-domain] --a <admin-address>",
		Short: "Get the current allowance of the user to consume from the provider",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			providerDomain := args[0]
			adminAddress, _ := cmd.Flags().GetString("a")

			err := commands.AllowanceCommand(cfg, adminAddress, providerDomain)
			handleError(err)
		},
	}
	allowanceCommand.Flags().StringP("a", "k", "", "admin wallet address")
	allowanceCommand.MarkFlagRequired("a")

	approveCommand := &cobra.Command{
		Use:   "approve [provider-domain] [amount] [price] --pk <admin-pk>",
		Short: "Approve the provider to consume USDC from the user's account",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			providerDomain := args[0]
			amount := args[1]
			price := args[2]
			adminPrivateKey, _ := cmd.Flags().GetString("pk")

			err := commands.ApproveCommand(cfg, amount, price, providerDomain, adminPrivateKey)
			handleError(err)
		},
	}
	approveCommand.Flags().StringP("pk", "k", "", "admin wallet private key")
	approveCommand.MarkFlagRequired("pk")

	balanceCommand := &cobra.Command{
		Use:   "balance --a <address>",
		Short: "Get the current balance of the user in the escrow account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			address, _ := cmd.Flags().GetString("a")

			err := commands.BalanceCommand(cfg, address)
			handleError(err)
		},
	}
	balanceCommand.Flags().StringP("a", "k", "", "admin wallet address")
	balanceCommand.MarkFlagRequired("a")

	depositCommand := &cobra.Command{
		Use:   "deposit [amount] --pk <private key>",
		Short: "Deposit funds into the escrow account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			amount := args[0]
			adminPrivateKey, _ := cmd.Flags().GetString("pk")

			err := commands.DepositCommand(cfg, amount, adminPrivateKey)
			handleError(err)
		},
	}
	depositCommand.Flags().StringP("pk", "", "", "admin wallet private key")
	depositCommand.MarkFlagRequired("pk")

	withdrawCommand := &cobra.Command{
		Use:   "withdraw [amount] --pk <private key>",
		Short: "Withdraw funds from the escrow account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			amount := args[0]
			adminPrivateKey, _ := cmd.Flags().GetString("pk")

			err := commands.WithdrawCommand(cfg, amount, adminPrivateKey)
			handleError(err)
		},
	}
	withdrawCommand.Flags().StringP("pk", "", "", "admin wallet private key")

	resolveCmd := &cobra.Command{
		Use:   "resolve [domain]",
		Short: "Resolve a domain to get the provider",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]

			err := commands.ResolveCommand(cfg, domain)
			handleError(err)
		},
	}

	buildCmd := &cobra.Command{
		Use:   "build [service]",
		Short: "Build the artifact",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]

			err := commands.BuildCommand(cfg, serviceName)
			handleError(err)
		},
	}

	runCmd := &cobra.Command{
		Use:                "run [service] ...args",
		Short:              "Run the local build",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			serviceName := args[0]
			serviceArgs := args[1:]
			fmt.Println("Running", serviceName, "with args", strings.Join(serviceArgs, " "))

			err := commands.RunCommand(cfg, serviceName, strings.Join(serviceArgs, " "))
			handleError(err)
		},
	}

	publishCmd := &cobra.Command{
		Use:   "publish [service] [provider] --pk <admin private key>",
		Short: "Publish artifact to ipfs making it available for consumption",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			service := args[0]
			provider := args[1]
			adminPrivateKey, _ := cmd.Flags().GetString("pk")

			err := commands.PublishCommand(cfg, service, provider, adminPrivateKey)
			handleError(err)
		},
	}
	publishCmd.Flags().StringP("pk", "", "", "admin wallet private key")
	publishCmd.MarkFlagRequired("pk")

	// TODO: command to deploy resolver and register it
	// TODO: set global config

	// Add the commands to the root command
	rootCmd.AddCommand(
		allowanceCommand,
		approveCommand,
		publishCmd,
		resolveCmd,
		balanceCommand,
		depositCommand,
		withdrawCommand,
		buildCmd,
		runCmd,
	)

	return &CliHandler{
		RootCmd: rootCmd,
	}
}

func (c *CliHandler) Handle() {
	if err := c.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
}
