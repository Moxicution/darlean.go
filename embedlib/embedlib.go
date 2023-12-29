package main

// See: https://github.com/enobufs/go-calls-c-pointer

// typedef void (*invoke_cb)(_GoString_);
// void makeCallback(_GoString_ bufhandle, invoke_cb cb);
import "C"
import (
	"strings"

	"github.com/darlean-io/darlean.go/base/invoker"
)

var apiInstance *Api

//export Start
func Start(appId string, natsAddr string, hosts string) {
	hostsArray := strings.Split(hosts, ",")
	apiInstance = NewApi(appId, natsAddr, hostsArray)
}

//export Stop
func Stop() {
	apiInstance.Stop()
}

//export Invoke
func Invoke(cb C.invoke_cb, actorType string, actorId []string, actionName string, arguments string) {
	//goCb := func(bufhandle int) {
	//	C.makeCallback(C.int(bufhandle), cb)
	//}
	goCb := func(bufhandle string) {
		C.makeCallback(bufhandle, cb)
	}
	request := invoker.Request{
		ActorType:  actorType,
		ActorId:    actorId, //strings.Split(actorId, ","),
		ActionName: actionName,
		Parameters: []any{arguments},
	}
	apiInstance.Invoke(&request, goCb)
}
