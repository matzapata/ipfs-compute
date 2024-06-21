package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/cli/commands"
	"github.com/spf13/cobra"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	rootCmd := &cobra.Command{
		Use:   "ipfs-compute",
		Short: "Create and manage IPFS Compute deployments",
	}

	// Define the deploy command
	deployCmd := &cobra.Command{
		Use:   "deploy --private-key <owner private key> --provider <provider>",
		Short: "Deploy a new application to IPFS Compute",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// get variables from env
			rpc := os.Getenv("RPC")
			privateKey := os.Getenv("PRIVATE_KEY")
			pinataApiKey := os.Getenv("PINATA_APIKEY")
			pinataSecret := os.Getenv("PINATA_SECRET")

			// get variables from flags
			provider, _ := cmd.Flags().GetString("provider")

			commands.DeployCommand(privateKey, provider, pinataApiKey, pinataSecret, rpc)
		},
	}
	deployCmd.Flags().StringP("provider", "", "", "IPFS Compute provider domain")

	// Allowance command
	allowanceCommand := &cobra.Command{
		Use:   "allowance --address <address> --provider <provider>",
		Short: "Get the current allowance of the user to consume from the provider",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// get variables from env
			rpc := os.Getenv("RPC")

			// get variables from flags
			address, _ := cmd.Flags().GetString("address")
			provider, _ := cmd.Flags().GetString("provider")

			commands.AllowanceCommand(address, provider, rpc)
		},
	}
	allowanceCommand.Flags().StringP("address", "", "", "Address to get the allowance")

	// Deposit command
	depositCommand := &cobra.Command{
		Use:   "deposit --amount <amount>",
		Short: "Deposit funds into the escrow account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// get variables from env
			rpc := os.Getenv("RPC")
			privateKey := os.Getenv("PRIVATE_KEY")

			// get variables from flags
			amount, _ := cmd.Flags().GetUint("amount")

			commands.DepositCommand(privateKey, amount, rpc)
		},
	}
	depositCommand.Flags().UintP("amount", "", 0, "Amount to deposit")

	// Withdraw command

	// Approve command

	// Register command

	// Balance command
	balanceCommand := &cobra.Command{
		Use:   "balance --address <address>",
		Short: "Get the current balance of the user in the escrow account",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// get variables from env
			rpc := os.Getenv("RPC")

			// get variables from flags
			address, _ := cmd.Flags().GetString("address")

			commands.BalanceCommand(rpc, address)
		},
	}
	balanceCommand.Flags().StringP("address", "", "", "Address to get the balance")

	// Resolve domain command
	resolveCmd := &cobra.Command{
		Use:   "resolve --domain <domain>",
		Short: "Resolve a domain to get the provider",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// get variables from env
			rpc := os.Getenv("RPC")

			// get variables from flags
			domain, _ := cmd.Flags().GetString("domain")

			commands.ResolveCommand(rpc, domain)
		},
	}
	resolveCmd.Flags().StringP("domain", "", "", "Domain to resolve")

	// Add the commands to the root command
	rootCmd.AddCommand(deployCmd, resolveCmd, balanceCommand, allowanceCommand, depositCommand)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
