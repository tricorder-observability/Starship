package exec

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func Run(argv []string) (string, string, error) {
	cmd := NewCommand(argv)
	err := cmd.Start()
	msg := fmt.Sprintf("command=%v stdout=%s stderr=%s error: %v", argv, cmd.Stdout(), cmd.Stderr(), err)
	log.Infof(msg)
	if err != nil {
		return "", "", fmt.Errorf("start failed, message=%s", msg)
	}
	err = cmd.Wait()
	msg = fmt.Sprintf("command=%v stdout=%s stderr=%s error: %v", argv, cmd.Stdout(), cmd.Stderr(), err)
	if err != nil {
		return "", "", fmt.Errorf("wait failed, message=%s", msg)
	}
	return cmd.Stdout(), cmd.Stderr(), nil
}
