package contract

type ListQuery struct {
	PageNo   int    `form:"pageNo,default=1" json:"pageNo" binding:"gte=1" example:"1"`
	PageSize int    `form:"pageSize,default=50" json:"pageSize" binding:"gte=1,lte=100" example:"50"`
	OrderBy  string `form:"orderBy,default=desc" json:"orderBy" binding:"max=4" example:"desc"`
}
