package bitcask_go

import "errors"

var (
	ErrKeyIsEmpty        = errors.New("the key is empty")
	ErrIndexUpdataFailed = errors.New("failed to update index")
	ErrKeyNotFound       = errors.New("key not found in database")
	ErrDataFileNotFound  = errors.New("data file is not found")
	ErrDataDirectory     = errors.New("the database directory maybe corrupted")
)
