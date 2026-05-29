package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

// --- Hoard ------------------------------------------------------------------

func (m *MySQL) CreateHoard(ctx context.Context, h market.Hoard) (market.Hoard, error) {
	if err := h.Validate(); err != nil {
		return market.Hoard{}, err
	}
	if h.ID.IsZero() {
		h.ID = market.NewHoardID()
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO hoards (id, owner_id, amount, created_at) VALUES (?, ?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q, string(h.ID), string(h.OwnerID), int64(h.Amount), h.CreatedAt); err != nil {
		if isDuplicate(err) {
			return market.Hoard{}, ErrAlreadyExists
		}
		return market.Hoard{}, err
	}
	return h, nil
}

func (m *MySQL) GetHoard(ctx context.Context, id market.HoardID) (market.Hoard, error) {
	const q = `SELECT id, owner_id, amount, created_at FROM hoards WHERE id = ?`
	var h market.Hoard
	var idStr, ownerStr string
	var amount int64
	err := m.db.QueryRowContext(ctx, q, string(id)).Scan(&idStr, &ownerStr, &amount, &h.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Hoard{}, ErrNotFound
	}
	if err != nil {
		return market.Hoard{}, err
	}
	h.ID = market.HoardID(idStr)
	h.OwnerID = market.OwnerID(ownerStr)
	h.Amount = market.Pence(amount)
	return h, nil
}

func (m *MySQL) ListHoards(ctx context.Context) ([]market.Hoard, error) {
	const q = `SELECT id, owner_id, amount, created_at FROM hoards ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Hoard
	for rows.Next() {
		var h market.Hoard
		var idStr, ownerStr string
		var amount int64
		if err := rows.Scan(&idStr, &ownerStr, &amount, &h.CreatedAt); err != nil {
			return nil, err
		}
		h.ID = market.HoardID(idStr)
		h.OwnerID = market.OwnerID(ownerStr)
		h.Amount = market.Pence(amount)
		out = append(out, h)
	}
	return out, rows.Err()
}

// --- PaymentObligation -----------------------------------------------------

func (m *MySQL) CreatePaymentObligation(ctx context.Context, p market.PaymentObligation) (market.PaymentObligation, error) {
	if err := p.Validate(); err != nil {
		return market.PaymentObligation{}, err
	}
	if p.ID.IsZero() {
		p.ID = market.NewPaymentObligationID()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO payment_obligations
		(id, creditor_id, debtor_id, amount, created_at, due_at, paid_at)
		VALUES (?, ?, ?, ?, ?, ?, NULL)`
	if _, err := m.db.ExecContext(ctx, q,
		string(p.ID), string(p.CreditorID), string(p.DebtorID), int64(p.Amount), p.CreatedAt, p.DueAt,
	); err != nil {
		if isDuplicate(err) {
			return market.PaymentObligation{}, ErrAlreadyExists
		}
		return market.PaymentObligation{}, err
	}
	return p, nil
}

func (m *MySQL) GetPaymentObligation(ctx context.Context, id market.PaymentObligationID) (market.PaymentObligation, error) {
	const q = `SELECT id, creditor_id, debtor_id, amount, created_at, due_at, paid_at
		FROM payment_obligations WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	p, err := scanPaymentObligation(row)
	if errors.Is(err, sql.ErrNoRows) {
		return market.PaymentObligation{}, ErrNotFound
	}
	if err != nil {
		return market.PaymentObligation{}, err
	}
	return p, nil
}

func (m *MySQL) ListPaymentObligations(ctx context.Context) ([]market.PaymentObligation, error) {
	const q = `SELECT id, creditor_id, debtor_id, amount, created_at, due_at, paid_at
		FROM payment_obligations ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.PaymentObligation
	for rows.Next() {
		p, err := scanPaymentObligation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (m *MySQL) SettlePaymentObligation(ctx context.Context, id market.PaymentObligationID) (market.PaymentObligation, error) {
	now := m.now().UTC()
	const updateQ = `UPDATE payment_obligations SET paid_at = ? WHERE id = ? AND paid_at IS NULL`
	res, err := m.db.ExecContext(ctx, updateQ, now, string(id))
	if err != nil {
		return market.PaymentObligation{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return market.PaymentObligation{}, err
	}
	if rows == 0 {
		// Either the obligation is unknown or already settled — return whatever
		// the row currently looks like (or ErrNotFound for an unknown id).
		return m.GetPaymentObligation(ctx, id)
	}
	return m.GetPaymentObligation(ctx, id)
}

// rowScanner abstracts *sql.Row and *sql.Rows for the shared scanner below.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanPaymentObligation(r rowScanner) (market.PaymentObligation, error) {
	var p market.PaymentObligation
	var idStr, creditorStr, debtorStr string
	var amount int64
	var paidAt sql.NullTime
	if err := r.Scan(&idStr, &creditorStr, &debtorStr, &amount, &p.CreatedAt, &p.DueAt, &paidAt); err != nil {
		return market.PaymentObligation{}, err
	}
	p.ID = market.PaymentObligationID(idStr)
	p.CreditorID = market.OwnerID(creditorStr)
	p.DebtorID = market.OwnerID(debtorStr)
	p.Amount = market.Pence(amount)
	if paidAt.Valid {
		t := paidAt.Time.UTC()
		p.PaidAt = &t
	}
	return p, nil
}

// --- WorldMoneyTransfer ----------------------------------------------------

func (m *MySQL) CreateWorldMoneyTransfer(ctx context.Context, t market.WorldMoneyTransfer) (market.WorldMoneyTransfer, error) {
	if err := t.Validate(); err != nil {
		return market.WorldMoneyTransfer{}, err
	}
	if t.ID.IsZero() {
		t.ID = market.NewWorldMoneyTransferID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO world_money_transfers (id, sender_id, receiver_id, gold_mg, created_at)
		VALUES (?, ?, ?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q, string(t.ID), string(t.SenderID), string(t.ReceiverID), t.GoldMg, t.CreatedAt); err != nil {
		if isDuplicate(err) {
			return market.WorldMoneyTransfer{}, ErrAlreadyExists
		}
		return market.WorldMoneyTransfer{}, err
	}
	return t, nil
}

func (m *MySQL) ListWorldMoneyTransfers(ctx context.Context) ([]market.WorldMoneyTransfer, error) {
	const q = `SELECT id, sender_id, receiver_id, gold_mg, created_at
		FROM world_money_transfers ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.WorldMoneyTransfer
	for rows.Next() {
		var t market.WorldMoneyTransfer
		var idStr, sendStr, recvStr string
		if err := rows.Scan(&idStr, &sendStr, &recvStr, &t.GoldMg, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.ID = market.WorldMoneyTransferID(idStr)
		t.SenderID = market.OwnerID(sendStr)
		t.ReceiverID = market.OwnerID(recvStr)
		out = append(out, t)
	}
	return out, rows.Err()
}
