package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/matzapata/ipfs-compute/cli/helpers"

	cp "github.com/otiai10/copy"
	"github.com/wabarc/ipfs-pinner/pkg/pinata"
)

type DeploymentSpecification struct {
	Env []string `json:"env"`
}

type Deployment struct {
	DeploymentSpecification
	Owner          string `json:"owner"`
	OwnerSignature string `json:"owner_signature"`
	DeploymentCid  string `json:"deployment_cid"`
}

const TERMS_AND_CONDITIONS = "Deploy to IPFS Compute. Code will be runnable by everyone with the CID."
const IPFS_COMPUTE_PUBLIC_KEY = `-----BEGIN RSA PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAsvdSskRMJkBMNcBjPZU8
d8/PUTejbXOorzmipBC42RBEvBveLzBa76m7QOdiIlMIGsDWNXirE+5oYNLvw0XF
nRVcC8ahrm+FKvVsLfGCREF8Yc0WN6Gv11VtaEEMeY0CUsh1u8ItTT7ePQgXltzo
9qzfNW1DWxhB8YBi/W1Zy0R2k75aiBiAwuLbEnMnOviap9IbKGJaG/ZzNbxWOEfJ
2bZa1fva/xFYjh/pj03b5WaVcj1GK797twuU/LcmrUsaNjQLrBtbMNQonsfnIaO3
cS7uUv4KbSN5Y5H3hryY43kdaOmBIk//MZwGD/l3niAXi4F3q21coDomm8SSE5V7
b2+P6QOEnZT2AToPOCI4dYssV6FwZrimJIEAoupsyTRO5t7xN1Rxob/XeyX753Ib
zhxNJrOjxEN87GQluz6SRgcDfyRwttMGQWPuy9SViKVOuH6UNtOP+iUau40IBKeZ
4UBnvO1pu0V6B9sWUAS1lU8Fqz61DlDzIK/PvOExX7ClEwnWHjMMKHBTFU6MYiwZ
E7CV/AT45P3b34WJ6roC440bmUX1tkm6n2+kS92QM0noeGb7UzGidcYFvab9MrgL
axG80A+kA5yo5MhzdD2JAfgWMrPjPAvpzk5uhcKRq9Df9lJ4Og97eusAnbIAcOP7
jAlff9+PkjHViXkVH6M/HN8CAwEAAQ==
-----END RSA PUBLIC KEY-----`

// Paths
var (
	SRC_SPEC_FILE       = helpers.BuildCwdPath("deployment.json")
	SRC_BIN_FILE        = helpers.BuildCwdPath("main")
	SRC_PUBLIC_DIR      = helpers.BuildCwdPath("public")
	DIST_DIR            = helpers.BuildCwdPath("dist")
	DIST_SPEC_FILE      = helpers.BuildCwdPath("dist/deployment.json")
	DIST_BIN_FILE       = helpers.BuildCwdPath("dist/deployment/main")
	DIST_PUBLIC_DIR     = helpers.BuildCwdPath("dist/deployment/public")
	DIST_ZIP_FILE       = helpers.BuildCwdPath("dist/deployment.zip")
	DIST_DEPLOYMENT_DIR = helpers.BuildCwdPath("dist/deployment")
)

var pnt = pinata.Pinata{Apikey: os.Getenv("PINATA_API_KEY"), Secret: os.Getenv("PINATA_SECRET")}

func main() {
	// confirm deployment
	fmt.Println("IPFS Compute Deployment")
	fmt.Println(TERMS_AND_CONDITIONS)
	if !confirm("Do you want to deploy? (y/n): ") {
		return // exit
	}

	// create signature
	if len(os.Args) < 2 {
		log.Fatal("Usage: ipfs-compute <private_key>")
	}
	privateKey, err := crypto.HexToECDSA(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	signature, err := helpers.SignMessage([]byte(TERMS_AND_CONDITIONS), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// zip deployment files and pin it to ipfs
	err = buildDeploymentZip()
	if err != nil {
		log.Fatal(err)
	}
	cidDeploymentZip, err := pnt.PinFile(DIST_ZIP_FILE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deployment CID:", cidDeploymentZip)

	// build deployment specification and pin it to ipfs
	err = buildDeploymentSpecification(cidDeploymentZip, signature)
	if err != nil {
		log.Fatal(err)
	}
	cidDeploymentSpec, err := pnt.PinFile(DIST_SPEC_FILE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deployment Spec CID:", cidDeploymentSpec)
}

func buildDeploymentZip() error {
	err := os.Mkdir(DIST_DIR, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// copy binary and public folder to dist folder
	err = os.Mkdir(DIST_DEPLOYMENT_DIR, 0755)
	if err != nil {
		return err
	}
	err = cp.Copy(SRC_BIN_FILE, DIST_BIN_FILE)
	if err != nil {
		return err
	}
	err = cp.Copy(SRC_PUBLIC_DIR, DIST_PUBLIC_DIR)
	if err != nil {
		return err
	}

	// zip dist folder
	err = helpers.ZipFolder(DIST_DEPLOYMENT_DIR, DIST_ZIP_FILE)
	if err != nil {
		return err
	}

	return nil
}

func buildDeploymentSpecification(deploymentZipCid string, signature *helpers.Signature) error {
	// load public key
	ipfsComputePublicKey, err := helpers.LoadPublicKeyFromString(IPFS_COMPUTE_PUBLIC_KEY)
	if err != nil {
		log.Fatal(err)
	}

	// read spec json file
	deploymentSpecJson, err := os.ReadFile(SRC_SPEC_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var deploymentSpec DeploymentSpecification
	err = json.Unmarshal(deploymentSpecJson, &deploymentSpec)
	if err != nil {
		log.Fatal(err)
	}

	// add signature to json
	deployment := Deployment{
		DeploymentSpecification: deploymentSpec,
		Owner:                   signature.Address,
		OwnerSignature:          signature.Signature,
		DeploymentCid:           deploymentZipCid,
	}

	// encrypt it with public key
	deploymentJson, err := json.Marshal(deployment)
	if err != nil {
		return err
	}
	encDeploymentJson, err := helpers.EncryptBytes(ipfsComputePublicKey, deploymentJson)
	if err != nil {
		return err
	}

	// save it to a file
	err = os.WriteFile(DIST_SPEC_FILE, encDeploymentJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Trim whitespace and convert to lower case
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}
}
