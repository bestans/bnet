package protoc

import (
	"bnet/com/session/msgfactory"
	"reflect"
)

func registerMessage(msgType int32, msgReflectType reflect.Type)  {
	msgfactory.RegisterPBMessage(msgType, msgReflectType)
}
