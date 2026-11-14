package whatsapp

import (
	"context"
	"sync"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
)

// MemoryContactStore implements store.ContactStore purely in RAM without SQLite persistence.
type MemoryContactStore struct {
	mu       sync.RWMutex
	contacts map[types.JID]*types.ContactInfo
}

// NewMemoryContactStore creates an in-memory contact store.
func NewMemoryContactStore() *MemoryContactStore {
	return &MemoryContactStore{
		contacts: make(map[types.JID]*types.ContactInfo),
	}
}

func (s *MemoryContactStore) PutPushName(ctx context.Context, user types.JID, pushName string) (bool, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user = user.ToNonAD()
	contact, ok := s.contacts[user]
	if !ok {
		contact = &types.ContactInfo{Found: true, PushName: pushName}
		s.contacts[user] = contact
		return true, "", nil
	}

	old := contact.PushName
	if old == pushName {
		return false, old, nil
	}
	contact.PushName = pushName
	contact.Found = true
	return true, old, nil
}

func (s *MemoryContactStore) PutBusinessName(ctx context.Context, user types.JID, businessName string) (bool, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user = user.ToNonAD()
	contact, ok := s.contacts[user]
	if !ok {
		contact = &types.ContactInfo{Found: true, BusinessName: businessName}
		s.contacts[user] = contact
		return true, "", nil
	}

	old := contact.BusinessName
	if old == businessName {
		return false, old, nil
	}
	contact.BusinessName = businessName
	contact.Found = true
	return true, old, nil
}

func (s *MemoryContactStore) PutContactName(ctx context.Context, user types.JID, fullName, firstName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user = user.ToNonAD()
	contact, ok := s.contacts[user]
	if !ok {
		contact = &types.ContactInfo{Found: true, FullName: fullName, FirstName: firstName}
		s.contacts[user] = contact
		return nil
	}

	contact.FullName = fullName
	contact.FirstName = firstName
	contact.Found = true
	return nil
}

func (s *MemoryContactStore) PutAllContactNames(ctx context.Context, contacts []store.ContactEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range contacts {
		user := entry.JID.ToNonAD()
		contact, ok := s.contacts[user]
		if !ok {
			contact = &types.ContactInfo{Found: true, FullName: entry.FullName, FirstName: entry.FirstName}
			s.contacts[user] = contact
		} else {
			contact.FullName = entry.FullName
			contact.FirstName = entry.FirstName
			contact.Found = true
		}
	}
	return nil
}

func (s *MemoryContactStore) PutManyRedactedPhones(ctx context.Context, entries []store.RedactedPhoneEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, entry := range entries {
		user := entry.JID.ToNonAD()
		contact, ok := s.contacts[user]
		if !ok {
			contact = &types.ContactInfo{Found: true, RedactedPhone: entry.RedactedPhone}
			s.contacts[user] = contact
		} else {
			contact.RedactedPhone = entry.RedactedPhone
			contact.Found = true
		}
	}
	return nil
}

func (s *MemoryContactStore) GetContact(ctx context.Context, user types.JID) (types.ContactInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user = user.ToNonAD()
	if contact, ok := s.contacts[user]; ok && contact != nil {
		return *contact, nil
	}
	return types.ContactInfo{Found: false}, nil
}

func (s *MemoryContactStore) GetAllContacts(ctx context.Context) (map[types.JID]types.ContactInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make(map[types.JID]types.ContactInfo, len(s.contacts))
	for k, v := range s.contacts {
		if v != nil {
			res[k] = *v
		}
	}
	return res, nil
}

var _ store.ContactStore = (*MemoryContactStore)(nil)
