package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

type Config struct {
	Login string
	Oauth string
}

func start(client *twitch.Client) {
	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func main() {
	// read in config.json file
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error when opening config file: ", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Println("Error during Unmarshal(): ", err)
	}

	client := twitch.NewClient(config.Login, config.Oauth)

	client.OnConnect(func() {
		fmt.Println("Connected to Twitch IRC")
	})

	client.OnSelfJoinMessage(func(message twitch.UserJoinMessage) {
		fmt.Println("Joined " + message.Channel)
	})

	client.OnSelfPartMessage(func(message twitch.UserPartMessage) {
		fmt.Println("Parted " + message.Channel)
	})

	// Used for testing purposes
	// client.OnPrivateMessage(func(message twitch.PrivateMessage) {
	// 	fmt.Println(message.Channel, message.Message)
	// })

	// Since Connect() is blocking we start new go routine
	go start(client)

	// Open channels.txt for parsing
	channelFile, err := os.Open("channels.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer channelFile.Close()

	scanner := bufio.NewScanner(channelFile)
	channels := []string{}
	for i := 1; scanner.Scan(); i++ {
		channels = append(channels, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Open banlist.txt for parsing
	banFile, err := os.Open("banlist.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer banFile.Close()

	scanner = bufio.NewScanner(banFile)
	banList := []string{}
	for i := 1; scanner.Scan(); i++ {
		banList = append(banList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	// Connect to channels and ban users
	for _, c := range channels {
		client.Join(c)
		for i, u := range banList {
			if i == 0 {
				fmt.Println("Banning users...")
			}
			client.Ban(c, u, "banned via TwitchCrossMod")
			time.Sleep(500 * time.Millisecond)
		}

		client.Depart(c)
		time.Sleep(5 * time.Second)
	}

	time.Sleep(5 * time.Second)
	fmt.Println("TwitchCrossMod finished")

}
