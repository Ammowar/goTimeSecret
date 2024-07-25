package storages

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"sync"
	"time"
)

type Secret struct {
	Content   string
	ExpiresAt time.Time
}

type MemoryStorage struct {
	sync.RWMutex
	secrets map[string]*Secret
}

func NewMemoryStorage() *MemoryStorage {

	var m = &MemoryStorage{
		secrets: make(map[string]*Secret),
	}
	go m.cleanupExpiredSecrets()
	return m
}

func (s *MemoryStorage) Create(content string, expiration time.Duration) (string, time.Time, error) {
	s.Lock()
	defer s.Unlock()
	id := generateID()
	expiresAt := time.Now().Add(expiration)
	s.secrets[id] = &Secret{Content: content, ExpiresAt: expiresAt}
	return id, expiresAt, nil
}

func (s *MemoryStorage) View(id string) (*Secret, error) {
	s.RLock() // Read lock
	secret, ok := s.secrets[id]
	s.RUnlock()
	if !ok {
		return nil, errors.New("Secret not found or already viewed")
	}
	defer s.Delete(id)
	if time.Now().After(secret.ExpiresAt) {
		return nil, errors.New("Secret has expired")
	}
	return secret, nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.secrets, id)
	return nil
}

func generateID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (s *MemoryStorage) cleanupExpiredSecrets() {
	for {
		time.Sleep(1 * time.Hour)
		for key := range s.secrets {
			if time.Now().After(s.secrets[key].ExpiresAt) {
				s.Delete(key)
			}
		}
	}
}
