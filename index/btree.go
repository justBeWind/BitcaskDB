package index

import (
	"bitcask-go/data"
	"github.com/google/btree"
	"sync"
)

// BTree 封装谷歌的btree 库
type BTree struct {
	tree *btree.BTree
	lock *sync.RWMutex //由于Btree中 写的时候是并发不安全的，所以要手动加锁
}

func NewBTree() *BTree {
	return &BTree{
		tree: btree.New(32), //32表示btree中最多的节点树
		lock: new(sync.RWMutex),
	}
}

func (bt *BTree) Put(key []byte, pos *data.LogRecordPos) bool {
	it := &Item{key: key, pos: pos}
	//先在调用B树来存取时，先加锁
	bt.lock.Lock()
	bt.tree.ReplaceOrInsert(it)
	bt.lock.Unlock()
	return true

}
func (bt *BTree) Get(key []byte) *data.LogRecordPos {
	it := &Item{key: key}
	btreeItem := bt.tree.Get(it)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pos
}

func (bt *BTree) Delete(key []byte) bool {
	it := &Item{key: key}
	bt.lock.Lock()
	oldItem := bt.tree.Delete(it) //返回一个删除是否成功的值
	bt.lock.Unlock()
	if oldItem == nil {
		return false
	}
	return true
}
