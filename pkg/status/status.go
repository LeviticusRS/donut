package status

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    okayDescriptor = message.Descriptor{
        Id:       0,
        Size:     0,
        Provider: message.ProvideSingleton(Okay),
    }

    invalidCredentialsDescriptor = message.Descriptor{
        Id:       3,
        Size:     0,
        Provider: message.ProvideSingleton(InvalidCredentials),
    }

    accountDisabledDescriptor = message.Descriptor{
        Id:       4,
        Size:     0,
        Provider: message.ProvideSingleton(AccountDisabled),
    }

    alreadyLoggedInDescriptor = message.Descriptor{
        Id:       5,
        Size:     0,
        Provider: message.ProvideSingleton(AlreadyLoggedIn),
    }

    unsupportedVersionDescriptor = message.Descriptor{
        Id:       6,
        Size:     0,
        Provider: message.ProvideSingleton(UnsupportedVersion),
    }

    fullDescriptor = message.Descriptor{
        Id:       7,
        Size:     0,
        Provider: message.ProvideSingleton(Full),
    }

    loginLimitExceededDescriptor = message.Descriptor{
        Id:       9,
        Size:     0,
        Provider: message.ProvideSingleton(LoginLimitExceeded),
    }

    serverBeingUpdatedDescriptor = message.Descriptor{
        Id:       14,
        Size:     0,
        Provider: message.ProvideSingleton(ServerBeingUpdated),
    }

    worldRunningClosedBetaDescriptor = message.Descriptor{
        Id:       19,
        Size:     0,
        Provider: message.ProvideSingleton(WorldRunningClosedBeta),
    }

    profileTransferDescriptor = message.Descriptor{
        Id:       21,
        Size:     1,
        Provider: newProfileTransfer,
    }

    malformedLoginPacketDescriptor = message.Descriptor{
        Id:       22,
        Size:     0,
        Provider: message.ProvideSingleton(MalformedLoginPacket),
    }

    errorLoadingProfileDescriptor = message.Descriptor{
        Id:       24,
        Size:     0,
        Provider: message.ProvideSingleton(ErrorLoadingProfile),
    }

    blockedComputerAddressDescriptor = message.Descriptor{
        Id:       26,
        Size:     0,
        Provider: message.ProvideSingleton(BlockedComputerAddress),
    }

    serviceUnavailableDescriptor = message.Descriptor{
        Id:       27,
        Size:     0,
        Provider: message.ProvideSingleton(ServiceUnavailable),
    }

    customRejectionDescriptor = message.Descriptor{
        Id:       29,
        Size:     message.SizeVariableShort,
        Provider: newCustomRejection,
    }

    enterSixDigitDescriptor = message.Descriptor{
        Id:       56,
        Size:     0,
        Provider: message.ProvideSingleton(EnterSixDigitPinCode),
    }

    invalidSixDigitPinCodeDescriptor = message.Descriptor{
        Id:       57,
        Size:     0,
        Provider: message.ProvideSingleton(InvalidSixDigitPinCode),
    }

    Okay                   = okay{}
    InvalidCredentials     = invalidCredentials{}
    AccountDisabled        = accountDisabled{}
    AlreadyLoggedIn        = alreadyLoggedIn{}
    UnsupportedVersion     = unsupportedVersion{}
    Full                   = full{}
    LoginLimitExceeded     = loginLimitExceeded{}
    ServerBeingUpdated     = serverBeingUpdated{}
    WorldRunningClosedBeta = worldRunningClosedBeta{}
    MalformedLoginPacket   = malformedLoginPacket{}
    ErrorLoadingProfile    = errorLoadingProfile{}
    BlockedComputerAddress = blockedComputerAddress{}
    ServiceUnavailable     = serviceUnavailable{}
    EnterSixDigitPinCode   = enterSixDigitPinCode{}
    InvalidSixDigitPinCode = invalidSixDigitPinCode{}
)

type okay struct{}

func (okay) Descriptor() message.Descriptor { return okayDescriptor }

func (okay) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (okay) Encode(b *buffer.ByteBuffer) error { return nil }

type unsupportedVersion struct{}

func (unsupportedVersion) Descriptor() message.Descriptor { return unsupportedVersionDescriptor }

func (unsupportedVersion) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (unsupportedVersion) Encode(b *buffer.ByteBuffer) error { return nil }

type full struct{}

func (full) Descriptor() message.Descriptor { return fullDescriptor }

