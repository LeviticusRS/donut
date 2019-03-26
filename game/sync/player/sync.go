package player

import (
    "github.com/sprinkle-it/donut/game2ile"
    "github.com/sprinkle-it/donut/game2yer"
    "github.com/sprinkle-it/donut/game2c"
)

type State uint8

const (
    Active   State = 0x0
    Inactive State = 0x1
    Skipped  State = 0x2
)

type Sync struct {
    local   uint16
    nearby  []uint16
    distant []uint16
    updated []uint16
    states  []State
    blocks  sync.BlockList
    counter int
}

func NewSync(local uint16, capacity int) Sync {
    syn := Sync{
        local:   local,
        nearby:  make([]uint16, 0, capacity),
        distant: make([]uint16, 0, capacity),
        updated: make([]uint16, 0, capacity),
        states:  make([]State, capacity),
        blocks:  make(sync.BlockList, 0, capacity),
        counter: -1,
    }

    // The local player identifier is the first element in the nearby list. It will never be removed because the local
    // player can never leave the scene since the scene is centered at the player.
    syn.nearby = append(syn.nearby, local)

    // Append each player identifier to the distant list skipping over the local player identifier. Afterward
    // distant players can be moved to the nearby list once they are within the local players scene.
    for i := 1; i < capacity; i++ {
        id := uint16(i)

        if id == local {
            continue
        }

        syn.distant = append(syn.distant, id)
    }

    return syn
}

func (s *Sync) Process(list *mobile.List) {
    s.blocks.Clear()

    s.processNearby(list, Active)
    s.processNearby(list, Inactive)
    s.processDistant(list, Inactive)
    s.processDistant(list, Active)

    s.nearby = s.nearby[:0]
    s.distant = s.distant[:0]
    s.updated = s.updated[:0]

    for i := 1; i < list.Capacity(); i++ {
        id := uint16(i)

        if list.Exists(id) {
            s.nearby = append(s.nearby, id)
        } else {
            s.distant = append(s.distant, id)
        }

        s.states[i] >>= 1
    }
}

func (s *Sync) processNearby(list *mobile.List, state State) {
    s.blocks.Push(sync.StartList)

    loc := list.Get(s.local).(*player.Player)

    for _, id := range s.nearby {
        if (s.states[id] & 0x1) != state {
            continue
        }

        /*
        // Check if the player still exists. If not then we can alert the client to remove the player from the game.
        if !list.Exists(id) {
            s.finishSkip()
            s.blocks.Push(Remove)
            continue
        }

        plr := list.Get(id).(*player.Player)

        // Check that the player is still within the scene.
        if plr.Position.Distance(loc.Position) >= 16 {
            s.finishSkip()
            s.blocks.Push(Remove)
            continue
        }

        // Check that the player is active. If the player is not active then it can be skipped.
        if !plr.Active() {
            s.skip(id)
            continue
        }

        if plr.Updated() {
            s.blocks.Push(Updated)
            s.updated = append(s.updated, id)
        } else {
            s.blocks.Push(NotUpdated)
        }

        if plr.Moved() {

        }
        */
    }

    s.finishSkip()

    s.blocks.Push(sync.EndList)
}

func (s *Sync) processDistant(list *mobile.List, state State) {
    s.blocks.Push(sync.StartList)

    counter := -1
    for _, id := range s.distant {
        if (s.states[id] & 0x1) != state {
            continue
        }

        s.states[id] |= Skipped
        counter++
    }

    if counter >= 0 {
        s.blocks.Push(sync.Skip{Count: counter})
    }

    s.blocks.Push(sync.EndList)
}

func (s *Sync) skip(id uint16) {
    s.states[id] |= Skipped
    s.counter++
}

func (s *Sync) finishSkip() {
    if s.counter >= 0 {
        s.blocks.Push(sync.Skip{Count: s.counter})
        s.counter = -1
    }
}
