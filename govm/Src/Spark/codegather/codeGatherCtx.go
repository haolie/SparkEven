package codegather

import (
	"SparkEven/govm/Src/Model"
)

type codeGatherCtx struct {
	// 日期
	dateStr string

	// 总数（最终数量可能小于 totalCount
	totalCount int32

	// 当前采集页
	pageIndex int32

	// 总页数
	pageCount int32

	// 数据列
	clMap map[string]*Col

	// 已采集记录
	faceMap map[int]struct{}

	// 采集结果
	faceList []*Model.CodeFace
}
