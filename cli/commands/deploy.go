package commands

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matzapata/ipfs-compute/cli/config"
	"github.com/matzapata/ipfs-compute/cli/helpers"
	"github.com/matzapata/ipfs-compute/cli/services"
)

func DeployCommand(privateKey string, provider string, pinataApiKey string, pinataSecret string, rpc string) {
	ethclient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	deploymentService := services.NewDeploymentService()
	ipfsService := services.NewIpfsService(pinataApiKey, pinataSecret)
	registryService := services.NewRegistryService(ethclient, config.REGISTRY_ADDRESS)

	// confirm deployment
	fmt.Println("IPFS Compute Deployment")
	fmt.Println(config.TERMS_AND_CONDITIONS)
	if !helpers.Confirm("Do you want to deploy? (y/n): ") {
		return // exit
	}

	// create signature
	signature, err := helpers.SignMessage([]byte(config.TERMS_AND_CONDITIONS), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// zip deployment files and pin it to ipfs
	err = deploymentService.BuildDeploymentZip()
	if err != nil {
		log.Fatal(err)
	}
	cidDeploymentZip, err := ipfsService.PinFile(config.DIST_ZIP_FILE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deployment CID:", cidDeploymentZip)

	// grab public key from resolver
	resolver, err := registryService.ResolveDomain(provider)
	if err != nil {
		log.Fatal(err)
	}
	providerPublicKey, err := resolver.Pubkey(nil)
	if err != nil {
		log.Fatal(err)
	}

	// build deployment specification and pin it to ipfs
	err = deploymentService.BuildDeploymentSpecification(cidDeploymentZip, signature, providerPublicKey)
	if err != nil {
		log.Fatal(err)
	}
	cidDeploymentSpec, err := ipfsService.PinFile(config.DIST_SPEC_FILE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deployment Spec CID:", cidDeploymentSpec)
}
