// Package window implements a sliding time-window event counter.
//
// It divides a configurable period into a fixed number of buckets.
// Events are recorded into the current bucket; Total() sums only
// buckets whose timestamp falls within the window period, providing
// an approximate rolling count useful for rate-limiting and alerting.
//
// Example:
//
//	ctr := window.New(time.Minute, 6)
//	ctr.Add(1)
//	fmt.Println(ctr.Total()) // 1
package window
