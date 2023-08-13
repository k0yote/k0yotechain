package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/k0yote/privatechain/types"
)

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrInsufficientBalance = errors.New("insufficient account balance")
)

type Account struct {
	Address types.Address
	Balance uint64
}

func (a *Account) String() string {
	return fmt.Sprintf("%d", a.Balance)
}

type AccountState struct {
	mu       sync.RWMutex
	accounts map[types.Address]*Account
}

func NewAccountState() *AccountState {
	return &AccountState{
		accounts: make(map[types.Address]*Account),
	}
}

func (s *AccountState) CreateAccount(address types.Address) *Account {
	s.mu.Lock()
	defer s.mu.Unlock()

	acc := &Account{Address: address}
	s.accounts[address] = acc
	return acc
}

func (s *AccountState) CreateAccountWithBalance(address types.Address, balance uint64) *Account {
	s.mu.Lock()
	defer s.mu.Unlock()

	acc := &Account{Address: address, Balance: balance}
	s.accounts[address] = acc
	return acc
}

func (s *AccountState) GetAccount(address types.Address) (*Account, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.getAccountWithoutLock(address)
}

func (s *AccountState) getAccountWithoutLock(address types.Address) (*Account, error) {
	account, ok := s.accounts[address]
	if !ok {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (s *AccountState) GetBalance(address types.Address) (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	account, err := s.getAccountWithoutLock(address)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}

func (s *AccountState) Transfer(from, to types.Address, amount uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fromAccount, err := s.getAccountWithoutLock(from)
	if err != nil {
		return err
	}

	if fromAccount.Balance < amount {
		return ErrInsufficientBalance
	}

	fromAccount.Balance -= amount

	if s.accounts[to] == nil {
		s.accounts[to] = &Account{
			Address: to,
		}
	}

	s.accounts[to].Balance += amount

	return nil
}
