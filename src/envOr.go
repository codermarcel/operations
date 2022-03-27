package src

import "os"

func envOr(envName string, defaultValue string) string {
	v, found := os.LookupEnv(envName)

	if !found {
		return defaultValue
	}

	return v
}
