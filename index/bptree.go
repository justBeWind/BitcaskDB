package index

import (
	"bitcask-go/data"
	"go.etcd.io/bbolt"
	_ "go.etcd.io/bbolt"
	"path/filepath"
)

// B+树索引 存储到磁盘的索引
// 主要封装了 go.etcd.io/bbolt
type BPlusTree struct {
	tree *bbolt.DB
}

const bptreeIndexFileName = "bptree-index"

var indexBucketName = []byte("bitcask-index")

// NewBPlusTree 初始化B+树索引
func NewBPlusTree(dirPath string, syncWrites bool) *BPlusTree {
	opts := bbolt.DefaultOptions
	opts.NoSync = !syncWrites

	bptree, err := bbolt.Open(filepath.Join(dirPath, bptreeIndexFileName), 0644, opts)
	if err != nil {
		panic("failed to ope bptree")
	}

	if err := bptree.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(indexBucketName)
		return err
	}); err != nil {
		panic("failed to create bucket in bptree")
	}

	return &BPlusTree{tree: bptree}
}

// Put 向索引中存储key对应的数据位置信息
func (bpt *BPlusTree) Put(key []byte, pos *data.LogRecordPos) bool {
	if err := bpt.tree.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(indexBucketName)
		return bucket.Put(key, data.EncodeLogRecordPos(pos))
	}); err != nil {
		panic("falied to put value in bptree")
	}
	return true
}

// Get 根据key 取出对应的索引位置信息
func (bpt *BPlusTree) Get(key []byte) *data.LogRecordPos {
	var pos *data.LogRecordPos
	if err := bpt.tree.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(indexBucketName)
		value := bucket.Get(key)
		if len(value) != 0 {
			pos = data.DecodeLogRecordPos(value)
		}
		return nil
	}); err != nil {
		panic("failed to get value in bptree")
	}
	return pos
} //返回数据位置信息

// Delete 根据 key 删除对应的索引位置信息
func (bpt *BPlusTree) Delete(key []byte) bool {
	var ok bool
	if err := bpt.tree.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(indexBucketName)
		if value := bucket.Get(key); len(value) != 0 {
			ok = true
			return bucket.Delete(key)
		}
		return nil
	}); err != nil {
		panic("falied to delete value in bptree")
	}
	return ok
}

// Size 索引中存在的数据量
func (bpt *BPlusTree) Size() int {
	var size int
	if err := bpt.tree.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(indexBucketName)
		size = bucket.Stats().KeyN

		return nil
	}); err != nil {
		panic("failed to get size in bptree")
	}
	return size
}

// Iterator 索引迭代器
func (bpt *BPlusTree) Iterator(reverse bool) Iterator {
	return newBptreeIterator(bpt.tree, reverse)
}

func (bpt *BPlusTree) Close() error {
	return bpt.tree.Close()
}

// B+树迭代器
type bptreeIterator struct {
	tx        *bbolt.Tx
	cursor    *bbolt.Cursor
	reverse   bool
	currKey   []byte
	currValue []byte
}

func newBptreeIterator(tree *bbolt.DB, reverse bool) *bptreeIterator {
	tx, err := tree.Begin(false)
	if err != nil {
		panic("failed to begin a transaction")
	}

	bpi := &bptreeIterator{
		tx:      tx,
		cursor:  tx.Bucket(indexBucketName).Cursor(),
		reverse: reverse,
	}
	//为了给 currKey 和currValue赋值
	bpi.Rewind()
	return bpi

}

// Rewind 重新回到迭代器的起点，即第一个数据
func (bpi *bptreeIterator) Rewind() {
	if bpi.reverse {
		bpi.currKey, bpi.currValue = bpi.cursor.Last()
	} else {
		bpi.currKey, bpi.currValue = bpi.cursor.First()
	}
}

// Seek 根据传入的key查找到第一个大于（或小于）等于目标的key，根据这个key开始遍历
func (bpi *bptreeIterator) Seek(key []byte) {
	bpi.currKey, bpi.currValue = bpi.cursor.Seek(key)
}

// Next 跳转到下一个key
func (bpi *bptreeIterator) Next() {
	if bpi.reverse {
		bpi.currKey, bpi.currValue = bpi.cursor.Prev()
	} else {
		bpi.currKey, bpi.currValue = bpi.cursor.Next()
	}
}

// Valid 是否有效，即是否已经遍历完了所有的key，用于退出遍历
func (bpi *bptreeIterator) Valid() bool {
	return len(bpi.currKey) != 0
}

// Key 当前遍历位置的Key数据
func (bpi *bptreeIterator) Key() []byte {
	return bpi.currKey
}

// Value 当前遍历位置的Value数据
func (bpi *bptreeIterator) Value() *data.LogRecordPos {
	return data.DecodeLogRecordPos(bpi.currValue)
}

// Close 关闭迭代器，释放资源
func (bpi *bptreeIterator) Close() {
	_ = bpi.tx.Commit()
}
