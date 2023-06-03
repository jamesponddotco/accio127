package model

// Status represents the health status of the service.
type Status struct {
	// Server is the server status.
	Server string `json:"server"`

	// Database is the database status.
	Database string `json:"database"`
}

// NewStatus creates a new Status.
func NewStatus(server, database string) *Status {
	return &Status{
		Server:   server,
		Database: database,
	}
}
