// Package backoff implements exponential backoff with optional jitter,
// suitable for retrying transient failures such as webhook delivery or
// port scan errors in portwatch.
//
// Usage:
//
//	b := backoff.NewDefault()
//	for {
//		err := doSomething()
//		if err == nil {
//			b.Reset()
//			break
//		}
//		wait := b.Next()
//		time.Sleep(wait)
//	}
package backoff
