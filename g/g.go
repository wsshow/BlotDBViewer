package g

import (
	"BBoltViewer/cmd"
	"BBoltViewer/db"
	"BBoltViewer/log"
	"BBoltViewer/utils"
	"os"
)

var (
	Log         = new(log.Log)
	Conf        = new(cmd.Command)
	SignalExit  = make(chan os.Signal)
	CacheDBConn = make(map[string]*db.Visiter)
)

func NewLog() *log.Log {
	utils.NotExistToMkdir("./log")
	return log.New("BBlotViewer", Conf.Loglevel)
}
