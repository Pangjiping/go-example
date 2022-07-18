package subject

type Subject interface {
	Register(Observer observer)
	Deregister(Observer observer)
	NotidyAll()
}
