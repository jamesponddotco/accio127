// Package endpoint contains the endpoint paths for each handler.
package endpoint

import "git.sr.ht/~jamesponddotco/accio127/internal/build"

const (
	// Slash is the endpoint for the root handler.
	Slash string = "/"

	// IP is the endpoint for the IP handler.
	IP string = Slash + build.APIVersion + "/ip"

	// IPAnonymized is the endpoint for the IPAnonymize handler.
	IPAnonymize string = Slash + build.APIVersion + "/ip/anonymized"

	// IPHashed is the endpoint for the IPHashed handler.
	IPHashed string = Slash + build.APIVersion + "/ip/hashed"

	// Metrics is the endpoint for the Metrics handler.
	Metrics string = Slash + build.APIVersion + "/metrics"

	// Status is the endpoint for the Status handler.
	Status string = Slash + build.APIVersion + "/status"

	// Ping is the endpoint for the Heartbeat handler.
	Ping string = Slash + build.APIVersion + "/ping"
)
