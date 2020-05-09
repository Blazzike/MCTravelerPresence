package main

import (
	"github.com/hugolgst/rich-go/client"
	"github.com/sqweek/dialog"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var minecraftDirectoryPath string

type MinecraftState struct {
	Away        bool
	PlayerCount int
	Region      string
	Timestamp   time.Time
}

func main() {
	appDataPath := os.Getenv("APPDATA")
	minecraftDirectoryPath = appDataPath + "\\.minecraft\\"
	startupDirectoryPath := appDataPath + "\\Microsoft\\Windows\\Start Menu\\Programs\\Startup"

	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}

	if filepath.Dir(executable) != startupDirectoryPath { // Not running from startup. Install!
		executableDestination := startupDirectoryPath + "\\MCTravelerPresence.exe"
		if err := CopyFile(executable, executableDestination); err != nil {
			panic(err)
		}

		if err := os.Chdir(startupDirectoryPath); err != nil {
			panic(err)
		}

		cmd := exec.Command("cmd.exe", "/C", "start", "/b", "MCTravelerPresence.exe")
		if err := cmd.Run(); err != nil {
			panic(err)
		}

		dialog.Message("MCTraveler's Discord Rich Presence has been configured.").Title("Installation successful!").Info()

		return
	}

	minecraftState := MinecraftState{
		Away:        false,
		PlayerCount: 0,
		Region:      "",
		Timestamp:   time.Time{},
	}

	for {
		launcherPid := waitUntilMinecraftRunning()

		uuid := getMinecraftUuid()
		if uuid == "" {
			dialog.Message("Your UUID could not be detected.").Title("Uh-oh!").Error()

			os.Exit(1)
		}

		log.Println("Connecting to MCTraveler...")
		webSocketConnect(uuid, func(closeSocket func()) {
			log.Println("Connected!")

			go func() {
				waitUntilMinecraftStopsRunning(launcherPid)
				log.Println("Minecraft stopped running. Closing socket...")
				closeSocket()
			}()
		}, func(json map[string]interface{}) {
			payloadType := json["type"]
			switch payloadType {
			case "online":
				minecraftState.Timestamp = time.Now()

				break
			case "playerCount":
				minecraftState.PlayerCount = int(json["playerCount"].(float64))

				break
			case "region":
				region := json["region"]
				if region == nil {
					minecraftState.Region = ""

					break
				}

				minecraftState.Region = region.(string)

				break
			case "offline":
				minecraftState.Timestamp = time.Time{}

				break
			case "away":
				minecraftState.Away = json["away"].(bool)

				break
			default:
				log.Printf("Invalid type '%s' sent.", payloadType)

				return
			}

			if minecraftState.Timestamp.IsZero() {
				removePresence()

				return
			}

			playState := "Playing"
			if minecraftState.Away {
				playState = "Away"
			}

			state := "In the wilderness"
			if minecraftState.Region != "" {
				state = "In " + minecraftState.Region
			}

			otherPlayerCount := minecraftState.PlayerCount - 1
			details := playState + " with " + strconv.Itoa(otherPlayerCount) + " other player"
			if otherPlayerCount != 1 {
				details += "s"
			}

			err := updatePresence(PresenceState{
				State:     state,
				Details:   details,
				Timestamp: &minecraftState.Timestamp,
			})

			if err != nil {
				log.Println("Discord presence could not be set. Is Discord running?")
			}
		})

		log.Println("Connection closed!")

		client.Logout() // Safe to call even if logged out, internal checks
	}
}
