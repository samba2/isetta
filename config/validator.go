package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/3th1nk/cidr"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"

	ut "github.com/go-playground/universal-translator"
)

var humanReadableValidationMessages = map[string]string{
	"required": "{0} is missing",
	"url":      "{0}: {1} is an an invalid URL",
	"cidrv4":   "{0}: {1} is not a valid CIDR address",
	"ip4_addr": "{0}: {1} is not a valid IPv4 address",
	"alpha":    "{0}: {1} is not a letters-only string",
}

type MyValidator struct {
	Validate       *validator.Validate
	Trans          ut.Translator
	Config         *Config
	ValidLogLevels []string
}

func NewValidator(conf *Config, validLogLevels []string) MyValidator {
	uni := ut.New(en.New())
	trans, _ := uni.GetTranslator("en")

	myValidator := MyValidator{
		Validate:       validator.New(),
		Trans:          trans,
		Config:         conf,
		ValidLogLevels: validLogLevels,
	}

	myValidator.registerHumanReadableErrorMessages()
	return myValidator
}

func (v MyValidator) registerHumanReadableErrorMessages() {
	for key, value := range humanReadableValidationMessages {
		// extra assignment was needed to ensure tag and msg
		// are accessed as expected by the anonymous functions.
		// still have not fully understood why this is necessary
		tag := key
		msg := value
		v.Validate.RegisterTranslation(tag, v.Trans,
			// map tag to translation aka the human friendly validation message
			func(ut ut.Translator) error {
				return ut.Add(tag, msg, true)
			},
			// provide placeholders with values from fieldError
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(tag, fe.StructNamespace(), fmt.Sprint(fe.Value()))
				return t
			},
		)
	}
}

func (v MyValidator) DoValidate() error {
	err := v.Validate.Struct(v.Config)
	if err != nil {
		return v.buildHumanReadableValidationErrorMessage(err)
	}
	err = v.validateSubnetSize()
	if err != nil {
		return err
	}

	return v.validateLogLevel()
}

func (v MyValidator) validateLogLevel() error {
	logLevel := v.Config.General.LogLevel
	if slices.Contains(v.ValidLogLevels, logLevel) {
		return nil
	} else {		
		return fmt.Errorf("log level '%v' is invalid. Valid log levels: %v", logLevel, strings.Join(v.ValidLogLevels, ", "))
	}
}

func (v MyValidator) validateSubnetSize() error {
	subnet, err := cidr.Parse(v.Config.Network.WslToWindowsSubnet)
	if err != nil {
		return err
	}

	if isSubnetSizeTooSmall(subnet) {
		return errors.New("configured subnet in wsl_to_windows_subnet is too small. Smallest allowed size is a /30 network")
	}

	return nil
}

func isSubnetSizeTooSmall(subnet *cidr.CIDR) bool {
	minimumIpAddressInSubnetCnt := 4 // all IPs of a /30 network. Only 2 IPs are usable for routing
	return subnet.IPCount().Uint64() < uint64(minimumIpAddressInSubnetCnt)
}

func (v MyValidator) buildHumanReadableValidationErrorMessage(err error) error {
	// produce human readable output for first validation error
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for i, e := range errs {
			if i == 0 {
				humanReadableMessage := e.Translate(v.Trans)
				return errors.New(humanReadableMessage)
			}
		}
	}
	return nil
}
