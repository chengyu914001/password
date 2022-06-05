package main

import (
	"log"
	"os"
	"io"
	"findPassword/pkg/password"
	"gopkg.in/ini.v1"
	"path/filepath"
	"runtime"
)

func main(){
	_, err := os.Stat("unzip.log")
	if !os.IsNotExist(err) {
		os.Remove("unzip.log")
	}
	logFile, err := os.OpenFile("unzip.log", os.O_CREATE | os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic("can not create log file! ")
	}
	defer logFile.Close()
	logger := log.New(io.MultiWriter(os.Stdout, logFile), "", log.Ldate | log.Ltime)

	passwordFinder :=  new(passwordGenerator.PasswordFinder)
	passwordFinder.Init()
	passwordFinder.SetTableClear()

	cfg, err := ini.Load("setting.ini")
	if err != nil {
		logger.Fatal("Fail to read file: ", err)
	}

	filename := cfg.Section("").Key("filename").String()
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		logger.Fatal(filename, " not exist!")
	}
	filenameExt := filepath.Ext(filename)
	if filenameExt == "" {
		logger.Fatal(filename, " is no extension")
	}
	cmdExe := ""
	if filenameExt == ".7z" || filenameExt == ".tar" || filenameExt == ".wim" || filenameExt == ".zip" {
		cmdExe = cfg.Section("install").Key("7z").String()
		if cmdExe == "" {
			logger.Fatal("Please write 7z command path!")
		}
	}
	if cmdExe == "" {
		logger.Fatal("Can't support", filenameExt, "File!")
	}
	passwordGenerator.Set7zPath(cmdExe)
	logger.Println("filename:", filename)
	
	passwordCfg := cfg.Section("password")
	
	isAllNumber, _ := passwordCfg.Key("AllNumber").Bool()
	if isAllNumber {
		passwordFinder.AddTableAllNumber()
	}
	isAllLowercase, _ := passwordCfg.Key("AllLowercase").Bool()
	if isAllLowercase {
		passwordFinder.AddTableAllLowercase()
	}
	isAllUppercase, _ := passwordCfg.Key("AllUppercase").Bool()
	if isAllUppercase {
		passwordFinder.AddTableAllUppercase()
	}
	isAllSpecial, _ := passwordCfg.Key("AllSpecial").Bool()
	if isAllSpecial {
		passwordFinder.AddTableAllSpecial()
	}
	passwordFinder.AddTable(passwordCfg.Key("defineOthers").String())
	{
		table := passwordFinder.GetTable()
		if table[0] == ' '{
			table = "<space>" + table[1:]
		}
		logger.Println("Table:", table)
	}

	passwordFinder.SetReg(passwordCfg.Key("regexp").String())
	logger.Println("regexp:",passwordFinder.GetReg())

	startLen, err := passwordCfg.Key("startLen").Int()
	if err != nil {
		logger.Fatal("startLen:", err)
	}
	endLen, err := passwordCfg.Key("endLen").Int()
	if err != nil {
		logger.Fatal("endLen:", err)
	}
	logger.Println("Password len is", startLen, "to", endLen, ".")

	cmdType := filepath.Base(cmdExe)
	logger.Println("Cmd type", cmdType)
	
	logger.Println("start!")

	ans := passwordFinder.Find7ZPassword(uint8(startLen), uint8(endLen), filename, uint8(runtime.NumCPU() - 1))
	if ans == "" {
		logger.Println("password no answer!")
	}else{
		logger.Println("password is :", ans)
	}
}