package commands

type CommandReplyOutcome string

const (
	FAILED  CommandReplyOutcome = "failed"
	SUCCESS CommandReplyOutcome = "success"
)

type Command interface {
}
