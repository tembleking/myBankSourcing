//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Event struct {
	StreamName    string
	StreamVersion int32
	EventName     string
	EventData     []byte
	HappenedOn    time.Time
	ContentType   string
}
