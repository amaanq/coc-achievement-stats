package log

import (
	"os"

	"github.com/withmandala/go-log"
)

var (
	Log = log.New(os.Stdout).WithColor().WithDebug().WithTimestamp()
)
