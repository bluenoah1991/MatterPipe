package main

import (
	"fmt"
	"os"

	"encoding/json"

	"io/ioutil"

	"path"

	"github.com/mattermost/platform/model"
)

// Configuration file format
type Configuration struct {
	Address  string `json:"address"`
	Clientid string `json:"clientid"`
	Password string `json:"password"`
	Team     string `json:"team"`
	Channel  string `json:"channel"`
}

func main() {
	var inputBytes []byte

	if os.Getenv("MATTERPIPE_DEBUG") == "true" {
		inputBytes = []byte("Debug data.")
	} else {
		fi, err := os.Stdin.Stat()
		if err != nil {
			fmt.Println("Stream Error.")
			os.Exit(1)
		}
		if (fi.Mode() & os.ModeNamedPipe) == 0 {
			fmt.Println("No access to a named pipe.")
			os.Exit(1)
		}
		// if fi.Size() == 0 {
		// 	fmt.Println("No input data.")
		// 	os.Exit(1)
		// }

		inputBytes, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Failed to read data.")
			os.Exit(1)
		}
	}

	configPath := "/etc/matterpipe.json"
	if os.Getenv("MATTERPIPE_OS") == "windows" {
		pwd, _ := os.Getwd()
		configPath = path.Join(pwd, "matterpipe.json")
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Configuration file not found.")
		os.Exit(1)
	}

	configDecoder := json.NewDecoder(configFile)
	config := Configuration{}
	err = configDecoder.Decode(&config)
	if err != nil {
		fmt.Println("Incorrect configuration file format.")
		os.Exit(1)
	}

	client := model.NewClient(config.Address)

	if loginResult, err := client.Login(config.Clientid, config.Password); err != nil {
		fmt.Printf("Logon failure: %s\n", err.Message)
		os.Exit(1)
	} else {
		botUser := loginResult.Data.(*model.User)
		if os.Getenv("MATTERPIPE_OS") != "linux" && os.Getenv("MATTERPIPE_DEBUG") == "true" {
			fmt.Printf("Logon to: %s\n", botUser.Username)
		}
	}

	var initialLoad *model.InitialLoad

	if initialLoadResults, err := client.GetInitialLoad(); err != nil {
		fmt.Printf("Initial failure: %s\n", err.Message)
		os.Exit(1)
	} else {
		initialLoad = initialLoadResults.Data.(*model.InitialLoad)
	}

	var currentTeam *model.Team

	for _, team := range initialLoad.Teams {
		if team.Name == config.Team {
			currentTeam = team
			break
		}
	}
	if currentTeam == nil {
		fmt.Println("Special team not found.")
		os.Exit(1)
	}

	client.SetTeamId(currentTeam.Id)

	var currentChannel *model.Channel

	if channelsResult, err := client.GetChannels(""); err != nil {
		fmt.Println("Failed to get channel data.")
		os.Exit(1)
	} else {
		channelList := channelsResult.Data.(*model.ChannelList)
		for _, channel := range *channelList {
			if channel.Name == config.Channel {
				currentChannel = channel
				break
			}
		}
	}

	if currentChannel == nil {
		fmt.Println("Special channel not found.")
		os.Exit(1)
	}

	post := &model.Post{}
	post.ChannelId = currentChannel.Id
	post.Message = string(inputBytes)
	post.RootId = ""

	if _, err := client.CreatePost(post); err != nil {
		fmt.Println("Failed to send.")
		os.Exit(1)
	}

}
