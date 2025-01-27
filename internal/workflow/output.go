package workflow

import (
	"encoding/json"
	"fmt"
)

type outputType interface {
	SSEParams
}

type Output struct {
	Name        string
	Destination Destination
	Condition   *Condition
	params      any
}

func UnmarshallOutput(name, destination string, condition *Condition, params any) (*Output, error) {
	if name == "" {
		return nil, fmt.Errorf("output name is empty")
	}

	dest, err := getDestination(destination)
	if err != nil {
		return nil, err
	}

	o := Output{
		Name:        name,
		Destination: dest,
		Condition:   condition,
		params:      params,
	}

	err = o.unmarshallParams()
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (o *Output) unmarshallParams() error {
	marshalled, err := json.Marshal(o.params)
	if err != nil {
		return fmt.Errorf("error marshalling output params: %w", err)
	}

	var params parameter
	switch o.Destination {
	case SSE:
		var sse SSEParams
		err = unmarshalParams(marshalled, &sse)
		if err != nil {
			return err
		}

		params = sse
	default:
		return fmt.Errorf("unknown output type: %s", o.Destination)
	}

	o.params = params
	return params.validate()
}

func (o *Output) SSEParams() (SSEParams, error) {
	params, ok := o.params.(SSEParams)
	if !ok {
		return SSEParams{}, fmt.Errorf("output params are not of type SSEParams")
	}

	return params, nil
}

func (o *Output) MQTTParams() (MQTTParams, error) {
	if o.Destination != MQTTPub {
		return MQTTParams{}, fmt.Errorf("output source is not MQTT")
	}

	params, ok := o.params.(MQTTParams)
	if !ok {
		return MQTTParams{}, fmt.Errorf("input params are not of type MQTTParams")
	}

	return params, nil
}
