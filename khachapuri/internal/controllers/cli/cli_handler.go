package cli_controller

import (
	"fmt"
	"os"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/controllers/cli/commands"
	"github.com/spf13/cobra"
)

type CliHandler struct {
	RootCmd *cobra.Command
}

func NewCliHandler(cfg *config.Config) *CliHandler {
	// Define the root command
	rootCmd := &cobra.Command{
		Use:   "khachapuri",
		Short: "Create and manage kachapuri deployments",
	}

	deployCmd := &cobra.Command{
		Use:   "publish --pk <admin private key>",
		Short: "Publish artifact to ipfs making it available for consumption",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			adminPrivateKey, _ := cmd.Flags().GetString("private-key")

			err := commands.PublishCommand(cfg, adminPrivateKey)
			handleError(err)
		},
	}
	deployCmd.Flags().StringP("pk", "", "", "admin wallet private key")
	deployCmd.MarkFlagRequired("pk")

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

	// Balance command
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

	// Deposit command
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

	// Withdraw command
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

	// Resolve domain command
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
		Use:   "build",
		Short: "Build the artifact",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := commands.BuildCommand(cfg)
			handleError(err)
		},
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the local build",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			err := commands.RunCommand(cfg)
			handleError(err)
		},
	}

	// TODO: command to deploy resolver and register it
	// TODO: publish command (publish build to ipfs)
	// TODO: set global config

	// Add the commands to the root command
	rootCmd.AddCommand(
		allowanceCommand,
		approveCommand,
		deployCmd,
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
