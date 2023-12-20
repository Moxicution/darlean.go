package transporthandler

import (
	"core"
	"core/invoke"
	"core/wire"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type pendingCall struct {
	finished chan<- *invoke.InvokeResponse
}

type TransportHandler struct {
	appId        string
	transport    core.Transport
	pendingCalls map[string]pendingCall
	mutex        sync.Mutex
	dispatcher   InwardCallDispatcher
}

type InwardCallDispatcher interface {
	Dispatch(tags *wire.Tags) error
}

func (invoker *TransportHandler) Listen() {
	for tags := range invoker.transport.GetInputChannel() {
		switch tags.Remotecall_Kind {
		case "call":
			if invoker.dispatcher == nil {
				fmt.Println("transporthandler: Ignore incoming message: no dispatcher assigned")
				continue
			}

			fmt.Println("transporthandler: Ignore incoming message: 'call' is not yet implemented")
			continue
		case "return":
			invoker.handleReturnMessage(tags)
		}
	}
}

func (invoker *TransportHandler) handleReturnMessage(tags *wire.Tags) {
	invoker.mutex.Lock()
	call, found := invoker.pendingCalls[tags.Remotecall_Id]
	if found {
		delete(invoker.pendingCalls, tags.Remotecall_Id)
	}
	invoker.mutex.Unlock()

	if !found {
		fmt.Println("Received value without matching call")
		return
	}

	call.finished <- &invoke.InvokeResponse{
		Value: tags.Value,
		Error: tags.Error,
	}
}

func New(transport core.Transport, dispatcher InwardCallDispatcher, appId string) *TransportHandler {
	invoker := TransportHandler{
		appId:        appId,
		transport:    transport,
		dispatcher:   dispatcher,
		pendingCalls: make(map[string]pendingCall),
	}

	go invoker.Listen()

	return &invoker
}

func (invoker *TransportHandler) Invoke(req *invoke.TransportHandlerInvokeRequest) *invoke.InvokeResponse {
	id := uuid.NewString()

	tags := wire.Tags{}
	tags.Transport_Receiver = req.Receiver
	tags.Transport_Return = invoker.appId
	tags.Remotecall_Id = id
	tags.Remotecall_Kind = "call"
	tags.ActorType = req.ActorType
	tags.ActorId = req.ActorId
	tags.ActionName = req.ActionName
	tags.Arguments = req.Parameters

	response := make(chan *invoke.InvokeResponse)

	invoker.mutex.Lock()
	invoker.pendingCalls[id] = pendingCall{
		finished: response,
	}
	invoker.mutex.Unlock()

	err := invoker.transport.Send(tags)
	if err != nil {
		panic(err)
	}

	return <-response
}