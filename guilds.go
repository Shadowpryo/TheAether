package main

import (

	"sync"
	"errors"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type GuildsManager struct {

	db          *DBHandler
	querylocker sync.RWMutex

}


type GuildRecord struct {

	ID 			string `storm:"id"` // primary key
	Name 		string

	Region 		string
	Icon		string

	AFKChannel 	string
	AFKTimeout 	int

	OwnerID		string
	RoleIDs		[]string
	UserIDs		[]string

	AdminID		string
	ModeratorID	string
	BuilderID	string
	EveryoneID	string
}



func (h *GuildsManager) SaveGuildToDB(guild GuildRecord) (err error) {
	h.querylocker.Lock()
	defer h.querylocker.Unlock()

	db := h.db.rawdb.From("Guilds")
	err = db.Save(&guild)
	return err
}

func (h *GuildsManager) RemoveGuildFromDB(guild GuildRecord) (err error) {
	h.querylocker.Lock()
	defer h.querylocker.Unlock()

	db := h.db.rawdb.From("Guilds")
	err = db.DeleteStruct(&guild)
	return err
}

func (h *GuildsManager) RemoveGuildByID(guildID string) (err error) {

	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}

	err = h.RemoveGuildFromDB(guild)
	if err != nil {
		return err
	}

	return nil
}

func (h *GuildsManager) GetGuildByID(guildID string) (guild GuildRecord, err error) {

	guilds, err := h.GetAllGuilds()
	if err != nil{
		return guild, err
	}

	for _, i := range guilds {
		if i.ID == guildID{
			return i, nil
		}
	}

	return guild, errors.New("No guild record found")
}

func (h *GuildsManager) GetGuildByName(guildname string, guildID string) (guild GuildRecord, err error) {

	guilds, err := h.GetAllGuilds()
	if err != nil{
		return guild, err
	}

	for _, i := range guilds {
		if i.Name == guildname && i.ID == guildID{
			return i, nil
		}
	}

	return guild, errors.New("No guild record found")
}


// GetAllRooms function
func (h *GuildsManager) GetAllGuilds() (guildlist []GuildRecord, err error) {
	h.querylocker.Lock()
	defer h.querylocker.Unlock()

	db := h.db.rawdb.From("Guilds")
	err = db.All(&guildlist)
	if err != nil {
		return guildlist, err
	}

	return guildlist, nil
}



func (h *GuildsManager) AddRoleToGuild(guildID string, roleID string) (err error) {

	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}

	for _, role := range guild.RoleIDs {
		if role == roleID {
			return nil
		}
	}

	guild.RoleIDs = append(guild.RoleIDs, roleID)

	err = h.SaveGuildToDB(guild)
	if err != nil {
		return err
	}

	return nil
}

func (h *GuildsManager) RemoveRoleFromGuild(guildID string, roleID string) (err error) {
	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}

	guild.RoleIDs = RemoveStringFromSlice(guild.RoleIDs, roleID)

	err = h.SaveGuildToDB(guild)
	if err != nil {
		return err
	}

	return nil
}



func (h *GuildsManager) IsGuildRegistered(guildID string) (valid bool) {

	guildlist, err := h.GetAllGuilds()
	if err != nil {
		return false
	}

	for _, guild := range guildlist {
		if guild.ID == guildID {
			return true
		}
	}
	return false
}

func (h *GuildsManager) IsGuildIDValid(guildID string, s *discordgo.Session) (valid bool) {

	_, err := s.Guild(guildID)
	if err != nil{
		return false
	}
	return false
}



func (h *GuildsManager) AddUserToGuild(guildID string, userID string) (err error) {

	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}

	for _, user := range guild.UserIDs {
		if user == userID {
			return nil
		}
	}

	guild.UserIDs = append(guild.UserIDs, userID)

	err = h.SaveGuildToDB(guild)
	if err != nil {
		return err
	}

	return nil
}

func (h *GuildsManager) RemoveUserFromGuild(guildID string, userID string) (err error) {
	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}

	guild.UserIDs = RemoveStringFromSlice(guild.UserIDs, userID)

	err = h.SaveGuildToDB(guild)
	if err != nil {
		return err
	}

	return nil
}



func (h *GuildsManager) SetAdminID(guildID string, adminID string) (err error) {

	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}
	guild.AdminID = adminID

	err = h.SaveGuildToDB(guild)
	if err != nil {
		return err
	}
	return nil
}

