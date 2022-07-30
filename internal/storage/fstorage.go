package storage

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Fstorage struct {
	mu    sync.RWMutex
	file  *os.File
	index map[string]bool
}

func (f *Fstorage) Append(email string) error {
	if f.index[email] {
		return errors.New("email is already in the list")
	}
	str := fmt.Sprintf("%s\n", email)
	f.mu.Lock()
	_, err := f.file.WriteString(str)
	f.index[email] = true
	f.mu.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func buildIndex(file *os.File) map[string]bool {
	index := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		index[scanner.Text()] = true
	}
	return index
}
func (f *Fstorage) GetEmails() []string {
	f.mu.RLock()
	emails := make([]string, 0, len(f.index))
	for k := range f.index {
		emails = append(emails, k)
	}
	f.mu.RUnlock()
	return emails
}

func Open(filename string) *Fstorage {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic("Storage running failed")
	}
	return &Fstorage{
		file:  file,
		index: buildIndex(file),
	}
}
