package runner

import (
	"os/exec"
	"strconv"
	"strings"

	log "github.com/golang/glog"
)

const SANDBOX_USER = "sandbox"
const SANDBOX_USER_UID = 65537

var SANDBOX_GROUP_ID = 0

func init() {
	// create sandbox user
	user := SANDBOX_USER
	uid := SANDBOX_USER_UID

	// check if user exists
	_, err := exec.Command("id", user).Output()
	if err != nil {
		// create user
		output, err := exec.Command("bash", "-c", "useradd -u "+strconv.Itoa(uid)+" "+user).Output()
		if err != nil {
			log.Errorf("failed to create user: %v, %v", err, string(output))
		}
	}

	// get gid of sandbox user and setgid
	gid, err := exec.Command("id", "-g", SANDBOX_USER).Output()
	if err != nil {
		log.Errorf("failed to get gid of user: %v", err)
	}

	SANDBOX_GROUP_ID, err = strconv.Atoi(strings.TrimSpace(string(gid)))
	if err != nil {
		log.Errorf("failed to convert gid: %v", err)
	}
}
