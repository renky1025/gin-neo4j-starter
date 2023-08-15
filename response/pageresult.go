package response

// PageResult represents a 响应格式结构
//
// swagger:model
type PageResult struct {
	Total   int64       `json:"total"`
	Records interface{} `json:"records"`
	Size    int64       `json:"size"`
	Current int64       `json:"current"`
}
