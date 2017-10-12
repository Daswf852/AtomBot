/*
Copyright 2017 daswf852@outlook.com

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package Core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Prefix                 string
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
	parser.SetPrefix(c.Prefix)

	logger := MakeLogger()
	logger.ReadFromFile(c.NlDatabaseFile)

	return parser, logger
}

func (c *Config) End(file string, p *Parser, l *Logger) {
	c.Prefix = p.GetPrefix()
	l.OutToFile(c.NlDatabaseFile)
	c.OutToFile(file)
}
