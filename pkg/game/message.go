package game

import (
	"fmt"
	"github.com/sprinkle-it/donut/pkg/account"
	"github.com/sprinkle-it/donut/pkg/buffer"
	"github.com/sprinkle-it/donut/pkg/message"
	"reflect"
	"strings"
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

	ClientPerformanceMeasuredConfig = message.Config{
		Id:   111,
		Size: 10,
		New:  func() message.Message { return &ClientPerformanceMeasured{} },
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

	SetPlayerContextMenuOptionConfig = message.Config{
		Id:   66,
		Size: message.SizeVariableByte,
		New:  func() message.Message { return &SetPlayerContextMenuOption{} },
	}

	OpenChildInterfaceConfig = message.Config{
		Id:   77,
		Size: 7,
		New:  func() message.Message { return &OpenChildInterface{} },
	}

	RelocateChildInterfaceConfig = message.Config{
		Id:   82,
		Size: 8,
		New:  func() message.Message { return &RelocateChildInterface{} },
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

	DisplaySystemMessageConfig = message.Config{
		Id:   3,
		Size: message.SizeVariableByte,
		New:  func() message.Message { return &DisplaySystemMessage{} },
	}

	SetEnergyConfig = message.Config{
		Id:   60,
		Size: 1,
		New:  func() message.Message { return &SetEnergy{} },
	}

	SetWeightConfig = message.Config{
		Id:   71,
		Size: 2,
		New:  func() message.Message { return &SetWeight{} },
	}

	SetMinimapStateConfig = message.Config{
		Id:   74,
		Size: 1,
		New:  func() message.Message { return &SetMinimapState{} },
	}

	SetSystemUpdateTimerConfig = message.Config{
		Id:   72,
		Size: 2,
		New:  func() message.Message { return &SetSystemUpdateTimer{} },
	}

	ClearPerspectiveCameraConfig = message.Config{
		Id:   2,
		Size: 0,
		New:  func() message.Message { return &ClearPerspectiveCamera{} },
	}

	ClearInventoryConfig = message.Config{
		Id:   7,
		Size: 4,
		New:  func() message.Message { return &ClearInventory{} },
	}

	LogoutConfig = message.Config{
		Id:   1,
		Size: 0,
		New:  func() message.Message { return &Logout{} },
	}

	TargetPatchConfig = message.Config{
		Id:   64,
		Size: 2,
		New:  func() message.Message { return &TargetPatch{} },
	}

	ClearPatchConfig = message.Config{
		Id:   25,
		Size: 2,
		New:  func() message.Message { return &ClearPatch{} },
	}

	Set32BitVariableConfig = message.Config{
		Id:   4,
		Size: 6,
		New:  func() message.Message { return &Set32BitVariable{} },
	}

	Set8BitVariableConfig = message.Config{
		Id:   63,
		Size: 3,
		New:  func() message.Message { return &Set8BitVariable{} },
	}

	ClearVariablesConfig = message.Config{
		Id:   78,
		Size: 0,
		New:  func() message.Message { return &ClearVariables{} },
	}

	RevertVariablesConfig = message.Config{
		Id:   73,
		Size: 0,
		New:  func() message.Message { return &RevertVariables{} },
	}

	SetSkillConfig = message.Config{
		Id:   22,
		Size: 6,
		New:  func() message.Message { return &SetSkill{} },
	}

	ModifyLabelTextConfig = message.Config{
		Id:   19,
		Size: message.SizeVariableShort,
		New:  func() message.Message { return &ModifyLabelText{} },
	}

	ModifyLabelColourConfig = message.Config{
		Id:   24,
		Size: 6,
		New:  func() message.Message { return &ModifyLabelColour{} },
	}

	InvokeInterfaceScriptConfig = message.Config{
		Id:   62,
		Size: message.SizeVariableShort,
		New:  func() message.Message { return &InvokeInterfaceScript{} },
	}

	ToggleComponentVisibilityConfig = message.Config{
		Id:   21,
		Size: 5,
		New:  func() message.Message { return &ToggleComponentVisibility{} },
	}

	RequestClientPerformanceConfig = message.Config{
		Id:   69,
		Size: 8,
		New:  func() message.Message { return &RequestClientPerformance{} },
	}

	GroupedEntityUpdateConfig = message.Config{
		Id:   17,
		Size: message.SizeVariableShort,
		New:  func() message.Message { return &GroupedEntityUpdate{} },
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
	SubComponent uint16 // TODO is truly a subcomponent? need confirmation or proper identification
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

type ClientPerformanceMeasured struct {
	FirstKey  uint32
	SecondKey uint32
	FPS       uint8
	GCTime    uint8
}

func (ClientPerformanceMeasured) Config() message.Config { return ClientPerformanceMeasuredConfig }

func (c ClientPerformanceMeasured) Decode(buf *buffer.ByteBuffer, length int) (err error) {
	if c.GCTime, err = buf.GetUint8(); err != nil {
		return
	}

	if c.FPS, err = buf.GetUint8(); err != nil {
		return
	}

	if c.FirstKey, err = buf.GetUint32(); err != nil {
		return
	}

	if c.SecondKey, err = buf.GetUint32(); err != nil {
		return
	}

	return
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

type GroupedEntityUpdate struct {
	X        uint8
	Z        uint8
	Children []message.Outbound
}

func (*GroupedEntityUpdate) Config() message.Config { return GroupedEntityUpdateConfig }

func (g *GroupedEntityUpdate) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint8(g.Z); err != nil {
		return err
	}

	if err := buf.PutUint8(g.X); err != nil {
		return err
	}

	for _, msg := range g.Children {
		if err := buf.PutUint8(msg.Config().Id); err != nil {
			return err
		}

		if err := msg.Encode(buf); err != nil {
			return err
		}
	}

	return nil
}

type ModifyLabelText struct {
	Parent uint32
	Text   string
}

func (*ModifyLabelText) Config() message.Config { return ModifyLabelTextConfig }

func (m *ModifyLabelText) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint32(m.Parent); err != nil {
		return err
	}

	if err := buf.PutCString(m.Text); err != nil {
		return err
	}

	return nil
}

type ModifyLabelColour struct {
	Parent uint32
	R      int
	G      int
	B      int
}

func (*ModifyLabelColour) Config() message.Config { return ModifyLabelColourConfig }

func (m *ModifyLabelColour) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint32(m.Parent); err != nil {
		return err
	}

	if err := buf.PutUint16(uint16(m.R<<10&31 | m.G<<5&31 | m.B&31)); err != nil {
		return err
	}

	return nil
}

type SetPlayerContextMenuOption struct {
	Slot        uint8
	Label       string
	Prioritized bool
}

func (*SetPlayerContextMenuOption) Config() message.Config { return SetPlayerContextMenuOptionConfig }

func (s SetPlayerContextMenuOption) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutBool(s.Prioritized); err != nil {
		return err
	}

	if err := buf.PutUint8(s.Slot); err != nil {
		return err
	}

	if err := buf.PutCString(s.Label); err != nil {
		return err
	}

	return nil
}

type ClearInputBox struct{}

func (*ClearInputBox) Config() message.Config { return ClearInputBoxConfig }

func (ClearInputBox) Encode(buf *buffer.ByteBuffer) error {
	return nil
}

type ToggleComponentVisibility struct {
	Parent uint32
	Hidden bool
}

func (*ToggleComponentVisibility) Config() message.Config { return ToggleComponentVisibilityConfig }

func (t ToggleComponentVisibility) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint32(t.Parent); err != nil {
		return err
	}

	if err := buf.PutBool(t.Hidden); err != nil {
		return err
	}

	return nil
}

