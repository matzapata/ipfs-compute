package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/matzapata/ipfs-compute/cli/commands"
	"github.com/spf13/cobra"
)

func main() {
	err := godotenv.Load()
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
			privateKey := os.Getenv("PRIVATE_KEY")
			pinataApiKey := os.Getenv("PINATA_APIKEY")
			pinataSecret := os.Getenv("PINATA_SECRET")
			rpc := os.Getenv("RPC")

			// get variables from flags
			provider, _ := cmd.Flags().GetString("provider")

			commands.DeployCommand(privateKey, provider, pinataApiKey, pinataSecret, rpc)
		},
	}
	deployCmd.Flags().StringP("provider", "", "", "IPFS Compute provider domain")

	// Allowance command

	// Deposit command

	// Withdraw command

	// Balance command

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
	rootCmd.AddCommand(deployCmd, resolveCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
