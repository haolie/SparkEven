package DBSupport

import (
	"context"

	"SparkEven/govm/Src/Common/Def"
	"SparkEven/govm/Src/Interface"
	"SparkEven/govm/Src/Spark/MPark"
)

func init() {
	implMap = make(map[string]Interface.IDbModelFill, 2)
	implMap[Def.Model_Key_CodeFace] = &CodeFaceImpl{}

	MPark.DbSupport = &mysqlImpl{}
	MPark.RegisterLoad(mPageKey, func(ctx context.Context) []string {
		conn, errStr := CreateConn()
		if errStr != "" {
			return []string{errStr}
		}

		errStr = InitSupport(conn)
		if errStr != "" {
			return []string{errStr}
		}
		return []string{}
	})
}
