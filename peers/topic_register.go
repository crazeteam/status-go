package peers

import (
	"sync"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/p2p/discv5"

	"github.com/status-im/status-go/common"
	"github.com/status-im/status-go/discovery"
	"github.com/status-im/status-go/logutils"
)

// Register manages register topic queries
type Register struct {
	discovery discovery.Discovery
	topics    []discv5.Topic

	wg   sync.WaitGroup
	quit chan struct{}
}

// NewRegister creates instance of topic register
func NewRegister(discovery discovery.Discovery, topics ...discv5.Topic) *Register {
	return &Register{discovery: discovery, topics: topics}
}

// Start topic register query for every topic
func (r *Register) Start() error {
	if !r.discovery.Running() {
		return ErrDiscv5NotRunning
	}
	r.quit = make(chan struct{})
	for _, topic := range r.topics {
		r.wg.Add(1)
		go func(t discv5.Topic) {
			defer common.LogOnPanic()
			logutils.ZapLogger().Debug("v5 register topic", zap.String("topic", string(t)))
			if err := r.discovery.Register(string(t), r.quit); err != nil {
				logutils.ZapLogger().Error("error registering topic", zap.String("topic", string(t)), zap.Error(err))
			}
			r.wg.Done()
		}(topic)
	}
	return nil
}

// Stop all register topic queries and waits for them to exit
func (r *Register) Stop() {
	if r.quit == nil {
		return
	}
	select {
	case <-r.quit:
		return
	default:
		close(r.quit)
	}
	logutils.ZapLogger().Debug("waiting for register queries to exit")
	r.wg.Wait()
}
