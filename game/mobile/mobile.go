package mobile

import "github.com/sprinkle-it/donut/game2ity"

const(
    walkingSteps = 1
    runningSteps = 2
)

type Update uint8

const (
    GraphicPlayed Update = 0x1
    Teleported    Update = 0x2
    Walked        Update = 0x4
)

func (u *Update) Set(update Update) { *u |= update }

func (u Update) Test(update Update) bool { return u&update != 0 }

func (u *Update) Clear() { *u = 0 }

type Mobile struct {
    entity.Entity

    id uint16

    position         Position
    previousPosition Position

    path    Path
    running bool

    traveled Path

    graphic Graphic
    updates Update
}

func (m *Mobile) Id() uint16 {
    return m.id
}

func (m *Mobile) Clear() {
    m.traveled.Clear()
    m.updates.Clear()
}

// Sets the path that the mobile will traverse. This function will reset the path that the mobile is currently
// traversing.
func (m *Mobile) StartPath(path Path) {
    m.path.Clear()
    m.path.PushAll(path)
}

// Walks the current path. The number of steps that will be walked will be the minimum of either the number of steps
// remaining in the path the mobile is currently walking or the maximum number of steps the mobile can take per turn.
// This function will immediately return if there are no steps in the path. All of the traversed steps will be pushed
// to the traveled path. The previous position will be set to be the position of where the mobile was currently located
// before this function is called.
func (m *Mobile) Walk() {
    if len(m.path) < 1 {
        return
    }

    count := m.stepsPerTurn()
    if count > len(m.path) {
        count = len(m.path)
    }

    m.markPosition()

    for count > 0 {
        dir := m.path.Poll()
        m.position.X += dir.X()
        m.position.Z += dir.Z()
        m.traveled.Push(dir)
        count--
    }

    m.updates.Set(Walked)
}

// Gets the number of steps per turn that the mobile should traverse when walking a path.
func (m *Mobile) stepsPerTurn() int {
    if m.running {
        return runningSteps
    }
    return walkingSteps
}

// Teleports the mobile to the provided world coordinates. Teleporting the mobile will clear the mobile's current path.
// The previous position will be set to be the position of where the mobile was currently located before this function
// was called.
func (m *Mobile) Teleport(level int8, x, z int16) {
    m.path.Clear()

    m.markPosition()

    m.position.Level = level
    m.position.X = x
    m.position.Z = z

    m.updates.Set(Teleported)
}

// Sets the previous position as the current position. Any operation that updates the mobiles current position must
// perform this action.
func (m *Mobile) markPosition() {
    m.previousPosition.Level = m.position.Level
    m.previousPosition.X = m.position.X
    m.previousPosition.Z = m.position.Z
}

func (m *Mobile) PlayGraphic(graphic Graphic) {
    m.graphic = graphic
    m.updates.Set(GraphicPlayed)
}

func (m *Mobile) Graphic() Graphic {
    return m.graphic
}
