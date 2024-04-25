package painter

import (
	"kpi-apz-lab-3/ui"
	"sync"
)

// Receiver отримує стан вікна, що був підготовлений в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(ui.Getter)
}

// Loop реалізує цикл подій для формування стану вікна через виконання операцій, отриманих із внутрішньої черги.
type Loop struct {
	Receiver Receiver

	st ui.State

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start() {
	l.st = ui.InitState()

	l.stop = make(chan struct{})

	go func() {
		for !l.stopReq || !l.mq.empty() {
			op := l.mq.pull()
			if update := op.Do(l.st); update {
				l.Receiver.Update(l.st)
			}
		}
		close(l.stop)
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(ui.Setter) {
		l.stopReq = true
	}))
	<-l.stop
}

type messageQueue struct {
	ops     []Operation
	mu      sync.Mutex
	blocked chan struct{}
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.ops = append(mq.ops, op)

	if mq.blocked != nil {
		close(mq.blocked)
		mq.blocked = nil
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for len(mq.ops) == 0 {
		mq.blocked = make(chan struct{})
		mq.mu.Unlock()
		<-mq.blocked
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
