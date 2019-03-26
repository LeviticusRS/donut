package game

import (
    "github.com/sprinkle-it/donut/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    HandshakeConfig = message.Config{
        Id:   14,
        Size: 0,
        New:  message.Singleton(Handshake),
    }

    AuthenticateConfig = message.Config{
        Id:   16,
        Size: message.SizeVariableShort,
        New:  func() message.Message { return &Authenticate{} },
    }

    WindowUpdateConfig = message.Config{
        Id:   35,
        Size: 5,
        New:  func() message.Message { return &WindowUpdate{} },
    }

    FocusChangedConfig = message.Config{
        Id:   73,
        Size: 1,
        New:  func() message.Message { return &FocusChanged{} },
    }

    SceneRebuiltConfig = message.Config{
        Id:   76,
        Size: 0,
        New:  message.Singleton(SceneRebuilt),
    }

    HeartbeatConfig = message.Config{
        Id:   122,
        Size: 0,
        New:  message.Singleton(Heartbeat),
    }

    KeyTypedConfig = message.Config{
        Id: 67,
        Size: message.SizeVariableShort,
        New: func() message.Message { return &KeyTyped{} },
    }

    CameraRotatedConfig = message.Config{
        Id: 39,
        Size: 4,
        New: func() message.Message { return &CameraRotated{} },
    }

    ReadyConfig = message.Config{
        Id:   0,
        Size: 8,
        New:  func() message.Message { return &Ready{} },
    }

    InitializeSceneConfig = message.Config{
        Id:   0,
        Size: message.SizeVariableShort,
        New:  func() message.Message { return &InitializeScene{} },
    }

    SuccessConfig = message.Config{
        Id:   2,
        Size: message.SizeVariableByte,
        New:  func() message.Message { return &Success{} },
    }

    SetHudConfig = message.Config{
        Id:   84,
        Size: 2,
        New:  func() message.Message { return &SetHud{} },
    }

    PlayerUpdateConfig = message.Config{
        Id:   79,
        Size: message.SizeVariableShort,
        New:  func() message.Message { return &PlayerUpdate{} },
    }

    Handshake    = &handshake{}
    Heartbeat    = &heartbeat{}
    SceneRebuilt = &sceneRebuilt{}
)

type handshake struct{}

func (handshake) Config() message.Config { return HandshakeConfig }

func (handshake) Decode(buf *buffer.ByteBuffer, length int) error { return nil }

type Ready struct {
    AuthenticationKey uint64
}

type Authenticate struct {
}

func (Authenticate) Config() message.Config {
    return AuthenticateConfig
}

func (Authenticate) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}

func (*Ready) Config() message.Config { return ReadyConfig }

func (r *Ready) Encode(buf *buffer.ByteBuffer) error { return buf.PutUint64(r.AuthenticationKey) }

type Success struct {
    UserGroup uint8
    Moderator bool
    PlayerId  uint16
    Members   bool
}

func (*Success) Config() message.Config { return SuccessConfig }

func (s *Success) Encode(buf *buffer.ByteBuffer) error {
    if err := buf.PutUint8(s.UserGroup); err != nil {
        return err
    }

    if err := buf.PutBool(s.Moderator); err != nil {
        return err
    }

    if err := buf.PutUint16(s.PlayerId); err != nil {
        return err
    }

    if err := buf.PutBool(s.Members); err != nil {
        return err
    }

    return nil
}

type InitializeScene struct {
    Position        Position
    PlayerPositions [2046]Position
}

func (*InitializeScene) Config() message.Config { return InitializeSceneConfig }

func (r *InitializeScene) Encode(buf *buffer.ByteBuffer) error {
    buf.StartBitAccess()

    r.Position.EncodeHash(buf)

    for i := 0; i < len(r.PlayerPositions); i++ {
        r.PlayerPositions[i].EncodeBlockHash(buf)
    }

    buf.FinishBitAccess()

    if err := buf.PutUint16(r.Position.ChunkX()); err != nil {
        return err
    }

    if err := buf.PutUint16(r.Position.ChunkZ()); err != nil {
        return err
    }

    return nil
}

type WindowUpdate struct {
}

func (WindowUpdate) Config() message.Config { return WindowUpdateConfig }

func (WindowUpdate) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}

type heartbeat struct {
}

func (heartbeat) Config() message.Config { return HeartbeatConfig }

func (heartbeat) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}

// Inbound message from the client that lets the server know that the scene has been successfully rebuilt.
type sceneRebuilt struct{}

func (sceneRebuilt) Config() message.Config { return SceneRebuiltConfig }

func (sceneRebuilt) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}

type FocusChanged struct {
    Focused bool
}

func (*FocusChanged) Config() message.Config { return FocusChangedConfig }

func (f *FocusChanged) Decode(buf *buffer.ByteBuffer, length int) error {
    var err error
    if f.Focused, err = buf.GetBool(); err != nil {
        return err
    }
    return nil
}

type SetHud struct {
    Id uint16
}

func (*SetHud) Config() message.Config { return SetHudConfig }

func (s *SetHud) Encode(buf *buffer.ByteBuffer) error {
    return buf.PutUint16(s.Id)
}

type KeyTyped struct {

}

func (KeyTyped) Config() message.Config {
    return KeyTypedConfig
}

func (KeyTyped) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}

type CameraRotated struct {

}

func (CameraRotated) Config() message.Config {
    return CameraRotatedConfig
}

func (CameraRotated) Decode(buf *buffer.ByteBuffer, length int) error {
    return nil
}


