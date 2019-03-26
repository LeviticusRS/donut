package status

import (
    "github.com/sprinkle-it/donut/pkg/buffer"
    "github.com/sprinkle-it/donut/pkg/message"
)

var (
    okayConfig = message.Config{
        Id:   0,
        Size: 0,
        New:  message.Singleton(Okay),
    }

    invalidCredentialsConfig = message.Config{
        Id:   3,
        Size: 0,
        New:  message.Singleton(InvalidCredentials),
    }

    accountDisabledConfig = message.Config{
        Id:   4,
        Size: 0,
        New:  message.Singleton(AccountDisabled),
    }

    alreadyOnlineConfig = message.Config{
        Id:   5,
        Size: 0,
        New:  message.Singleton(AlreadyOnline),
    }

    unsupportedVersionConfig = message.Config{
        Id:   6,
        Size: 0,
        New:  message.Singleton(UnsupportedVersion),
    }

    fullConfig = message.Config{
        Id:   7,
        Size: 0,
        New:  message.Singleton(Full),
    }

    loginLimitExceededConfig = message.Config{
        Id:   9,
        Size: 0,
        New:  message.Singleton(LoginLimitExceeded),
    }

    serverUpdateConfig = message.Config{
        Id:   14,
        Size: 0,
        New:  message.Singleton(ServerUpdate),
    }

    closedBetaConfig = message.Config{
        Id:   19,
        Size: 0,
        New:  message.Singleton(ClosedBeta),
    }

    profileTransferConfig = message.Config{
        Id:   21,
        Size: 1,
        New:  func() message.Message { return &ProfileTransfer{} },
    }

    malformedLoginPacketConfig = message.Config{
        Id:   22,
        Size: 0,
        New:  message.Singleton(MalformedLoginPacket),
    }

    errorLoadingProfileConfig = message.Config{
        Id:   24,
        Size: 0,
        New:  message.Singleton(ErrorLoadingProfile),
    }

    blockedAddressConfig = message.Config{
        Id:   26,
        Size: 0,
        New:  message.Singleton(BlockedAddress),
    }

    serviceUnavailableConfig = message.Config{
        Id:   27,
        Size: 0,
        New:  message.Singleton(ServiceUnavailable),
    }

    customRejectionConfig = message.Config{
        Id:   29,
        Size: message.SizeVariableShort,
        New:  func() message.Message { return &CustomRejection{} },
    }

    enterPinConfig = message.Config{
        Id:   56,
        Size: 0,
        New:  message.Singleton(EnterPin),
    }

    invalidPinConfig = message.Config{
        Id:   57,
        Size: 0,
        New:  message.Singleton(InvalidPin),
    }

    Okay                 = okay{}
    InvalidCredentials   = invalidCredentials{}
    AccountDisabled      = accountDisabled{}
    AlreadyOnline        = alreadyOnline{}
    UnsupportedVersion   = unsupportedVersion{}
    Full                 = full{}
    LoginLimitExceeded   = loginLimitExceeded{}
    ServerUpdate         = serverUpdate{}
    ClosedBeta           = closedBeta{}
    MalformedLoginPacket = malformedLoginPacket{}
    ErrorLoadingProfile  = errorLoadingProfile{}
    BlockedAddress       = blockedAddress{}
    ServiceUnavailable   = serviceUnavailable{}
    EnterPin             = enterPin{}
    InvalidPin           = invalidPin{}
)

type okay struct{}

func (okay) Config() message.Config { return okayConfig }

func (okay) Encode(b *buffer.ByteBuffer) error { return nil }

type unsupportedVersion struct{}

func (unsupportedVersion) Config() message.Config { return unsupportedVersionConfig }

func (unsupportedVersion) Encode(b *buffer.ByteBuffer) error { return nil }

type full struct{}

func (full) Config() message.Config { return fullConfig }

func (full) Encode(b *buffer.ByteBuffer) error { return nil }

type invalidCredentials struct{}

func (invalidCredentials) Config() message.Config { return invalidCredentialsConfig }

func (invalidCredentials) Encode(b *buffer.ByteBuffer) error { return nil }

type accountDisabled struct{}

func (accountDisabled) Config() message.Config { return accountDisabledConfig }

func (accountDisabled) Encode(b *buffer.ByteBuffer) error { return nil }

type alreadyOnline struct{}

func (alreadyOnline) Config() message.Config { return alreadyOnlineConfig }

func (alreadyOnline) Encode(b *buffer.ByteBuffer) error { return nil }

type loginLimitExceeded struct{}

func (loginLimitExceeded) Config() message.Config { return loginLimitExceededConfig }

func (loginLimitExceeded) Encode(b *buffer.ByteBuffer) error { return nil }

type serverUpdate struct{}

func (serverUpdate) Config() message.Config { return serverUpdateConfig }

func (serverUpdate) Encode(b *buffer.ByteBuffer) error { return nil }

type ProfileTransfer struct {
    Delay uint8
}

func (ProfileTransfer) Config() message.Config { return profileTransferConfig }

func (msg ProfileTransfer) Encode(b *buffer.ByteBuffer) error { return b.PutUint8(msg.Delay) }

type malformedLoginPacket struct{}

func (malformedLoginPacket) Config() message.Config { return malformedLoginPacketConfig }

func (malformedLoginPacket) Encode(b *buffer.ByteBuffer) error { return nil }

type closedBeta struct{}

func (closedBeta) Config() message.Config { return closedBetaConfig }

func (closedBeta) Encode(b *buffer.ByteBuffer) error { return nil }

type errorLoadingProfile struct{}

func (errorLoadingProfile) Config() message.Config { return errorLoadingProfileConfig }

func (errorLoadingProfile) Encode(b *buffer.ByteBuffer) error { return nil }

type blockedAddress struct{}

func (blockedAddress) Config() message.Config { return blockedAddressConfig }

func (blockedAddress) Encode(b *buffer.ByteBuffer) error { return nil }

type serviceUnavailable struct{}

func (serviceUnavailable) Config() message.Config { return serviceUnavailableConfig }

func (serviceUnavailable) Encode(b *buffer.ByteBuffer) error { return nil }

type CustomRejection struct {
    TopLabel    string
    CenterLabel string
    BottomLabel string
}

func (CustomRejection) Config() message.Config { return customRejectionConfig }

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

type enterPin struct{}

func (enterPin) Config() message.Config { return enterPinConfig }

func (enterPin) Encode(b *buffer.ByteBuffer) error { return nil }

type invalidPin struct{}

func (invalidPin) Config() message.Config { return invalidPinConfig }

func (invalidPin) Encode(b *buffer.ByteBuffer) error { return nil }