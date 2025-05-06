package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"sshtui/config"
)

func Connect(server config.Server) error {
	addr := fmt.Sprintf("%s@%s", server.User, server.Host)
	port := fmt.Sprintf("-p%d", server.Port)

	cmd := exec.Command("ssh", port, addr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

