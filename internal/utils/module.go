package utils

import (
	"os"
	"os/exec"
	"strings"
)

func SliceContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func GetRunningContainerNames() ([]string, error) {
	c := exec.Command("docker", "ps", "--format", "{{.Names}}")
	c.Env = os.Environ()
	stdout, err := c.Output()
	if err != nil {
		return []string{}, err
	}
	return strings.Split(strings.Trim(string(stdout), " \n\t\r"), "\n"), nil
}
