package bitcask_go

// Options 提供给用户的配置项
type Options struct {
	DirPath string //数据库数据目录

	//数据文件的大小
	DataFileSize int64

	//每次写数据是否要持久化
	SyncWrites bool

	//索引类型
	IndexerType IndexerType
}

type IndexerType = int8

const (
	//BTree 索引
	BTree IndexerType = iota + 1

	//ART Adpative Radix Tree
	ART
)
