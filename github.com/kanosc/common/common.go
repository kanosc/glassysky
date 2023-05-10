package common

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func GetDateStr() string {
	t := time.Now()
	weekday := []string{"日", "一", "二", "三", "四", "五", "六"}
	date := fmt.Sprintf("%v年%v月%v日 星期%v", t.Year(), int(t.Month()), t.Day(), weekday[int(t.Weekday())])
	return date
}

func RollInt(end int64) int {
	ret, _ := rand.Int(rand.Reader, big.NewInt(end+1))
	return int(ret.Int64())
}

func GetDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != "" {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func GetFileNameList(path string) []string {
	filenames := make([]string, 0)
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		log.Println("###########", len(info.Name()))
		if !info.IsDir() && len(info.Name()) != 0 {
			filenames = append(filenames, info.Name())
		}
		return err
	})
	log.Println(len(filenames))
	for _, j := range filenames {
		log.Println("[" + j + "]")
	}
	return filenames
}

func FileExistInDir(path, filename string) (bool, error) {
	//	var fileExistFlag = false
	/* err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info.Name() == filename {
			fileExistFlag = true
		}
		return err
	}) */
	fullFileName := path + filename
	if _, err := os.Stat(fullFileName); os.IsNotExist(err) {
		return false, err
	} else {
		return true, err
	}
}