type RequestClientPerformance struct {
	FirstKey  uint32
	SecondKey uint32
}

func (*RequestClientPerformance) Config() message.Config { return RequestClientPerformanceConfig }

func (r RequestClientPerformance) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint32(r.FirstKey); err != nil {
		return err
	}

	if err := buf.PutUint32(r.SecondKey); err != nil {
		return err
	}

	return nil
}

type ScriptArgument interface{}

type InvokeInterfaceScript struct {
	Id        uint32
	Arguments []ScriptArgument
}

func (i InvokeInterfaceScript) Identifiers() (string, error) {
	var bldr strings.Builder
	for j := len(i.Arguments) - 1; j >= 0; j-- {
		switch value := i.Arguments[j].(type) {
		case string:
			bldr.WriteRune('s')
		case int:
			bldr.WriteRune('i')
		default:
			return "", fmt.Errorf("given type of %v is unsupported in interface script", reflect.TypeOf(value))
		}
	}

	return bldr.String(), nil
}

func (*InvokeInterfaceScript) Config() message.Config { return InvokeInterfaceScriptConfig }

func (i InvokeInterfaceScript) Encode(buf *buffer.ByteBuffer) error {
	scriptIdentifiers, err := i.Identifiers()
	if err != nil {
		return err
	}

	if err := buf.PutCString(scriptIdentifiers); err != nil {
		return err
	}

	for argumentIdx, character := range scriptIdentifiers {
		var err error
		var value = i.Arguments[argumentIdx]

		if string(character) == "s" {
			err = buf.PutCString(value.(string))
		} else {
			err = buf.PutUint32(uint32(value.(int)))
		}

		if err != nil {
			return err
		}
	}

	if err := buf.PutUint32(i.Id); err != nil {
		return err
	}

	return nil
}

type DisplaySystemMessage struct {
	Type            uint8
	InteractingWith *account.DisplayName
	Text            string
}

func (*DisplaySystemMessage) Config() message.Config { return DisplaySystemMessageConfig }

