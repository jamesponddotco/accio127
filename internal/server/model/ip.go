package model

import "net"

// Address represents an IP address.
type IP struct {
	V4 string `json:"ipv4,omitempty"`
	V6 string `json:"ipv6,omitempty"`
}

// NewIP creates a new IP address.
func NewIP(ip string) *IP {
	address := net.ParseIP(ip)

	var (
		ipv4 = address.To4()
		ipv6 = address.To16()
	)

	if ipv4 != nil {
		return &IP{
			V4: ipv4.String(),
		}
	}

	if ipv6 != nil {
		return &IP{
			V6: ipv6.String(),
		}
	}

	return nil
}
