package v1

type StaroidSke struct {
	ID     string `json:"name"`
	Cloud  string `json:"cloud"`
	Region string `json:"region"`
}

type StaroidCluster struct {
	Name  string     `json:"name"`
	Ske   StaroidSke `json:"ske"`
	OrgID int64      `json:"orgId"`
}

type StaroidOrg struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	ID       int64  `json:"id"`
}
