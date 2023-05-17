package model

import "net"

// Address represents an IP address.
type IP struct {
	V4 string `json:"ipv4"`
	V6 string `json:"ipv6"`
}

// NewIP creates a new IP address.
func NewIP(ip string) *IP {
	address := net.ParseIP(ip)
	if address == nil {
		return nil
	}

	switch len(address) {
	case net.IPv4len:
		return &IP{
			V4: address.String(),
		}
	case net.IPv6len:
		return &IP{
			V6: address.String(),
		}
	default:
		return nil
	}
}
