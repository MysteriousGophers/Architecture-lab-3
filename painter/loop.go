package painter

import (
	"golang.org/x/exp/shiny/screen"
	"image"
	"sync"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver

	next screen.Texture
	prev screen.Texture

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(800, 800)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.stop = make(chan struct{})
	go func() {
		for !l.stopReq || !l.mq.empty() {
			op := l.mq.pull()
			if update := op.Do(l.next); update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
		close(l.stop)
	}()
}

func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(t screen.Texture) {
		l.stopReq = true
	}))
	<-l.stop
}

type messageQueue struct {
	ops        []Operation
	mu         sync.Mutex
	pushSignal chan struct{}
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.ops = append(mq.ops, op)
	if mq.pushSignal != nil {
		close(mq.pushSignal)
		mq.pushSignal = nil
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	for len(mq.ops) == 0 {
		mq.pushSignal = make(chan struct{})
		mq.mu.Unlock()
		<-mq.pushSignal
		mq.mu.Lock()
	}
	op := mq.ops[0]
	mq.ops[0] = nil
	mq.ops = mq.ops[1:]
	return op
}

func (mq *messageQueue) empty() bool {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	return len(mq.ops) == 0
}
