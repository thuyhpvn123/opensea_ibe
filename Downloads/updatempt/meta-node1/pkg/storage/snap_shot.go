package storage

type SnapShot interface {
	GetIterator() IIterator
}
