package websocket

import (
	"Chatapp/pkg/response"
	"encoding/json"
	"log"
)

// TODO: почему нельзя было передать в параметре пользователя
func (ws *Ws) ConnectUser() {
	// почему response ?
	user := ws.User
	res_user := response.NewUser(user, 1)

	log.Printf("%s Connected\n", user.Username)
	ws.Conns.AddUser(user.ID.Hex(), ws)

	res_channels := []response.Channel{}
	channels := ws.Db.GetChannels(user)
	for _, channel := range channels {
		recipients := []response.User{}
		for _, recipient := range channel.Recipients {
			// почему самому себе не отправляем только в личных чатах?
			if channel.Type == 1 && recipient.Hex() == user.ID.Hex() {
				continue
			}
			recipient, _ := ws.Db.GetUser(recipient.Hex())
			recipients = append(recipients, response.NewUser(recipient, ws.Conns.GetUserStatus(recipient.ID.Hex())))
		}
		res_channels = append(res_channels, response.NewChannel(&channel, recipients))

		status := response.Status{
			UserID:    user.ID.Hex(),
			Status:    1,
			Type:      1,
			ChannelID: channel.ID.Hex(),
		}
		// наверное, по поводу онлайна можно отправлять STATUS_UPDATE по каждому сотруднику всем сотрудникам
		ws.Conns.BroadcastToChannel(channel.ID.Hex(), "STATUS_UPDATE", status)
		ws.Conns.AddUserToChannel(user.ID.Hex(), channel.ID.Hex())
	}

	relationships := ws.Db.GetRelationships(user.ID)
	for _, relationship := range relationships {
		// TODO: где enum ??? relationship.Type
		if relationship.Type != 1 {
			continue
		}
		status := response.Status{
			UserID: user.ID.Hex(),
			Status: 1,
			Type:   0,
		}
		ws.Conns.SendToUser(relationship.ToUserID.Hex(), "STATUS_UPDATE", status)
	}

	ws_msg := WS_Message{
		Event: "READY",
		Data: Ready{
			User:     res_user,
			Channels: res_channels,
		},
	}

	ws_res, _ := json.Marshal(ws_msg)
	ws.Write(ws_res)
}
