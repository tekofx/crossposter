package types

type SocialNetWork int

const (
	Bluesky SocialNetWork = iota
	Instagram
	Telegram
	Twitter
)

func (d SocialNetWork) String() string {
	return [...]string{"Bluesky", "Instagram", "Telegram", "Twitter"}[d]
}
