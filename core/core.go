package core

import (
	"bibt-SpeedSkat/backup/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var (
	INIcCfg    *ini.File
	BackDir    string
	BackTmpDir string
	Wg         sync.WaitGroup
	//--Pg
	PgConn  *gorm.DB
	PgLimit int
)

func Run() {
	InIConfig()
	if service, err := INIcCfg.Section("").Key("service").Int(); service == 1 && err == nil {
		go InIWebSite()
		for {
			Wg.Add(1)
			TaskRun()
			Wg.Done()
			Wg.Wait()

			//sleep
			_sleepTime()
		}
	} else {
		TaskRun()
	}
}

func InIConfig() {
	var err error
	INIcCfg, err = ini.Load("./cfg/backup.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	BackDir = INIcCfg.Section("").Key("backupdir").String()
}

func InIWebSite() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		fileList := utils.GetFileList(BackDir)
		htmlStr := "<a>"
		if len(fileList) > 0 {
			for _, v := range fileList {
				htmlStr += fmt.Sprintf("<font color='green'>%s</font>&nbsp;&nbsp; <a href='download/%s'> %s</a>&nbsp;&nbsp;<font color='red'>%s byte</font>&nbsp;&nbsp;<font color='black'>%s</font><p>",v["mode"], v["name"], v["name"],v["size"],v["modetime"])
			}
		}
		c.String(200, htmlStr)
	})

	r.Static("download", BackDir)
	r.Run(":17894")
}

func TaskRun() {

	tmpDir := _getTmpDir()
	for _, item := range INIcCfg.Sections() {
		Task(item, tmpDir)
	}

	//zip
	_compressTarGz(tmpDir)
	//remove
	os.RemoveAll(tmpDir)
}

func Task(item *ini.Section, backDir string) {
	BackTmpDir = backDir + "/" + item.Name()
	if strings.ToLower(item.Name()) == "default" {
		return
	}

	fmt.Println(fmt.Sprintf(":: Begin Execute Task 【%s】::", item.Name()))
	if ok, _ := PathExists(BackTmpDir); !ok {
		Mkdir(BackTmpDir)
	}

	switch item.Key("mode").String() {
	case "file":
		CopyDir(item.Key("dir").String(), BackTmpDir)
	case "pgsql":
		limit, _ := item.Key("limit").Int()
		PgLimit = cast.ToInt(utils.If(limit > 0, limit, 20000))

		PgSqldump(item.Key("dsn").String(), item.Key("schema").String())
	case "ssh":
		SshCopy(item)
	}
}

//-----------------local func
func _getTmpDir() (tmpDir string) {
	backDir := INIcCfg.Section("").Key("backupdir").String()
	tmpDir = path.Join(backDir, "tmp"+utils.New())
	os.RemoveAll(tmpDir)
	Mkdir(tmpDir)
	return
}

func _sleepTime() {
	sleepTime := INIcCfg.Section("").Key("sleeptime").String()
	times := strings.Split(sleepTime, "/")
	timeNums := 60
	timeM := "m"
	if len(times) > 0 {
		if cast.ToInt(times[0]) > 0 {
			timeNums = cast.ToInt(times[0])
		}

		if times[1] == "" || !utils.InArray(timeM, []string{"m", "s"}) {
			timeM = "m"
		} else {
			timeM = times[1]
		}
	}
	if timeM == "s" {
		time.Sleep(time.Duration(timeNums) * time.Second)
	} else {
		time.Sleep(time.Duration(timeNums) * time.Minute)
	}
}

func _compressTarGz(tmpDir string) {
	fmt.Println(fmt.Sprintf(":: Begin Compress BackupDataDict 【%s】::", tmpDir))
	f, err := os.Open(tmpDir)
	if err != nil {
		fmt.Println("Compress Backup File Fail,", err.Error())
	}
	uid, _ := utils.GenerateUUID()
	dest := fmt.Sprintf("%s/%s@BackupFile%s.tar.gz", INIcCfg.Section("").Key("backupdir").String()+"/", utils.GetDate(), uid)
	if err = Compress(f, dest); err != nil {
		fmt.Println("######## Compress Backup Data Is Error ", err.Error())
	} else {
		fmt.Println("################## Compress Backup Data Is Success! ################")
	}
}

func _cleanTar() {
	fmt.Println("clean Tar is Success")
}
