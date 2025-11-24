package stagen

import "time"

type FakeTimeSpecImpl struct {
	time time.Time
}

func NewFakeTimeSpec(time time.Time) *FakeTimeSpecImpl {
	return &FakeTimeSpecImpl{
		time: time,
	}
}

func (t *FakeTimeSpecImpl) ModTime() time.Time {
	return t.time
}

func (t *FakeTimeSpecImpl) AccessTime() time.Time {
	return t.time
}

func (t *FakeTimeSpecImpl) ChangeTime() time.Time {
	return t.time
}

func (t *FakeTimeSpecImpl) BirthTime() time.Time {
	return t.time
}

func (t *FakeTimeSpecImpl) HasChangeTime() bool {
	return true
}

func (t *FakeTimeSpecImpl) HasBirthTime() bool {
	return true
}
