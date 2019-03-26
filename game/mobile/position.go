package mobile

type Position struct {
    Level int8
    X     int16
    Z     int16
}

func (p Position) Distance(o Position) int16 {
    dx := o.X - p.X
    dz := o.Z - p.Z

    if dx < 0 {
        dx = -dx
    }

    if dz < 0 {
        dz = -dz
    }

    if dx > dz {
        return dx
    }

    return dz
}
