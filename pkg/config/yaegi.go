package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const (
	FilenameFunctionIdentifier  = "config.Filename"
	ConfigureFunctionIdentifier = "config.Configure"
)

const initialConfigFileContent = `package config

func Filename(
	ip string,
	macAddress string,
	arch string,
	archID int,
) string {
	switch arch {
	case "x64 UEFI":
		return "ipxe.efi"
	default:
		return "undionly.kpxe"
	}
}

func Configure() map[string]string {
	return map[string]string{
		"useStdlib": "false",
	}
}
`

func GetFileName(
	configFileLocation string,
	ip string,
	macAddress string,
	archID int,
	pure bool,
) (string, error) {
	// Read the config file (we are re-reading each time so that a server restart is unnecessary)
	src, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		return "", err
	}

	// Configure the interpreter
	useStdlib := false
	{
		// Start the interpreter (for configuration)
		i := interp.New(interp.Options{})
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
	}

	// Manually prevent stdlib use if set to pure
	if pure {
		useStdlib = false
	}

	// Start the interpreter (for file name)
	e := interp.New(interp.Options{})
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
	return getFileName(
		ip,
		macAddress,
		GetNameForArchId(archID),
		archID,
	), nil
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
	}

	return nil
}
