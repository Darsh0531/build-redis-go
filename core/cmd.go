package core

// RedisCmd represents a parsed Redis command
type RedisCmd struct {
	Cmd  string
	Args []string
}
