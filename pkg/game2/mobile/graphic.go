package mobile

var StopGraphic = Graphic{Id: 65535}

type Graphic struct {
    Id    uint16
    Y     uint16
    Delay uint16
}
