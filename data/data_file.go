package data

import (
	"bitcask-go/fio"
	"fmt"
	"path/filepath"
)

const DataFileNameSuffix = ".data"

type DataFile struct {
	FileId   uint32       // 文件id
	WriteOff int64        //文件写到了哪个位置
	IoManger fio.IOManger //io 读写管理
}

// OpenDataFile 打开新的数据文件
func OpenDataFile(dirPath string, fileId uint32) (*DataFile, error) {
	//根据dirPath和fileId生成完整的文件名称
	fileName := filepath.Join(dirPath, fmt.Sprintf("%09d", fileId)+DataFileNameSuffix)
	//初始化IOManger管理器接口
	ioManager, err := fio.NewIOManger(fileName)
	if err != nil {
		return nil, err
	}

	return &DataFile{
		FileId:   fileId,
		WriteOff: 0,
		IoManger: ioManager,
	}, nil

}

// ReadLogRecord 根据offset从数据文件中LogRecord 读取文件数据 并返回LogRecord的大小，方便对offset的更新
func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, int64, error) {

	return nil, 0, nil
}

// 写入文件
func (df *DataFile) Write(buf []byte) error {
	return nil
}

// Sync 持久化文件
func (df *DataFile) Sync() error {
	return nil
}
