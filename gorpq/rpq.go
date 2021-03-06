package gorpq

//KeyType type for rpq index key
type KeyType = int

//ValueType type for rpq values
type ValueType = int

type node struct {
	key    KeyType
	val    ValueType
	left   *node
	next   *node
	parent *node
	rank   int
}

type stack struct {
	internal []*node
	length   int
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func (st *stack) push(node *node) {
	st.internal = append(st.internal, node)
	st.length++
}
func (st *stack) top() *node {
	if st.empty() {
		return nil
	}
	return st.internal[st.length-1]

}
func (st *stack) pop() {
	if !st.empty() {
		st.internal = st.internal[:st.length-1]
		st.length--
	}
}
func (st *stack) size() int {
	return st.length
}
func (st *stack) empty() bool {
	return st.length == 0
}

//RPQ Is the rank pairing queue.
type RPQ struct {
	head  *node
	size  int
	index map[KeyType]*node
}

//Empty returns bool
func (q *RPQ) Empty() bool {
	return q.size == 0
}

//Size return size
func (q *RPQ) Size() int {
	return q.size
}

//Top returns top value.
func (q *RPQ) Top() (KeyType, ValueType) {
	return q.head.key, q.head.val
}

//Push inserts a value, only unique keys.
func (q *RPQ) Push(key KeyType, val ValueType) {
	if _, ok := q.index[key]; ok {
		return
	}
	n := &node{key: key, val: val}
	q.insert(n)
	q.index[key] = n
	q.size++
}

//Pop is delete min
func (q *RPQ) Pop() { //ok
	if q.Empty() {
		return
	}
	bucket := make([]*node, maxBucketSize(q.size))
	for ptr := q.head.left; ptr != nil; {
		next := ptr.next
		ptr.next = nil
		ptr.parent = nil
		q.multipass(bucket, ptr)
		ptr = next
	}
	for ptr := q.head.next; ptr != q.head; {
		next := ptr.next
		ptr.next = nil
		q.multipass(bucket, ptr)
		ptr = next
	}
	delete(q.index, q.head.key)
	q.head = nil
	q.size--

	for _, b := range bucket {
		if b != nil {
			q.insert(b)
		}
	}

}

//Clear does post order traversal to delete.
func (q *RPQ) Clear() {
	for !q.Empty() {
		q.Pop()
	}
}

//Decrease key value. Type 2
func (q *RPQ) Decrease(key KeyType, val ValueType) {
	ptr, ok := q.index[key]
	if !ok {
		return
	}
	if val < ptr.val {
		ptr.val = val
	}
	if ptr == q.head {
		return
	}
	if ptr.parent == nil {
		if ptr.val < q.head.val {
			q.head = ptr
		}
	} else {
		parentPtr := ptr.parent
		if ptr == parentPtr.left {
			parentPtr.left = ptr.next
			if parentPtr.left != nil {
				parentPtr.left.parent = parentPtr
			}
		} else {
			parentPtr.next = ptr.next
			if parentPtr.next != nil {
				parentPtr.next.parent = parentPtr
			}
		}
		ptr.parent = nil
		ptr.next = nil
		if ptr.left != nil {
			ptr.rank = ptr.left.rank + 1
		} else {
			ptr.rank = 0
		}
		q.insert(ptr)
		if parentPtr.parent == nil {
			if parentPtr.left != nil {
				parentPtr.rank = parentPtr.left.rank + 1
			} else {
				parentPtr.rank = 0
			}
		} else {
			for parentPtr.parent != nil {
				var i, j int = -1, -1
				if parentPtr.left != nil {
					i = parentPtr.left.rank
				}
				if parentPtr.next != nil {
					j = parentPtr.next.rank
				}
				var k int
				if abs(i-j) > 1 {
					k = max(i, j)
				} else {
					k = max(i, j) + 1
				}
				if k >= parentPtr.rank {
					break
				}
				parentPtr.rank = k
				parentPtr = parentPtr.parent

			}
		}
	}
}
func (q *RPQ) insert(ptr *node) {
	if q.head == nil {
		q.head = ptr
		ptr.next = ptr
	} else {
		ptr.next = q.head.next
		q.head.next = ptr
		if ptr.val < q.head.val {
			q.head = ptr
		}
	}
}
func (q *RPQ) link(left *node, right *node) *node {
	if right == nil {
		return left
	}
	var winner, loser *node
	if right.val < left.val {
		winner = right
		loser = left
	} else {
		winner = left
		loser = right
	}
	loser.parent = winner
	if winner.left != nil {
		loser.next = winner.left
		loser.next.parent = loser
	}
	winner.left = loser
	winner.rank = loser.rank + 1

	return winner
}

//ceil(log2()) + 1
func maxBucketSize(size int) int {
	var bit, count int = 1, size
	count >>= 1
	for count != 0 {
		bit++
		count >>= 1
	}
	return bit + 1
}
func (q *RPQ) multipass(bucket []*node, ptr *node) {
	for bucket[ptr.rank] != nil {
		var rank uint = uint(ptr.rank)
		ptr = q.link(ptr, bucket[rank])
		bucket[rank] = nil
	}
	bucket[ptr.rank] = ptr
}

//NewRPQ generates RPQ structs
func NewRPQ() *RPQ {
	return &RPQ{index: make(map[KeyType]*node)}
}
