package filesystem

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

func getFileOwner(path string) (owner string, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	
	stat := fileInfo.Sys().(*syscall.Stat_t)
	uid := stat.Uid
	
	o, err := user.LookupId(fmt.Sprintf("%d", uid))
	if err != nil {
		return "", err
	}
	
	return o.Username, nil
}

func getFileInfo(path string) (owner string, permissions string, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", "", err
	}
	
	permissions = info.Mode().Perm().String()
	
	owner, err = getFileOwner(path)
	if err != nil {
		return "", "", err
	}
	return
}
