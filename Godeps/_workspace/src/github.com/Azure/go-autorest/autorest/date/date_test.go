package date

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func ExampleParseDate() {
	d, _ := ParseDate("2001-02-03")
	fmt.Println(d)
	// Output: 2001-02-03
}

func ExampleDate() {
	d, _ := ParseDate("2001-02-03")
	t, _ := time.Parse(time.RFC3339, "2001-02-04T00:00:00Z")

	// Date acts as time.Time when the receiver
	if d.Before(t) {
		fmt.Printf("Before ")
	} else {
		fmt.Printf("After ")
	}

	// Convert Date when needing a time.Time
	if t.After(d.ToTime()) {
		fmt.Printf("After")
	} else {
		fmt.Printf("Before")
	}
	// Output: Before After
}

func ExampleDate_MarshalBinary() {
	d, _ := ParseDate("2001-02-03")
	t, _ := d.MarshalBinary()
	fmt.Println(string(t))
	// Output: 2001-02-03
}

func ExampleDate_UnmarshalBinary() {
	d := Date{}
	t := "2001-02-03"

	err := d.UnmarshalBinary([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: 2001-02-03
}

func ExampleDate_MarshalJSON() {
	d, _ := ParseDate("2001-02-03")
	j, _ := json.Marshal(d)
	fmt.Println(string(j))
	// Output: "2001-02-03"
}

func ExampleDate_UnmarshalJSON() {
	var d struct {
		Date Date `json:"date"`
	}
	j := `{
    "date" : "2001-02-03"
  }`

	err := json.Unmarshal([]byte(j), &d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d.Date)
	// Output: 2001-02-03
}

func ExampleDate_MarshalText() {
	d, _ := ParseDate("2001-02-03")
	t, _ := d.MarshalText()
	fmt.Println(string(t))
	// Output: 2001-02-03
}

func ExampleDate_UnmarshalText() {
	d := Date{}
	t := "2001-02-03"

	err := d.UnmarshalText([]byte(t))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)
	// Output: 2001-02-03
}

func TestDateString(t *testing.T) {
	d, _ := ParseDate("2001-02-03")
	if d.String() != "2001-02-03" {
		t.Errorf("date: String failed (%v)", d.String())
	}
}

func TestDateBinaryRoundTrip(t *testing.T) {
	d1, err := ParseDate("2001-02-03")
	t1, err := d1.MarshalBinary()
	if err != nil {
		t.Errorf("date: MarshalBinary failed (%v)", err)
	}

	d2 := Date{}
	err = d2.UnmarshalBinary(t1)
	if err != nil {
		t.Errorf("date: UnmarshalBinary failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip Binary failed (%v, %v)", d1, d2)
	}
}

func TestDateJSONRoundTrip(t *testing.T) {
	type s struct {
		Date Date `json:"date"`
	}
	var err error
	d1 := s{}
	d1.Date, err = ParseDate("2001-02-03")
	j, err := json.Marshal(d1)
	if err != nil {
		t.Errorf("date: MarshalJSON failed (%v)", err)
	}

	d2 := s{}
	err = json.Unmarshal(j, &d2)
	if err != nil {
		t.Errorf("date: UnmarshalJSON failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip JSON failed (%v, %v)", d1, d2)
	}
}

func TestDateTextRoundTrip(t *testing.T) {
	d1, err := ParseDate("2001-02-03")
	t1, err := d1.MarshalText()
	if err != nil {
		t.Errorf("date: MarshalText failed (%v)", err)
	}

	d2 := Date{}
	err = d2.UnmarshalText(t1)
	if err != nil {
		t.Errorf("date: UnmarshalText failed (%v)", err)
	}

	if !reflect.DeepEqual(d1, d2) {
		t.Errorf("date: Round-trip Text failed (%v, %v)", d1, d2)
	}
}

func TestDateToTime(t *testing.T) {
	var d Date
	d, err := ParseDate("2001-02-03")
	if err != nil {
		t.Errorf("date: ParseDate failed (%v)", err)
	}
	var v interface{} = d.ToTime()
	switch v.(type) {
	case time.Time:
		return
	default:
		t.Errorf("date: ToTime failed to return a time.Time")
	}
}

func TestDateUnmarshalJSONReturnsError(t *testing.T) {
	var d struct {
		Date Date `json:"date"`
	}
	j := `{
    "date" : "February 3, 2001"
  }`

	err := json.Unmarshal([]byte(j), &d)
	if err == nil {
		t.Error("date: Date failed to return error for malformed JSON date")
	}
}

func TestDateUnmarshalTextReturnsError(t *testing.T) {
	d := Date{}
	txt := "February 3, 2001"

	err := d.UnmarshalText([]byte(txt))
	if err == nil {
		t.Error("date: Date failed to return error for malformed Text date")
	}
}
