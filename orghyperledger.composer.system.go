package main

import "time"

type Asset struct {
}

type Concept struct {
}

type Participant struct {
}

type Event struct {
	EventID   string    `json:"eventId"`
	Timestamp time.Time `json:"timestamp"`
}
type Registry struct {
	Asset
	RegistryID string `json:"registryId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	System     bool   `json:"system"`
}
type AssetRegistry struct {
	Registry
}
type ParticipantRegistry struct {
	Registry
}
type TransactionRegistry struct {
	Registry
}
type Network struct {
	Asset
	NetworkID      string `json:"networkId"`
	RuntimeVersion string `json:"runtimeVersion"`
}
type NetworkAdmin struct {
	Participant
	ParticipantID string `json:"participantId"`
}
type HistorianRecord struct {
	Asset
	TransactionID        string    `json:"transactionId"`
	TransactionType      string    `json:"transactionType"`
	EventsEmitted        []Event   `json:"eventsEmitted"`
	TransactionTimestamp time.Time `json:"transactionTimestamp"`
}

type IdentityState int

const (
	ISSUED IdentityState = 1 + iota
	BOUND
	ACTIVATED
	REVOKED
)

type Identity struct {
	Asset
	IdentityID  string        `json:"identityId"`
	Name        string        `json:"name"`
	Issuer      string        `json:"issuer"`
	Certificate string        `json:"certificate"`
	State       IdentityState `json:"state"`
}
