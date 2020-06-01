package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	// Len of list
	Len() int
	// Front - get first item from list
	Front() *Item
	// Back - get last item from list
	Back() *Item
	// PushFront - add item to beginning
	PushFront(v interface{}) *Item
	// PushBack - add item to the end
	PushBack(v interface{}) *Item
	// Remove item from list
	Remove(i *Item)
	// MoveToFront - move item to the beginning
	MoveToFront(i *Item)
	// Get item by value
	Get(v interface{}) *Item
}

type value interface{}

type Item struct {
	Value value
	Next  *Item
	Prev  *Item
}

type list struct {
	items map[value]*Item
	first *Item
	last  *Item
}

func NewList() List {
	return &list{
		items: make(map[value]*Item),
	}
}

func (l *list) Len() int {
	return len(l.items)
}

func (l *list) Get(v interface{}) *Item {
	value, ok := v.(value)
	if !ok {
		return nil
	}

	i, ok := l.items[value]
	if !ok {
		return nil
	}
	return i
}

func (l *list) Front() *Item {
	return l.first
}

func (l *list) Back() *Item {
	return l.last
}

func (l *list) PushFront(v interface{}) *Item {
	i := &Item{
		Value: v,
	}

	if (l.Len()) > 0 {
		l.first.Prev = i
		i.Next = l.first
	}

	l.first = i
	if l.Len() == 0 {
		l.last = i
	}

	l.items[i.Value] = i

	return i
}

func (l *list) PushBack(v interface{}) *Item {
	i := &Item{
		Value: v,
	}

	if l.Len() > 0 {
		l.last.Next = i
		i.Prev = l.last
	}

	l.last = i
	if l.Len() == 0 {
		l.first = i
	}

	l.items[i.Value] = i

	return i
}

func (l *list) Remove(i *Item) {
	delete(l.items, i.Value)

	if i == l.first {
		l.first.Next.Prev = nil
		l.first = l.first.Next
		return
	}

	if i == l.last {
		l.last.Prev.Next = nil
		l.last = l.last.Prev
		return
	}

	i.Next.Prev = i.Prev
	i.Prev.Next = i.Next
}

func (l *list) MoveToFront(i *Item) {
	l.Remove(i)
	l.PushFront(i.Value)
}
