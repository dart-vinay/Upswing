# Retrial Library Documentation

## 1. Project Overview

The Retrial library provides a flexible mechanism for retrying functions in Go. It allows you to control the retry logic using different strategies, making it easier to handle transient errors in your applications. The library is designed to be simple to integrate and use.

The core feature of the library is to determine whether a retry should be allowed based on a configurable strategy and a maximum number of retries. It currently supports two strategies:

*   **Linear Backoff**: Retries after a fixed time interval.
*   **Exponential Backoff**: Retries after an exponentially increasing time interval.

## 2. Architecture

The library is built around a few key interfaces and structs that work together to provide the retry functionality.

### 2.1. `RetrialInterface`

This is the main interface of the library. It defines the operations for managing a retry session.

```go
type RetrialInterface interface {
	Close() error
	AllowRetry() (bool, error)
}
```

*   `AllowRetry() (bool, error)`: This method is called to check if a retry is permitted. It returns `true` if a retry should be attempted, and `false` otherwise. It may also return an error if the retry session is invalid or if the maximum number of retries has been exceeded.
*   `Close() error`: This method closes the retry session, cleaning up any associated resources.

### 2.2. `RetrialRequest`

This struct is the default implementation of the `RetrialInterface`.

```go
type RetrialRequest struct {
	MaxRetries           int
	SessionIdentifier    string // unique key
	RetrialStrategy      StrategyInterface
	LastRequestTimestamp *time.Time
	mutex                *sync.Mutex
}
```

*   `MaxRetries`: The maximum number of retries allowed.
*   `SessionIdentifier`: A unique identifier for the retry session.
*   `RetrialStrategy`: The strategy to be used for determining retry intervals.
*   `LastRequestTimestamp`: The timestamp of the last retry attempt.
*   `mutex`: A mutex to ensure thread safety.

### 2.3. `StrategyInterface`

This interface defines the contract for retry strategies.

```go
type StrategyInterface interface {
	SetStrategyParams(val float64)
	GetStrategyType() StrategyIdentifer
	GetParamValue() float64
}
```

*   `SetStrategyParams(val float64)`: Sets the parameters for the strategy (e.g., the retry window for linear backoff or the exponent for exponential backoff).
*   `GetStrategyType() StrategyIdentifer`: Returns the type of the strategy.
*   `GetParamValue() float64`: Returns the parameter value for the strategy.

### 2.4. Strategy Implementations

The library provides two implementations of the `StrategyInterface`:

*   **`LinearStrategy`**: Implements a linear backoff strategy. It retries after a fixed time window.
    ```go
    type LinearStrategy struct {
        Strategy
        RetryWindow float64
    }
    ```
*   **`ExponentStrategy`**: Implements an exponential backoff strategy. The retry interval increases exponentially with each attempt.
    ```go
    type ExponentStrategy struct {
        Strategy
        Exponent float64
    }
    ```

## 3. Usage

### 3.1. Integrating the Library

To use the Retrial library in your project, you need to create a `RetrialRequest`. Here's a typical usage pattern:

1.  **Choose a strategy**: Decide whether you want to use `LinearStrategy` or `ExponentStrategy`.
2.  **Create a `RetrialRequest`**: Call `RetrialUpswing.CreateNewRetryRequest`, passing in the chosen strategy and the maximum number of retries.
3.  **Use a loop for retries**: In your code, create a loop that calls `retryRequest.AllowRetry()` to check if you should proceed with an operation.
4.  **Close the session**: Once you are done, call `retryRequest.Close()` to clean up the session.

Here is a code example:

```go
package main

import (
	"Upswing/RetrialUpswing"
	"errors"
	"log"
)

func main() {
	// 1. Choose and configure a strategy
	retrialStrategy := new(RetrialUpswing.LinearStrategy)
	retrialStrategy.SetStrategyParams(2) // 2-second window

	// 2. Create a new retry request
	retryRequest, err := RetrialUpswing.CreateNewRetryRequest(retrialStrategy, 5) // Max 5 retries
	if err != nil {
		log.Fatalf("Error creating new retry request: %v", err)
	}
	defer retryRequest.Close()

	// 3. Use a loop for retries
	for {
		retry, err := retryRequest.AllowRetry()
		if err != nil {
			if err == RetrialUpswing.ErrMaxRetryExceeds {
				log.Println("Max retries exceeded!")
				break
			}
			log.Printf("An error occurred: %v", err)
			break
		}

		if retry {
			log.Println("Attempting operation...")
			if err := SomeFallibleOperation(); err != nil {
				log.Printf("Operation failed: %v. Retrying...", err)
			} else {
				log.Println("Operation succeeded!")
				break
			}
		} else {
            // Not allowed to retry yet, wait before next AllowRetry call
            // Note: The current implementation of AllowRetry is blocking,
            // so this 'else' block may not be hit as expected.
            // A delay might be needed here in a real application.
		}
	}
}

// An example of an operation that might fail
func SomeFallibleOperation() error {
	return errors.New("something went wrong")
}
```

### 3.2. Using the `main` Executable

The repository includes a `main` executable that demonstrates the library's functionality. You can run it from the command line with the following flags:

*   `-strategy`: The retry strategy to use. `0` for Linear, `1` for Exponential. (Default: `0`)
*   `-param`: The parameter for the strategy. For Linear, this is the time window in seconds. For Exponential, this is the base of the exponent. (Default: `2`)
*   `-max-retries`: The maximum number of retries. (Default: `5`)

**Example:**

To run with an exponential backoff strategy with an exponent of 3 and a maximum of 4 retries:

```bash
./main -strategy 1 -param 3 -max-retries 4
```

## 4. Error Handling

The library defines several custom errors to indicate specific failure conditions:

*   `ErrBadRequest`: A general error for bad requests (not currently used in the library).
*   `ErrMaxRetryExceeds`: Returned by `AllowRetry` when the maximum number of retries has been reached.
*   `ErrRetryWindowAbsent`: Returned by `AllowRetry` if the linear strategy is used with a retry window of 0 or less.
*   `ErrInvalidRetryStrategy`: Returned if an unknown or unsupported strategy is used.
*   `ErrInvalidRetrialSession`: Returned by `AllowRetry` if the session is invalid or has been closed.
