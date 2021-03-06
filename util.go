package main

// Utility Functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	// See http://golang.org/pkg/time/#Parse
	globalTimeFormat = "2006-01-02 15:04 MST"
)

// minDuration and maxDuration const for rounding
const (
	minTimeDuration time.Duration = -1 << 63
	maxTimeDuration time.Duration = 1<<63 - 1
)

// RemoveStringFromSlice function
func RemoveStringFromSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// SafeInput function
func SafeInput(s *discordgo.Session, m *discordgo.MessageCreate, conf *Config) bool {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return false
	}

	// Ignore bots
	if m.Author.Bot {
		return false
	}

	if !strings.HasPrefix(m.Content, conf.MainConfig.CP) {
		return false
	}

	// Set our command prefix to the default one within our config file
	message := strings.Fields(m.Content)
	if len(message) < 1 {
		return false
	}

	return true
}

// CleanCommand function
func CleanCommand(input string, conf *Config) (command string, message []string) {

	// Set our command prefix to the default one within our config file
	cp := conf.MainConfig.CP
	message = strings.Fields(input)

	// Remove the prefix from our command
	message[0] = strings.Trim(message[0], cp)
	command = message[0]
	message = RemoveStringFromSlice(message, command)

	return command, message

}

// SplitPayload function
func SplitPayload(input []string) (command string, message []string) {

	// Remove the prefix from our command
	command = input[0]
	message = RemoveStringFromSlice(input, command)

	return command, message

}

// GetArgumentAndFlags function
func GetArgumentAndFlags(input []string) (argument string, flags []string) {

	// Remove the prefix from our command
	command := input[0]
	message := RemoveStringFromSlice(input, command)
	argument, flags = SplitPayload(message)

	return argument, flags
}

// RemoveFromString function
func RemoveFromString(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// SliceToString function
func SliceToString(s []string, separator string) (formatted string) {
	for _, element := range s {
		formatted = formatted + element + separator
	}
	return formatted
}

// CleanChannel function
func CleanChannel(mention string) string {

	mention = strings.TrimPrefix(mention, "<#")
	mention = strings.TrimSuffix(mention, ">")
	return mention

}

// MentionChannel function
func MentionChannel(channelid string, s *discordgo.Session) (mention string, err error) {
	dgchannel, err := s.Channel(channelid)
	if err != nil {
		return "", err
	}

	return "<#" + dgchannel.ID + ">", nil
}

// CheckPermissions function
func CheckPermissions(command string, channelid string, user *User, s *discordgo.Session, com *CommandHandler) bool {

	usergroups, err := com.user.GetGroups(user.ID, s, channelid)
	if err != nil {
		//fmt.Println("Error Retrieving User Groups for " + usermanager.ID)
		return false
	}

	commandgroups, err := com.registry.GetGroups(command)
	if err != nil {
		//fmt.Println("Error Retrieving Registry Groups for " + command)
		return false
	}

	commandchannels, err := com.registry.GetChannels(command)
	if err != nil {
		//fmt.Println("Error Retrieving Channels for " + command)
		return false
	}

	commandusers, err := com.registry.GetUsers(command)
	if err != nil {
		//fmt.Println("Error Retrieving Users for " + command)
		return false
	}

	// Verify our channel is valid
	_, err = s.Channel(channelid)
	if err != nil {
		return false
	}

	// Look to see if the provided channel id matches one in the command's channel list
	match := false
	for _, commandchannelid := range commandchannels {
		if commandchannelid == channelid {
			match = true
		}
	}
	// If command is not in channel list we return false
	if !match {
		return false
	}

	// Look to see if our usermanager ID is in the users list for the command.
	for _, commanduser := range commandusers {
		if commanduser == user.ID {
			return true
		}
	}

	// Finally we want to try to check the usermanager group list
	for _, usergroup := range usergroups {
		for _, commandgroup := range commandgroups {
			if usergroup == commandgroup {
				return true
			}
		}
	}

	return false
}

/*
// MentionOwner function
func MentionOwner(conf *Config, s *discordgo.Session, m *discordgo.MessageCreate) (mention string, err error) {
	usermanager, err := s.User(conf.MainConfig.AdminID)
	if err != nil {
		return mention, err
	}

	return usermanager.Mention(), nil
}

// OwnerName function
func OwnerName(conf *Config, s *discordgo.Session, m *discordgo.MessageCreate) (name string, err error) {
	usermanager, err := s.User(conf.DiscordConfig.AdminID)
	if err != nil {
		return name, err
	}

	return usermanager.Username, nil
}

*/

// IsVoiceChannelEmpty function
func IsVoiceChannelEmpty(s *discordgo.Session, channelid string, botid string) bool {

	channel, err := s.Channel(channelid)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	guild, err := s.Guild(channel.GuildID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if len(guild.VoiceStates) > 0 {
		for _, state := range guild.VoiceStates {
			if state.ChannelID == channelid && state.UserID != botid {
				return false
			}
		}
		return true
	}

	return true

}

// TruncateTime returns the result of rounding d toward zero to a multiple of m.
// If m <= 0, Truncate returns d unchanged.
func TruncateTime(d time.Duration, m time.Duration) time.Duration {
	if m <= 0 {
		return d
	}
	return d - d%m
}

// RoundTime returns the result of rounding d to the nearest multiple of m.
// The rounding behavior for halfway values is to round away from zero.
// If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration,
// Round returns the maximum (or minimum) duration.
// If m <= 0, Round returns d unchanged.
func RoundTime(d time.Duration, m time.Duration) time.Duration {
	if m <= 0 {
		return d
	}
	r := d % m
	if d < 0 {
		r = -r
		if r+r < m {
			return d + r
		}
		if d1 := d - m + r; d1 < d {
			return d1
		}
		return minTimeDuration // overflow
	}
	if r+r < m {
		return d - r
	}
	if d1 := d + m - r; d1 > d {
		return d1
	}
	return maxTimeDuration // overflow
}

// truncateString function Shorten a string to num characters
func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

// GetGuildID function
func getGuildID(s *discordgo.Session, channelID string) (guildID string, err error) {

	channel, err := s.Channel(channelID)
	if err != nil {
		return "", err
	}

	return channel.GuildID, nil
}

// getGuildEveryoneRoleID function
func getGuildEveryoneRoleID(s *discordgo.Session, guildID string) (everyoneid string, err error) {

	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return "", err
	}

	for _, role := range roles {
		if role.Name == "@everyone" {
			return role.ID, nil
		}
	}

	return "", errors.New("Everyone Role ID Not Found")

}

// getGuildChannelIDByName function
func getGuildChannelIDByName(s *discordgo.Session, guildID string, name string) (channelid string, err error) {

	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Name == name {
			return channel.ID, nil
		}
	}

	return "", errors.New("Channel ID Not Found: " + name)
}

