package gw

import (
	"log"
	"time"
)

type Heartbeat struct {
	IsDead bool
	MsToDie time.Duration

	lastBeatTime time.Time
	isChecking bool
}

func NewHeartbeat(msToDie time.Duration) *Heartbeat {
	return &Heartbeat{
		MsToDie: msToDie,
	}
}

func (hb *Heartbeat) StartHealthCheck() {
	hb.isChecking = true
	hb.lastBeatTime = time.Now()
	hb.heathCheck()
}

func (hb *Heartbeat) heathCheck() {
	time.AfterFunc(time.Second, func() {
		if hb.checkIsDead() {
			hb.IsDead = true
			hb.isChecking = false
			log.Println("Die", hb)
		} else {
			hb.heathCheck()
			log.Println("checking..")
		}
	})
}

func (hb *Heartbeat) Beat() {
	hb.lastBeatTime = time.Now()
}

func (hb *Heartbeat) checkIsDead() bool {
	return makeTimestamp(time.Now()) - makeTimestamp(hb.lastBeatTime) > hb.MsToDie
}

func makeTimestamp(t time.Time) time.Duration {
	return time.Duration(t.UnixNano() / int64(time.Millisecond)) * time.Millisecond
}
