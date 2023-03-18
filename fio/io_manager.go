package fio

const DataFilePerm = 0644

// IOManger 抽象IO管理接口，可以接入不同的IO类型，目前支持标准文件 IO
type IOManger interface {
	// Read 从文件的给定位置读取对应的数据
	Read([]byte, int64) (int, error)
	// Write 写入字节数组到文件中
	Write([]byte) (int, error)

	// Sync 持久化数据
	Sync() error

	// Close 关闭文件
	Close() error
}

// 初始化IOManger，目前只支持标准FileIO
func NewIOManger(fileName string) (IOManger, error) {
	return newFileIOManager(fileName)
}
