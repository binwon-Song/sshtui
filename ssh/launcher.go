package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"sshtui/config"
	"sshtui/crypto"
)

func Connect(server config.Server) error {
	addr := fmt.Sprintf("%s@%s", server.User, server.Host)
	port := fmt.Sprintf("-p%d", server.Port)

	// 비밀번호가 저장되어 있는 경우
	if server.EncryptedPass != "" {
		password, err := crypto.Decrypt(server.EncryptedPass)
		if err != nil {
			return fmt.Errorf("비밀번호 복호화 실패: %v", err)
		}
		// sshpass를 사용하여 자동으로 비밀번호 입력
		cmd := exec.Command("sshpass", "-p", password, "ssh", port, addr)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// 비밀번호가 없는 경우 기존 방식대로 연결
	cmd := exec.Command("ssh", port, addr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
