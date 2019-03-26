package game

import (
    "github.com/sprinkle-it/donut/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

type SkippedFlag uint8

func (f SkippedFlag) Test(v bool) bool {
    switch v {
    case true:
        return f&0x1 != 0
    default:
        return f&0x1 == 0
    }
}

func (f *SkippedFlag) Mark() {
    *f |= 0x2
}

func (f *SkippedFlag) Update() {
    *f >>= 1
}

type PlayerSync struct {
    nearby          []uint16
    distant         []uint16
    skipped         []SkippedFlag
    nearbyActive    BlockList
    nearbyInactive  BlockList
    distantActive   BlockList
    distantInactive BlockList
}

func NewPlayerSync(local uint16, capacity int) PlayerSync {
    sync := PlayerSync{
        nearby:          make([]uint16, 0, capacity),
        distant:         make([]uint16, 0, capacity),
        skipped:         make([]SkippedFlag, capacity),
        nearbyActive:    make(BlockList, 0, capacity),
        nearbyInactive:  make(BlockList, 0, capacity),
        distantActive:   make(BlockList, 0, capacity),
        distantInactive: make(BlockList, 0, capacity),
    }

    // The local player identifier is the first element in the nearby list. It will never be removed because the local
    // player can never leave the scene since it is centered around the player.
    sync.nearby = append(sync.nearby, local)

    // Append each player identifier to the distant list skipping over the local player identifier. Afterward
    // distant players can be moved to the nearby list once they are within the local players scene.
    for i := 1; i < capacity; i++ {
        id := uint16(i)

        if id == local {
            continue
        }

        sync.distant = append(sync.distant, id)
    }

    return sync
}

func (s *PlayerSync) Process(w *World) *PlayerUpdate {
    s.nearbyActive = s.nearbyActive[:0]
    s.nearbyInactive = s.nearbyInactive[:0]
    s.distantActive = s.distantActive[:0]
    s.distantInactive = s.distantInactive[:0]

    s.encodeNearbyPlayers(w, &s.nearbyActive, false)
    s.encodeNearbyPlayers(w, &s.nearbyInactive, true)
    s.encodeDistantPlayers(w, &s.distantActive, false)
    s.encodeDistantPlayers(w, &s.distantInactive, true)

    for i := 0; i < len(s.skipped); i++ {
        s.skipped[i].Update()
    }

    return &PlayerUpdate{
        ActiveNearbyBlocks:    s.nearbyActive.Copy(),
        InactiveNearbyBlocks:  s.nearbyInactive.Copy(),
        ActiveDistantBlocks:   s.distantActive.Copy(),
        InactiveDistantBlocks: s.distantInactive.Copy(),
    }
}

func (s *PlayerSync) encodeNearbyPlayers(w *World, blocks *BlockList, active bool) {
    counter := -1
    for _, id := range s.nearby {
        if s.skipped[id].Test(active) {
            continue
        }
        s.skipped[id].Mark()
        counter++
    }

    if counter >= 0 {
        *blocks = append(*blocks, SkipBlock{Count: counter})
    }
}

func (s *PlayerSync) encodeDistantPlayers(w *World, blocks *BlockList, skipped bool) {
    counter := -1
    for _, id := range s.distant {
        if s.skipped[id].Test(skipped) {
            continue
        }
        s.skipped[id].Mark()
        counter++
    }

    if counter >= 0 {
        *blocks = append(*blocks, SkipBlock{Count: counter})
    }
}

type PlayerUpdate struct {
    ActiveNearbyBlocks    BlockList
    InactiveNearbyBlocks  BlockList
    ActiveDistantBlocks   BlockList
    InactiveDistantBlocks BlockList
}

func (*PlayerUpdate) Config() message.Config {
    return PlayerUpdateConfig
}

func (u *PlayerUpdate) Encode(buf *buffer.ByteBuffer) error {
    if err := u.ActiveNearbyBlocks.Encode(buf); err != nil {
        return err
    }

    if err := u.InactiveNearbyBlocks.Encode(buf); err != nil {
        return err
    }

    if err := u.InactiveDistantBlocks.Encode(buf); err != nil {
        return err
    }

    if err := u.ActiveDistantBlocks.Encode(buf); err != nil {
        return err
    }

    return nil
}

type BlockList []SyncBlock

func (l BlockList) Encode(buf *buffer.ByteBuffer) error {
    buf.StartBitAccess()

    for _, block := range l {
        if err := block.Encode(buf); err != nil {
            return err
        }
    }

    buf.FinishBitAccess()

    return nil
}

func (l BlockList) Copy() BlockList {
    arr := make([]SyncBlock, len(l))
    copy(arr, l)
    return arr
}

type SyncBlock interface {
    Encode(*buffer.ByteBuffer) error
}

// When players do not need to be updated they will be skipped. A skip block instructs the client to skip over a certain
// number of players in the descriptor list it is currently decoding.
type SkipBlock struct {
    Count int
}

func (b SkipBlock) Encode(buf *buffer.ByteBuffer) error {
    buf.PutBits(0, 1)
    switch {
    case b.Count < 1:
        buf.PutBits(uint32(b.Count), 2)
        return nil
    case b.Count >= 1 && b.Count <= 31:
        buf.PutBits(1, 2)
        buf.PutBits(uint32(b.Count), 5)
    case b.Count >= 32 && b.Count <= 255:
        buf.PutBits(2, 2)
        buf.PutBits(uint32(b.Count), 8)
    default:
        buf.PutBits(3, 2)
        buf.PutBits(uint32(b.Count), 11)
    }
    return nil
}
