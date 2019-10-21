package common

import "testing"

func TestArithmeticGetGroupOrder(t *testing.T) {
	a := &ArithmeticGroup{3}

	var tests = []struct {
		Input     int
		Output    int
		ShouldErr bool
	}{
		{0, 0, true},
		{1, 1, false},

		{2, 2, false},
		{3, 2, false},
		{4, 2, false},

		{5, 3, false},
		{6, 3, false},
		{7, 3, false},

		{8, 4, false},
		{9, 4, false},
		{10, 4, false},
	}

	for i, tt := range tests {
		result, err := a.GetGroupOrder(tt.Input)
		if err != nil && !tt.ShouldErr {
			t.Errorf("test #%d: unexpected error: %v", i, err)
		}
		if err == nil {
			if tt.ShouldErr {
				t.Errorf("test #%d: expected error, got none", i)
			}
			if result != tt.Output {
				t.Errorf("test #%d: wrong result. expected %v, but %v", i, tt.Output, result)
			}
		}
	}
}

func TestArithmeticGetGroupRange(t *testing.T) {
	a := &ArithmeticGroup{3}

	var tests = []struct {
		GroupOrder int
		StartOrder int
		LastOrder  int
		ShouldErr  bool
	}{
		{0, 0, 0, true},
		{1, 1, 1, false},
		{2, 2, 4, false},
		{3, 5, 7, false},
		{4, 8, 10, false},
	}

	for i, tt := range tests {
		startOrder, lastOrder, err := a.GetGroupRange(tt.GroupOrder)

		if err != nil && !tt.ShouldErr {
			t.Errorf("test #%d: unexpected error: %v", i, err)
		}
		if err == nil {
			if tt.ShouldErr {
				t.Errorf("test #%d: expected error, got none", i)
			}

			if startOrder != tt.StartOrder || lastOrder != tt.LastOrder {
				t.Errorf("test #%d: wrong result. expected start : %d, last : %d but start : %d, last : %d",
					i, tt.StartOrder, tt.LastOrder, startOrder, lastOrder)
			}
		}
	}
}
