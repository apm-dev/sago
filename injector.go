package sago

import (
	"context"
	"time"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/pkg/errors"
)

type Injector interface {
	InjectEventToFlow(uniqueId string, s Saga, e Event) error
}

type injector struct {
	zb zbc.Client
}

func NewInjector(c zbc.Client) Injector {
	return &injector{c}
}

type Event struct {
	Name string
	Vars map[string]interface{}
	TTL  time.Duration
}

func (i *injector) InjectEventToFlow(uniqueId string, s Saga, e Event) error {
	const op string = "sago.injector.InjectEventToFlow"
	// build zeebe message of event
	req := i.zb.NewPublishMessageCommand().
		MessageName(e.Name).
		CorrelationKey(buildSagoKey(s.SagaType(), uniqueId)).
		TimeToLive(e.TTL)
	// set variables of message if there was any
	if e.Vars != nil && len(e.Vars) > 0 {
		removeReservedVariableKeys(e.Vars)
		var err error
		req, err = req.VariablesFromMap(e.Vars)
		if err != nil {
			return errors.Wrapf(err, op)
		}
	}
	// send message to zeebe
	_, err := req.Send(context.Background())
	if err != nil {
		return errors.Wrapf(err, op)
	}
	return nil
}
