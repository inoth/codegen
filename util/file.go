package util

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func CreateFile(path, content string) error {
	fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("%v/文件创建失败, %v", path, err.Error())
		return err
	}
	defer fw.Close()
	_, err = fw.WriteString(content)
	if err != nil {
		log.Errorf("%v/文件创建失败, %v", path, err.Error())
		return err
	}
	return nil
}

func CreateFileBytes(path string, write func(*os.File) error) error {
	fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Errorf("%v/文件创建失败, %v", path, err.Error())
		return err
	}
	defer fw.Close()
	// _, err = fw.Write(content)
	err = write(fw)
	if err != nil {
		log.Errorf("%v/文件创建失败, %v", path, err.Error())
		return err
	}
	return nil
}
