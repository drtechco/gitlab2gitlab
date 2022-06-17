package configs

import "time"

var (
	SqliteDsn   = ""
	ServerTZ, _ = time.LoadLocation(
		"Asia/Shanghai",
	)
	TempDir = ""
)
