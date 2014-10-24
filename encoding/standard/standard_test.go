package standard

import "testing"

var encodingTests = []struct {
	in, out string
}{
	{"opentext", "opentext"},
	{"\nblue\t63", "\nblue\t63"},
}

func TestEncoding(t *testing.T) {
	for _, tt := range encodingTests {
		encoded, err := Encode([]byte(tt.in))
		if err != nil {
			t.Errorf(err.Error())
		}
		decoded, err := Decode(encoded)
		if err != nil {
			t.Errorf(err.Error())
		}
		if tt.out != string(decoded) {
			t.Errorf("want %s, got %s", tt.out, decoded)
		}
	}
}