func (h *GuildsManager) SetModeratorID(guildID string, moderatorID string) (err error) {

	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}
	guild.ModeratorID = moderatorID

	err = h.SaveGuildToDB(guild)
	if err != nil {
		return err
	}
	return nil
}

func (h *GuildsManager) SetBuilderID(guildID string, builderID string) (err error) {

	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}
	guild.BuilderID = builderID

	err = h.SaveGuildToDB(guild)
	if err != nil {
		return err
	}
	return nil
}

func (h *GuildsManager) SetEveryoneID(guildID string, everyoneID string) (err error) {

	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return err
	}
	guild.EveryoneID = everyoneID

	err = h.SaveGuildToDB(guild)
	if err != nil {
		return err
	}
	return nil
}




func (h *GuildsManager) GetGuildAdminID(guildID string) (adminID string, err error) {
	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return "", err
	}
	return guild.AdminID, nil
}


func (h *GuildsManager) GetGuildModeratorID(guildID string) (moderatorID string, err error) {
	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return "", err
	}
	return guild.ModeratorID, nil
}

func (h *GuildsManager) GetGuildBuilderID(guildID string) (builderID string, err error) {
	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return "", err
	}
	return guild.BuilderID, nil
}

func (h *GuildsManager) GetGuildEveryoneID(guildID string) (builderID string, err error) {
	guild, err := h.GetGuildByID(guildID)
	if err != nil {
		return "", err
	}
	return guild.EveryoneID, nil
}


func (h *GuildsManager) GetGuildDiscordAdminID(guildID string, s *discordgo.Session) (adminID string, err error) {
	adminID, err = getRoleIDByName(s, guildID, "Admin")
	if err != nil {
		return "", err
	}
	return adminID, nil
}

func (h *GuildsManager) GetGuildDiscordModeratorID(guildID string, s *discordgo.Session) (moderatorID string, err error) {
	moderatorID, err = getRoleIDByName(s, guildID, "Moderator")
	if err != nil {
		return "", err
	}
	return moderatorID, nil
}

func (h *GuildsManager) GetGuildDiscordBuilderID(guildID string, s *discordgo.Session) (builderID string, err error) {
	builderID, err = getRoleIDByName(s, guildID, "Builder")
	if err != nil {
		return "", err
	}
	return builderID, nil
}

func (h *GuildsManager) GetGuildDiscordEveryoneID(guildID string, s *discordgo.Session) (everyoneid string, err error) {
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



func (h *GuildsManager) RegisterGuild(guildID string, s *discordgo.Session) (err error) {

	guildRecord, err := h.GetGuildByID(guildID)
	if err != nil {
		if strings.Contains(err.Error(), "No guild record found"){

			discordguild, err := s.Guild(guildID)
			if err != nil {
				return err
			}

			guildRecord = GuildRecord{ID: guildID, Name: discordguild.Name }
			roles, err := s.GuildRoles(guildID)
			if err != nil {
				return err
			}

			guildRecord.Name = discordguild.Name
			guildRecord.OwnerID = discordguild.OwnerID

			guildRecord.AFKChannel = discordguild.AfkChannelID
			guildRecord.AFKTimeout = discordguild.AfkTimeout
			guildRecord.Icon = discordguild.Icon

			adminID, err := h.GetGuildDiscordAdminID(guildID, s)
			if err != nil {
				return err
			}
			guildRecord.AdminID = adminID

			builderID, err := h.GetGuildDiscordBuilderID(guildID, s)
			if err != nil {
				return err
			}
			guildRecord.BuilderID = builderID

			moderatorID, err := h.GetGuildDiscordModeratorID(guildID, s)
			if err != nil {
				return err
			}
			guildRecord.ModeratorID = moderatorID

			everyoneID, err := h.GetGuildDiscordEveryoneID(guildID, s)
			if err != nil {
				return err
			}
			guildRecord.EveryoneID = everyoneID

			err = h.SaveGuildToDB(guildRecord)
			if err != nil {
				return err
			}

			for _, role := range roles {
				// We should not a get a duplicate record here because this is only for first-time guild registration
				err = h.AddRoleToGuild(guildID, role.ID)
				if err != nil {
					return err
				}
			}

			for _, member := range discordguild.Members {
				err = h.AddUserToGuild(guildID, member.User.ID)
				if err != nil {
					return err
				}
			}

			return nil

		} else {
			return err
		}
	} else {
		return nil
	}
}

