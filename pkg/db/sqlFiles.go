package db

import "embed"

//go:embed queries/*/*.sql
var SqlFiles embed.FS
