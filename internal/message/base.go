package message

type ListQueryBase struct {
	PageNo   int    `form:"pageNo"   json:"pageNo"   binding:"gte=1" example:"1"`          // 当前页码
	PageSize int    `form:"pageSize" json:"pageSize" binding:"gte=1,lte=100" example:"20"` // 每页数量
	OrderBy  string `form:"orderBy"  json:"orderBy"  binding:"max=4" example:"desc"`       // 排序方式[desc, asc]
}

type ListResponseBase struct {
	ListQueryBase
	Total int         `json:"total" example:"20"` // 总数
	Data  interface{} `json:"data,omitempty"`
}
