//go:build mage

package main

import (
	"github.com/joho/godotenv"
	"github.com/magefile/mage/sh"
)

var profile = "ots"

func Apply() error {
	godotenv.Load()
	return sh.RunV("aws-vault", "exec", "--no-session", profile, "--", "terraform", "apply")
}

func Applya() error {
	godotenv.Load()
	return sh.RunV("aws-vault", "exec", "--no-session", profile, "--", "terraform", "apply", "-auto-approve")
}

func Plan() error {
	godotenv.Load()
	return sh.RunV("aws-vault", "exec", "--no-session", profile, "--", "terraform", "plan")
}

func Init() error {
	godotenv.Load()

	return sh.RunV("aws-vault", "exec", "--no-session", profile, "--", "terraform", "init")
}

func Initu() error {
	godotenv.Load()

	return sh.RunV("aws-vault", "exec", "--no-session", profile, "--", "terraform", "init", "-upgrade")
}
