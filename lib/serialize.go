package lib

import (
	"agus/configs"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type Quoted struct {
	Message     *waE2E.Message
	ID          types.MessageID
	Participant types.JID
}

type M struct {
	Full        *events.Message
	ID          string
	Chat        types.JID
	Sender      types.JID
	IsFromMe    bool
	IsGroup     bool
	IsOwner     bool
	PushName    string
	Mentions    []types.JID
	Quoted      *Quoted
	ContextInfo *waE2E.ContextInfo
	Text        string
	Prefix      string
	IsCmd       bool
	Command     string
	Args        []string
	Query       string
	Reply       func(string)
}

func Includes(text string, substr ...string) bool {
	for _, str := range substr {
		if strings.Contains(text, str) {
			return true
		}
	}
	return false
}

func getText(str ...string) string {
	for _, s := range str {
		if s != "" {
			return s
		}
	}
	return ""
}
func getCtxInfo(ctxinfo ...*waE2E.ContextInfo) *waE2E.ContextInfo {
	for _, q := range ctxinfo {
		if q != nil {
			return q
		}
	}
	return &waE2E.ContextInfo{}
}
func getMentions(s ...string) []types.JID {
	var jids []types.JID
	for _, id := range s {
		jid, _ := types.ParseJID(id)
		jids = append(jids, jid)
	}
	return jids
}

func Serialize(msg *events.Message, conn *IClient) *M {
	var m = &M{}
	message := msg.Message
	conversation := message.GetConversation()
	extend := message.GetExtendedTextMessage()
	image := message.GetImageMessage()
	video := message.GetVideoMessage()
	m.ContextInfo = getCtxInfo(
		extend.GetContextInfo(),
		image.GetContextInfo(),
		video.GetContextInfo())

	if m.ContextInfo.QuotedMessage != nil {
		p, _ := types.ParseJID(m.ContextInfo.GetParticipant())
		m.Quoted = &Quoted{
			Message:     m.ContextInfo.GetQuotedMessage(),
			ID:          m.ContextInfo.GetStanzaID(),
			Participant: p,
		}
	}
	m.Mentions = getMentions(m.ContextInfo.MentionedJID...)

	m.Full = msg
	m.Text = getText(
		conversation,
		extend.GetText(),
		image.GetCaption(),
		video.GetCaption(),
	)

	m.Chat = msg.Info.Chat
	m.Sender = msg.Info.Sender
	m.ID = msg.Info.ID
	m.IsGroup = msg.Info.IsGroup
	m.IsFromMe = msg.Info.IsFromMe
	m.IsOwner = Includes(m.Sender.User, configs.CONFIG.Owner...) || m.IsFromMe
	m.PushName = msg.Info.PushName
	m.Prefix = configs.CONFIG.Prefix
	m.Args = strings.Split(m.Text, " ")[1:]
	m.IsCmd = strings.HasPrefix(strings.Split(m.Text, " ")[0], m.Prefix)
	m.Query = strings.Join(m.Args, " ")
	if m.IsCmd {
		m.Command = strings.Replace(strings.Split(m.Text, " ")[0], m.Prefix, "", 1)
	} else {
		m.Command = ""
	}
	m.Reply = func(s string) {
		conn.SendText(m.Chat, s, &waE2E.ContextInfo{
			QuotedMessage: m.Full.Message,
			StanzaID:      &m.ID,
			Participant:   proto.String(m.Sender.String()),
			RemoteJID:     proto.String(m.Chat.String()),
			MentionedJID:  []string{m.Sender.String()},
			// ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
			// 	Title:     proto.String("H A L O"),
			// 	Body:      proto.String("Simple Bot WhatsApp Using Go Lang"),
			// 	SourceURL: proto.String("https://web.whatsapp.com/"),
			// },
		})
	}

	return m
}
