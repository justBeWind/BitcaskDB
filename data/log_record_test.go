package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeLogRecord(t *testing.T) {
	//正常情况
	rec1 := &LogRecord{
		Key:   []byte("name"),
		Value: []byte("bitcask-go"),
		Type:  LogRecordNormal,
	}
	res1, n1 := EncodeLogRecord(rec1)
	t.Log(res1)
	t.Log(n1)
	assert.NotNil(t, res1)
	assert.Greater(t, n1, int64(5))

	//value为空的情况

	//对Deleted情况的测试

}

func TestDecoderLogRecordHeader(t *testing.T) {

}
