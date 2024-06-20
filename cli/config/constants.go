package config

import (
	"github.com/matzapata/ipfs-compute/cli/helpers"
)

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

const TERMS_AND_CONDITIONS = "Deploy to IPFS Compute. Code will be runnable by everyone with the CID."

const REGISTRY_ADDRESS = "0xdb42A86B1bfe04E75B2A5F2bF7a3BBB52D7FFD2F"
const ESCROW_ADDRESS = "0x"
