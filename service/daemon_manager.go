package service

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// systemDRecord - standard record (struct) for linux systemD version of daemon package
type SystemDRecord struct {
	Name string
	Path string
}

const (
	success = "[  OK  ]"
	failed  = "[  FAILED  ]"
)

const (
	statNotInstalled = "Service not installed"
)

var (
	// ErrUnsupportedSystem appears if try to use service on system which is not supported by this release
	ErrUnsupportedSystem = errors.New("Unsupported system")

	// ErrRootPrivileges appears if run installation or deleting the service without root privileges
	ErrRootPrivileges = errors.New("You must have root user privileges. Possibly using 'sudo' command should help")

	// ErrNotInstalled appears if try to delete service which was not been installed
	ErrNotInstalled = errors.New("Service is not installed")

	// ErrAlreadyRunning appears if try to start already running service
	ErrAlreadyRunning = errors.New("Service is already running")

	// ErrAlreadyStopped appears if try to stop already stopped service
	ErrAlreadyStopped = errors.New("Service has already been stopped")
)

func checkPrivileges() (bool, error) {

	if output, err := exec.Command("id", "-g").Output(); err == nil {
		if gid, parseErr := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 32); parseErr == nil {
			if gid == 0 {
				return true, nil
			}
			return false, ErrRootPrivileges
		}
	}
	return false, ErrUnsupportedSystem
}

// service path for systemD daemons
func (linux *SystemDRecord) servicePath() string {
	return linux.Path + linux.Name + ".service"
}

// Is a service installed
func (linux *SystemDRecord) isInstalled() bool {
	if _, err := os.Stat(linux.servicePath()); err == nil {
		return true
	}
	return false
}

func (linux *SystemDRecord) checkRunning() (string, bool) {
	output, err := exec.Command("systemctl", "status", linux.Name+".service").Output()
	if err == nil {
		return string(output), true
	}
	return "Service " + linux.Name + " is stopped", false
}

// Start the service
func (linux *SystemDRecord) Start() (string, error) {
	startAction := "Starting " + linux.Name + ":"

	if ok, err := checkPrivileges(); !ok {
		return startAction + failed, err
	}

	if !linux.isInstalled() {
		return startAction + failed, ErrNotInstalled
	}

	if _, ok := linux.checkRunning(); ok {
		return startAction + failed, ErrAlreadyRunning
	}

	if err := exec.Command("systemctl", "start", linux.Name+".service").Run(); err != nil {
		return startAction + failed, err
	}

	return startAction + success, nil
}

// Status - Get service status
func (linux *SystemDRecord) Status() (string, error) {

	if ok, err := checkPrivileges(); !ok {
		return "", err
	}

	if !linux.isInstalled() {
		return statNotInstalled, ErrNotInstalled
	}

	statusAction, _ := linux.checkRunning()

	return statusAction, nil
}

// Stop the service
func (linux *SystemDRecord) Stop() (string, error) {
	stopAction := "Stopping " + linux.Name + ":"

	if ok, err := checkPrivileges(); !ok {
		return stopAction + failed, err
	}

	if !linux.isInstalled() {
		return stopAction + failed, ErrNotInstalled
	}

	if _, ok := linux.checkRunning(); !ok {
		return stopAction + failed, ErrAlreadyStopped
	}

	if err := exec.Command("systemctl", "stop", linux.Name+".service").Run(); err != nil {
		return stopAction + failed, err
	}

	return stopAction + success, nil
}

// Restart the service
func (linux *SystemDRecord) Restart() (string, error) {
	restartAction := "Restart " + linux.Name + ":"

	if ok, err := checkPrivileges(); !ok {
		return restartAction + failed, err
	}

	if !linux.isInstalled() {
		return restartAction + failed, ErrNotInstalled
	}

	if _, ok := linux.checkRunning(); !ok {
		return restartAction + failed, ErrAlreadyStopped
	}

	if err := exec.Command("systemctl", "restart", linux.Name+".service").Run(); err != nil {
		return restartAction + failed, err
	}

	return restartAction + success, nil
}
