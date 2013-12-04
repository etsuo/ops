package ops

import (
	"errors"
	"fmt"
	"github.com/etsuo/log"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
)

type MountedSmbVolume struct {
	userName      string
	password      string
	serverAddress string
	remotePath    string
	localPath     string
	mounted       bool
}

// Mount mounts a remote SMB share locally. It is cross platform
// compatible with Windows and Darwin
func (m *MountedSmbVolume) Mount(userName string, password string, server string, remotePath string, localPath string) error {
	m.userName = userName
	m.password = password
	m.serverAddress = server
	m.remotePath = remotePath
	m.localPath = localPath

	var err error = nil

	switch runtime.GOOS {
	case "darwin":
		err = m.mount_Darwin()
	case "windows":
		err = m.mount_Windows()
	case "linux":
		fallthrough
	default:
		err = errors.New(fmt.Sprintf("OS %s not supported.", runtime.GOOS))
	}

	return err
}

func (m *MountedSmbVolume) mount_Darwin() error {
	m.createLocalTmpDir()

	files, _ := ioutil.ReadDir(m.localPath)
	if len(files) > 0 {
		log.Send.Warningf("Path '%s' was already mounted - it's likely from a failed prior run. Dismounting prior to continuing.", m.localPath)
		m.dismount_Darwin()
	}

	var cmd string
	if m.password == "" {
		cmd = fmt.Sprintf("mount -t smbfs //%s@%s/%s %s", m.userName, m.serverAddress, m.remotePath, m.localPath)

	} else {
		cmd = fmt.Sprintf("mount -t smbfs //%s:%s@%s/%s %s", m.userName, url.QueryEscape(m.password), m.serverAddress, m.remotePath, m.localPath)
	}

	if err := RunCommand(cmd); err != nil {
		return err
	}

	m.mounted = true
	return nil
}

func (m *MountedSmbVolume) mount_Windows() error {
	return errors.New("Not implemented")
}

// Dismount removes the previously created mount. It is cross platform
// compatible with Windows and Darwin
func (m *MountedSmbVolume) Dismount() error {
	var err error = nil

	if !m.IsMounted() {
		log.Send.Infof("Dismount requested for path '%s', which is not mounted. No action took place.", m.localPath)
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		err = m.dismount_Darwin()
	case "windows":
		err = m.mount_Windows()
	case "linux":
		fallthrough
	default:
		err = errors.New(fmt.Sprintf("OS %s not supported.", runtime.GOOS))
	}

	m.removeLocalTmpDir()
	return err
}

func (m *MountedSmbVolume) dismount_Darwin() error {
	cmd := fmt.Sprintf("diskutil unmount %s", m.localPath)
	return RunCommand(cmd)
}

func (m *MountedSmbVolume) dismount_Windows() error {
	return errors.New("Not implemented")
}

func (m *MountedSmbVolume) createLocalTmpDir() error {
	log.Send.Debugf("Creating temp mount directory %s", m.localPath)

	err := os.MkdirAll(m.localPath, 0777)
	if err != nil {
		log.Send.Fatalf("Unable to create required directory. Error: %s", err.Error())
	}
	return err
}

func (m *MountedSmbVolume) removeLocalTmpDir() error {
	log.Send.Debugf("Clearning up temp mount directory %s", m.localPath)

	files, _ := ioutil.ReadDir(m.localPath)
	expectedFileCount := 0

	if len(files) > expectedFileCount {
		log.Send.Warningf("Expected directory '%s' to have %d files in it prior to cleaning it up post dismount.",
			" Intead, %d files were found. An attempt to delete the directory was aborted.", m.localPath, expectedFileCount, len(files))
	}

	err := os.RemoveAll(m.localPath)
	if err != nil {
		log.Send.Warningf("Unable to remove directory '%s'. Error: %s", m.localPath, err.Error())
	} else {
		log.Send.Infof("Removed path: %s", m.localPath)
	}

	return err
}

func (m *MountedSmbVolume) IsMounted() bool {
	return m.mounted
}
