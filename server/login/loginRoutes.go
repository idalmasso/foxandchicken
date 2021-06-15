package login

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type UserAddWebSocketInterface interface {
	UpdateWebSocket(*websocket.Conn)
}

type user struct {
	Username string `json:"username"`
}
type userResponse struct {
	Id string `json:"id"`
}
type UserAuthRouter struct {
	LoginUser func(username string) (UserAddWebSocketInterface, error)
}

func (userRouter UserAuthRouter) ManageRequest(w http.ResponseWriter, r *http.Request) {
	var u user
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user, err := userRouter.LoginUser(u.Username); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.UpdateWebSocket(ws)

	}
}
