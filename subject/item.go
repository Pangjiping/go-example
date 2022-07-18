package subject

import "fmt"

// subject的具体实现
type SubjectImpl struct {
	ObserverList []observer
	Name         string
	InStock      bool
}

func NewSubjectImpl(name string) *SubjectImpl {
	return &SubjectImpl{
		Name: name,
	}
}

func (s *SubjectImpl) updateAvailability() {
	fmt.Printf("SubjectImpl %s is now in stock\n", s.Name)
	s.InStock = true
	s.NotifyAll()
}

func (s *SubjectImpl) Register(o observer) {
	s.ObserverList = append(s.ObserverList, o)
}

func (s *SubjectImpl) Deregister(o observer) {
	s.ObserverList = removeFromslice(s.ObserverList, o)
}

func (s *SubjectImpl) NotifyAll() {
	for _, observer := range s.ObserverList {
		observer.update(s.Name)
	}
}

func removeFromslice(observerList []observer, observerToRemove observer) []observer {
	observerListLength := len(observerList)
	for i, observer := range observerList {
		if observerToRemove.getID() == observer.getID() {
			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
			return observerList[:observerListLength-1]
		}
	}
	return observerList
}
