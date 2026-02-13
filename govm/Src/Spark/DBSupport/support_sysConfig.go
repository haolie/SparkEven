package DBSupport

import (
	"database/sql"
	"fmt"

	"SparkEven/govm/Src/Model"
)

func getSysConfigs(conn *sql.DB) (sysModels []*Model.SysConfig, errStr Model.Err) {

	sqlStr := "select `Key`,`Value` from sys_config;"

	rows, err := conn.Query(sqlStr)
	if err != nil {
		errStr = Model.Err(fmt.Sprintf(" sys_config 失败%v", err))
		return
	}

	sysModels = make([]*Model.SysConfig, 0, 8)
	for rows.Next() {
		item := &Model.SysConfig{}
		err := rows.Scan(&item.Key, &item.Value)

		if err != nil {
			errStr = Model.Err(fmt.Sprintf("GetSysConfigs err=%v", err))
			return
		}

		sysModels = append(sysModels, item)
	}

	return
}

func updateSysConfig(conn *sql.DB, k, v string) (errStr Model.Err) {
	str := fmt.Sprintf("REPLACE INTO sys_config(`key`,`value`) VALUES(\"%s\",\"%s\");", k, v)
	_, err := conn.Exec(str)
	if err != nil {
		errStr = Model.Err(fmt.Sprintf("UpdateSysConfig err=%v", err))
	}

	return
}
