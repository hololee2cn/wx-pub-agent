package vo

const (
	// SendPending 发送判定中
	SendPending = iota

	// Sending 发送中
	Sending

	// SendSuccess 发送成功
	SendSuccess

	// SendFailure 发送失败
	SendFailure
)

// MsgStatus 消息发送状态
type MsgStatus struct {
	Val int
}

func (m *MsgStatus) GetPending() int {
	return SendPending
}

func (m *MsgStatus) GetSending() int {
	return Sending
}

func (m *MsgStatus) GetSuccess() int {
	return SendSuccess
}

func (m *MsgStatus) GetFailure() int {
	return SendFailure
}
