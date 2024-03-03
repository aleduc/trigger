package call

type ID string

type Meta struct {
	PhoneNumber    string
	VirtualAgentID string
	ID             ID
}

type Body struct {
	PhoneNumber    string `json:"phone_number"`
	VirtualAgentID string `json:"virtual_agent_id"`
}
