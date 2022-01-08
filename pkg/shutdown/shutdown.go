package shutdown

import (
	"io"
	"os"
	"os/signal"
	"project/pkg/logging"
)

func Graceful(signals []os.Signal, closeItems ...io.Closer) {
	logger := logging.GetLogger()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)
	sig := <-sigChan
	logger.Infof("Caught signal %s. Shutting down."+"..", sig)

	// Here we can do graceful shutdown (close connections and etc)
	for _, closer := range closeItems {
		if err := closer.Close(); err != nil {
			logger.Errorf("failed to close %v: %v", closer, err)
		}
	}
}
