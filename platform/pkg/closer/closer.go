package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/clava1096/rocket-service/platform/pkg/logger"
	"go.uber.org/zap"
)

const shutdownTimeout = 5 * time.Second

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}
type Closer struct {
	mu     sync.Mutex
	once   sync.Once
	done   chan struct{}
	funcs  []func(ctx context.Context) error
	logger Logger
}

var globalCloser = NewWithLogger(&logger.NoopLogger{})

func AddNamed(name string, f func(ctx context.Context) error) {
	globalCloser.AddNamed(name, f)
}

func Add(f ...func(ctx context.Context) error) {
	globalCloser.Add(f...)
}

func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

func SetLogger(logger Logger) {
	globalCloser.SetLogger(logger)
}

func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

func New(signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Logger(), signals...)
}

func NewWithLogger(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}
	return c
}

func (c *Closer) AddNamed(name string, f func(ctx context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("Close %s...", name))

		err := f(ctx)
		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("Error while closing %s: %v (tooked %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("Close %s took %s sec", name, duration))
		}
		return err
	})
}

func (c *Closer) SetLogger(logger Logger) {
	c.logger = logger
}

func (c *Closer) Add(f ...func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info(context.Background(), "signal received, shutting down...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "failed to shutdown gracefully, error while close resource: %v", zap.Error(err))
		}
	case <-c.done:

	}
}

func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) > 0 {
			c.logger.Info(ctx, "No functions to close.")
		}

		c.logger.Info(ctx, "Begin shutdown gracefully...")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		for i := 0; i < len(funcs); i++ {
			f := funcs[i]
			wg.Add(1)
			go func(f func(ctx context.Context) error) {
				defer wg.Done()

				go func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "Context cancelled while shutdown gracefully")

				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, "All resources have been closed.")
					return
				}

				c.logger.Error(ctx, "Error while closing resource: ", zap.Error(err))

				if err == nil {
					result = err
				}
			}
		}
	})

	return result
}
