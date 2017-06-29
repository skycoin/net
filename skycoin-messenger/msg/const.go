package msg

const (
	MSG_OP_SIZE = 1

	MSG_HEADER_BEGIN = 0
	MSG_OP_BEGIN
	MSG_OP_END       = MSG_OP_BEGIN + MSG_OP_SIZE

	MSG_HEADER_END
)

const (
	OP_REG  = iota
	OP_SEND
	OP_SIZE
)

const (
	PUSH_MSG  = iota
	PUSH_SIZE
)
