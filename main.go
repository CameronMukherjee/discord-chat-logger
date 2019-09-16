/*Didatic Discord Bot v2
Allow logging to be toggled on and off per user.

Collection of all userids - with option of logging as on or off.
If on - log
If off - ignore all message (no not log)
*/

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"

	"./config"
)

//BotID - type: string
var BotID string

//InsertMessage - The insert struct for MongoDB
type InsertMessage struct {
	Date    string
	Time    string
	Message string
}

func main() {
	color.Yellow.Println("Attempting to connect to Discord API...")
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatalln(err.Error())
		color.Red.Println("(-) :: Could not establish a connection to the Discord API!")
	}

	//TODO: Figure out what discord.User does and change the error
	user, err := discord.User("@me")
	if err != nil {
		log.Fatalln(err.Error())
		color.Red.Println("(-) :: Could not set BotID as UserID")
	}

	BotID = user.ID

	discord.AddHandler(messageHandler)
	err = discord.Open()
	if err != nil {
		log.Fatalln(err.Error())
	}
	printTable()
	<-make(chan struct{})
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*messageHandler
	This code is ran for every message sent.
	*/
}

func mongoConnect() {
	/*mongoConnect
	Establishing a connection with MongoDB to save logs.
	*/
}

func printTable() {
	color.Green.Println("Discord Token Verified!")
	time.Sleep(time.Second)
	color.Blue.Println("Logging has started.")
	time.Sleep(time.Second)
	color.Cyan.Printf("-       STATUS        ")
	fmt.Print("-")
	color.Yellow.Printf("       DATETIME      ")
	fmt.Print("-")
	color.Blue.Printf(" USERNAME ")
	fmt.Print("-")
	color.Magenta.Printf(" MESSAGE\n")
}
