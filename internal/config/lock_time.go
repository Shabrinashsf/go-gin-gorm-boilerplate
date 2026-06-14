package config

import "os"

type ApiLockRange struct {
	Start   string
	End     string
	Message string
}

var (
	ApiLockRanges map[string]ApiLockRange
)

func InitApiLockRanges() {
	if os.Getenv("APP_ENV") == "development" {
		ApiLockRanges = ApiLockRangesDev
	} else {
		ApiLockRanges = ApiLockRangesProd
	}
}

// how to use: "lock-name": {Start: "01-06-2026 10:00:00", End: "30-06-2026 23:59:59", Message: "api closed for maintenance"},

var ApiLockRangesDev = map[string]ApiLockRange{}

var ApiLockRangesProd = map[string]ApiLockRange{}
