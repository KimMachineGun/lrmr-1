package lrmr

import (
	"fmt"
	"github.com/airbloc/logger"
	"github.com/therne/lrmr/coordinator"
	"github.com/therne/lrmr/transformation"
	"github.com/therne/lrmr/worker"
	"os"
	"os/signal"
)

var (
	log = logger.New("lrmr")
)

func init() {
	transformation.Register(&LocalInput{})
}

func RunMaster(optionalOpt ...*Options) (*Master, error) {
	opt := DefaultOptions()
	if len(optionalOpt) > 0 {
		opt = optionalOpt[0]
	}

	etcd, err := coordinator.NewEtcd(opt.EtcdEndpoints)
	if err != nil {
		return nil, fmt.Errorf("connect etcd: %w", err)
	}
	return NewMaster(etcd, opt)
}

func RunWorker(optionalOpt ...*Options) error {
	opt := DefaultOptions()
	if len(optionalOpt) > 0 {
		opt = optionalOpt[0]
	}

	etcd, err := coordinator.NewEtcd(opt.EtcdEndpoints)
	if err != nil {
		return fmt.Errorf("connect etcd: %w", err)
	}
	w, err := worker.New(etcd, opt.Worker)
	if err != nil {
		return fmt.Errorf("init worker: %w", err)
	}
	go func() {
		if err := w.Start(); err != nil {
			log.Wtf("failed to start worker", err)
			return
		}
	}()

	waitForExit := make(chan os.Signal)
	signal.Notify(waitForExit, os.Interrupt, os.Kill)
	<-waitForExit

	if err := w.Stop(); err != nil {
		log.Error("failed to shutdown historical node", err)
	}
	log.Info("Bye")
	return nil
}

func Input(m *Master, path string) Session {
	sess := NewSession(m)
	sess.AddStage("__input", NewLocalInput(path))
	return sess
}