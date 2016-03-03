package ssh

import (
	"fmt"
	"log"
	// "os"
	"os/exec"
)

func Exec(user, ip string, port int, command string) ([]byte, error) {
	sshPath, err := exec.LookPath("ssh")
	if err != nil {
		log.Println("ressources/ssh: cannot find ssh in $PATH, can't connect to nodes.")
		return nil, err
	}

	privateKey := PrivateKeyPath()

	command = fmt.Sprintf("%s -o BatchMode=yes -o ConnectTimeout=5 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i %s -p %d %s@%s -- %s", sshPath, privateKey, port, user, ip, command)

	log.Println("ressources/ssh: command:", command)

	cmd := exec.Command("sh", "-c", command)

	out, err := cmd.CombinedOutput()

	return out, err
}

func CheckSsh(user, ip string, port int) error {

	_, err := Exec(user, ip, port, "exit")
	return err

}
