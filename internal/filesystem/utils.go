package filesystem

import "os"

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func GetEntryType(isDir bool) EntryType {
	if isDir {
		return EntryTypeDir
	}
	return EntryTypeFile
}
