package utils

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func GetProcessID(port string) (int, error) {
	cmd := exec.Command("lsof", "-t", "-i:"+port)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	pidString := strings.TrimSpace(string(output))
	pid, err := strconv.Atoi(pidString)
	if err != nil {
		return 0, err
	}

	return pid, nil
}

func KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Kill()
	if err != nil {
		return err
	}

	return nil
}

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if strings.Contains(item, v) {
			return true
		}
	}
	return false
}

func IsPortInUse(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func GeneratePatterns(cmdLine string) []string {
	words := strings.Fields(cmdLine)

	if len(words) > 1 {
		goFile := words[2]

		fileName := filepath.Base(goFile)
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

		return []string{
			cmdLine,
			"/tmp/go-build.*exe/" + fileName,
		}
	}

	return []string{}
}
