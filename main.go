package main

import (
	"Upswing/RetrialUpswing"
	"errors"
	"flag"
	"github.com/labstack/gommon/log"
)

//
const LINEAR = 0
const EXPONENTIAL = 1

func main() {

	strategy := flag.Float64("strategy", LINEAR, "")
	strategyParams := flag.Float64("param", 2, "")
	maxRetriesFlag := flag.Float64("max-retries", 5, "")

	flag.Parse()
	strategyString := "LINEAR"
	if *strategy == LINEAR {
		strategyString = "LINEAR"
	} else if *strategy == EXPONENTIAL {
		strategyString = "EXPONENTIAL"
	} else {
		log.Errorf("Invalid strategy chosen %v", *strategy)
		return
	}
	log.Infof("Strategy %v, Param Value %v, MaxRetries %v", strategyString, *strategyParams, *maxRetriesFlag)

	var retrialStrategy RetrialUpswing.StrategyInterface
	if *strategy == LINEAR {
		retrialStrategy = new(RetrialUpswing.LinearStrategy)
	} else if *strategy == EXPONENTIAL {
		retrialStrategy = new(RetrialUpswing.ExponentStrategy)
	} else {
		log.Errorf("Invalid Strategy Flag %v", strategy)
		return
	}

	retrialStrategy.SetStrategyParams(*strategyParams)
	maxRetries := int(*maxRetriesFlag)

	retryRequest, err := RetrialUpswing.CreateNewRetryRequest(retrialStrategy, maxRetries)
	if err != nil {
		log.Errorf("Error creating new retry request %v", err)
		return
	}
	defer retryRequest.Close()

	for {
		if retry, err := retryRequest.AllowRetry(); retry && err == nil {
			if err2 := BigOperation(); err2 != nil {
				// Log the error
			} else {
				break
			}
		} else if err != nil {
			if err == RetrialUpswing.ErrMaxRetryExceeds {
				log.Infof("Max retries exceeded!!")
				break
			} else if err == RetrialUpswing.ErrInvalidRetrialSession {
				log.Infof("Session Invalidated!")
			}
		} else if !retry {
			//log.Infof("Denied retrial!!")
		}
	}
}

func BigOperation() error {
	log.Info("Retrying Function: BigOperation")
	return errors.New("Self Created Error")
}