func (s *DisplaySystemMessage) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint8(s.Type); err != nil { // TODO encode this as a smart
		return nil
	}

	if err := buf.PutBool(s.InteractingWith != nil); err != nil {
		return err
	}

	if s.InteractingWith != nil {
		if err := buf.PutCString(string(*s.InteractingWith)); err != nil {
			return err
		}
	}

	if err := buf.PutCString(s.Text); err != nil {
		return err
	}

	return nil
}

type SetSystemUpdateTimer struct {
	Ticks uint16
}

func (*SetSystemUpdateTimer) Config() message.Config { return SetSystemUpdateTimerConfig }

func (s *SetSystemUpdateTimer) Encode(buf *buffer.ByteBuffer) error {
	return buf.PutUint16(s.Ticks)
}

type ClearPerspectiveCamera struct{}

func (*ClearPerspectiveCamera) Config() message.Config { return ClearPerspectiveCameraConfig }

func (s *ClearPerspectiveCamera) Encode(buf *buffer.ByteBuffer) error {
	return nil
}

type ClearInventory struct {
	Parent uint32
}

func (*ClearInventory) Config() message.Config { return ClearInventoryConfig }

func (s *ClearInventory) Encode(buf *buffer.ByteBuffer) error {
	return buf.PutUint32(s.Parent)
}

type Logout struct{}

func (*Logout) Config() message.Config { return LogoutConfig }

func (s *Logout) Encode(buf *buffer.ByteBuffer) error {
	return nil
}

type TargetPatch struct {
	X uint8
	Z uint8
}

func (*TargetPatch) Config() message.Config { return TargetPatchConfig }

func (s *TargetPatch) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint8(s.Z); err != nil {
		return err
	}

	if err := buf.PutUint8(s.X); err != nil {
		return err
	}

	return nil
}

type ClearPatch struct {
	X uint8
	Z uint8
}

func (*ClearPatch) Config() message.Config { return ClearPatchConfig }

func (s *ClearPatch) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint8(s.Z); err != nil {
		return err
	}

	if err := buf.PutUint8(s.X); err != nil {
		return err
	}

	return nil
}

type SetSkill struct {
	Id         uint8
	Level      uint8
	Experience uint32
}

func (*SetSkill) Config() message.Config { return SetSkillConfig }

func (s *SetSkill) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint8(s.Id); err != nil {
		return err
	}

	if err := buf.PutUint8(s.Level); err != nil {
		return err
	}

	if err := buf.PutUint32(s.Experience); err != nil {
		return err
	}

	return nil
}

type Set32BitVariable struct {
	Id    uint16
	Value uint32
}

func (*Set32BitVariable) Config() message.Config { return Set32BitVariableConfig }

func (s *Set32BitVariable) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint16(s.Id); err != nil {
		return err
	}

	if err := buf.PutUint32(s.Value); err != nil {
		return err
	}

	return nil
}

type Set8BitVariable struct {
	Id    uint16
	Value uint8
}

func (*Set8BitVariable) Config() message.Config { return Set8BitVariableConfig }

func (s *Set8BitVariable) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint16(s.Id); err != nil {
		return err
	}

	if err := buf.PutUint8(s.Value); err != nil {
		return err
	}

	return nil
}

type ClearVariables struct{}

func (*ClearVariables) Config() message.Config { return ClearVariablesConfig }

func (s *ClearVariables) Encode(buf *buffer.ByteBuffer) error {
	return nil
}

type RevertVariables struct{}

func (*RevertVariables) Config() message.Config { return RevertVariablesConfig }

func (s *RevertVariables) Encode(buf *buffer.ByteBuffer) error {
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

type SetEnergy struct {
	Percentage uint8
}

func (*SetEnergy) Config() message.Config { return SetEnergyConfig }

func (s *SetEnergy) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint8(s.Percentage); err != nil {
		return err
	}

	return nil
}

type SetWeight struct {
	Kilograms uint16
}

func (*SetWeight) Config() message.Config { return SetWeightConfig }

func (s *SetWeight) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint16(s.Kilograms); err != nil {
		return err
	}

	return nil
}

type SetMinimapState struct {
	Id uint8
}

func (*SetMinimapState) Config() message.Config { return SetMinimapStateConfig }

func (s *SetMinimapState) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint8(s.Id); err != nil {
		return err
	}

	return nil
}

type RelocateChildInterface struct {
	ParentFrom uint32
	ParentTo   uint32
}

func (*RelocateChildInterface) Config() message.Config { return RelocateChildInterfaceConfig }

func (r *RelocateChildInterface) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint32(r.ParentTo); err != nil {
		return err
	}

	if err := buf.PutUint32(r.ParentFrom); err != nil {
		return err
	}

	return nil
}

type CloseChildInterface struct {
	Parent uint32
}

func (*CloseChildInterface) Config() message.Config { return CloseChildInterfaceConfig }

func (c *CloseChildInterface) Encode(buf *buffer.ByteBuffer) error {
	if err := buf.PutUint32(c.Parent); err != nil {
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
