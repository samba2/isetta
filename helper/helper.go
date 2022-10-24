package helper

import (
	"errors"
	"fmt"
	"time"

	log "org.samba/isetta/simplelogger"
)

type RetryParams struct {
	Description string
	Attempts    int
	Sleep       time.Duration // time between function execution, doubles with each retry 
	Func        func() bool   // function which should be retried until success
}

func Retry(p RetryParams) error {
	for i := 0; i < p.Attempts; i++ {
		if i > 0 {
			log.Logger.Trace("%v: Trying %vst time, backing off for %v", p.Description, i, p.Sleep)
			time.Sleep(p.Sleep)
			p.Sleep *= 2
		}

		isSuccessful := p.Func()
		if isSuccessful {
			log.Logger.Debug("%v: Success", p.Description)
			return nil
		}
	}
	errMsg := fmt.Sprintf("%v: failed after %d attempts", p.Description, p.Attempts)
	return errors.New(errMsg)
}


func AssertNoError(err error, format string, v ...any) {
    if err != nil {
        errMsg := fmt.Sprintf(format, v...)
        newErr := fmt.Errorf("%v, error was: %w", errMsg, err)
		log.Logger.Error2(newErr)
	}
}

func AssertNoError2(err error) {
    if err != nil {
		log.Logger.Error2(err)
	}
}