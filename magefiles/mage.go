//go:build mage

package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func RunAir() error {
	godotenv.Load()
	return sh.RunV("aws-vault", "exec", "--no-session", "ots", "--", "air")
}

const lambdaPrefix = "./cmd/lambda"

func Build() error {
	entries, err := os.ReadDir(lambdaPrefix)
	if err != nil {
		return err
	}
	for _, entry := range entries {

		scripts := []string{
			"CGO_ENABLED=0 go build -o ./bin/{{.lambda}}/bootstrap {{.prefix}}/{{.lambda}}",
		}

		err := executeScript(scripts, map[string]string{
			"lambda": entry.Name(),
			"prefix": lambdaPrefix,
		})
		if err != nil {
			return err
		}

	}

	// // using the os package from the std lib, list all files with a file extension of .go in the cmd/photos directory and append them to an array
	// entries, _ := os.ReadDir(fmt.Sprintf("./cmd/%s", lambda))
	// var files = []string{
	// 	"build", "-o", fmt.Sprintf("./bin/%s/bootstrap", lambda),
	// }
	// for _, entry := range entries {
	// 	if filepath.Ext(entry.Name()) == ".go" {
	// 		files = append(files, filepath.Join(fmt.Sprintf("./cmd/%s", lambda), entry.Name()))
	// 	}
	// }

	// return sh.RunV("go", files...)
	return nil
}

func Run() error {
	mg.Deps(Build)
	godotenv.Load()

	script := []string{
		"aws-vault exec --no-session ots -- go run ./cmd/photos-handler/*.go",
	}

	return executeScript(script, nil)
}

func Tidy() error {
	return sh.RunV("go", "mod", "tidy")
}

func Vendor() error {
	return sh.RunV("go", "mod", "vendor")
}

func Edge() error {
	cwd, _ := os.Getwd()
	scripts := []string{
		"cd ./cmd/edge/photos-edge-validation",
		"zip photos-edge-validation.zip *.js",
		fmt.Sprintf("mv photos-edge-validation.zip %s/terraform/assets", cwd),
	}

	return executeScript(scripts, nil)
}

func Deploy() error {
	mg.Deps(Build)
	entries, err := os.ReadDir(lambdaPrefix)
	if err != nil {
		return err
	}

	var wg = new(sync.WaitGroup)

	for _, entry := range entries {
		wg.Add(1)
		go func(entry fs.DirEntry) {
			defer wg.Done()

			logEntry := logrus.NewEntry(logger).WithField("entry", entry.Name())

			script := []string{
				// cd into the directory
				"cd ./bin/{{.lambda}}",
				// create a zip file that is the same name as the directory and add the bootstrap file to the zip
				"zip -q {{.lambda}}.zip bootstrap",
				// Upload the lambda function code
				"aws-vault exec --no-session ots -- aws lambda update-function-code --function-name {{.lambda}} --zip-file fileb://{{.lambda}}.zip --query 'Handler'",
			}

			// logEntry.Info("Upload Lambda Function Code")

			err := executeScript(script, map[string]string{
				"lambda": entry.Name(),
			})
			if err != nil {
				logEntry.WithError(err).Fatal("failed to upload lambda function code")
				return
			}

			// logEntry.Info("Upload Lambda Function Configuration")
		}(entry)
	}

	wg.Wait()

	// for _, entry := range entries {

	// 	scripts := []string{
	// 		"aws-vault exec --no-session ots -- aws lambda update-function-code --function-name {{.lambda}} --zip-file fileb://./bin/{{.lambda}}.zip",
	// 	}

	// 	err := executeScript(scripts, map[string]string{
	// 		"lambda": entry.Name(),
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}

	// }
	return nil
}

func executeScript(script []string, args any) error {

	cmd := new(bytes.Buffer)
	err := template.Must(template.New("").Parse(strings.Join(script, " && "))).Execute(cmd, args)
	if err != nil {
		return err
	}

	return sh.RunV("/bin/bash", "-c", cmd.String())

}
