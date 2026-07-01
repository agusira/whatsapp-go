package lib

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/purpshell/meowcaller"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

type IClient struct {
	WA     *whatsmeow.Client
	Caller *meowcaller.Client
}

var background = context.Background()

func SerializeClient(sock *whatsmeow.Client, caller *meowcaller.Client) *IClient {
	return &IClient{
		WA:     sock,
		Caller: caller,
	}
}

func generateMessageID() types.MessageID {
	id := make([]byte, 15)
	_, err := rand.Read(id)
	if err != nil {
		panic(err)
	}
	return "AC" + strings.ToUpper(hex.EncodeToString(id))
}
func (c *IClient) ParseMention(text string) []string {
	res := []string{}
	matches := regexp.MustCompile("@([0-9]{5,16}|0)").FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		res = append(res, match[1]+"@s.whatsapp.net")
	}
	return res
}

// Call func
func (conn *IClient) Call(jid string) (*meowcaller.Call, error) {
	return conn.Caller.Call(background, jid)
}

func (conn *IClient) SendInteractive(from types.JID, text string, opt *waE2E.ContextInfo) {
	msg := &waE2E.Message{
		InteractiveMessage: &waE2E.InteractiveMessage{
			Body: &waE2E.InteractiveMessage_Body{
				Text: &text,
			},
			Footer: &waE2E.InteractiveMessage_Footer{
				Text: proto.String("Agus"),
			},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons: []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
						// {
						// 	Name:             proto.String("quick_reply"),
						// 	ButtonParamsJSON: proto.String("{}"),
						// },
					},
				},
			},
			ContextInfo: opt,
		},
	}

	conn.WA.SendMessage(background, from, msg, whatsmeow.SendRequestExtra{
		ID: generateMessageID(),
	})
}

func (conn *IClient) SendText(from types.JID, txt string, opts *waE2E.ContextInfo, optn ...whatsmeow.SendRequestExtra) (whatsmeow.SendResponse, error) {
	ok, er := conn.WA.SendMessage(background, from, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        &txt,
			ContextInfo: opts,
		},
	}, whatsmeow.SendRequestExtra{
		ID: generateMessageID(),
	})
	if er != nil {
		return whatsmeow.SendResponse{}, er
	}
	return ok, nil
}

func (conn *IClient) SendImage(from types.JID, data []byte, caption string, opts *waE2E.ContextInfo) (whatsmeow.SendResponse, error) {
	uploaded, err := conn.WA.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return whatsmeow.SendResponse{}, err
	}
	resultImg := &waE2E.Message{
		ImageMessage: &waE2E.ImageMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Caption:       proto.String(caption),
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo:   opts,
		},
	}
	ok, _ := conn.WA.SendMessage(background, from, resultImg, whatsmeow.SendRequestExtra{ID: generateMessageID()})
	return ok, nil
}

func (conn *IClient) SendVideo(from types.JID, data []byte, caption string, opts *waE2E.ContextInfo) (whatsmeow.SendResponse, error) {
	uploaded, err := conn.WA.Upload(background, data, whatsmeow.MediaVideo)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return whatsmeow.SendResponse{}, err
	}
	resultVideo := &waE2E.Message{
		VideoMessage: &waE2E.VideoMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Caption:       proto.String(caption),
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo:   opts,
		},
	}
	ok, er := conn.WA.SendMessage(context.Background(), from, resultVideo, whatsmeow.SendRequestExtra{ID: generateMessageID()})
	if er != nil {
		return whatsmeow.SendResponse{}, er
	}
	return ok, nil
}

func (conn *IClient) SendSticker(from types.JID, data []byte, opts *waE2E.ContextInfo) (whatsmeow.SendResponse, error) {
	uploaded, err := conn.WA.Upload(context.Background(), data, whatsmeow.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return whatsmeow.SendResponse{}, err
	}
	resultImg := &waE2E.Message{
		StickerMessage: &waE2E.StickerMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo:   opts,
		},
	}
	ok, _ := conn.WA.SendMessage(background, from, resultImg, whatsmeow.SendRequestExtra{ID: generateMessageID()})
	return ok, nil
}

func (conn *IClient) AddParticipant(groupJID types.JID, participantJID types.JID) ([]types.GroupParticipant, error) {
	bareJID := types.JID{
		User:   participantJID.User,
		Server: participantJID.Server,
		Device: 0,
	}

	participants, err := conn.WA.UpdateGroupParticipants(background, groupJID, []types.JID{bareJID}, whatsmeow.ParticipantChangeAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to add participant %s to group %s: %w", bareJID.String(), groupJID.String(), err)
	}
	return participants, err
}

func (conn *IClient) RemoveParticipant(groupJID types.JID, participantJID []types.JID) ([]types.GroupParticipant, error) {

	participants, err := conn.WA.UpdateGroupParticipants(background, groupJID, participantJID, whatsmeow.ParticipantChangeRemove)
	if err != nil {
		return nil, fmt.Errorf("failed to remove participant %s from group %s: %w", participantJID[0:], groupJID.String(), err)
	}
	return participants, err
}

func (conn *IClient) PromoteParticipant(groupJID types.JID, participantJID types.JID) ([]types.GroupParticipant, error) {
	bareJID := types.JID{
		User:   participantJID.User,
		Server: participantJID.Server,
		Device: 0,
	}

	participants, err := conn.WA.UpdateGroupParticipants(background, groupJID, []types.JID{bareJID}, whatsmeow.ParticipantChangePromote)
	if err != nil {
		return nil, fmt.Errorf("failed to promote participant %s in group %s: %w", bareJID.String(), groupJID.String(), err)
	}
	return participants, nil
}

func (conn *IClient) DemoteParticipant(groupJID types.JID, participantJID types.JID) ([]types.GroupParticipant, error) {
	bareJID := types.JID{
		User:   participantJID.User,
		Server: participantJID.Server,
		Device: 0,
	}

	participants, err := conn.WA.UpdateGroupParticipants(background, groupJID, []types.JID{bareJID}, whatsmeow.ParticipantChangeDemote)
	if err != nil {
		return nil, fmt.Errorf("failed to demote participant %s in group %s: %w", bareJID.String(), groupJID.String(), err)
	}
	return participants, nil
}

func (conn *IClient) CopyNForward(JID types.JID, message *waE2E.Message) {
	if message.ImageMessage != nil && *message.ImageMessage.ViewOnce {
		*message.ImageMessage.ViewOnce = false
		conn.WA.SendMessage(background, JID, &waE2E.Message{
			ImageMessage: message.ImageMessage,
		})
		return
	} else if message.VideoMessage != nil && *message.VideoMessage.ViewOnce {
		*message.VideoMessage.ViewOnce = false
		conn.WA.SendMessage(background, JID, &waE2E.Message{
			VideoMessage: message.VideoMessage,
		})
		return
	} else if message.AudioMessage != nil && *message.AudioMessage.ViewOnce {
		*message.AudioMessage.ViewOnce = false
		conn.WA.SendMessage(background, JID, &waE2E.Message{
			AudioMessage: message.AudioMessage,
		})
		return
	}
}
