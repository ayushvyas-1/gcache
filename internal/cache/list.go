package cache

type DoublyLinkedList struct {
	head  *DoublyNode
	last  *DoublyNode
	count int
}

func NewDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{}
}

func (l *DoublyLinkedList) AttachNode(node *DoublyNode) {
	if l.head == nil {
		l.head = node
	} else {
		l.last.SetNext(node)
		node.SetPrev(l.last)
	}
	l.last = node
	l.count++
}

func (l *DoublyLinkedList) Add(data interface{}) {
	l.AttachNode(NewDoublyNode(data))
}

func (l *DoublyLinkedList) Count() int {
	return l.count
}

// func (l *DoublyLinkedList) GetNext() (*DoublyNode, error) {
// 	if l.head == nil {
// 		return nil, errors.New("list is empty")
// 	}

// 	return l.head, nil
// }

// func (l *DoublyLinkedList) GetPrev() (*DoublyNode, error) {
// 	if l.last == nil {
// 		return nil, errors.New("list is empty")
// 	}

// 	return l.last, nil
// }

// func (l *DoublyLinkedList) GetByIndex(index int) (*DoublyNode, error) {
// 	if l.head == nil {
// 		return nil, errors.New("list is emptly")
// 	}

// 	if index+1 > l.count {
// 		return nil, errors.New("index out of range")
// 	}

// 	node := l.head
// 	for i := 0; i < index; i++ {
// 		node = node.GetNext()
// 	}
// 	return node, nil
// }

func (l *DoublyLinkedList) Remove(node *DoublyNode) {
	if node == nil {
		return
	}

	if node == l.head {
		l.head = node.next
	}

	if node == l.last {
		l.last = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	node.prev = nil
	node.next = nil
	l.count--
}

func (l *DoublyLinkedList) InsertAtFront(node *DoublyNode) {
	node.prev = nil
	node.next = l.head

	if l.head != nil {
		l.head.prev = node
	} else {
		l.last = node
	}

	l.head = node
	l.count++
}
