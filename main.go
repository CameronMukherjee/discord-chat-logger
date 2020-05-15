package main

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	Date     string
	Time     string
	User     string
	Chatroom string
	Message  string
}

func main() {

	fmt.Println("Attempting to connect to Discord API.")

	discord, err := discordgo.New("Bot " + getEnvVar("DiscordToken"))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to Discord API established.")
	}

	dbSetup()

	_, err = discord.User("@me")
	if err != nil {
		log.Fatal(err)
	}

	discord.AddHandler(messageHandler)
	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Content == "" {
		return
	}

	username := m.Author.String()
	chatroom := "https://discordapp.com/channels/" + m.GuildID + "/" + m.ChannelID

	insert := Message{currentDate(), currentTime(), username, chatroom, m.Content}
	go dbInsert(insert)

}

func dbInsert(message Message) {

	db, err := sql.Open("sqlite3", getEnvVar("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	query, err := db.Prepare("INSERT INTO Messages(username, chatroom, message, date, time) " +
		"VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = query.Exec(message.User, message.Chatroom, message.Message, message.Date, message.Time)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Message Added: ", message)
	}
}

//noinspection GoNilness
func dbSetup() {

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}

	setup, err := db.Prepare("CREATE TABLE IF NOT EXISTS Messages(id INTEGER PRIMARY KEY AUTOINCREMENT," +
		" username varchar(255), chatroom varchar(255), message varchar(255), date varchar(255), time varchar(255))")

	_, err = setup.Exec()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Database Initialization Complete.")
	}

	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}

}

func currentDate() string {

	dt := time.Now()
	return dt.Format("02-01-2006")

}

func currentTime() string {

	dt := time.Now()
	return dt.Format("15:04:05")

}

func getEnvVar(key string) string {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)

}
