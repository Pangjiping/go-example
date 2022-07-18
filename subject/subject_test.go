package subject

import "testing"

func Test_SubjectAndObserver(t *testing.T) {
	shirtItem := NewSubjectImpl("Nike Shirt")
	observerFirst := &customer{id: "abc@gmail.com"}
	observerSecond := &customer{id: "xyz@gmail.com"}

	shirtItem.Register(observerFirst)
	shirtItem.Register(observerSecond)

	shirtItem.updateAvailability()
}
