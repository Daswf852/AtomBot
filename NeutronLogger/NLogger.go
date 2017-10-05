/*
Copyright 2017 daswf852@outlook.com

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package NeutronLogger

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"time"
)

type UserInfo struct {
	LastSeen    time.Time
	LastMessage string
	LastChannel string
	LastGame    string
	dogeCoins   int
}

type Logger struct {
	users map[string]*UserInfo
}

func MakeLogger() Logger {
	return Logger{make(map[string]*UserInfo)}
}
func (l *Logger) MakeUser(id string) {
	_, exists := l.users[id]
	if !exists {
		l.users[id] = &UserInfo{
			LastSeen:    time.Now(),
			LastMessage: "Last message not recorded",
			LastGame:    "Last played game not recorded"}
	}
}

func (l *Logger) GetInfo(id string) (*UserInfo, int) {
	user, exists := l.users[id]
	if exists {
		return user, 0
	} else {
		return &UserInfo{}, 1
	}
}

func (l *Logger) OutToFile(f string) {
	j, jerr := json.Marshal(l.users)
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

func (l *Logger) ReadFromFile(f string) {
	j, ioerr := ioutil.ReadFile(f)
	if ioerr != nil {
		fmt.Println("error while reading from file: \n", ioerr)
		return
	}
	jerr := json.Unmarshal(j, &l.users)
	if jerr != nil {
		fmt.Println("error while parsing json: \n", jerr)
		return
	}
	fmt.Println("Done reading from file ", f)
}

func (l *Logger) EntryExists(id string) bool {
	_, exists := l.users[id]
	return exists
}

func (l *Logger) UpdateEntryMsg(id string, m *discordgo.MessageCreate) {
	if l.EntryExists(id) {
		l.users[id].LastSeen = time.Now()
		l.users[id].LastMessage = m.Content
		l.users[id].LastChannel = m.ChannelID
	} else {
		l.MakeUser(id)
		l.UpdateEntryMsg(id, m)
	}
}

func (l *Logger) UpdateEntryPresence(id string, p *discordgo.PresenceUpdate) {
	if l.EntryExists(id) {
		l.users[id].LastSeen = time.Now()
	} else {
		l.MakeUser(id)
		l.UpdateEntryPresence(id, p)
	}
}

func (l *Logger) DeleteEntry(id string) {
	if l.EntryExists(id) {
		delete(l.users, id)
	}
}
