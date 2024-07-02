package config

import "github.com/matzapata/ipfs-compute/provider/pkg/system"

var (
	SRC_SPEC_FILE       = system.BuildCwdPath("deployment.json")
	DIST_DIR            = system.BuildCwdPath("dist")
	DIST_SPEC_FILE      = system.BuildCwdPath("dist/deployment.json")
	DIST_BIN_FILE       = system.BuildCwdPath("dist/deployment/main")
	DIST_PUBLIC_DIR     = system.BuildCwdPath("dist/deployment/public")
	DIST_ZIP_FILE       = system.BuildCwdPath("dist/deployment.zip")
	DIST_DEPLOYMENT_DIR = system.BuildCwdPath("dist/deployment")
)

const TERMS_AND_CONDITIONS = "Deploy to IPFS Compute. Code will be runnable by everyone with the CID."
