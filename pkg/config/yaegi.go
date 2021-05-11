package config

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/codeclysm/extract/v3"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const (
	FilenameFunctionIdentifier  = "config.Filename"
	ConfigureFunctionIdentifier = "config.Configure"
)

const initialConfigFileContent = `package config

import "log"

func Filename(
	ip string,
	macAddress string,
	arch string,
	archID int,
) string {
	log.Println("You did not set up boot files yet!")

	return "changeme"
}

func Configure() map[string]string {
	return map[string]string{
		"useStdlib": "true",
	}
}
`

func GetFileName(
	configFileLocation string,
	ip string,
	macAddress string,
	arch string,
	archID int,
	pure bool,
	handleOutput func(string),
) (string, error) {
	// Read the config file (we are re-reading each time so that a server restart is unnecessary)
	src, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		return "", err
	}

	// Configure the interpreter
	useStdlib := false
	{
		// Setup stdout/stderr handling
		outputReader, outputWriter, err := os.Pipe()
		if err != nil {
			return "", err
		}

		// Start the interpreter (for configuration)
		i := interp.New(interp.Options{
			Stdout: outputWriter,
			Stderr: outputWriter,
		})
		i.Use(stdlib.Symbols)

		// "Run" the config file, exporting the config function identifier
		if _, err := i.Eval(string(src)); err != nil {
			return "", err
		}

		// Get the config function by it's identifier
		v, err := i.Eval(ConfigureFunctionIdentifier)
		if err != nil {
			return "", err
		}

		// Cast the function
		configure, ok := v.Interface().(func() map[string]string)
		if !ok {
			return "", errors.New("could not parse config function: invalid config function signature")
		}

		// Run the function
		configParameters := configure()
		for key, value := range configParameters {
			if key == "useStdlib" && value == "true" {
				useStdlib = true
			}
		}

		// Close the output pipe
		if err := outputWriter.Close(); err != nil {
			return "", err
		}

		// Read & handle output
		out, err := ioutil.ReadAll(outputReader)
		if err != nil {
			return "", err
		}

		handleOutput(string(out))
	}

	// Setup stdout/stderr handling
	outputReader, outputWriter, err := os.Pipe()
	if err != nil {
		return "", err
	}

	// Manually prevent stdlib use if set to pure
	if pure {
		useStdlib = false
	}

	// Start the interpreter (for file name)
	e := interp.New(interp.Options{
		Stdout: outputWriter,
		Stderr: outputWriter,
	})
	if useStdlib {
		e.Use(stdlib.Symbols)
	}

	// "Run" the config file, exporting the file name function identifier
	if _, err := e.Eval(string(src)); err != nil {
		return "", err
	}

	// Get the file name function by it's identifier
	w, err := e.Eval(FilenameFunctionIdentifier)
	if err != nil {
		return "", err
	}

	// Cast the function
	getFileName, ok := w.Interface().(func(
		ip string,
		macAddress string,
		arch string,
		archID int,
	) string)
	if !ok {
		return "", errors.New("could not parse file name function: invalid file name function signature")
	}

	// Run the function
	rv, err := getFileName(
		ip,
		macAddress,
		arch,
		archID,
	), nil

	// Close the output pipe
	if err := outputWriter.Close(); err != nil {
		return "", err
	}

	// Read & handle output
	out, err := ioutil.ReadAll(outputReader)
	if err != nil {
		return "", err
	}

	handleOutput(string(out))

	return rv, err
}

func CreateConfigIfNotExists(configFileLocation string) error {
	// If config file does not exist, create and write to it
	if _, err := os.Stat(configFileLocation); os.IsNotExist(err) {
		// Create leading directories
		leadingDir, _ := filepath.Split(configFileLocation)
		if err := os.MkdirAll(leadingDir, os.ModePerm); err != nil {
			return err
		}

		// Create file
		out, err := os.Create(configFileLocation)
		if err != nil {
			return err
		}
		defer out.Close()

		// Write to file
		if err := ioutil.WriteFile(configFileLocation, []byte(initialConfigFileContent), os.ModePerm); err != nil {
			return err
		}

		return nil
	}

	return nil
}

func GetStarterIfNotExists(configFileLocation string, starterURL string, outDir string) error {
	// If config file does not exist, get and extract starter
	if _, err := os.Stat(configFileLocation); os.IsNotExist(err) {
		// Create directory to extract to
		if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
			return err
		}

		// Download .tar.gz
		resp, err := http.Get(starterURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Extract .tar.gz
		return extract.Gz(context.Background(), resp.Body, outDir, nil)
	}

	return nil
}
