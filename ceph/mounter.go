package ceph

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/zv0n/ceph-proxy/configuration"
	"k8s.io/utils/mount"

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

func Umount(targetPath string) error {
	notMnt, err := mount.New("").IsLikelyNotMountPoint(targetPath)

	if err != nil {
		if os.IsNotExist(err) {
			return status.Error(codes.NotFound, "Targetpath not found")
		} else {
			return status.Error(codes.Internal, err.Error())
		}
	}
	if notMnt {
		return status.Error(codes.NotFound, "Volume not mounted")
	}

	err = mount.New("").Unmount(targetPath)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func Mount(input MountInput, config *configuration.Configuration) (err error, uidMap string, gidMap string) {
	uidMap, err = writeIDMapping(input.UidLocal, input.UidRemote)
	if err != nil {
		return err, "", ""
	}
	gidMap, err = writeIDMapping(input.GidLocal, input.GidRemote)
	if err != nil {
		return err, "", ""
	}

	mountCmd := "ceph-fuse"
	mountArgs := []string{}

	// TODO maybe make conf/keyring configurable through some sort of conf file
	clientName := fmt.Sprintf("client.%s", input.Client)
	configPath := fmt.Sprintf("%s/%s.conf", config.ClientConfPath, input.Client)
	keyringPath := fmt.Sprintf("%s/%s.keyring", config.ClientKeyringPath, input.Client)
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
	err = os.MkdirAll(input.TargetPath, 0750)
	if err != nil {
		return err, "", ""
	}

	log.Printf("executing mount command cmd=%s, args=%s", mountCmd, mountArgs)

	out, err := exec.Command(mountCmd, mountArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("mounting failed: %v cmd: '%s %s' output: %q",
			err, mountCmd, strings.Join(mountArgs, " "), string(out)), "", ""
	}

	return nil, uidMap, gidMap
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
