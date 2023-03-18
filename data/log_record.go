package data

type LogRecordType = byte //枚举类型

const (
	LoRecordNormal LogRecordType = iota
	LogRecordDeleted
)

// LogRecord 写入到数据文件的记录
// 之所以叫日志，是因为数据文件中的数据是追加写入的，类似日志的格式
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

type LogRecordPos struct {
	Fid    uint32 //文件id ，表示将数据存储到了哪个文件中
	Offset int64  // 偏移， 表示将数据存储到了数据文件的哪个位置中
}

// EncodeLogRecord 对 LogRecord进行编码，返回字节数组及长度
func EncodeLogRecord(logRecord *LogRecord) ([]byte, int64) {
	return nil, 0
}
