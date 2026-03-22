package core

import (
	"errors"
	"net"
)

// EvalAndRespond figures out which command to run based on the Cmd string
func EvalAndRespond(cmd *RedisCmd, c net.Conn) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	default:
		// Fallback for right now (Arpit did this in the video)
		return evalPING(cmd.Args, c)
	}
}

// evalPING handles the PING command logic
func evalPING(args []string, c net.Conn) error {
	var b []byte

	// PING can take a max of 1 argument. E.g., "PING hello"
	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'ping' command")
	}

	// If it's just "PING", respond with simple string "+PONG\r\n"
	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		// If it's "PING hello", respond with bulk string "$5\r\nhello\r\n"
		b = Encode(args[0], false)
	}

	// Write the encoded RESP bytes back to the client
	_, err := c.Write(b)
	return err
}
