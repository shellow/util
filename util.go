package util

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type FileInfo struct {
	Filename string `form:"filename" json:"filename" xml:"filename"`
	Size     int    `form:"size" json:"size" xml:"size"`
	FileTime string `form:"time" json:"time" xml:"time"`
}

func CmdLinuxExc(cmdstr string) (string, error) {
	sysType := runtime.GOOS
	if !strings.EqualFold(sysType, "linux") {
		return "", errors.New("system not support")
	}
	cmd := exec.Command("/bin/bash", "-c", cmdstr)
	buf, err := cmd.CombinedOutput()
	return string(buf), err
}

func ListFiles(dir string, start int, num int) ([]FileInfo, error) {
	fileinfos := make([]FileInfo, 0)
	if start < 0 {
		return nil, errors.New("start error")
	}
	if num < 0 || num > 50 {
		return nil, errors.New("num error")
	}

	realnum := num + 5

	excstr := "ls -lt \"" + dir + "\" --time-style '+%Y-%m-%d %H:%M:%S' | awk '{print $5,$6,$7,$8}' | tail -n +" + strconv.Itoa(start+2) + " | head -n " + strconv.Itoa(realnum)
	//log.Println(excstr)
	ret, err := CmdLinuxExc(excstr)
	if err != nil {
		return nil, err
	}
	if strings.EqualFold("", ret) {
		return fileinfos, nil
	}
	lines := strings.Split(ret, "\n")
	for _, line := range lines {
		if len(fileinfos) >= num {
			break
		}
		var fileInfo FileInfo
		infos := strings.Split(line, " ")
		if len(infos) != 4 {
			continue
		}
		fileInfo.Filename = infos[3]
		fileInfo.FileTime = infos[1] + " " + infos[2]
		size, err := strconv.Atoi(infos[0])
		if err != nil {
			continue
			//return nil, errors.New("size error")
		}
		fileInfo.Size = size
		fileinfos = append(fileinfos, fileInfo)
	}
	//fmt.Println(len(lines))
	return fileinfos, nil
}

func IsFileDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func GetDirSize(dir string) int {
	if IsFileDir(dir) {
		excstr := "du -sb \"" + dir + "\" | awk '{print $1}'"
		//log.Println(excstr)
		ret, err := CmdLinuxExc(excstr)
		if err != nil {
			return -1
		}
		reg := regexp.MustCompile("[0-9]+")
		ret = reg.FindString(ret)
		size, err := strconv.Atoi(ret)
		if err != nil {
			return -1
		}
		return size
	}
	return -1
}
