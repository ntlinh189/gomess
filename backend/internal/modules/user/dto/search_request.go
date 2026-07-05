package dto

type SearchRequest struct {
	Provider string `form:"provider" binding:"required"`
	Keyword  string `form:"keyword" binding:"required"`
	Skip     int    `form:"skip"`
	Limit    int    `form:"limit"`
}
