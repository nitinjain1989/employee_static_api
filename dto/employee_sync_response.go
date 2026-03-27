package dto

type SyncResponse struct {
	Employees  []Employee `json:"employees"`
	NextCursor Cursor     `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
}
