package market

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// WorldMoneyTransferID uniquely identifies a WorldMoneyTransfer.
type WorldMoneyTransferID string

// NewWorldMoneyTransferID returns a new random 96-bit hex transfer ID.
func NewWorldMoneyTransferID() WorldMoneyTransferID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("market: rand: %w", err))
	}
	return WorldMoneyTransferID(hex.EncodeToString(b[:]))
}

// IsZero reports whether the ID is the empty string.
func (id WorldMoneyTransferID) IsZero() bool { return id == "" }

// WorldMoneyTransfer is a cross-border movement of money in its world-money
// function — settlement between national markets.
//
// "It is only in the world-market that money acquires to the full extent the
// character of the commodity whose bodily form is also the immediate social
// incarnation of human labour in the abstract. Its mode of existence becomes
// adequate to its idea." (Ch. 3 §3.c)
//
// At this level money sheds its national coinage and acts as bullion: the
// quantity exchanged is measured in milligrams of gold, not in the
// denomination of any particular state.
type WorldMoneyTransfer struct {
	ID         WorldMoneyTransferID `json:"id"`
	SenderID   OwnerID              `json:"sender_id"`
	ReceiverID OwnerID              `json:"receiver_id"`
	GoldMg     int64                `json:"gold_mg"`
	CreatedAt  time.Time            `json:"created_at"`
}

// Validate checks that t is well-formed.
func (t WorldMoneyTransfer) Validate() error {
	if t.SenderID.IsZero() {
		return errors.New("market: world_money_transfer sender_id is required")
	}
	if t.ReceiverID.IsZero() {
		return errors.New("market: world_money_transfer receiver_id is required")
	}
	if t.SenderID == t.ReceiverID {
		return errors.New("market: world_money_transfer sender and receiver cannot be the same owner")
	}
	if t.GoldMg <= 0 {
		return errors.New("market: world_money_transfer gold_mg must be positive")
	}
	return nil
}
