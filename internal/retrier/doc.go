// Package retrier provides a simple retry mechanism with exponential backoff.
//
// Use New or NewFromConfig to create a Retrier, then call Do with a context
// and a function to retry. Wrap errors with Permanent to prevent retrying.
//
// Example:
//
//	r, _ := retrier.New(retrier.DefaultConfig())
//	err := r.Do(ctx, func() error {
//		return callExternalService()
//	})
package retrier
