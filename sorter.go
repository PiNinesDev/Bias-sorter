package main

import (
	"example.com/bias-sorter/db"
)

type Sorter interface {
	Select(int64)
	GetCurrentPair() []db.Entry
	GetQuiz() db.Quiz
}

/*
type MockSorter struct {
	Quiz            db.Quiz
	entries         []db.Entry
	internalCounter int64
}

func (s MockSorter) GetQuiz() db.Quiz {
	return s.Quiz
}

func (s MockSorter) GetCurrentPair() []db.Entry {
	return s.entries
}

func (s MockSorter) Select(index int64) {
	s.entries[0] = s.entries[1]
	s.entries[1] = s.nextEntry()
}

func (s MockSorter) nextEntry() db.Entry {
	return db.Entry{s.internalCounter, 2, fmt.Sprintf("test%v", s.internalCounter)}
}

func NewMockSorter() MockSorter {
	ms := MockSorter{db.Quiz{2, "title"}, make([]db.Entry, 2, 2), 0}
	ms.entries[0] = ms.nextEntry()
	ms.entries[1] = ms.nextEntry()
	return ms
}

func (s MockSorter) serialize() string {
	serializedSorter, _ := json.Marshal(s)
	return string(serializedSorter)
}

func unserializeMockSorter(str string, s *MockSorter) error {
	err := json.Unmarshal([]byte(str), s)
	if err != nil {
		return err
	}
	return nil
}
*/
