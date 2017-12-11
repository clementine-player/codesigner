package codesigner

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	unsignedPath = "/some/path/to/a/clementine.dmg"
)

func fakeExecCommand(exitCode int) func(string, ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1", fmt.Sprintf("GO_WANT_EXIT_CODE=%d", exitCode)}
		return cmd
	}
}

func TestUnlockKeychain(t *testing.T) {
	Convey("Unlock keychain success", t, func() {
		execCommand = fakeExecCommand(0)
		defer func() { execCommand = exec.Command }()
		_, err := unlockKeychain("foo", "bar")
		So(err, ShouldBeNil)
	})
	Convey("Unlock keychain fails", t, func() {
		execCommand = fakeExecCommand(1)
		defer func() { execCommand = exec.Command }()
		_, err := unlockKeychain("foo", "bar")
		So(err, ShouldNotBeNil)
	})
}

func TestHelperProcess(t *testing.T) {
	Convey("Test helper process", t, func() {
		if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
			return
		}
		exitEnv, ok := os.LookupEnv("GO_WANT_EXIT_CODE")
		if !ok {
			defer os.Exit(0)
		} else {
			exitCode, err := strconv.Atoi(exitEnv)
			if err != nil {
				defer os.Exit(0)
			} else {
				defer os.Exit(exitCode)
			}
		}

		args := os.Args
		for len(args) > 0 {
			if args[0] == "--" {
				args = args[1:]
				break
			}
			args = args[1:]
		}

		cmd, args := args[0], args[1:]
		switch cmd {
		case "codesign":
			So(args[1], ShouldEqual, "-fv")
			So(args[2], ShouldEqual, "-s")
			So(args[3], ShouldEqual, unsignedPath)
		}
	})
}
