package mobile

type Direction int

const (
    East      Direction = 0
    NorthEast Direction = 1
    North     Direction = 2
    NorthWest Direction = 3
    West      Direction = 4
    SouthWest Direction = 5
    South     Direction = 6
    SouthEast Direction = 7
)

func (d Direction) X() int16 {
    switch d {
    case East, NorthEast, SouthEast:
        return 1
    case West, NorthWest, SouthWest:
        return -1
    default:
        return 0
    }
}

func (d Direction) Z() int16 {
    switch d {
    case North, NorthEast, NorthWest:
        return 1
    case South, SouthWest, SouthEast:
        return -1
    default:
        return 0
    }
}

type Path []Direction

func (p *Path) Poll() Direction {
    old := *p
    v := old[0]
    *p = old[1:]
    return v
}

func (p *Path) Push(dir Direction) {
    *p = append(*p, dir)
}

func (p *Path) PushAll(path Path) {
    *p = append(*p, path...)
}

func (p *Path) Clear() {
    old := *p
    *p = old[:0]
}
