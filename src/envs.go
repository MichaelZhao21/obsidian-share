package src

import (
	"log"
	"os"
)

var reqEnvs = []string{
	"SSH_PRIVATE_KEY",
	"REPO_URL",
	"MONGODB_URI",
}

func CheckEnvs() {
	for _, env := range reqEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("Environment variable %s is required", env)
		}
	}
}
