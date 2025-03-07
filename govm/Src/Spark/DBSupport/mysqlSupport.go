package DBSupport

import (
	"database/sql"
	"fmt"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Interface"
	"SparkEven/govm/Src/Model"
)

var (
	implMap  map[string]Interface.IDbModelFill
	mPageKey = "mysqlSupport"
)

func InsertList(conn *sql.DB, modelList []Interface.IDbModel) string {

	fillImpl := implMap[Def.Model_Key_CodeFace]
	i, count := 0, 300
	for {
		index := i + count
		if index > len(modelList) {
			index = len(modelList)
		}

		array := modelList[i:index]

		var vStr = fillImpl.GetInsertPerStr()
		for i, item := range array {

			if i > 0 {
				vStr += ","
			}

			codeFace := item.(*Model.CodeFace)
			noId, errStr := GetOrAddId(codeFace.Code, conn)
			if errStr != "" {
				return errStr
			}

			codeFace.ID = int(getNewFaceId())
			vStr += fillImpl.GetInsertValueStr(item, noId)
		}

		vStr = vStr + ";"
		_, err := conn.Exec(vStr)
		if err != nil {
			return fmt.Sprintf("mysqlSupport.InsertList err=%v sql=%s", err, vStr)
		}
		i += count
		if i >= len(modelList) {
			break
		}
	}

	return ""
}

func getDateCodeFace(conn *sql.DB, code int, date string) (face *Model.CodeFace, errStr Model.Err) {
	list, er := GetCodeFace(conn, code, date)
	errStr = Model.Err(er)
	if errStr.Exists() {
		return
	}

	if len(list) != 1 {
		return
	}

	face = list[0].(*Model.CodeFace)
	return
}

func GetCodeFace(conn *sql.DB, code int, date string) (faceList []Interface.IDbModel, errStr string) {

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
