/*
Copyright 2017 daswf852@outlook.com

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
for porting to different protocols, edit:
	'Command' in struct 'Command'

*/
package Core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type Command struct {
	Name              string
	ArgumentCount     int //0 for functions that ignore arguments or can take variable arguments
	HelpMsg           string
	UsageMsg          string
	IsDisplayedOnHelp bool
	Command           func(Arguments, *discordgo.Session, *discordgo.MessageCreate) string
}
type Arguments struct {
	Args  []string
	Count int
}
type Parser struct {
	commands   map[string]Command
	prefix     string
	unknownMsg bool
}

func MakeParser() Parser {
	return Parser{
		make(map[string]Command), "", true}
}
func (p *Parser) SetPrefix(pr string) {
	p.prefix = pr

}
func (p *Parser) GetPrefix() string {
	return p.prefix
}
func makeArguments(s string) Arguments {
	parsed := strings.Split(s, " ")
	return Arguments{parsed, len(parsed) - 1}
}

func (p *Parser) Register(c *Command) {
	if c != nil {
		fmt.Println("Registered command: ", c.Name)
		p.commands[c.Name] = *c
	}
}

func (p *Parser) Execute(s *discordgo.Session, m *discordgo.MessageCreate) string {
	arguments := makeArguments(m.Content)
	valid := strings.HasPrefix(arguments.Args[0], p.prefix)
	function, exists := p.commands[strings.TrimLeft(arguments.Args[0], p.prefix)]
	if valid && exists {
		if function.ArgumentCount <= arguments.Count {
			return function.Command(arguments, s, m)
		} else {
			return fmt.Sprintln("Minimum argument requirement not met, it needs to be atleast ", function.ArgumentCount, "but is ", arguments.Count)
		}
	} else if (strings.TrimLeft(arguments.Args[0], p.prefix) == "help") && valid {
		if len(arguments.Args) > 1 {
			return p.help(arguments.Args[1])
		} else {
			return p.help("")
		}
	} else {
		if p.unknownMsg && valid {
			return p.help("")
		}
	}
	return ""
}

func (p *Parser) help(cmd string) string {
	var retStr string

	retStr = fmt.Sprintf("%sThe current prefix is: %s\n", retStr, p.prefix)

	foundCmd, isCmdFound := p.commands[cmd]

	if isCmdFound {
		if foundCmd.IsDisplayedOnHelp {
			retStr = fmt.Sprintf("%sCommand: %s\n", retStr, foundCmd.Name)
			retStr = fmt.Sprintf("%sMinimum argument count: %v\n", retStr, foundCmd.ArgumentCount)
			retStr = fmt.Sprintf("%sHelp message:\n\t%s\n", retStr, foundCmd.HelpMsg)
			retStr = fmt.Sprintf("%sUsage: %s", retStr, foundCmd.UsageMsg)
		}
	} else {
		retStr = fmt.Sprintf("%sCommand list:\n\t", retStr)
		for _, v := range p.commands {
			retStr = fmt.Sprintf("%s**%s** - %s\n\t", retStr, v.Name, v.HelpMsg)
		}
	}
	return retStr
}