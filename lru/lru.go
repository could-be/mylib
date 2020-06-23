package lru

import (
	"fmt"
	"sync"
)

type LRU struct {
	mu    sync.RWMutex          // 保护
	table map[int]*HashLinkList // 防止内存逃逸?
	cap   int                   // 容积
	len   int                   // 当前元素个数
	head  *HashLinkList         // 头
	tail  *HashLinkList         // 尾
}

func NewLRU(cap int) *LRU {
	return &LRU{
		mu:    sync.RWMutex{},
		table: make(map[int]*HashLinkList, cap),
		cap:   cap,
	}
}

// k-v
type HashLinkList struct {
	pre  *HashLinkList
	Next *HashLinkList
	// Key int
	Value int
}

func newHashLinkList(v int) *HashLinkList {
	return &HashLinkList{
		pre:   nil,
		Next:  nil,
		Value: v,
	}
}

// 当只有一个元素的时候，head != nil, tail == nil

// 最右端添加新元素, 添加到末尾
func (h *LRU) Assign(k, v int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 表头是否为空
	if h.head == nil {
		h.head = newHashLinkList(v)
		h.tail = h.head
		return
	}

	// 是否已存在
	item, ok := h.table[k]
	switch ok {
	case true:
		// 更新
		// 存在, 更新对应项，并且移动到表位
		item.Value = v
		// item.pre.Next = item.Next
		// 移动当前元素到表末尾

		// 大于1个
		switch {
		case h.tail == item:
			// 表尾, 直接返回
			return
		case h.head == item:
			// 表头, 更新表头
			h.head = item.Next
			h.head.pre = nil
		default:
			// 中间, 大于三个元素
			item.pre.Next = item.Next
		}

		// 表尾特征
		h.tail.Next = item
		item.pre = h.tail
		h.tail = item
		h.tail.Next = nil

	default:
		// 新增
		h.len++
		// 不存在
		// 未达到上限
		if h.len < h.cap {
			// 容积是否已经达到上限
			item := newHashLinkList(v)
			// 这种情况一定是, 初始化, head和tail同为空或同为不为空
			if h.tail == nil {
				h.head, h.tail = item, item
			} else {
				// 更新表尾
				h.tail.Next = item
				item.pre = h.tail
				h.tail = item
			}
		} else {
			// 达到上限，淘汰链表头元素
			lastHead, lastTail := h.head, h.tail
			h.head = lastHead.Next
			h.head.pre = nil

			lastTail.Next = lastHead
			lastHead.pre = h.tail
			h.tail = lastHead
			h.tail.Next = nil
		}
	}
}

func (h *LRU) Remove(k int) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 不存在
	item, ok := h.table[k]
	if !ok {
		return
	}

	h.len--
	delete(h.table, k)

	// 是表头
	if h.head == item {
		switch {
		// 只有一个元素, 更新表头和表尾
		case h.head == h.tail:
			h.head = nil
			h.tail = nil
		default:
			// 更新表头元素
			h.head = item.Next
		}
		return
	}

	// 表尾, 这里不会出现只有一个元素的情况了，也就是tail一定有前驱
	if h.tail == item {
		item.pre.Next = item.Next // 即 item.pre.Next = nil
		return
	}

	// 非表头 非表尾, 即至少是三个
	item.pre.Next = item.Next
	return
}

func (h *LRU) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.len
}

func (h *LRU) Print() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	item := h.head

	for item != nil {
		fmt.Println(item.Value)
	}
}
