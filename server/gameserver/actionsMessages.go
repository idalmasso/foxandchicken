package gameserver

type actionMessageTypes string

const (
	//ActionMessageCreateRoom is the action to create a room
	ActionMessageCreateRoom actionMessageTypes = "CREATEROOM"
	//ActionMessageLeaveRoom is the action to leave a room
	ActionMessageLeaveRoom actionMessageTypes = "LEAVEROOM"
	//ActionMessageJoinRoom is the action to join a room
	ActionMessageJoinRoom actionMessageTypes = "JOINROOM"
	//ActionMessageMovement is the action to send a movement
	ActionMessageMovement actionMessageTypes = "MOVEMENT"
	//ActionListRooms is the action to request a list for all rooms
	ActionMessageListRooms actionMessageTypes = "LISTROOMS"
	//ActionMovesRoom is the message with the movements for all players
	ActionMessageMovesRoom actionMessageTypes = "MOVES"
)

type movementStruct struct {
	PositionX     float64 `json:"position_x"`
	PositionY     float64 `json:"position_y"`
	VelocityX     float64 `json:"velocity_x"`
	VelocityY     float64 `json:"velocity_y"`
	AccelerationX float64 `json:"a_x"`
	AccelerationY float64 `json:"a_y"`
	Rotation      float64 `json:"rotation"`
}

//movemementMessage is the message sent to and from the client with the position/velocity of the player
type movemementMessage struct {
	Action  actionMessageTypes `json:"action"`
	Message movementStruct     `json:"message"`
}

//message is an action/message type of message from/to the client
type message struct {
	Action  actionMessageTypes `json:"action"`
	Message string             `json:"message"`
}

//GetAction returns the action of the message
func (m message) GetAction() actionMessageTypes {
	return m.Action
}

type singleStringReturnMessage struct {
	Message string `json:"message"`
}

//genericMessage is a type for the other messages
type genericMessage struct {
	Action  actionMessageTypes `json:"action"`
	Message interface{}        `json:"message"`
}

//GetAction returns the action of the message
func (m genericMessage) GetAction() actionMessageTypes {
	return m.Action
}
