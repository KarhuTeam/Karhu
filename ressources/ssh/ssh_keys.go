package ssh

import (
	"fmt"
	"github.com/gotoolz/env"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	SSH_KEYS_DIR             = "ssh"
	SSH_KEY_NAME             = "karhu"
	SSH_PUBLIC_KEY_NAME      = SSH_KEY_NAME + ".pub"
	SSH_AUTHORIZED_KEYS_FILE = "authorized_keys"
	SSH_AUTHORIZED_KEYS_DIR  = "~/.ssh"
)

// Check for ssh keys and generate them if needed
func init() {

	generateSSHKeyPair()
}

func keyDir() string {
	return path.Clean(fmt.Sprintf("%s/%s", env.GetDefault("DATA_DIR", "./data"), SSH_KEYS_DIR))
}

func privateKeyPath() string {
	return path.Clean(fmt.Sprintf("%s/%s", keyDir(), SSH_KEY_NAME))
}

func publicKeyPath() string {
	return path.Clean(fmt.Sprintf("%s/%s", keyDir(), SSH_PUBLIC_KEY_NAME))
}

func AuthorizedKeysPath() string {
	return path.Join(SSH_AUTHORIZED_KEYS_DIR, SSH_AUTHORIZED_KEYS_FILE)
}

func generateSSHKeyPair() error {

	targetDir := keyDir()
	targetFile := privateKeyPath()

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

func GetPublicKey() ([]byte, error) {

	targetFile := publicKeyPath()

	data, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return nil, err
	}

	return []byte(strings.Trim(string(data), " \t\n")), err
}
