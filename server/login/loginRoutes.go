package login

import (
	"encoding/json"
	"net/http"
)

type user struct {
	Username string `json:"username"`
}
type userResponse struct {
	Id string `json:"id"`
}
type UserAuthRouter struct {
	LoginUser func(username string) (string, error)
}

func (userRouter UserAuthRouter) ManageRequest(w http.ResponseWriter, r *http.Request) {
	var u user
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id, err := userRouter.LoginUser(u.Username); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		if err = json.NewEncoder(w).Encode(userResponse{Id: id}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
