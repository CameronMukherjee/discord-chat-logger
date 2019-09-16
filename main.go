/*Didatic Discord Bot v2
Allow logging to be toggled on and off per user.

Collection of all userids - with option of logging as on or off.
If on - log
If off - ignore all message (no not log)
*/

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/color"
	"go.mongodb.org/mongo-driver/mongo/options"

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
	color.Yellow.Println("(/) :: Attempting to connect to Discord API...")
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatalln(err.Error())
		color.Red.Println("(-) :: Connection to Discord API could not be established!")
	}
	color.Green.Println("(+) :: Connection to Discord API Successful!")

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

	//For every message this code runs, move client to main.
	insert := InsertMessage{currentDate(), currentTime(), m.Content}
	mongoClient := mongoConnect()
	go insertToMongo(mongoClient, insert)
}

func mongoConnect() *mongo.Client {
	/*mongoConnect
	Establishing a connection with MongoDB to save logs.
	*/
	clientOptions := options.Client().ApplyURI(config.MongoURL)
	color.Yellow.Println("(/) :: Attempting to connect to MongoDB...")

	//TODO: What is context TODO?
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalln(err.Error())
		color.Red.Println("(-) :: Connection to MongoDB could not be established!")
	}
	color.Green.Println("(+) :: Connection to MongoDB Successful!")

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln(err.Error())
		color.Red.Println("(-) :: Connection to MongoDB could not be established!")
	}
	return client
}

func insertToMongo(client *mongo.Client, message InsertMessage) {
	/*insertToMongo
	Inserts message from 'messageHandler' to MongoDB
	Utilising the mongoConnect function.
	*/
	collection := client.Database("didaticDBv2").Collection("messages")
	_, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		color.Red.Println("(-) :: Failed to write message to database!")
	} else {
		color.Green.Println("(+) :: Successfully added new log!")
	}
	return
}

func printTable() {
	//TODO: Update whole table
	color.Green.Println("(+) :: Discord Token Verified!")
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

func currentDate() string {
	/*currentDate
	Returns date in format "00-00-0000" (British Format)
	This return is a string.
	*/
	dt := time.Now()
	return dt.Format("02-01-2006")
}

func currentTime() string {
	/*currentTime
	Returns time in format "00:00:00".
	This return is a string.
	*/
	dt := time.Now()
	return dt.Format("15:04:05")
}
