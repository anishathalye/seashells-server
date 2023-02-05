package datamanager

import (
	"fmt"
	"github.com/anishathalye/seashells-server/slidingbuffer"
	mrand "math/rand"
	"sync"
	"time"
)

type Session struct {
	user      string // ip address
	id        string
	data      *slidingbuffer.SlidingBuffer
	finalized bool
	finished  time.Time
	mu        *sync.Mutex
	cond      *sync.Cond
}

type DataManager struct {
	// mapping from identifier -> buffer
	byId map[string]*Session
	// mapping from ip string -> identifiers
	byUser       map[string][]*Session
	limit        int
	perUserLimit int
	gcTime       time.Duration
	mu           sync.Mutex
	dead         bool
}

// returns a random number in [low, high)
func rand(low, high int64) int64 {
	return low + mrand.Int63n(high-low)
}

func (sess *Session) String() string {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	var finished string
	if sess.finalized {
		finished = fmt.Sprintf("%v", sess.finished)
	} else {
		finished = ""
	}
	return fmt.Sprintf(
		"Session{ip=%s, id=%s, len=%d, finalized=%t, finished=%s}",
		sess.user,
		sess.id,
		sess.data.Len(),
		sess.finalized,
		finished,
	)
}

// returns false if the session is dead
func (sess *Session) Append(data []byte) bool {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	if sess.finalized {
		return false
	}
	sess.data.Append(data)
	sess.cond.Broadcast()
	return true
}

func (sess *Session) Finalize() {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	// just in case Finalize() is called multiple times
	if !sess.finalized {
		sess.finished = time.Now()
	}
	sess.finalized = true
	sess.cond.Broadcast()
}

func (sess *Session) isFinalized() bool {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	return sess.finalized
}

// returns nil if it doesn't exist
// channel is closed once we're sure there will be no more data
func (sess *Session) Subscribe() (chan []byte, func()) {
	done := make(chan bool)
	dataChan := make(chan []byte)
	go func() {
		index := 0
		sess.mu.Lock()
		for {
			// this select block is so we can exit quickly when the
			// client is done; otherwise, we might have to wait
			// till we get new data or the session is finalized
			sess.mu.Unlock()
			select {
			case <-done:
				close(dataChan)
				return
			default:
				sess.mu.Lock()
			}
			if sess.data.Len() > index {
				// more data to send
				_, newData, nextIndex := sess.data.Get(index)
				index = nextIndex
				sess.mu.Unlock() // unlock before potentially blocking operation
				select {
				case <-done:
					// we've already released the lock
					close(dataChan)
					return
				case dataChan <- newData:
					// grab the lock again
					sess.mu.Lock()
				}
			} else if sess.finalized {
				// we're done, we've sent all the data already
				sess.mu.Unlock()
				close(dataChan)
				<-done // wait for receiver to acknowledge
				return
			} else {
				sess.cond.Wait()
			}
		}
	}()
	doneFunc := func() {
		go func() {
			sess.mu.Lock()
			sess.cond.Broadcast()
			sess.mu.Unlock()
			done <- true
		}()
	}
	return dataChan, doneFunc
}

func newSession(user, id string, limit int) *Session {
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	sess := &Session{
		user:      user,
		id:        id,
		data:      slidingbuffer.New(limit),
		finalized: false,
		mu:        mu,
		cond:      cond,
	}
	return sess
}

func (mgr *DataManager) Create(user string, id string) *Session {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if sess := mgr.byId[id]; sess != nil {
		return nil // can't create same session twice
	}

	// see if we need to delete user's old sessions
	userSessions := mgr.byUser[user]
	if len(userSessions) > mgr.perUserLimit-1 {
		toDelete := userSessions[:len(userSessions)-mgr.perUserLimit+1]
		for _, sess := range toDelete {
			mgr.destroy(sess.id)
		}
	}

	sess := newSession(user, id, mgr.limit)

	mgr.byId[id] = sess
	mgr.byUser[user] = append(mgr.byUser[user], sess)

	return sess
}

// need to be holding lock before calling this
func (mgr *DataManager) destroy(id string) {
	sess := mgr.byId[id]
	if sess == nil {
		return
	}
	sess.Finalize()
	delete(mgr.byId, sess.id)
	var newUserSessions []*Session
	// sess.user and sess.id are set at creation time, so it's okay to
	// access even when we're not holding the lock
	for _, otherSess := range mgr.byUser[sess.user] {
		if otherSess.id != sess.id {
			newUserSessions = append(newUserSessions, otherSess)
		}
	}
	if len(newUserSessions) > 0 {
		mgr.byUser[sess.user] = newUserSessions
	} else {
		delete(mgr.byUser, sess.user)
	}
}

func (sess *Session) Dump() []byte {
	sess.mu.Lock()
	defer sess.mu.Unlock()
	_, data, _ := sess.data.Get(0)
	return data
}

func (mgr *DataManager) Get(id string) *Session {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	return mgr.byId[id]
}

func (mgr *DataManager) All() []*Session {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	var sessions []*Session
	for _, sess := range mgr.byId {
		sessions = append(sessions, sess)
	}
	return sessions
}

func (mgr *DataManager) isDead() bool {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	return mgr.dead
}

func (mgr *DataManager) gc() {
	for !mgr.isDead() {
		mgr.mu.Lock()
		var toDelete []string = nil // ids
		for id, sess := range mgr.byId {
			// it's dead if it's greater than the gc time
			// and if the client has disconnected
			if sess.isFinalized() {
				// need to grab lock to read .finished
				sess.mu.Lock()
				if time.Since(sess.finished) > mgr.gcTime {
					toDelete = append(toDelete, id)
				}
				sess.mu.Unlock()
			}
		}
		for _, id := range toDelete {
			mgr.destroy(id)
		}
		mgr.mu.Unlock()
		// GC check interval is between 25% and 50% of gcTime
		time.Sleep(time.Duration(rand(int64(mgr.gcTime/4), int64(mgr.gcTime/2))))
	}
}

func (mgr *DataManager) Kill() {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.dead = true
}

func New(limit int, perUserLimit int, gcTime time.Duration) *DataManager {
	manager := &DataManager{
		byId:         map[string]*Session{},
		byUser:       map[string][]*Session{},
		limit:        limit,
		perUserLimit: perUserLimit,
		gcTime:       gcTime,
		dead:         false,
	}
	go manager.gc()
	return manager
}
