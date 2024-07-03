package cli_controller

import (
	"fmt"
	"os"

	"github.com/matzapata/ipfs-compute/provider/internal/config"
	"github.com/matzapata/ipfs-compute/provider/internal/controllers/cli/commands"
	"github.com/spf13/cobra"
)

type CliHandler struct {
}

func NewCliHandler() *CliHandler {
	return &CliHandler{}
}

func (c *CliHandler) Handle() {
	// load config
	cfg := config.ReadConfig("")

	// Define the root command
	rootCmd := &cobra.Command{
		Use:   "khachapuri",
		Short: "Create and manage kachapuri deployments",
	}

	// Define the deploy command
	deployCmd := &cobra.Command{
		Use:   "deploy --pk <admin private key>",
		Short: "Deploy a new application to kachapuri",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// get variables from flags
			adminPrivateKey, _ := cmd.Flags().GetString("private-key")

			commands.DeployCommand(cfg, adminPrivateKey)
		},
	}
	deployCmd.Flags().StringP("private-key", "pk", "", "admin wallet private key")

	// Allowance command
	allowanceCommand := &cobra.Command{
		Use:   "allowance --admin-address <address> --provider-address <provider>",
		Short: "Get the current allowance of the user to consume from the provider",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// get variables from flags
			adminAddress, _ := cmd.Flags().GetString("admin-address")
			providerAddress, _ := cmd.Flags().GetString("provider-address")

			commands.AllowanceCommand(cfg, adminAddress, providerAddress)
		},
	}
	allowanceCommand.Flags().StringP("admin-address", "aa", "", "address of the admin")
	allowanceCommand.Flags().StringP("provider-address", "pa", "", "address of the provider")

	// Deposit command
	depositCommand := &cobra.Command{
		Use:   "deposit [amount] --pk <private key>",
		Short: "Deposit funds into the escrow account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			amount := args[0]
			adminPrivateKey, _ := cmd.Flags().GetString("private-key")

			commands.DepositCommand(cfg, amount, adminPrivateKey)
		},
	}
	depositCommand.Flags().UintP("amount", "", 0, "amount to deposit")
	depositCommand.Flags().StringP("private-key", "pk", "", "admin wallet private key")

	// Withdraw command
	withdrawCommand := &cobra.Command{
		Use:   "withdraw [amount] --pk <private key>",
		Short: "Withdraw funds from the escrow account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			amount := args[0]
			adminPrivateKey, _ := cmd.Flags().GetString("private-key")

			commands.WithdrawCommand(cfg, amount, adminPrivateKey)
		},
	}
	depositCommand.Flags().StringP("private-key", "pk", "", "admin wallet private key")

	// Approve command
	approveCommand := &cobra.Command{
		Use:   "approve --amount <amount> --price <price> --provider-address <provider address>",
		Short: "Approve the provider to consume USDC from the user's account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			// get variables from flags
			amount, _ := cmd.Flags().GetString("amount")
			price, _ := cmd.Flags().GetString("price")
			providerAddress, _ := cmd.Flags().GetString("provider-address")
			adminPrivateKey, _ := cmd.Flags().GetString("private-key")

			commands.ApproveCommand(cfg, amount, price, providerAddress, adminPrivateKey)
		},
	}
	approveCommand.Flags().UintP("amount", "", 0, "amount to approve as a bignumber string")
	approveCommand.Flags().UintP("price", "", 0, "limit price per request as a bignumber string")
	approveCommand.Flags().StringP("provider-address", "", "", "provider address")
	depositCommand.Flags().StringP("private-key", "pk", "", "admin wallet private key")

	// Balance command
	balanceCommand := &cobra.Command{
		Use:   "balance [address]",
		Short: "Get the current balance of the user in the escrow account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			address := args[0]

			commands.BalanceCommand(cfg, address)
		},
	}
	balanceCommand.Flags().StringP("address", "", "", "Address to get the balance")

	// Resolve domain command
	resolveCmd := &cobra.Command{
		Use:   "resolve [domain]",
		Short: "Resolve a domain to get the provider",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			domain := args[0]

			commands.ResolveCommand(cfg, domain)
		},
	}

	// TODO: command to deploy resolver and register it
	// TODO: set global config

	// Add the commands to the root command
	rootCmd.AddCommand(
		deployCmd,
		resolveCmd,
		balanceCommand,
		allowanceCommand,
		depositCommand,
		withdrawCommand,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
