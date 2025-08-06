package main

import (
	"errors"
	"fmt"
)

type DoublyNode struct {
	data interface{}
	next *DoublyNode
	prev *DoublyNode
}

func NewDoublyNode(data interface{}) *DoublyNode {
	return &DoublyNode{data: data}
}

func (n *DoublyNode) SetData(data interface{}) {
	n.data = data
}

func (n *DoublyNode) GetData() interface{} {
	return n.data
}

func (n *DoublyNode) SetNext(next *DoublyNode) {
	n.next = next
}

func (n *DoublyNode) GetNext() *DoublyNode {
	return n.next
}

func (n *DoublyNode) SetPrev(prev *DoublyNode) {
	n.prev = prev
}

func (n *DoublyNode) GetPrev() (*DoublyNode, error) {
	if n.prev == nil {
		return nil, errors.New("no previous node")
	}
	return n.prev, nil
}

func (n *DoublyNode) ToString() string {
	return fmt.Sprintf("%v", n.data)
}
