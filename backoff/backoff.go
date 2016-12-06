package backoff

import (
	"math"
	"time"
)

// Strategy has a means to calculate the duration of a wait period based on a given number of attempts
type Strategy interface {
	// WaitDuration calculates a wait period based on the given attempt.
	// In practice, no wait period is used on the 0th attempt.
	Duration(attempt int) time.Duration
}

// ConstantStrategy always waits a set amount of time
type constantStrategy time.Duration

// Constant returns a backoff strategy that always waits same duration
func Constant(wait time.Duration) Strategy {
	return constantStrategy(wait)
}

// Duration returns the predetermined constant duration
func (c constantStrategy) Duration(attempt int) time.Duration {
	return time.Duration(c)
}

// ExponentialBackoffStrategy increases the wait period exponentially using the formula:
// 	min(max, initial * factor^attempt)
type exponentialStrategy struct {
	// initial is the time to wait after the first failed attempt
	initial float64

	// max is the hard ceiling for the calculated wait
	max float64

	// factor is the base of the exponent to exponentially increase the wait period
	factor float64
}

// NewExponentialBackoffStrategy creates a Strategy whose WaitDuration increases using the formula:
// 	min(max, initial * factor^attempt)
func Exponential(initial, max time.Duration, factor float64) Strategy {
	return &exponentialStrategy{
		initial: float64(initial),
		max:     float64(max),
		factor:  factor,
	}
}

// Duration returns the expontentially increasing duration
func (c *exponentialStrategy) Duration(attempt int) time.Duration {
	wait := math.Pow(c.factor, float64(attempt-1))
	wait = math.Min(wait*c.initial, c.max)
	return time.Duration(wait)
}
