/*
Copyright 2017 daswf852@outlook.com

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	Core "AtomBot/Core"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	Token  string
	DBFile string
	Parser Core.Parser
	Logger Core.Logger
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&DBFile, "d", "", "'Database' file")
	flag.Parse()
	Parser = Core.MakeParser()

	//register commands
	echoCmd := Core.Command{
		"echo",
		1,
		"A command to *echo* what you say", "echo [message]",
		true,
		echo}
	Parser.Register(&echoCmd)

	pingCmd := Core.Command{
		"Ping!",
		0,
		"Ping the bot! Or maybe a website in future...", "Ping!",
		true,
		ping}
	Parser.Register(&pingCmd)

	succCmd := Core.Command{
		"succ",
		0,
		"succ someone", "succ (opt.)[user mention]",
		true,
		succ}
	Parser.Register(&succCmd)

	fuccCmd := Core.Command{
		"fucc",
		0,
		"fucc someone", "fucc (opt.)[user mention]",
		true,
		fucc}
	Parser.Register(&fuccCmd)

	whoamiCmd := Core.Command{
		"whoami",
		0,
		"A command to get your user info", "whoami",
		true,
		whoami}
	Parser.Register(&whoamiCmd)

	seenCmd := Core.Command{
		"seen",
		1,
		"A command to see the last seen info of a user mentioned", "seen [user mention]",
		true,
		seen}
	Parser.Register(&seenCmd)
	//register commands end

	Parser.SetPrefix(",")
	Logger = Core.MakeLogger()
	Logger.ReadFromFile("db.json")
}

func main() {
	bot, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session:\n\t", err)
		return
	}

	bot.AddHandler(onMessage)
	bot.AddHandler(onStatusUpdate)

	err = bot.Open()
	if err != nil {
		fmt.Println("Error opening connection:\n\t", err)
		return
	}

	fmt.Println("Bot is up and running!")
	bot.UpdateStatus(0, "github.com/Daswf852/AtomBot")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Close()
	Logger.OutToFile("db.json")
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		Logger.UpdateEntryMsg(m.Author.ID, m)
		return
	}
	s.ChannelMessageSend(m.ChannelID, Parser.Execute(s, m))
	if strings.Contains(m.Content, "🅱") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🅱")
	}

	Logger.UpdateEntryMsg(m.Author.ID, m)
}

func onStatusUpdate(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	Logger.UpdateEntryPresence(p.Presence.User.ID, p)
}

//commands

func echo(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	var retString string
	for i := 1; len(args.Args) > i; i++ {
		retString = fmt.Sprintln(retString, args.Args[i])
	}
	fmt.Println(retString)
	return retString
}

func ping(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	return "Pong!"
}

func succ(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	if args.Count >= 1 {
		return fmt.Sprintf("***%s succs %s***", m.Author.Mention(), args.Args[1])
	} else {
		return fmt.Sprintf("***%s succs %s***", s.State.User.Mention(), m.Author.Mention())
	}
}

func fucc(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	if args.Count >= 1 {
		return fmt.Sprintf("***%s fuccs %s***", m.Author.Mention(), args.Args[1])
	} else {
		return fmt.Sprintf("***%s fuccs %s***", s.State.User.Mention(), m.Author.Mention())
	}
}

func whoami(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	author := discordgo.MessageEmbedAuthor{
		Name:    fmt.Sprintln("User info of: ", m.Author.String()),
		IconURL: m.Author.AvatarURL("")}

	retEmbed := discordgo.MessageEmbed{
		Author: &author,
		Color:  0x43c605}

	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Is bot?", Value: fmt.Sprintln(m.Author.Bot), Inline: false})
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Username", Value: m.Author.String(), Inline: false})
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Mention", Value: m.Author.Mention(), Inline: false})
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "ID", Value: m.Author.ID, Inline: false})
	s.ChannelMessageSendEmbed(m.ChannelID, &retEmbed)
	return fmt.Sprintln("Here you go ", m.Author.Mention(), "!")
}

func seen(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	if len(m.Mentions) > 0 {
		if Logger.EntryExists(m.Mentions[0].ID) {
			user, _ := Logger.GetInfo(m.Mentions[0].ID)
			author := discordgo.MessageEmbedAuthor{
				Name:    fmt.Sprintln("Last seen info of ", m.Mentions[0].String()),
				IconURL: m.Mentions[0].AvatarURL("")}
			retEmbed := discordgo.MessageEmbed{
				Author: &author,
				Color:  0x05c699}
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Last seen on", Value: fmt.Sprintln(user.LastSeen)})

			LastChanName := "Invalid Channel"
			LastGuildName := "Invalid Guild"

			LastChan, cerr := s.Channel(user.LastChannel)

			if cerr != nil {
				//Logger.DeleteEntry(m.Mentions[0].ID)
				fmt.Println("Invalid channel detected on ", m.Mentions[0].ID)
			} else {
				LastGuild, _ := s.Guild(LastChan.GuildID)
				LastChanName = LastChan.Name
				LastGuildName = LastGuild.Name
			}

			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Last messaged server & channel", Value: fmt.Sprintln(LastGuildName, ", ", LastChanName)})
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Last message", Value: user.LastMessage})
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Last played game", Value: user.LastGame})
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, &retEmbed)
			if err != nil {
				fmt.Println(err)
			}
			return fmt.Sprintln("Here you go ", m.Author.Mention(), "!")
		} else {
			return "User not yet registered."
		}
	} else {
		return "Invalid mention!"
	}
	return "error?"
}
