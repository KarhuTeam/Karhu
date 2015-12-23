package ssh

import (
	"fmt"
	"github.com/gotoolz/env"
	"log"
	"os"
	"os/exec"
	"path"
)

const (
	SSH_KEYS_DIR = "ssh"
	SSH_KEY_NAME = "karhu"
)

// Check for ssh keys and generate them if needed
func init() {

	GenerateSSHKeyPair()
}

func GenerateSSHKeyPair() error {

	targetDir := path.Clean(fmt.Sprintf("%s/%s", env.GetDefault("DATA_DIR", "./data"), SSH_KEYS_DIR))
	targetFile := path.Clean(fmt.Sprintf("%s/%s", targetDir, SSH_KEY_NAME))

	// Check if key already exist
	if _, err := os.Stat(targetFile); !os.IsNotExist(err) {
		return nil
	}

	log.Println("ressources/ssh: generating new ssh key-pair...")

	sshKeyGenPath, err := exec.LookPath("ssh-keygen")
	if err != nil {
		log.Println("ressources/ssh: cannot find ssh-keygen in $PATH, can't generate ssh key-pair.")
		return err
	}

	// Check for target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Println("ressources/ssh: failed to generate ssh directory:", err)
		return err
	}

	command := fmt.Sprintf(`%s -q -t rsa -b 4096 -N "" -C karhu@karhu-master -f %s`, sshKeyGenPath, targetFile)

	log.Println("ressources/ssh: exec:", command)
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println("ressources/ssh: failed to generate ssh key-pair:", err)
		return err
	}

	return nil
}
