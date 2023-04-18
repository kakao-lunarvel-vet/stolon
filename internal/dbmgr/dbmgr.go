package dbmgr

import (
	"github.com/mitchellh/copystructure"
	"github.com/sorintlab/stolon/internal/common"
	"reflect"
)

type SystemData struct {
	SystemID   string
	TimelineID uint64
	XLogPos    uint64
}

type TimelineHistory struct {
	TimelineID  uint64
	SwitchPoint uint64
	Reason      string
}

type InitConfig struct {
	Locale        string
	Encoding      string
	DataChecksums bool
}

type Manager interface {
	GetTimelinesHistory(timeline uint64) ([]*TimelineHistory, error)
}

type RecoveryMode int

const (
	RecoveryModeNone RecoveryMode = iota
	RecoveryModeStandby
	RecoveryModeRecovery
)

type RecoveryOptions struct {
	RecoveryMode       RecoveryMode
	RecoveryParameters common.Parameters
}

func (r *RecoveryOptions) DeepCopy() *RecoveryOptions {
	nr, err := copystructure.Copy(r)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(r, nr) {
		panic("not equal")
	}
	return nr.(*RecoveryOptions)
}

func NewRecoveryOptions() *RecoveryOptions {
	return &RecoveryOptions{RecoveryParameters: make(common.Parameters)}
}
