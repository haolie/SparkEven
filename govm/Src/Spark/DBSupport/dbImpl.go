package DBSupport

import (
	"database/sql"

	"SparkEven/govm/Src/Interface"
	"SparkEven/govm/Src/Model"
)

type mysqlImpl struct{}

func (impl *mysqlImpl) SearchCodeFace(coon *sql.DB, date string, code int) (list []*Model.CodeFace, err Model.Err) {
	itemList, errStr := GetCodeFace(coon, code, date)
	err = Model.Err(errStr)

	if err.Exists() {
		return list, err
	}

	list = make([]*Model.CodeFace, 0, len(itemList))
	for _, item := range itemList {
		list = append(list, item.(*Model.CodeFace))
	}

	return
}

func (impl *mysqlImpl) SaveFaceList(conn *sql.DB, list []*Model.CodeFace) (err Model.Err) {

	insertList := make([]Interface.IDbModel, 0, len(list))
	for _, item := range list {
		insertList = append(insertList, item)
	}

	er := InsertList(conn, insertList)
	err = Model.Err(er)

	return ""
}

func (impl *mysqlImpl) GetConn() (conn *sql.DB, errStr Model.Err) {
	conn, err := CreateConn()
	errStr = Model.Err(err)

	return
}

func (impl *mysqlImpl) SaveFacePrices(conn *sql.DB, face *Model.CodeFace, priceList []*Model.CodePrice) Model.Err {
	return SaveTimePrices(conn, face, priceList)
}

func (impl *mysqlImpl) GetDateCodePrice(conn *sql.DB, date string, code int) ([]*Model.CodePrice, Model.Err) {

	return searchPrice(conn, date, code)
}

//
//SearchCodeFace(date time.Time, code int) (list []*Model.CodeFace, err error)
//SaveFaceList([]*Model.CodeFace) error
