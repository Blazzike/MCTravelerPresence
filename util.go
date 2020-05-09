package main

import (
	"encoding/json"
	"github.com/hugolgst/rich-go/client"
	"github.com/mitchellh/go-ps"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type LauncherProfiles struct {
	AuthenticationDatabase map[string]struct {
		Profiles map[string]struct {
			DisplayName string `json:"displayName"`
		} `json:"profiles"`
	} `json:"authenticationDatabase"`
	SelectedUser struct {
		Account string `json:"account"`
	}
}

func getMinecraftLauncherPid() int {
	processes, err := ps.Processes()
	if err != nil {
		panic(err)
	}

	for _, process := range processes {
		if process.Executable() == "MinecraftLauncher.exe" {
			return process.Pid()
		}
	}

	return -1
}

func getMinecraftUuid() string {
	bytes, err := ioutil.ReadFile(minecraftDirectoryPath + "launcher_profiles.json")
	if err != nil {
		return ""
	}

	launcherProfiles := LauncherProfiles{}
	err = json.Unmarshal(bytes, &launcherProfiles)
	if err != nil {
		return ""
	}

	for id, authenticationDatabaseItem := range launcherProfiles.AuthenticationDatabase {
		if id == launcherProfiles.SelectedUser.Account {
			if len(authenticationDatabaseItem.Profiles) != 0 {
				for uuid := range authenticationDatabaseItem.Profiles {
					return uuid
				}
			}

			break
		}
	}

	return ""
}

func waitUntilMinecraftRunning() int {
	for {
		launcherPid := getMinecraftLauncherPid()
		if launcherPid != -1 {
			return launcherPid
		}

		time.Sleep(time.Second * 10)
	}
}

func waitUntilMinecraftStopsRunning(pid int) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return
	}

	_, err = process.Wait()
	if err != nil {
		panic(err)
	}

	//for {
	//	killErr := syscall.Kill(pid, syscall.Signal(0))
	//	if killErr != nil {
	//		break
	//	}
	//
	//	time.Sleep(time.Second * 10)
	//}
}

type PresenceState struct {
	State     string
	Details   string
	Timestamp *time.Time
}

func removePresence() {
	client.Logout()
}

func updatePresence(state PresenceState) error {
	err := client.Login("708032641261371424")
	if err != nil {
		return err
	}

	err = client.SetActivity(client.Activity{
		State:      state.State,
		Details:    state.Details,
		LargeImage: "logo",
		LargeText:  "play.MCTraveler.eu",
		Timestamps: &client.Timestamps{
			Start: state.Timestamp,
		},
	})

	return err
}

func CopyFile(src, destination string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}
