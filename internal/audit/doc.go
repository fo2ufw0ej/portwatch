// Package audit provides structured audit logging for portwatch.
//
// Events are written in either human-readable text or JSON format.
// Each event carries a timestamp, kind, optional port number, and message.
//
// Usage:
//
//	f, _ := os.OpenFile("audit.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
//	logger := audit.New(f, "json")
//	logger.Log(audit.EventPortOpened, "new port detected", 8080)
package audit
