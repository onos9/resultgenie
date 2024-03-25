package bot

import (
	"context"
	"fmt"
	"os"
	"repot/pkg/model"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var dclient *discordgo.Session
var userID string

type Discord struct {
	id        string
	channelID string
	botPrefix string
	*discordgo.Session
}

func NewDiscord() (*Discord, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	return nil, err
	// }

	token := os.Getenv("DISCORD_TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	u, err := session.User("@me")
	if err != nil {
		return nil, err
	}
	dclient = session
	userID = u.ID
	return &Discord{
		id:        u.ID,
		channelID: os.Getenv("CHANNEL_ID"),
		Session:   dclient,
		botPrefix: os.Getenv("BOT_PREFIX"),
	}, nil
}

func GetDiscordInstance() (*Discord, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	return nil, err
	// }

	return &Discord{
		id:        userID,
		channelID: os.Getenv("CHANNEL_ID"),
		Session:   dclient,
		botPrefix: os.Getenv("BOT_PREFIX"),
	}, nil
}

func (d *Discord) Start(ctx context.Context) error {

	d.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == d.id || m.Author.Bot {
			return
		}
		if m.Content == "ping" {
			embed := &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{{
					Type:        discordgo.EmbedTypeRich,
					Title:       "Failed to generate result",
					Description: "Something went wrong while generating your result. Please try again.",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Student Name",
							Value:  "Godsgrace Brown",
							Inline: false,
						},
						{
							Name:   "AdminNo",
							Value:  "00876",
							Inline: true,
						},
						{
							Name:   "ID",
							Value:  "76",
							Inline: true,
						},
						{
							Name:   "URL",
							Value:  "https://google.com",
							Inline: false,
						},
					},
				},
				},
			}
			s.ChannelMessageSendComplex(m.ChannelID, embed)
		}
	})

	err := d.Open()
	if err != nil {
		return err
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	return nil
}

func (d *Discord) SendComplex(msg string, std model.Student) error {
	embed := parseComplexMessage(msg, std)

	_, err := d.ChannelMessageSendComplex(d.channelID, embed)
	if err != nil {
		return fmt.Errorf("failed to send message to Discord channel '%s'", err)
	}

	return nil
}

func (d *Discord) Send(subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	_, err := d.ChannelMessageSend(d.channelID, fullMessage)
	if err != nil {
		return fmt.Errorf("failed to send message to Discord channel '%s'", err)
	}

	return nil
}

func parseComplexMessage(msg string, data model.Student) *discordgo.MessageSend {
	id := int(data.ID)
	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Result Genaration Error",
			Description: msg,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Student Name",
					Value:  data.FullName,
					Inline: false,
				},
				{
					Name:   "AdminNo",
					Value:  fmt.Sprintf("%04d", int(data.AdmissionNo)),
					Inline: true,
				},
				{
					Name:   "ID",
					Value:  fmt.Sprintf("%d", int(id)),
					Inline: true,
				},
				{
					Name:   "URL",
					Value:  "https://llacademy.ng/student-view/" + strconv.Itoa(id),
					Inline: false,
				},
			},
		},
		},
	}

	return embed
}
