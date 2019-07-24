package common

import (
	"log"
	"os"
)

//Log var for log
var Log *log.Logger

//NewLog create log File
func NewLog(logpath string) {
	println("LogFile: " + logpath)
	file, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

//LogToFile save log to file
func LogToFile(row string) {
	f, err := os.OpenFile("ut_amplio_api_log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
}
