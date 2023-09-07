package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var s *discordgo.Session

var (
	commands = []discordgo.ApplicationCommand{
		{
			Name:        "mfuse",
			Description: "Test your morfuse scripts.",
		},
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mfuse": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseModal,
				Data: &discordgo.InteractionResponseData{
					CustomID: "mfuse_" + i.Interaction.Member.User.ID,
					Title:    "Morfuse",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.TextInput{
									CustomID:  "mfuse",
									Label:     "Write here your morfuse script to test.",
									Style:     discordgo.TextInputParagraph,
									Required:  true,
									MaxLength: 2000,
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}
)

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Bot parameters
	var (
		TOKEN     = os.Getenv("TOKEN")
		CLIENT_ID = os.Getenv("CLIENT_ID")
		GUILD_ID  = os.Getenv("GUILD_ID")
	)

	s, err = discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionModalSubmit:
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "**Sent :white_check_mark:**",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			data := i.ModalSubmitData()

			if !strings.HasPrefix(data.CustomID, "mfuse") {
				return
			}

			textValue := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			morfuseResponse := morfuse(textValue)
			username := i.Member.User.Username
			_, err = s.ChannelMessageSendEmbed(i.ChannelID, &discordgo.MessageEmbed{
				Title:       ":green_circle: Script executed",
				Description: fmt.Sprintf("```\n%s```", morfuseResponse),
				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("Execution time: %s\nRequested By: %s • %s", elapsed, username, time.Now().Format("03:04 PM")),
				},
			})

			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		}
	})

	cmdIDs := make(map[string]string, len(commands))

	for _, cmd := range commands {
		rcmd, err := s.ApplicationCommandCreate(CLIENT_ID, GUILD_ID, &cmd)
		if err != nil {
			log.Fatalf("Cannot create slash command %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name
	}

	err = s.Open()
	s.UpdateGameStatus(0, "/mfuse")

	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")

	for id, name := range cmdIDs {
		err := s.ApplicationCommandDelete(CLIENT_ID, GUILD_ID, id)
		if err != nil {
			log.Fatalf("Cannot delete slash command %q: %v", name, err)
		}
	}

}

var start time.Time
var elapsed time.Duration

func morfuse(script string) string {
	// compile and execute the script using the Morfuse executable
	start = time.Now() // get current time
	os.Setenv("LD_LIBRARY_PATH", "morfuse")
	cmd := exec.Command("morfuse/mfuse", script)

	// Redirect the command's output to stdout
	out := catchOutput(cmd)
	return out
}

func catchOutput(cmd *exec.Cmd) string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	elapsed = time.Since(start) // get elapsed time
	if err == nil {
		return stdout.String() // print morfuse execution result
	} else {
		return stderr.String() // print morfuse execution result error
	}

}
