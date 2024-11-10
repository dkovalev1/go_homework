package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first *ListItem
	last  *ListItem
	count int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	elem := new(ListItem)
	elem.Value = v
	elem.Next = l.first
	if l.first != nil {
		l.first.Prev = elem
	}

	l.first = elem

	if l.last == nil {
		l.last = elem
	}
	l.count++
	return elem
}

func (l *list) PushBack(v interface{}) *ListItem {
	elem := new(ListItem)
	elem.Value = v

	elem.Prev = l.last

	if l.last != nil {
		l.last.Next = elem
	}

	l.last = elem
	if l.first == nil {
		l.first = elem
	}

	l.count++
	return elem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.first = i.Next
	}

	if i.Next == nil {
		l.last = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = nil
	}
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		// It's already at front
		return
	}

	i.Prev.Next = i.Next

	if i.Next == nil {
		// Last element
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	i.Next = l.first
	i.Prev = nil
	l.first = i
}
