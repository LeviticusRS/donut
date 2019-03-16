package game

import "github.com/sprinkle-it/donut/pkg/client"

type PlayerList struct {
    new    PlayerFactory
    arr    []*Player
    unused []uint16
    active []uint16
}

func NewPlayerList(capacity int, new PlayerFactory) PlayerList {
    unused := make([]uint16, 0, capacity)
    for i := 1; i < cap(unused); i++ {
        unused = append(unused, uint16(i))
    }

    return PlayerList{
        new:    new,
        arr:    make([]*Player, capacity),
        unused: unused,
        active: make([]uint16, 0, capacity),
    }
}

func (l *PlayerList) New(cli *client.Client) (*Player, bool) {
    if len(l.unused) == 0 {
        return nil, false
    }

    // Pop an id from the list of unused identifiers.
    var id uint16
    id, l.unused = l.unused[0], l.unused[1:]

    // Create and store the player.
    player := l.new(cli, id)
    l.arr[player.Id()] = player

    // Append the players identifier to the active list.
    l.active = append(l.active, player.Id())

    return player, true
}

func (l *PlayerList) Capacity() int {
    return cap(l.arr)
}

func (l *PlayerList) Get(n uint16) *Player {
    return l.arr[n]
}