// getRoleIDByName function
func getRoleIDByName(s *discordgo.Session, guildID string, name string) (roleid string, err error) {

	name = strings.Title(name)
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return "", err
	}

	for _, role := range roles {
		if role.Name == name {
			return role.ID, nil
		}
	}

	return "", errors.New("Role ID Not Found: " + name)
}

// getRoleNameByID function
func getRoleNameByID(roleID string, guildID string, s *discordgo.Session) (rolename string, err error) {

	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return "", err
	}

	for _, role := range roles {

		if role.ID == roleID {
			return role.Name, nil
		}
	}

	return "", errors.New("Role " + roleID + " not found in guild " + guildID)
}

// getGuildOwnerID function
func getGuildOwnerID(s *discordgo.Session, channelID string) (ownerID string, err error) {

	guildID, err := getGuildID(s, channelID)
	if err != nil {
		return "", err
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		return "", err
	}

	return guild.OwnerID, nil

}

// FlushMessages function
func FlushMessages(s *discordgo.Session, channelID string, count int) (err error) {

	if count <= 0 {
		return errors.New(":rotating_light: Invalid message count supplied")
	}

	if count > 100 {
		return errors.New(":rotating_light: Count must be less than or equal to 100")

	}

	_, err = s.Channel(channelID)
	if err != nil {
		return err
	}

	messages, err := s.ChannelMessages(channelID, count, "", "", "")
	if err != nil {
		return err
	}

	if count > len(messages) {
		count = len(messages)
	}

	var flushmessages []string

	for i := 0; i < count; i++ {
		flushmessages = append(flushmessages, messages[i].ID)
	}

	err = s.ChannelMessagesBulkDelete(channelID, flushmessages)
	if err != nil {
		return err
	}

	return nil
}

// RollDiceString function
func RollDiceString(faces int, count int) (rolls []string) {

	// Always 1 or higher!
	if faces == 0 {
		faces = 1
	}
	if count == 0 {
		count = 1
	}
	// We're not terribly concerned with the randomness here as these aren't secrets
	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	for i := 0; i < count; i++ {
		roll := randomGen.Intn(faces)
		rolls = append(rolls, strconv.Itoa(roll))
	}

	return rolls
}

// RollDice function
func RollDice(faces int, count int) (rolls []int) {

	// Always 1 or higher!
	if faces == 0 {
		faces = 1
	}
	if count == 0 {
		count = 1
	}
	// We're not terribly concerned with the randomness here as these aren't secrets
	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	for i := 0; i < count; i++ {
		roll := randomGen.Intn(faces)
		rolls = append(rolls, roll)
	}

	return rolls
}

// RollDiceAndAdd function
func RollDiceAndAdd(faces int, count int) (total int) {
	// Always 1 or higher!
	if faces == 0 {
		faces = 1
	}
	if count == 0 {
		count = 1
	}
	// We're not terribly concerned with the randomness here as these aren't secrets
	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	for i := 0; i < count; i++ {
		roll := randomGen.Intn(faces)
		total = total + roll
	}

	return total
}

// IsJSON function
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
