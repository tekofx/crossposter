package types

type PostStatus int

const (
	Created PostStatus = iota
	Scheduled
	Posted
)

func (d PostStatus) String() string {
	return [...]string{"Creado", "Programado", "Publicado"}[d]
}
