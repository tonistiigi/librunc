package librunc

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/opencontainers/specs"
)

func TestSimpleRun(t *testing.T) {
	spec := getBlankSpec()
	tmpdir, err := ioutil.TempDir("", "librunc")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	spec.Process.Terminal = false
	spec.Process.Args = []string{"echo", "ok!"}

	curdir, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	spec.Root.Path = filepath.Join(curdir, "fixtures/busybox")

	c, err := New(tmpdir, spec)
	if err != nil {
		t.Fatal(err)
	}

	r, w := io.Pipe()
	err = c.Start(nil, w, nil)
	if err != nil {
		t.Fatal(err)
	}

	out := make([]byte, 3)
	io.ReadFull(r, out)
	actual := string(out)
	expected := "ok!"
	if expected != actual {
		t.Fatalf("Wrong output. Expected %v, got %v.", expected, actual)
	}

}

func getBlankSpec() specs.LinuxSpec {
	return specs.LinuxSpec{
		Spec: specs.Spec{
			Version: specs.Version,
			Platform: specs.Platform{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
			Root: specs.Root{
				Path:     "rootfs",
				Readonly: true,
			},
			Process: specs.Process{
				Terminal: true,
				User:     specs.User{},
				Args: []string{
					"sh",
				},
				Env: []string{
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"TERM=xterm",
				},
			},
			Hostname: "shell",
			Mounts: []specs.Mount{
				{
					Type:        "proc",
					Source:      "proc",
					Destination: "/proc",
					Options:     "",
				},
				{
					Type:        "tmpfs",
					Source:      "tmpfs",
					Destination: "/dev",
					Options:     "nosuid,strictatime,mode=755,size=65536k",
				},
				{
					Type:        "devpts",
					Source:      "devpts",
					Destination: "/dev/pts",
					Options:     "nosuid,noexec,newinstance,ptmxmode=0666,mode=0620,gid=5",
				},
				{
					Type:        "tmpfs",
					Source:      "shm",
					Destination: "/dev/shm",
					Options:     "nosuid,noexec,nodev,mode=1777,size=65536k",
				},
				{
					Type:        "mqueue",
					Source:      "mqueue",
					Destination: "/dev/mqueue",
					Options:     "nosuid,noexec,nodev",
				},
				{
					Type:        "sysfs",
					Source:      "sysfs",
					Destination: "/sys",
					Options:     "nosuid,noexec,nodev",
				},
				{
					Type:        "cgroup",
					Source:      "cgroup",
					Destination: "/sys/fs/cgroup",
					Options:     "nosuid,noexec,nodev,relatime,ro",
				},
			},
		},
		Linux: specs.Linux{
			Namespaces: []specs.Namespace{
				{
					Type: "pid",
				},
				{
					Type: "network",
				},
				{
					Type: "ipc",
				},
				{
					Type: "uts",
				},
				{
					Type: "mount",
				},
			},
			Capabilities: []string{
				"AUDIT_WRITE",
				"KILL",
				"NET_BIND_SERVICE",
			},
			Devices: []specs.Device{
				{
					Type:        'c',
					Path:        "/dev/null",
					Major:       1,
					Minor:       3,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/random",
					Major:       1,
					Minor:       8,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/full",
					Major:       1,
					Minor:       7,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/tty",
					Major:       5,
					Minor:       0,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/zero",
					Major:       1,
					Minor:       5,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/urandom",
					Major:       1,
					Minor:       9,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
			},
			Resources: specs.Resources{
				Memory: specs.Memory{
					Swappiness: -1,
				},
			},
		},
	}
}
