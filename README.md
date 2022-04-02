# Retrial Library

## About

The library provides an interface to allow the retrial of your function using the specified strategy. These are the two strategy we are using:
* **Exponential BackOff**: Contains DB connection related information
* **Linear BackOff**: Retry the function after the specified time again 

## Usage
The program contains an executable file called main. Run the file with the follwing allowed flags :
* **strategy**: 0 = Linear. 1 = Exponential (by default we would use linear)
* **param**: Integer value. Would be treated as window in case of linear and exponent in case of exponential strategy (by default we would use 2)
* **max-retries**: Max Retries we want (by default we would use 5 )

Example: The below command retires the function using exponential strategy with backing off exponent as 3 and max retry count as 4. 
```bash
./main -strategy 1 -param 3 -max-retries 4
```

## Explanation
The library provides `CreateNewRetryRequest` function which return a `RetryInterface` whose structure looks like :

```bash
type RetrialInterface interface {
	Close() error
	AllowRetry() (bool, error)
}

// RetrialRequest implements the RetrialInterface
type RetrialRequest struct {
	MaxRetries           int
	SessionIdentifier    string // unique key
	RetrialStrategy      StrategyInterface
	LastRequestTimestamp time.Time
	mutex                *sync.Mutex
}
```
Our interface allows two functionality 
* **AllowRetry**: This is called everytime we want to check the retriability of our function
* **Close**: This function closes our RetrialRequest once the requirement is completed

Moreover, there are two strategies, ExponentialStrategy and LinearStrategy which implements the StrategyInterface: 
```bash
type StrategyInterface interface {
	SetStrategyParams(val float64)
	GetStrategyType() StrategyIdentifer
	GetParamValue() float64
}

// Parent Strategy
type Strategy struct {
	Name StrategyIdentifer
}

// Using Object Composition for reusability
type ExponentStrategy struct {
	Strategy
	Exponent float64
}

// Using Object Composition for reusability
type LinearStrategy struct {
	Strategy
	RetryWindow float64
}
```
