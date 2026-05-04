// Package keyconv provides decoding and conversion for Amnezia VPN provisioning files.
package keyconv

import "fmt"

// ConnectionParams stores all parameters for generating the VLESS link.
type ConnectionParams struct {
	Address, UserID, Flow, Encryption string
	Port                              int
	Security, Network                 string
	SNI, Fingerprint, PublicKey       string
	ShortID, SpiderX                  string
}

// buildProfileLink assembles the final profile string.
func buildProfileLink(client string, p ConnectionParams) string {
	return fmt.Sprintf(
		"vless://%s@%s:%d?security=%s&sni=%s&fp=%s&pbk=%s&sid=%s&spx=%s&type=%s&flow=%s&encryption=%s#%s",
		p.UserID, p.Address, p.Port,
		p.Security, p.SNI, p.Fingerprint, p.PublicKey, p.ShortID, p.SpiderX,
		p.Network, p.Flow, p.Encryption, client,
	)
}
