package config

const (
	MB                    = 1 << 20 // 1 MB in bytes (1 << 20 is 2^20)
	MAX_ZIPPED_DEPLOYMENT = uint(50 * MB)
	DEPLOYMENT_SPEC_FILE  = "khachapuri.json"
)
