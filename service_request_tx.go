package main

// ServiceRequest containing all main fields for organization request
type InvCreateUpdateServiceRequestTransaction struct {
	InvIdentity        string  `json:"inv_identity" valid:"required"` //Key
	InvAmount          float64 `json:"inv_amount" valid:"required"`
	InvRemainingAmount float64 `json:"inv_remaining_amount" valid:"required"`
}
