package Interface

import (
	"database/sql"

	"SparkEven/govm/Src/Model"
)

type IDbSupport interface {
	SearchCodeFace(coon *sql.DB, date string, code int) (list []*Model.CodeFace, err Model.Err)
	SaveFaceList(*sql.DB, []*Model.CodeFace) Model.Err
	SaveFacePrices(*sql.DB, *Model.CodeFace, []*Model.CodePrice) Model.Err
	GetConn() (conn *sql.DB, errStr Model.Err)
	GetDateCodePrice(*sql.DB, string, int) ([]*Model.CodePrice, Model.Err)
	GetSysConfigList(conn *sql.DB) (list []*Model.SysConfig, err Model.Err)
	SaveSysConfig(conn *sql.DB, k, v string) (err Model.Err)
}
