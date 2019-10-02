package polishcalc

import "testing"

func TestSingleTerms(t *testing.T) {
	examples := []struct {
		in string
		want int
	}{
		{"+ 2 2", 4},
		{"* 2 3", 6},
		{"   - 6 3", 3},
		{"- 3 6     ", -3},
		{"/ 6     3", 2},
	}

	for _, v := range examples {
		got, err := Parse(v.in)

		if got != v.want || err != nil {
			t.Errorf("Solve(%q) == %q, want %q", v.in, got, v.want)
		}
	}
}

func TestMultipleTerms(t *testing.T) {
	examples := []struct {
		in string
		want int
	}{
		{"+ 2 + 2 2", 6},
		{"* 2 * 3 + 2 4", 6*3*2},
		{"- 6 * 3 - 6 / 2 * 2 9", 6 - (3 * (6 - (2 / (2 * 9))))},
		{"- 6 - 6 0", 0},
		{"- * / 15 - 7 + 1 1 3 + 2 + 1 1", 5}, // Wikipedia example
	}

	for _, v := range examples {
		got, err := Parse(v.in)

		if got != v.want || err != nil {
			t.Errorf("Solve(%q) == %q, want %q", v.in, got, v.want)
		}
	}
}

func TestBigValues(t *testing.T) {
	examples := []struct {
		in string
		want int
	}{
		{"* 581257125 6167616321", 581257125 * 6167616321},
		{"/ 25125215 2512521612216", 25125215 / 2512521612216},
		{"+ 25216 + 11245 9294", 25216 + 11245 + 9294},
	}

	for _, v := range examples {
		got, err := Parse(v.in)

		if got != v.want || err != nil {
			t.Errorf("Solve(%q) == %q, want %q", v.in, got, v.want)
		}
	}
}

func TestTooFewTermsError(t *testing.T) {
	_, err := Parse("+ 1")
	if err == nil {
		t.Errorf("The code did not error out")
	}
}

func TestTooMuchTermsError(t *testing.T) {
	_, err := Parse("+ 1 2 3")
	if err == nil {
		t.Errorf("The code did not error out")
	}
}

func TestTooMuchTermsErrorComplex(t *testing.T) {
	_, err := Parse("+ 1 / 2 3 2")
	if err == nil {
		t.Errorf("The code did not error out")
	}
}

func TestEmptyStringError(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Errorf("The code did not error out")
	}
}

func TestWrongOperationError(t *testing.T) {
	_, err := Parse("2 s 5")
	if err == nil {
		t.Errorf("The code did not error out")
	}

}

func TestWrongIntError(t *testing.T) {
	_, err := Parse("+ 2 5f")
	if err == nil {
		t.Errorf("The code did not error out")
	}

	_, err = Parse("+ 2 f")
	if err == nil {
		t.Errorf("The code did not error out")
	}

}

func TestOperationOnlyError(t *testing.T) {
	_, err := Parse("/")
	if err == nil {
		t.Errorf("The code did not error out")
	}

}