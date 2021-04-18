package config

import (
	"errors"
	"io/ioutil"

	"github.com/traefik/yaegi/interp"
)

func GetFileName(
	configFunctionIdentifier string,
	configFileLocation string,
	ip string,
	macAddress string,
	arch int,
	undi int,
) (string, error) {
	// Start the interpreter
	i := interp.New(interp.Options{})

	// Read the config file (we are re-reading each time so that a server restart is unnecessary)
	src, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		return "", err
	}

	// "Run" the config file, exporting the config function identifier
	if _, err := i.Eval(string(src)); err != nil {
		return "", err
	}

	// Get the config function by it's identifier
	v, err := i.Eval(configFunctionIdentifier)
	if err != nil {
		return "", err
	}

	// Cast the function
	getFileName, ok := v.Interface().(func(
		ip string,
		macAddress string,
		arch int,
		undi int,
	) string)
	if !ok {
		return "", errors.New("could not parse config function: invalid config function signature")
	}

	// Run the function
	return getFileName(
		ip,
		macAddress,
		arch,
		undi,
	), nil
}
