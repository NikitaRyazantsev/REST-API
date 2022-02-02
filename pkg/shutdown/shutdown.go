// Package shutdown - package for Graceful Shutdown
package shutdown

import (
	"io"
	"os"
	"os/signal"
	"project/pkg/logging"
)

// Graceful - func for shutdown
func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	logger := logging.GetLogger()

	// work wih os signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	// create signals channel
	sig := <-sigChan
	logger.Infof("Caught signal %s. Shutting down."+"..", sig)

	// here we can do graceful shutdown (close connections and etc)
	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logger.Errorf("failed to close %v: %v", closer, err)
		}
	}
}
