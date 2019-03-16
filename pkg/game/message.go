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

	MouseClickedConfig = message.Config{
		Id:   41,
		Size: 6,
		New:  func() message.Message { return &MouseClicked{} },
	}

	MouseActivityRecordedConfig = message.Config{
		Id:   34,
		Size: message.SizeVariableByte,
		New:  func() message.Message { return &MouseActivityRecorded{} },
	}

	KeyTypedConfig = message.Config{
		Id:   67,
		Size: message.SizeVariableShort,
		New:  func() message.Message { return &KeyTyped{} },
	}

	CameraRotatedConfig = message.Config{
		Id:   39,
		Size: 4,
		New:  func() message.Message { return &CameraRotated{} },
	}

	ButtonPressedConfig = message.Config{
		Id:   68,
		Size: 9,
		New:  func() message.Message { return &ButtonPressed{} },
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

	OpenChildInterfaceConfig = message.Config{
		Id:   77,
		Size: 7,
		New:  func() message.Message { return &OpenChildInterface{} },
	}

	CloseChildInterfaceConfig = message.Config{
		Id:   9,
		Size: 4,
		New:  func() message.Message { return &CloseChildInterface{} },
	}

	ClearInputBoxConfig = message.Config{
		Id:   52,
		Size: 0,
		New:  func() message.Message { return &ClearInputBox{} },
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

	buf.EndBitAccess()

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

type ButtonPressed struct {
	// NOTE(Sino): The option was previously defined by the id of the packet so packet 68
	// would represent button press with option 1 and another packet, the button press with
	// option 2. These options represent the right-click context menu options that are displayed
	// within various interfaces, such as the one for the bank.  By default, the given option is
	// the first one. However, I've simplified this by including the id of the option within the
	// payload encoded as an 8-bit integer value, thus reducing the amount of packets from 10 to 1.
	// The remaining packets are marked as no longer used within the client and can therefore
	// be used for something else. Something... custom.
	Option uint8

	Interface    int
	Button       int
	Component    uint16
	SubComponent uint16 // TODO identify this
}

func (ButtonPressed) Config() message.Config { return ButtonPressedConfig }

func (b *ButtonPressed) Decode(buf *buffer.ByteBuffer, length int) error {
	var err error

	if b.Option, err = buf.GetUint8(); err != nil {
		return err
	}

	var pack uint32
	if pack, err = buf.GetUint32(); err != nil {
		return err
	}

	b.Interface = int(pack >> 16)
	b.Button = int(pack & 0xFFFF)

	if b.Component, err = buf.GetUint16(); err != nil {
		return err
	}

	if b.SubComponent, err = buf.GetUint16(); err != nil {
		return err
	}

	return err
}

type MouseClicked struct {
	DeltaTime  int16
	RightClick bool
	X          uint16
	Y          uint16
}

func (MouseClicked) Config() message.Config { return MouseClickedConfig }

func (m *MouseClicked) Decode(buf *buffer.ByteBuffer, length int) error {
	var err error

	var pack uint16
	if pack, err = buf.GetUint16(); err != nil {
		return err
	}

	m.DeltaTime = int16(pack >> 1)
	m.RightClick = (pack & 1) == 1

	if m.X, err = buf.GetUint16(); err != nil {
		return err
	}

	if m.Y, err = buf.GetUint16(); err != nil {
		return err
	}

	return err
}

type MouseActivityRecorded struct {
}

func (MouseActivityRecorded) Config() message.Config { return MouseActivityRecordedConfig }

func (MouseActivityRecorded) Decode(buf *buffer.ByteBuffer, length int) error {
	// TODO for now just skip this. it is a quite complex packet

	return buf.Skip(length)
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

type ClearInputBox struct{}

func (*ClearInputBox) Config() message.Config { return ClearInputBoxConfig }

func (ClearInputBox) Encode(buf *buffer.ByteBuffer) error {
	return nil
}

type OpenChildInterface struct {
	Parent   uint32
	Id       uint16
	Behavior uint8
}

func (*OpenChildInterface) Config() message.Config { return OpenChildInterfaceConfig }

func (o *OpenChildInterface) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint32(o.Parent); err != nil {
		return err
	}

	if err := buf.PutUint16(o.Id); err != nil {
		return err
	}

	if err := buf.PutUint8(o.Behavior); err != nil {
		return err
	}

	return nil
}

type CloseChildInterface struct {
	Parent uint32
}

func (*CloseChildInterface) Config() message.Config { return CloseChildInterfaceConfig }

func (o *CloseChildInterface) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint32(o.Parent); err != nil {
		return err
	}

	return nil
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
