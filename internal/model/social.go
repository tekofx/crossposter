package model

type SocialNetWork int

// Declare related constants for each direction starting with index 1
const (
	Bluesky SocialNetWork = iota
	Instagram
	Telegram
	Twitter
)

// String - Creating common behavior - give the type a String function
func (d SocialNetWork) String() string {
	return [...]string{"Bluesky", "Instagram", "Telegram", "Twitter"}[d]
}
