// Package keyconv provides decoding and conversion for Amnezia VPN provisioning files.
package keyconv

import (
	"encoding/json"
	"fmt"
)

// parseInnerConfig extracts and parses the inner configuration from the full JSON.
func parseInnerConfig(fullJSON string) (map[string]any, error) {
	var top map[string]any
	if err := json.Unmarshal([]byte(fullJSON), &top); err != nil {
		return nil, fmt.Errorf("outer json: %w", err)
	}

	containersRaw, ok := top["containers"]
	if !ok {
		return nil, fmt.Errorf("containers key missing")
	}
	containers, ok := containersRaw.([]any)
	if !ok || len(containers) == 0 {
		return nil, fmt.Errorf("containers list empty or invalid")
	}
	lastRaw := containers[len(containers)-1]
	last, ok := lastRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("last container is not an object")
	}

	xrayRaw, ok := last["xray"]
	if !ok {
		return nil, fmt.Errorf("xray section missing")
	}
	xrayBlock, ok := xrayRaw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("xray section is not an object")
	}

	innerStrRaw, ok := xrayBlock["last_config"]
	if !ok {
		return nil, fmt.Errorf("last_config missing in xray")
	}
	innerStr, ok := innerStrRaw.(string)
	if !ok || innerStr == "" {
		return nil, fmt.Errorf("inner config missing or not a string")
	}

	var inner map[string]any
	if err := json.Unmarshal([]byte(innerStr), &inner); err != nil {
		return nil, fmt.Errorf("inner json: %w", err)
	}
	return inner, nil
}

// extractOutbound extracts the first outbound from the inner config.
func extractOutbound(inner map[string]any) (map[string]any, error) {
	outboundsRaw, ok := inner["outbounds"]
	if !ok {
		return nil, fmt.Errorf("outbounds missing")
	}
	outbounds, ok := outboundsRaw.([]any)
	if !ok || len(outbounds) == 0 {
		return nil, fmt.Errorf("no outbounds defined")
	}
	out0Raw := outbounds[0]
	out0, ok := out0Raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("first outbound is not an object")
	}
	return out0, nil
}

// extractServer extracts server details from outbound settings.
func extractServer(out0 map[string]any) (address string, port int, user map[string]any, err error) {
	settingsRaw, ok := out0["settings"]
	if !ok {
		err = fmt.Errorf("settings missing")
		return
	}
	settings, ok := settingsRaw.(map[string]any)
	if !ok {
		err = fmt.Errorf("settings is not an object")
		return
	}

	vnextRaw, ok := settings["vnext"]
	if !ok {
		err = fmt.Errorf("vnext missing")
		return
	}
	vnext, ok := vnextRaw.([]any)
	if !ok || len(vnext) == 0 {
		err = fmt.Errorf("vnext absent or invalid")
		return
	}
	serverRaw := vnext[0]
	server, ok := serverRaw.(map[string]any)
	if !ok {
		err = fmt.Errorf("first vnext entry is not an object")
		return
	}

	if addr, ok := server["address"].(string); ok {
		address = addr
	}
	if portVal, ok := server["port"].(float64); ok {
		port = int(portVal)
	}

	usersRaw, ok := server["users"]
	if !ok {
		err = fmt.Errorf("users missing")
		return
	}
	users, ok := usersRaw.([]any)
	if !ok || len(users) == 0 {
		err = fmt.Errorf("user list empty")
		return
	}
	userRaw := users[0]
	user, ok = userRaw.(map[string]any)
	if !ok {
		err = fmt.Errorf("first user is not an object")
		return
	}
	return
}

// extractStreamSettings extracts stream and reality settings from outbound.
func extractStreamSettings(out0 map[string]any) (security, network string, reality map[string]any, err error) {
	streamRaw, ok := out0["streamSettings"]
	if !ok {
		err = fmt.Errorf("stream settings missing")
		return
	}
	stream, ok := streamRaw.(map[string]any)
	if !ok {
		err = fmt.Errorf("streamSettings is not an object")
		return
	}
	if sec, ok := stream["security"].(string); ok {
		security = sec
	}
	if net, ok := stream["network"].(string); ok {
		network = net
	}

	realityRaw, ok := stream["realitySettings"]
	if !ok {
		err = fmt.Errorf("reality settings missing")
		return
	}
	reality, ok = realityRaw.(map[string]any)
	if !ok {
		err = fmt.Errorf("realitySettings is not an object")
		return
	}
	return
}

// extractParams extracts all necessary values from the inner configuration.
func extractParams(inner map[string]any) (params ConnectionParams, err error) {
	out0, err := extractOutbound(inner)
	if err != nil {
		return
	}

	addr, port, user, err := extractServer(out0)
	if err != nil {
		return
	}
	params.Address = addr
	params.Port = port
	if id, ok := user["id"].(string); ok {
		params.UserID = id
	}
	if flow, ok := user["flow"].(string); ok {
		params.Flow = flow
	}
	if enc, ok := user["encryption"].(string); ok {
		params.Encryption = enc
	}

	security, network, reality, err := extractStreamSettings(out0)
	if err != nil {
		return
	}
	params.Security = security
	params.Network = network

	if sni, ok := reality["serverName"].(string); ok {
		params.SNI = sni
	}
	if fp, ok := reality["fingerprint"].(string); ok {
		params.Fingerprint = fp
	}
	if pk, ok := reality["publicKey"].(string); ok {
		params.PublicKey = pk
	}
	if sid, ok := reality["shortId"].(string); ok {
		params.ShortID = sid
	}

	params.SpiderX = "/"
	if spx, ok := reality["spiderX"].(string); ok && spx != "" {
		params.SpiderX = spx
	}

	return params, nil
}
