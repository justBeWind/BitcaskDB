package index

import (
	"bitcask-go/data"
	"bytes"
	"github.com/google/btree"
	"sort"
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

func (bt *BTree) Size() int {
	return bt.tree.Len()
}

func (bt *BTree) Iterator(reverse bool) Iterator {
	if bt.tree == nil {
		return nil
	}
	bt.lock.RLock()
	defer bt.lock.RUnlock()
	return newBTreeIterator(bt.tree, reverse)
}

// BTree索引迭代器
type btreeIterator struct {
	currIndex int     //当前遍历的下标位置
	reverse   bool    //是否反向遍历
	values    []*Item //key + 位置索引信息
}

// BTree索引迭代器的构造方法
func newBTreeIterator(tree *btree.BTree, reverse bool) *btreeIterator {
	var idx int
	values := make([]*Item, tree.Len())

	//将所有的数据存放到数组中，
	//但是其实有一个潜在的问题，就是带来内存的急剧膨胀 （迭代器中22：16）
	saveValues := func(it btree.Item) bool {
		values[idx] = it.(*Item)
		idx++
		return true
	}
	if reverse {
		tree.Descend(saveValues)
	} else {
		tree.Ascend(saveValues)
	}
	return &btreeIterator{
		currIndex: 0,
		reverse:   reverse,
		values:    values,
	}
	//值得注意的是，BTree因为本事就是排好序的了，所以迭代器中的values数组也是排好序的

}

func (bti *btreeIterator) Rewind() {
	bti.currIndex = 0
}

func (bti *btreeIterator) Seek(key []byte) {
	if bti.reverse {
		bti.currIndex = sort.Search(len(bti.values), func(i int) bool {
			return bytes.Compare(bti.values[i].key, key) <= 0
		})
	} else {
		bti.currIndex = sort.Search(len(bti.values), func(i int) bool {
			return bytes.Compare(bti.values[i].key, key) >= 0
		})
	}
}

func (bti *btreeIterator) Next() {
	bti.currIndex += 1
}

func (bti *btreeIterator) Valid() bool {
	return bti.currIndex < len(bti.values)
}

func (bti *btreeIterator) Key() []byte {
	return bti.values[bti.currIndex].key
}

func (bti *btreeIterator) Value() *data.LogRecordPos {
	return bti.values[bti.currIndex].pos
}

func (bti *btreeIterator) Close() {
	bti.values = nil
}
