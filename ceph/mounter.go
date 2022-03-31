package ceph

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MountInput struct {
	Client     string
	SourcePath string
	TargetPath string
	UidLocal   int64
	UidRemote  int64
	GidLocal   int64
	GidRemote  int64
}

func Mount(input MountInput) error {
	uidMap, e := writeIDMapping(input.UidLocal, input.UidRemote)
	if e != nil {
		return e
	}
	gidMap, e := writeIDMapping(input.GidLocal, input.GidRemote)
	if e != nil {
		return e
	}

	mountCmd := "ceph-fuse"
	mountArgs := []string{}

	// TODO maybe make conf/keyring configurable through some sort of conf file
	clientName := fmt.Sprintf("client.%s", input.Client)
	configPath := fmt.Sprintf("/clients/conf/%s.conf", input.Client)
	keyringPath := fmt.Sprintf("/clients/keyring/%s.keyring", input.Client)
	mountArgs = append(
		mountArgs,
		"-n", clientName,
		"-c", configPath,
		"-k", keyringPath,
		"-o", "allow_other",
		"--uid-mapping", uidMap,
		"--gid-mapping", gidMap,
		"-r", input.SourcePath,
		input.TargetPath,
	)

	// create target, os.Mkdirall is noop if it exists
	err := os.MkdirAll(input.TargetPath, 0750)
	if err != nil {
		return err
	}

	glog.Infof("executing mount command cmd=%s, args=%s", mountCmd, mountArgs)

	out, err := exec.Command(mountCmd, mountArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("mounting failed: %v cmd: '%s %s' output: %q",
			err, mountCmd, strings.Join(mountArgs, " "), string(out))
	}

	return nil
}

func writeIDMapping(local int64, remote int64) (string, error) {
	f, e := ioutil.TempFile("", "pk-*")
	defer f.Close()
	if e != nil {
		return "", status.Errorf(codes.Internal, "can not create tmp file for ID mapping: %s", e)
	}

	idMap := fmt.Sprintf("%d:%d", local, remote)
	_, e = f.WriteString(idMap)
	if e != nil {
		return "", status.Errorf(codes.Internal, "can not create tmp file for ID mapping: %s", e)
	}
	e = f.Chmod(0600)
	if e != nil {
		return "", status.Errorf(codes.Internal, "can not change rights for ID mapping: %s", e)
	}
	return f.Name(), nil
}
