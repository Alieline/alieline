package main

import (
  "strings"

  "github.com/saintwish/core"
  "github.com/saintwish/core/log"
  "github.com/saintwish/core/conf"
  "github.com/saintwish/core/dyscor"

  "github.com/bwmarrin/discordgo"
)

func guildJoinEvent(s *discordgo.Session, m *discordgo.GuildCreate) {
  defer panicRecovery()

  //addServer(m.Name, m.Guild.ID)

  if m.Unavailable {
    log.Bot.Warnln("Joined unavailable guild: ", m.Guild.ID)
  }
  log.Bot.Printf("Sucessfully joined guild [%s:%s]", m.Name, m.Guild.ID)
}

func guildKickEvent(s *discordgo.Session, m *discordgo.GuildDelete) {
  defer panicRecovery()

  if m.Unavailable {
    log.Bot.Warnln("Kicked from unavailable guild: ", m.Guild.ID)
  }
  log.Bot.Println("Kicked from guild: ", m.Guild.ID, m.Name)
}

func guildRoleCreate(s *discordgo.Session, m *discordgo.GuildRoleCreate) {
  defer panicRecovery()

  //If the guild settings cannot be found then don't do anything.
  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnGuildRoleCreate {
      go hook.OnGuildRoleCreate(s, m)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("GuildRoleCreate: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}

func guildRoleDelete(s *discordgo.Session, m *discordgo.GuildRoleDelete) {
  defer panicRecovery()

  //If the guild settings cannot be found then don't do anything.
  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnGuildRoleDelete {
      go hook.OnGuildRoleDelete(s, m)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("GuildRoleDelete: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}

func messageCreateEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
  defer panicRecovery()

  //We don't want the bot responding to other bots.
  if m.Author.Bot {
    return
  }

  //Ignore message from a DM
  if ok, _ := dyscor.ComesFromDM(s, m); ok {
    return;
  }

  guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
  if err != nil {
    return
  }

  //If the guild settings cannot be found then don't do anything.
  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnMessageCreate {
      go hook.OnMessageCreate(s, m)
    }

    //Check if the channel is being ignored.
    if _, ok := key.Config.GetStringIntMap("chanignore")[m.ChannelID]; ok {
      return
    }

    //Parse command if it contains guild setting prefix.
    if strings.HasPrefix(m.Content, key.Config.GetString("prefix")) {
      core.ParseCommand(s, m, guildDetails, strings.TrimPrefix(m.Content, key.Config.GetString("prefix")))
      return
    }else if strings.HasPrefix(m.Content, conf.GetString("core.prefix")) {
      core.ParseCommand(s, m, guildDetails, strings.TrimPrefix(m.Content, conf.GetString("core.prefix")))
      return
    }

  }else{
    log.Bot.Errorf("MessageCreate: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}

func messageDeleteEvent(s *discordgo.Session, m *discordgo.MessageDelete) {
  defer panicRecovery()

  //If the guild settings cannot be found then don't do anything.
  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnMessageDelete {
      hook.OnMessageDelete(s, m)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("MessageDelete: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}

func messageReactAddEvent(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
  defer panicRecovery()

  //Ignore reaction from self.
  if m.UserID == s.State.User.ID {
		return
	}

  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnReactionAdd {
      go hook.OnReactionAdd(s, m)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("MessageReactionAdd: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}

func messageReactRemoveEvent(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
  defer panicRecovery()

  if m.UserID == s.State.User.ID {
		return
	}

  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnReactionRemove {
      go hook.OnReactionRemove(s, m)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("MessageReactionRemove: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}



func memberJoinEvent(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
  defer panicRecovery()

  //Ignore bots joining.
  if m.User.Bot == true {
    return
  }

  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnMemberAdd {
      go hook.OnMemberAdd(s, m)
    }

    var joinMsg = key.Config.GetString("joinmessage")
    if joinMsg != "" {
      go dyscor.SendUserDM(s, m.User.ID, joinMsg)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("GuildMemberAdd: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}

func memberRemoveEvent(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
  defer panicRecovery()

  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnMemberRemove {
      go hook.OnMemberRemove(s, m)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("GuildMemberRemove: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}

func channelCreateEvent(s *discordgo.Session, m *discordgo.ChannelCreate) {
  defer panicRecovery()

  if key, ok := core.GData[m.GuildID]; ok {
    //We call our hooks
    for _, hook := range key.Hooks.OnChannelCreate {
      go hook.OnChannelCreate(s, m)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("ChannelCreate: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}

func channelDeleteEvent(s *discordgo.Session, m *discordgo.ChannelDelete) {
  defer panicRecovery()

  if key, ok := core.GData[m.GuildID]; ok {
    if _, ok := key.Config.GetStringIntMap("chanignore")[m.ID]; ok {
      delete(key.Config.GetStringIntMap("chanignore"), m.ID)
    }

    //We call our hooks
    for _, hook := range key.Hooks.OnChannelDelete {
      go hook.OnChannelDelete(s, m)
    }

  }else{
    guildDetails, err := dyscor.GuildDetails(m.GuildID, s)
    if err != nil {
      return
    }

    log.Bot.Errorf("ChannelDelete: Error getting GData for [%s:%s]", m.GuildID, guildDetails.Name)
  }
}
