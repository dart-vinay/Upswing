package RetrialUpswing

import (
	"github.com/dchest/uniuri"
	"math"
	"sync"
	"time"
)

type StrategyInterface interface {
	SetStrategyParams(val float64)
	GetStrategyType() StrategyIdentifer
	GetParamValue() float64
}

type Strategy struct {
	Name StrategyIdentifer
}

type ExponentStrategy struct {
	Strategy
	Exponent float64
}

func (strategy *ExponentStrategy) SetStrategyParams(exponent float64) {
	strategy.Exponent = exponent
}

func (strategy *ExponentStrategy) GetStrategyType() StrategyIdentifer {
	return ExponentialBackOff
}

func (strategy *ExponentStrategy) GetParamValue() float64 {
	return strategy.Exponent
}

type LinearStrategy struct {
	Strategy
	RetryWindow float64
}

func (strategy *LinearStrategy) SetStrategyParams(window float64) {
	strategy.RetryWindow = window
}

func (strategy *LinearStrategy) GetStrategyType() StrategyIdentifer {
	return LinearBackOff
}

func (strategy *LinearStrategy) GetParamValue() float64 {
	return strategy.RetryWindow
}

var RetryCountMap map[string]int // Map to store the count the number of times the function has been retired

type StrategyIdentifer string

const (
	ExponentialBackOff = StrategyIdentifer("ExponentialBackOff")
	LinearBackOff      = StrategyIdentifer("LinearBackOff")
)

type RetrialInterface interface {
	Close() error
	AllowRetry() (bool, error)
}

type RetrialRequest struct {
	MaxRetries           int
	SessionIdentifier    string // unique key
	RetrialStrategy      StrategyInterface
	LastRequestTimestamp *time.Time
	mutex                *sync.Mutex
}

var (
	StdChars = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
)

func generateNewSessionKey(idLength int) string {
	return uniuri.NewLenChars(idLength, StdChars)
}

func CreateNewRetryRequest(strategy StrategyInterface, maxRetries int) (RetrialInterface, error) {

	InitializeLibraryVariables()

	mutex := new(sync.Mutex)
	sessionIdentifier := generateNewSessionKey(10)
	retryRequest := new(RetrialRequest)

	retryRequest.mutex = mutex
	retryRequest.SessionIdentifier = sessionIdentifier
	retryRequest.RetrialStrategy = strategy
	retryRequest.MaxRetries = maxRetries

	RetryCountMap[retryRequest.SessionIdentifier] = 0
	return retryRequest, nil
}

// Return if we are allowed to retry the request or not
func (retrialRequest *RetrialRequest) AllowRetry() (bool, error) {
	var err error
	currTime := time.Now()

	if _, ok := RetryCountMap[retrialRequest.SessionIdentifier]; !ok {
		return false, ErrInvalidRetrialSession
	}

	if RetryCountMap[retrialRequest.SessionIdentifier] >= retrialRequest.MaxRetries {
		return false, ErrMaxRetryExceeds
	} else if retrialRequest.RetrialStrategy.GetStrategyType() != "" {
		if retrialRequest.RetrialStrategy.GetStrategyType() == ExponentialBackOff {
			exponent := retrialRequest.RetrialStrategy.GetParamValue()

			if (retrialRequest.LastRequestTimestamp != nil && currTime.Sub(*retrialRequest.LastRequestTimestamp).Seconds() <= math.Pow(exponent, float64(RetryCountMap[retrialRequest.SessionIdentifier]))) {
				//retrialRequest.LastRequestTimestamp = currTime
				//RetryCountMap[retrialRequest.SessionIdentifier] += 1
				return false, err
			} else {
				//fmt.Println(currTime.Sub(retrialRequest.LastRequestTimestamp).Seconds())
				//fmt.Println(math.Pow(exponent, float64(RetryCountMap[retrialRequest.SessionIdentifier])))
				retrialRequest.LastRequestTimestamp = &currTime
				RetryCountMap[retrialRequest.SessionIdentifier] += 1
				return true, err
			}
		} else if retrialRequest.RetrialStrategy.GetStrategyType() == LinearBackOff {
			if retrialRequest.RetrialStrategy.GetParamValue() > 0 {
				if (retrialRequest.LastRequestTimestamp != nil &&  currTime.Sub(*retrialRequest.LastRequestTimestamp).Seconds() <= float64(retrialRequest.RetrialStrategy.GetParamValue())) {
					//RetryCountMap[retrialRequest.SessionIdentifier] += 1
					//retrialRequest.LastRequestTimestamp = currTime
					return false, err
				} else {
					retrialRequest.LastRequestTimestamp = &currTime
					RetryCountMap[retrialRequest.SessionIdentifier] += 1
					return true, err
				}
			} else {
				return false, ErrRetryWindowAbsent
			}
		} else {
			return false, ErrInvalidRetryStrategy
		}
	}
	return true, err
}

// TODO : Need to handle race condition while deleting record from RetryCountMap global map
func (retrialRequest *RetrialRequest) Close() error {
	var err error
	retrialRequest.mutex.Lock()
	defer retrialRequest.mutex.Unlock()
	delete(RetryCountMap, retrialRequest.SessionIdentifier)
	return err
}

func InitializeLibraryVariables() {
	if RetryCountMap == nil {
		RetryCountMap = make(map[string]int)
	}
}
