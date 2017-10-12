package main

import (
	Core "AtomBot/Core"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strconv"
)

//For commands

func registerCommands(p *Core.Parser) {
	echoCmd := Core.Command{
		"echo",
		1,
		"A command to *echo* what you say", "echo [message]",
		true,
		0,
		echo}
	p.Register(&echoCmd)

	pingCmd := Core.Command{
		"Ping!",
		0,
		"Ping the bot! Or maybe a website in future...", "Ping!",
		true,
		0,
		ping}
	p.Register(&pingCmd)

	succCmd := Core.Command{
		"succ",
		0,
		"succ someone", "succ (opt.)[user mention]",
		true,
		0,
		succ}
	p.Register(&succCmd)

	fuccCmd := Core.Command{
		"fucc",
		0,
		"fucc someone", "fucc (opt.)[user mention]",
		true,
		0,
		fucc}
	p.Register(&fuccCmd)

	whoamiCmd := Core.Command{
		"whoami",
		0,
		"A command to get your user info", "whoami",
		true,
		0,
		whoami}
	p.Register(&whoamiCmd)

	seenCmd := Core.Command{
		"seen",
		1,
		"A command to see the last seen info of a user mentioned", "seen [user mention]",
		true,
		0,
		seen}
	p.Register(&seenCmd)

	shutdownCmd := Core.Command{
		Name:              "shutdown",
		ArgumentCount:     0,
		HelpMsg:           "Shuts down the bot",
		UsageMsg:          "shutdown",
		IsDisplayedOnHelp: true,
		PermLevel:         3,
		Command:           shutdown}
	p.Register(&shutdownCmd)

	setPointsCmd := Core.Command{
		Name:              "setpts",
		ArgumentCount:     0,
		HelpMsg:           "Sets a user's points",
		UsageMsg:          "setpts <mention> <points>",
		IsDisplayedOnHelp: true,
		PermLevel:         3,
		Command:           setPoints}
	p.Register(&setPointsCmd)

	setPermCmd := Core.Command{
		Name:              "setperm",
		ArgumentCount:     0,
		HelpMsg:           "Sets a user's permission level",
		UsageMsg:          "setperm <mention> <permlevel>",
		IsDisplayedOnHelp: true,
		PermLevel:         3,
		Command:           setPerm}
	p.Register(&setPermCmd)
}

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

	user, _ := Logger.GetInfo(m.Author.ID)
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Is bot?", Value: fmt.Sprintln(m.Author.Bot), Inline: false})
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Username", Value: m.Author.String(), Inline: false})
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Mention", Value: m.Author.Mention(), Inline: false})
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "ID", Value: m.Author.ID, Inline: false})
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Permission level", Value: fmt.Sprintf("%v", user.PermLevel), Inline: false})
	retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Fancy points", Value: fmt.Sprintf("%v", user.FancyPoints), Inline: false})
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
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "ID", Value: m.Mentions[0].ID, Inline: false})
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Last messaged server & channel", Value: fmt.Sprintln(LastGuildName, ", ", LastChanName)})
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Last message", Value: user.LastMessage})
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Last played game", Value: user.LastGame})
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Permission level", Value: fmt.Sprintf("%v", user.PermLevel), Inline: false})
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{Name: "Fancy points", Value: fmt.Sprintf("%v", user.FancyPoints), Inline: false})
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

func shutdown(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	bot.Close()
	Config.End(CfgFile, &Parser, &Logger)
	os.Exit(0)
	return "this message shouldnt be seen"
}

func setPerm(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	if len(m.Mentions) > 0 {
		if Logger.EntryExists(m.Mentions[0].ID) {
			val, err := strconv.Atoi(args.Args[2])
			if err == nil {
				Logger.SetPerm(m.Mentions[0].ID, val)
				return "Done!"
			} else {
				return "Not a number!"
			}
		} else {
			return "User is not registered"
		}
	} else {
		return "Invalid mention!"
	}
	return ""
}

func setPoints(args Core.Arguments, s *discordgo.Session, m *discordgo.MessageCreate) string {
	if len(m.Mentions) > 0 {
		if Logger.EntryExists(m.Mentions[0].ID) {
			val, err := strconv.Atoi(args.Args[2])
			if err == nil {
				Logger.SetPoints(m.Mentions[0].ID, val)
				return "Done!"
			} else {
				return "Not a number!"
			}
		} else {
			return "User is not registered"
		}
	} else {
		return "Invalid mention!"
	}
	return ""
}
