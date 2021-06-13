package main

import(
  "syscall"
  "runtime"
  "os"
  "os/signal"

  "github.com/saintwish/core"
  "github.com/saintwish/core/log"
  "github.com/saintwish/core/conf"

  "github.com/bwmarrin/discordgo"
)

func main() {
  var err error

  conf.SetFile("session.yaml", "./runtime/", 0444) //0677
  conf.Read()

  log.InitDefaultCore("log.log", conf.GetString("log.dir"))
  log.InitDefaultBot("log.log", conf.GetString("log.dir"))

  core.S, err = discordgo.New("Bot " + conf.GetString("core.token"))
  if err != nil {
		log.Bot.Errorln("Error creating Discord session, ", err)
		os.Exit(3)
	}

  _, err = core.S.User("@me")
	if err != nil {
    log.Bot.Errorln("Error getting Discord user, ", err)
		os.Exit(3)
	}

  err = core.S.Open()
	if err != nil {
    log.Bot.Errorln("Error getting Discord web socket, ", err)
		os.Exit(3)
	}

  //Max threads allowed.
  if conf.GetInt("core.maxthreads") != 0 {
    runtime.GOMAXPROCS(conf.GetInt("core.maxthreads"))
  }

  //Event Handlers.
  core.S.AddHandler(guildJoinEvent)
  core.S.AddHandler(guildKickEvent)
  core.S.AddHandler(guildRoleCreate)
  core.S.AddHandler(guildRoleDelete)

  core.S.AddHandler(messageCreateEvent)
  core.S.AddHandler(messageDeleteEvent)
  core.S.AddHandler(messageReactAddEvent)
  core.S.AddHandler(messageReactRemoveEvent)

  core.S.AddHandler(memberJoinEvent)
  core.S.AddHandler(memberRemoveEvent)

  core.S.AddHandler(channelCreateEvent)
  core.S.AddHandler(channelDeleteEvent)

  //Run bot
  log.Bot.Println("Bot is now running. Press CTRL-C to exit.")
  sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

  core.S.Close()
}

func panicRecovery() {
  if panic := recover(); panic != nil {
    log.Bot.Panicln("Alieline has encountered an unrecoverable error and as crashed.")
    log.Bot.Panicln("Crash Information: " + panic.(error).Error())

    stack := make([]byte, 65536)
    l := runtime.Stack(stack, true)

    log.Bot.Panic("Stack trace:\n" + string(stack[:l]))

    os.Exit(1)
  }
}
