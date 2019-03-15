package game

import (
	"github.com/sprinkle-it/donut/pkg/buffer"
	"github.com/sprinkle-it/donut/pkg/message"
)

var (
	HandshakeConfig = message.Config{
		Id:   14,
		Size: 0,
		New:  message.Singleton(Handshake),
	}

	NewLoginConfig = message.Config{
		Id:   16,
		Size: message.SizeVariableShort,
		New:  func() message.Message { return &NewLogin{} },
	}

	ReconnectConfig = message.Config{
		Id:   18,
		Size: message.SizeVariableShort,
		New:  func() message.Message { return &Reconnect{} },
	}

	WindowUpdateConfig = message.Config{
		Id:   35,
		Size: 5,
		New:  func() message.Message { return &WindowUpdate{} },
	}

	FocusChangedConfig = message.Config{
		Id:   73,
		Size: 1,
		New:  message.Singleton(SceneRebuilt),
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

	ReadyConfig = message.Config{
		Id:   0,
		Size: 8,
		New:  func() message.Message { return &Ready{} },
	}

	RebuildSceneConfig = message.Config{
		Id:   0,
		Size: message.SizeVariableShort,
		New:  func() message.Message { return &RebuildScene{} },
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

type NewLogin struct {
	Authenticate
}

type Reconnect struct {
	Authenticate
}

type Authenticate struct {
	Email    string
	Password string

	Seeds             [4]uint32
	AuthenticationKey uint64

	ClientVersion    uint32
	ArchiveChecksums [18]uint32

	LowMemory     bool
	ResizableMode bool
}

func (NewLogin) Config() message.Config {
	return NewLoginConfig
}

func (n *NewLogin) Decode(buf *buffer.ByteBuffer, length int) error {
	return n.decodeAuthenticate(buf, n.Config().Id, length)
}

func (Reconnect) Config() message.Config {
	return ReconnectConfig
}

func (r *Reconnect) Decode(buf *buffer.ByteBuffer, length int) error {
	return r.decodeAuthenticate(buf, r.Config().Id, length)
}

func (a *Authenticate) decodeAuthenticate(buf *buffer.ByteBuffer, id uint8, length int) error {
	var err error

	if a.ClientVersion, err = buf.GetUint32(); err != nil {
		return err
	}

	for i := 0; i < len(a.Seeds); i++ {
		if a.Seeds[i], err = buf.GetUint32(); err != nil {
			return err
		}
	}

	if a.AuthenticationKey, err = buf.GetUint64(); err != nil {
		return err
	}

	reconnecting := id == ReconnectConfig.Id
	if !reconnecting {
		if a.Password, err = buf.GetCString(); err != nil {
			return err
		}
	}

	if a.Email, err = buf.GetCString(); err != nil {
		return err
	}

	screenPack, err := buf.GetUint8()
	if err != nil {
		return err
	}

	a.ResizableMode = screenPack>>1 == 1
	a.LowMemory = screenPack&1 == 1

	for i := 0; i < len(a.ArchiveChecksums); i++ {
		if a.ArchiveChecksums[i], err = buf.GetUint32(); err != nil {
			return err
		}
	}

	return err
}

func (*Ready) Config() message.Config { return ReadyConfig }

func (r *Ready) Encode(buf *buffer.ByteBuffer) error { return buf.PutUint64(r.AuthenticationKey) }

type Success struct {
}

func (*Success) Config() message.Config { return SuccessConfig }

func (s *Success) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint8(2); err != nil {
		return err
	}

	if err := buf.PutUint8(0); err != nil {
		return err
	}

	if err := buf.PutUint16(1); err != nil {
		return err
	}

	if err := buf.PutUint8(1); err != nil {
		return err
	}

	return nil
}

type RebuildScene struct {
	InitializePlayerPositions InitializePlayerPositions
	ChunkX                    uint16
	ChunkZ                    uint16
}

func (*RebuildScene) Config() message.Config { return RebuildSceneConfig }

func (r *RebuildScene) Encode(buf *buffer.ByteBuffer) error {
	if err := r.InitializePlayerPositions.Encode(buf); err != nil {
		return err
	}

	if err := buf.PutUint16(r.ChunkX); err != nil {
		return err
	}

	if err := buf.PutUint16(r.ChunkZ); err != nil {
		return err
	}

	return nil
}

type InitializePlayerPositions struct {
	LocalPosition Position
	Positions     [2046]Position
}

func (p *InitializePlayerPositions) Encode(buf *buffer.ByteBuffer) error {
	buf.StartBitAccess()

	p.LocalPosition.EncodeHash(buf)

	for i := 0; i < len(p.Positions); i++ {
		p.Positions[i].EncodeBlockHash(buf)
	}

	buf.EndBitAccess()
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

type sceneRebuilt struct{}

func (sceneRebuilt) Config() message.Config { return SceneRebuiltConfig }

func (sceneRebuilt) Decode(buf *buffer.ByteBuffer, length int) error {
	return nil
}

type FocusChanged struct {
}

func (FocusChanged) Config() message.Config { return FocusChangedConfig }

func (FocusChanged) Decode(buf *buffer.ByteBuffer, length int) error {
	return nil
}

type SetHud struct {
	Id uint16
}

func (*SetHud) Config() message.Config { return SetHudConfig }

func (s *SetHud) Encode(buf *buffer.ByteBuffer) error {
	return buf.PutUint16(s.Id)
}
