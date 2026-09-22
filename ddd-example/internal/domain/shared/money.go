package shared

import "errors"

// Value Object: immutable, disamakan lewat nilainya (bukan identitas).
type Money struct {
	amount   float64
	currency string
}

func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}
	if currency == "" {
		return Money{}, errors.New("currency is required")
	}
	return Money{amount: amount, currency: currency}, nil
}

func (m Money) Amount() float64 {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) Multiply(factor float64) Money {
	return Money{amount: m.amount * factor, currency: m.currency}
}
