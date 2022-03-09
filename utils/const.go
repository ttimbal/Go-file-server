package utils

import "regexp"

const (
	maxSizeMb   = 5
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
	SEND        = "send"
	HELP        = "help"
	EXIT        = "exit"
	START       = "server start"
	STOP        = "server stop"
	MAX_SIZE    = 1024 * 1024 * maxSizeMb
)

//Regex patterns
var RegexSubscribe, _ = regexp.Compile("\\s*^" + SUBSCRIBE + "\\s*channel:\\w*\\s*")
var RegexUnsubscribe, _ = regexp.Compile("\\s*^" + UNSUBSCRIBE + "\\s*channel:\\w*\\s*")
var RegexSend, _ = regexp.Compile("\\s*^" + SEND + "\\s*channel:\\w*\\s*file:.*\\s*")
var SingleSpacePattern = regexp.MustCompile(`\s+`)
