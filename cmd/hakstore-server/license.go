package main

import (
	"fmt"

	licensely "github.com/hakluke/licensely-client-go"
)

func checkLicense() (bool, string) {
	valid, err := licensely.VerifyKey(config.Server.HakstoreLicense)

	if err != nil {
		fmt.Println("Error checking licence key:", err)
		return false, "An error occured while attempting to validate your license key: " + err.Error()
	}
	if valid != true {
		return false, "Your license key is invalid."
	}
	return true, "Licence key validated, thank you for supporting my work!"
}
