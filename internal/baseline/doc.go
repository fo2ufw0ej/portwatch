// Package baseline provides storage and retrieval of the trusted port
// baseline for portwatch.
//
// The baseline represents the set of ports that are expected to be open on
// the monitored host. When portwatch detects a port that is not in the
// baseline, or a baseline port that has closed, an alert is raised.
//
// Usage:
//
//	store := baseline.New("/var/lib/portwatch/baseline.json")
//
//	// Capture current open ports as the trusted baseline
//	if err := store.Save(openPorts); err != nil {
//		log.Fatal(err)
//	}
//
//	// Later, load the baseline for comparison
//	entry, err := store.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Trusted ports:", entry.Ports)
package baseline
