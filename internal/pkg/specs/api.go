package specs

// MessageResponse represents a JSON response message as output.
type MessageResponse struct {
	Message string `json:"message"`
}

// MessageResponseWithID represents a JSON response message and ID as output.
type MessageResponseWithID struct {
	Message   string `json:"message"`
	ProfileID int    `json:"profile_id"`
}
