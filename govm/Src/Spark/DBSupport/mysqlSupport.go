package DBSupport

import (
	"context"
	"database/sql"
	"fmt"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Interface"
	"SparkEven/govm/Src/Spark/DBSupport/dbFillImpl"
)

var (
	implMap map[string]Interface.IDbModelFill
)

func init() {
	implMap = make(map[string]Interface.IDbModelFill, 2)
	implMap[Def.Model_Key_CodeFace] = &dbFillImpl.CodeFaceImpl{}
}

func InsertList(ctx context.Context, conn *sql.DB, modelList []Interface.IDbModel) error {
	i, count := 0, 300
	for {
		index := i + count
		if index > len(modelList) {
			index = len(modelList)
		}

		array := modelList[i:index]

		var vStr = ""
		for _, item := range array {

			if len(vStr) > 0 {
				vStr += ";"
			}
			vStr += item.CreateInsertStr()

		}

		_, err := conn.Exec(vStr)
		if err != nil {
			return err
		}
		i += count
		if i >= len(modelList) {
			break
		}
	}

	return nil
}

func GetCodeFace(ctx context.Context, conn *sql.DB, code int, date string) (faceList []Interface.IDbModel, errStr string) {

	if date == "" && code < 0 {
		errStr = fmt.Sprintf(" GetCodeFace 缺少查询条件")
		return
	}

	sqlStr := "select id,_date,no_id,_min,_max,_change,lastprice,startprice,dde_b,dde_s,face,volume,turnoverRate,turnover,state,per from codeface where"
	filter := ""
	if len(date) > 0 {
		filter = fmt.Sprintf(" _date='%s'", date)
	}

	if code >= 0 {
		if len(filter) > 0 {
			filter += " and"
		}

		noId, exists := GetIdbyNo(code)
		if !exists {
			return
		}

		filter = fmt.Sprintf("%s `no_id`=%d", filter, noId)
	}

	sqlStr = sqlStr + filter + ";"
	rows, err := conn.Query(sqlStr)
	if err != nil {
		errStr = fmt.Sprintf(" GetCodeFace 失败%v", err)
		return
	}

	faceList, errStr = getFaceFromRows(implMap[Def.Model_Key_CodeFace], rows)

	return
}

func getFaceFromRows(impl Interface.IDbModelFill, rows *sql.Rows) (faces []Interface.IDbModel, errStr string) {
	faces = make([]Interface.IDbModel, 0, 128)
	for rows.Next() {
		tempModel, err := impl.FillCodeFace(func(tempList []interface{}) error {
			return rows.Scan(tempList...)
		})

		if err != nil {
			errStr = fmt.Sprintf("getFaceFromRows err=%v", err)
			return
		}

		faces = append(faces, tempModel)
	}

	return
}
