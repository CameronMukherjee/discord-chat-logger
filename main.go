/*Didatic Discord Bot v2
Allow logging to be toggled on and off per user.

Collection of all userids - with option of logging as on or off.
If on - log
If off - ignore all message (no not log)
*/

package main

import (
	"context"
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
	Date string
	Time string
	User string
	Mess string
}

var client = mongoConnect()

func main() {
	color.Yellow.Println("(/) :: Attempting to connect to Discord API...")
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatalln(err.Error())
		color.Red.Println("(-) :: Connection to Discord API could not be established!")
	}
	color.Green.Println("(+) :: Connection to Discord API Successful!")
	color.Magenta.Printf("---------------------------------------------- \n")

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
	<-make(chan struct{})
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*messageHandler
	This code is ran for every message sent.
	*/

	if m.Content == "" {
		return
	}
	username := m.Author.String()
	insert := InsertMessage{currentDate(), currentTime(), username, m.Content}
	go insertToMongo(insert)
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

func insertToMongo(message InsertMessage) {
	/*insertToMongo
	Inserts message from 'messageHandler' to MongoDB
	Utilising the mongoConnect function.
	*/
	collection := client.Database("didaticDBv2").Collection("messages")
	_, err := collection.InsertOne(context.TODO(), message)
	if err != nil {
		color.Red.Println("(-) :: Failed to write message to database!")
	} else {
		color.Green.Println("(+) :: Successfully added new entry!")
	}
	return
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

// func encryptMessage(message string) []byte {
// 	preEncryptedMessage := []byte(message)
// 	c, err := aes.NewCipher(config.AESKey)
// 	if err != nil {
// 		color.Red.Println("(-) :: Could not generate new AES Cypher!")
// 	}

// 	gcm, err := cipher.NewGCM(c)
// 	if err != nil {
// 		color.Red.Println("(-) :: Could not generate new Galois Counter Mode operation!")
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
// 		color.Red.Println("(-) :: Could not secure memory!")
// 	}

// 	// fmt.Println(gcm.Seal(nonce, nonce, preEncryptedMessage, nil))
// 	postEncryptedMessage := gcm.Seal(nonce, nonce, preEncryptedMessage, nil)
// 	return postEncryptedMessage
// }
