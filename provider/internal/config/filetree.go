package config

import (
	files_helpers "github.com/matzapata/ipfs-compute/provider/pkg/helpers/files"
)

var (
	SRC_SPEC_FILE       = files_helpers.BuildCwdPath("deployment.json")
	SRC_BIN_FILE        = files_helpers.BuildCwdPath("main")
	SRC_PUBLIC_DIR      = files_helpers.BuildCwdPath("public")
	DIST_DIR            = files_helpers.BuildCwdPath("dist")
	DIST_SPEC_FILE      = files_helpers.BuildCwdPath("dist/deployment.json")
	DIST_BIN_FILE       = files_helpers.BuildCwdPath("dist/deployment/main")
	DIST_PUBLIC_DIR     = files_helpers.BuildCwdPath("dist/deployment/public")
	DIST_ZIP_FILE       = files_helpers.BuildCwdPath("dist/deployment.zip")
	DIST_DEPLOYMENT_DIR = files_helpers.BuildCwdPath("dist/deployment")
)

const TERMS_AND_CONDITIONS = "Deploy to IPFS Compute. Code will be runnable by everyone with the CID."
