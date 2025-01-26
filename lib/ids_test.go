package lib

import (
	"testing"
)

func TestBitset_HasID(t *testing.T) {
	tests := []struct {
		name     string
		bitset   Bitset
		id       ComponentID
		expected bool
	}{
		{"has id", Bitset(0b101), ComponentID(2), true},
		{"does not have id", Bitset(0b101), ComponentID(1), false},
		{"empty bitset", Bitset(0), ComponentID(0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bitset.HasID(tt.id)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBitset_AddID(t *testing.T) {
	tests := []struct {
		name     string
		bitset   Bitset
		id       ComponentID
		expected Bitset
	}{
		{"add id to empty set", Bitset(0), ComponentID(1), Bitset(0b10)},
		{"add id already present", Bitset(0b10), ComponentID(1), Bitset(0b10)},
		{"add id not present", Bitset(0b1), ComponentID(2), Bitset(0b101)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bitset.AddID(tt.id)
			if result != tt.expected {
				t.Errorf("expected %b, got %b", tt.expected, result)
			}
		})
	}
}

func TestBitset_RemoveID(t *testing.T) {
	tests := []struct {
		name     string
		bitset   Bitset
		id       ComponentID
		expected Bitset
	}{
		{"remove existing id", Bitset(0b101), ComponentID(0), Bitset(0b100)},
		{"remove id not present", Bitset(0b101), ComponentID(1), Bitset(0b101)},
		{"remove from empty set", Bitset(0), ComponentID(3), Bitset(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bitset.RemoveID(tt.id)
			if result != tt.expected {
				t.Errorf("expected %b, got %b", tt.expected, result)
			}
		})
	}
}

func TestBitset_With(t *testing.T) {
	tests := []struct {
		name     string
		bitset   Bitset
		withSet  Bitset
		expected bool
	}{
		{"with full subset", Bitset(0b111), Bitset(0b101), true},
		{"with empty subset", Bitset(0b111), Bitset(0), true},
		{"with mismatched subset", Bitset(0b111), Bitset(0b1000), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bitset.Has(tt.withSet)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBitset_Without(t *testing.T) {
	tests := []struct {
		name       string
		bitset     Bitset
		withoutSet Bitset
		expected   bool
	}{
		{"without full missing set", Bitset(0b101), Bitset(0b10), true},
		{"without empty set", Bitset(0b101), Bitset(0), true},
		{"without overlapping set", Bitset(0b101), Bitset(0b100), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bitset.DoesNotHave(tt.withoutSet)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
