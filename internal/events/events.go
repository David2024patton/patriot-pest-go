// Package events is the process-wide live feed hub behind the SSE HUD
// (/api/staff/events). Any module can Publish; the legacy module streams
// events to connected staff dashboards. Kept dependency-free so inbox,
// legacy, twilio and facebook can all emit without import cycles.
package events

import (
	"sync"
	"time"
)

// Event is one HUD notification.
type Event struct {
	Type string `json:"type"`
	Text string `json:"text"`
	At   string `json:"at"`
}

var (
	mu   sync.Mutex
	subs = map[chan Event]struct{}{}
	log  []Event // ring of last 50 for late joiners
)

// Publish fans an event out to every subscriber and keeps it in the ring.
func Publish(typ, text string) {
	ev := Event{Type: typ, Text: text, At: time.Now().Format("15:04:05")}
	mu.Lock()
	log = append(log, ev)
	if len(log) > 50 {
		log = log[len(log)-50:]
	}
	for ch := range subs {
		select {
		case ch <- ev:
		default: // slow client — drop rather than block the publisher
		}
	}
	mu.Unlock()
}

// Subscribe returns a channel that receives live events plus the backlog.
func Subscribe() chan Event {
	ch := make(chan Event, 16)
	mu.Lock()
	subs[ch] = struct{}{}
	backlog := append([]Event{}, log...)
	mu.Unlock()
	go func() {
		for _, ev := range backlog {
			select {
			case ch <- ev:
			default:
			}
		}
	}()
	return ch
}

// Unsubscribe detaches a subscriber channel.
func Unsubscribe(ch chan Event) {
	mu.Lock()
	delete(subs, ch)
	mu.Unlock()
}
