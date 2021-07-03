package gameserver

type actionMessageTypes string

const (
	ActionMessageCreateRoom actionMessageTypes = "CREATEROOM"
	ActionMessageLeaveRoom  actionMessageTypes = "LEAVEROOM"
	ActionMessageJoinRoom   actionMessageTypes = "JOINROOM"
	ActionMessageMovement 	actionMessageTypes="POSITION"
)
type movemementMessage struct {
	Action    actionMessageTypes `json:"action"`
	PositionX float32            `json:"position_x"`
	PositionY float32            `json:"position_y"`
	VelocityX float32            `json:"velocity_x"`
	VelocityY float32            `json:"velocity_y"`
	Rotation  float32            `json:"rotation"`
}
func (m *movemementMessage) GetAction()actionMessageTypes{
	return m.Action
}
type message struct {
	Action  actionMessageTypes `json:"action"`
	Message string             `json:"message"`
}
func (m *message) GetAction()actionMessageTypes{
	return m.Action
}

type singleStringReturnMessage struct {
	Message string `json:"message"`
}

type actionMessage interface{
	GetAction() actionMessageTypes
}
