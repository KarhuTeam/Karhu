package ssh

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func CheckSsh(user, ip string, port int) error {

	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		log.Println("ressources/ssh: cannot find ssh in $PATH, can't connect to nodes.")
		return err
	}

	privateKey := PrivateKeyPath()

	command := fmt.Sprintf("%s -o BatchMode=yes -o ConnectTimeout=5 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i %s -p %d %s@%s exit", sshPath, privateKey, port, user, ip)

	log.Println("ressources/ssh: command:", command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
