package util

import snowflake "github.com/yockii/snowflake_ext"

var worker *snowflake.Worker

func InitializeSnowflake() (err error) {
	worker, err = snowflake.NewSnowflake(1)
	return
}

func NextID() uint64 {
	return worker.NextId()
}
