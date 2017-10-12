package Core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	prefix                 string
	Playing                string
	NlDatabaseFile         string
	AnnouncementChannelIDs []string
	BootMessage            string
	ClosingMessage         string
}

func (c *Config) readFromFile(f string) {
	j, ioerr := ioutil.ReadFile(f)
	if ioerr != nil {
		fmt.Println("error while reading from file: \n", ioerr)
		return
	}
	jerr := json.Unmarshal(j, &c)
	if jerr != nil {
		fmt.Println("error while parsing json: \n", jerr)
		return
	}
	fmt.Println("Done reading from file ", f)
}

func (c *Config) OutToFile(f string) {
	j, jerr := json.MarshalIndent(*c, "", "\t")
	if jerr != nil {
		fmt.Println("error while generating json: \n", jerr)
		return
	}
	ioerr := ioutil.WriteFile(f, j, 0664)
	if ioerr != nil {
		fmt.Println("error while writing json to file: \n", ioerr)
		return
	}
	fmt.Println("Done writing to file ", f)
}

func MakeConfig() Config {
	return Config{
		Playing:        "github.com/Daswf852/AtomBot",
		NlDatabaseFile: "db.json",
		BootMessage:    "Bot is up!",
		ClosingMessage: "Bot is closing!"}
}

func (c *Config) Init(file string) (Parser, Logger) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		*c = MakeConfig()
	} else {
		c.readFromFile(file)
	}
	parser := MakeParser()
	parser.SetPrefix(c.prefix)

	logger := MakeLogger()
	logger.ReadFromFile(c.NlDatabaseFile)

	return parser, logger
}

func (c *Config) End(file string, p *Parser, l *Logger) {
	c.prefix = p.GetPrefix()
	l.OutToFile(c.NlDatabaseFile)
	c.OutToFile(file)
}
