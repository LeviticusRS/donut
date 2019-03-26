package mobile

type Node interface {
    Id() uint16
}

type List struct {
    new     func(uint16) Node
    entries []Node
    unused  []uint16
    active  []uint16
}

func NewList(capacity int, new func(uint16) Node) List {
    unused := make([]uint16, 0, capacity)
    for i := 1; i < cap(unused); i++ {
        unused = append(unused, uint16(i))
    }

    return List{
        new:     new,
        entries: make([]Node, capacity),
        unused:  unused,
        active:  make([]uint16, 0, capacity),
    }
}

func (l *List) Create() (Node, bool) {
    if len(l.unused) == 0 {
        return nil, false
    }

    var id uint16
    id, l.unused = l.unused[0], l.unused[1:]

    entry := l.new(id)
    l.entries[id] = entry

    l.active = append(l.active, id)

    return entry, true
}

func (l *List) Capacity() int {
    return cap(l.entries)
}

func (l *List) Exists(id uint16) bool {
    return l.entries[id] != nil
}

func (l *List) Get(id uint16) Node {
    return l.entries[id]
}