func (full) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (full) Encode(b *buffer.ByteBuffer) error { return nil }

type invalidCredentials struct{}

func (invalidCredentials) Descriptor() message.Descriptor { return invalidCredentialsDescriptor }

func (invalidCredentials) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (invalidCredentials) Encode(b *buffer.ByteBuffer) error { return nil }

type accountDisabled struct{}

func (accountDisabled) Descriptor() message.Descriptor { return accountDisabledDescriptor }

func (accountDisabled) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (accountDisabled) Encode(b *buffer.ByteBuffer) error { return nil }

type alreadyLoggedIn struct{}

func (alreadyLoggedIn) Descriptor() message.Descriptor { return alreadyLoggedInDescriptor }

func (alreadyLoggedIn) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (alreadyLoggedIn) Encode(b *buffer.ByteBuffer) error { return nil }

type loginLimitExceeded struct{}

func (loginLimitExceeded) Descriptor() message.Descriptor { return loginLimitExceededDescriptor }

func (loginLimitExceeded) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (loginLimitExceeded) Encode(b *buffer.ByteBuffer) error { return nil }

type serverBeingUpdated struct{}

func (serverBeingUpdated) Descriptor() message.Descriptor { return serverBeingUpdatedDescriptor }

func (serverBeingUpdated) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (serverBeingUpdated) Encode(b *buffer.ByteBuffer) error { return nil }

type ProfileTransfer struct {
    Delay uint8
}

func newProfileTransfer() message.Message { return &ProfileTransfer{} }

func (ProfileTransfer) Descriptor() message.Descriptor { return profileTransferDescriptor }

func (ProfileTransfer) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (msg ProfileTransfer) Encode(b *buffer.ByteBuffer) error { return b.PutUint8(msg.Delay) }

type malformedLoginPacket struct{}

func (malformedLoginPacket) Descriptor() message.Descriptor { return malformedLoginPacketDescriptor }

func (malformedLoginPacket) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (malformedLoginPacket) Encode(b *buffer.ByteBuffer) error { return nil }

type worldRunningClosedBeta struct{}

func (worldRunningClosedBeta) Descriptor() message.Descriptor { return worldRunningClosedBetaDescriptor }

func (worldRunningClosedBeta) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (worldRunningClosedBeta) Encode(b *buffer.ByteBuffer) error { return nil }

type errorLoadingProfile struct{}

func (errorLoadingProfile) Descriptor() message.Descriptor { return errorLoadingProfileDescriptor }

func (errorLoadingProfile) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (errorLoadingProfile) Encode(b *buffer.ByteBuffer) error { return nil }

type blockedComputerAddress struct{}

func (blockedComputerAddress) Descriptor() message.Descriptor { return blockedComputerAddressDescriptor }

func (blockedComputerAddress) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (blockedComputerAddress) Encode(b *buffer.ByteBuffer) error { return nil }

type serviceUnavailable struct{}

func (serviceUnavailable) Descriptor() message.Descriptor { return serviceUnavailableDescriptor }

func (serviceUnavailable) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (serviceUnavailable) Encode(b *buffer.ByteBuffer) error { return nil }

type CustomRejection struct {
    TopLabel    string
    CenterLabel string
    BottomLabel string
}

func newCustomRejection() message.Message { return &CustomRejection{} }

func (CustomRejection) Descriptor() message.Descriptor { return customRejectionDescriptor }

func (CustomRejection) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (msg CustomRejection) Encode(b *buffer.ByteBuffer) error {
    if err := b.PutCString(msg.TopLabel); err != nil {
        return err
    }

    if err := b.PutCString(msg.CenterLabel); err != nil {
        return err
    }

    if err := b.PutCString(msg.BottomLabel); err != nil {
        return err
    }

    return nil
}

type enterSixDigitPinCode struct{}

func (enterSixDigitPinCode) Descriptor() message.Descriptor { return enterSixDigitDescriptor }

func (enterSixDigitPinCode) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (enterSixDigitPinCode) Encode(b *buffer.ByteBuffer) error { return nil }

type invalidSixDigitPinCode struct{}

func (invalidSixDigitPinCode) Descriptor() message.Descriptor { return invalidSixDigitPinCodeDescriptor }

func (invalidSixDigitPinCode) Decode(b *buffer.ByteBuffer, length int) error { return message.ErrDecodeNotSupported }

func (invalidSixDigitPinCode) Encode(b *buffer.ByteBuffer) error { return nil }