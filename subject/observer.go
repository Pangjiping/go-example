package subject

type observer interface {
	update(s string)
	getID() string
}
