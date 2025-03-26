package models

import (
	"Group03-EX-StudentManagementAppBE/common"
	"strings"
)

type QuerySort struct {
	Origin string
	Sort   string `json:"sort" form:"sort"`
}

func (s QuerySort) Parse() string {
	return strings.ReplaceAll(s.Origin, ".", " ")
}

type QueryParams struct {
	Limit  int
	Offset int
	QuerySort
	Preload  []common.Preload
	Selected []string
}

type BaseRequestParamsUri struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Sort     string `form:"sort"`
}

type BaseListResponse struct {
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Items    interface{} `json:"items"`
	Extra    interface{} `json:"extra"`
}

type FileTypeQueryRequest struct {
	BaseRequestParamsUri
	Type string `form:"type" validate:"omitempty,oneof=csv json"`
}

func (q *FileTypeQueryRequest) GetFileType() string {
	if q.Type == "" {
		return common.FILE_TYPE_CSV
	}
	return q.Type
}

func (q *FileTypeQueryRequest) IsValidFileType() bool {
	return q.Type == "" || q.Type == common.FILE_TYPE_JSON || q.Type == common.FILE_TYPE_CSV
}
