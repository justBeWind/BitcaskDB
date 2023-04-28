package utils

import (
	"io/fs"
	"path/filepath"
	"syscall"
)

// DirSize 获取一个目录的大小
func DirSize(dirPath string) (int64, error) {
	var size int64
	//对目录进行递归遍历
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// AvailableDiskSize 获取磁盘剩余可用空间大小
func AvailableDiskSize() (uint64, error) {
	wd, err := syscall.Getwd() //wd是当前文件的绝对路径
	if err != nil {
		return 0, err
	}
	var stat syscall.Statfs_t
	if err = syscall.Statfs(wd, &stat); err != nil {
		return 0, err
	}
	return stat.Bavail * uint64(stat.Bsize), nil
	//h := syscall.MustLoadDLL("kernel32.dll")
	//c := h.MustFindProc("GetDiskFreeSpaceExW")
	//lpFreeBytesAvailable := int64(0)
	//lpTotalNumberOfBytes := int64(0)
	//lpTotalNumberOfFreeBytes := int64(0)
	//c.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(wd))),
	//	uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
	//	uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
	//	uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	//return uint64(lpTotalNumberOfFreeBytes), nil
}
