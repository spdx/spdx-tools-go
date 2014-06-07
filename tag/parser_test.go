package tag

import "strings"
import "testing"

func sameDoc(a, b []pair) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestSameDocFunc(t *testing.T) {
	a := []pair{
		{"one", "onev"},
		{"two", "twov"},
	}
	b := []pair{
		{"one", "onev"},
		{"two", "twov"},
	}

	if !sameDoc(a, b) {
		t.Error("Slices are the same and not detected")
	}

	if sameDoc(a, []pair{{"a", "c"}, {"two", "twov"}}) {
		t.Error("Slices are the same and detected as different")
	}
}

func TestEmptyDoc(t *testing.T) {
	r := strings.NewReader("")

	doc, err := parse(r)

	if len(doc) != 0 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}
}

func TestCommentsOnly(t *testing.T) {
	r := strings.NewReader("#this is a comment\n#this is another comment :)\n#whatever")

	doc, err := parse(r)

	if len(doc) != 0 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

}

func TestOneCommentNoNewline(t *testing.T) {
	r := strings.NewReader("#this is a comment")

	doc, err := parse(r)

	if len(doc) != 0 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

}

func TestOneCommentNewLine(t *testing.T) {
	r := strings.NewReader("#this is a comment\n")

	doc, err := parse(r)

	if len(doc) != 0 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

}

func TestDoubleEndlineAfterComment(t *testing.T) {
	r := strings.NewReader("#property:value\n\n")

	doc, err := parse(r)

	if len(doc) != 0 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}
}

func TestCommentsAndEmptyLines(t *testing.T) {
	r := strings.NewReader("#this is a comment\n\n#this is another comment :)\n#whatever\n\n\n#anoterOne")

	doc, err := parse(r)

	if len(doc) != 0 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}
}

func TestValidProperty(t *testing.T) {
	r := strings.NewReader("someKey:someValue")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "someValue"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestValidPropertyNewLine(t *testing.T) {
	r := strings.NewReader("someKey:someValue\n")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "someValue"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestClearWhitespaces(t *testing.T) {
	r := strings.NewReader("someKey  : someValue\n")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "someValue"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestValidPropertyNewLineCR(t *testing.T) {
	r := strings.NewReader("someKey:someValue\r\n")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "someValue"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestValidPropertyText(t *testing.T) {
	r := strings.NewReader("someKey:<text>someValue</text>")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "someValue"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestValidPropertyTextMultiline(t *testing.T) {
	r := strings.NewReader("someKey:<text>\nsomeValue\n123\n\n4\n</text>")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "someValue\n123\n\n4"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestValidPropertyTextNewLine(t *testing.T) {
	r := strings.NewReader("someKey:<text>someValue</text>\n")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "someValue"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestMoreInlineProperties(t *testing.T) {
	r := strings.NewReader("Property1:value1\nProperty2:value2\nProperty3:value3\n")

	doc, err := parse(r)

	if len(doc) != 3 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

	properties := []pair{
		{"Property1", "value1"},
		{"Property2", "value2"},
		{"Property3", "value3"},
	}

	if !sameDoc(properties, doc) {
		t.Errorf("Expected %s. Got %s", properties, doc)
	}
}

func TestInlinePropertiesAndComments(t *testing.T) {
	r := strings.NewReader("# comment\nProperty1:value1\nProperty2:value2\n# comment no two\nProperty3:value3\n#comm\n")

	doc, err := parse(r)

	if len(doc) != 3 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

	properties := []pair{
		{"Property1", "value1"},
		{"Property2", "value2"},
		{"Property3", "value3"},
	}

	if !sameDoc(properties, doc) {
		t.Errorf("Expected %s. Got %s", properties, doc)
	}
}

func TestInlinePropertiesCommentsAndNewlines(t *testing.T) {
	r := strings.NewReader("# comment\n\nProperty1:value1\n\n\nProperty2:value2\n# comment no two\nProperty3:value3\n#comm\n")

	doc, err := parse(r)

	if len(doc) != 3 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

	properties := []pair{
		{"Property1", "value1"},
		{"Property2", "value2"},
		{"Property3", "value3"},
	}

	if !sameDoc(properties, doc) {
		t.Errorf("Expected %s. Got %s", properties, doc)
	}
}

func TestMoreTextProperties(t *testing.T) {
	r := strings.NewReader("Property1:<text>value1</text>\nProperty2:<text>value2</text>\nProperty3:<text>value3</text>\n")

	doc, err := parse(r)

	if len(doc) != 3 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

	properties := []pair{
		{"Property1", "value1"},
		{"Property2", "value2"},
		{"Property3", "value3"},
	}

	if !sameDoc(properties, doc) {
		t.Errorf("Expected %s. Got %s", properties, doc)
	}
}

func TestMoreTextPropertiesAndComments(t *testing.T) {
	r := strings.NewReader("# this is a comment\nProperty1:<text>value1</text>\n#so is this\nProperty2:<text>value2</text>\nProperty3:<text>value3</text>\n#and this")

	doc, err := parse(r)

	if len(doc) != 3 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

	properties := []pair{
		{"Property1", "value1"},
		{"Property2", "value2"},
		{"Property3", "value3"},
	}

	if !sameDoc(properties, doc) {
		t.Errorf("Expected %s. Got %s", properties, doc)
	}
}

func TestMoreTextPropertiesCommentsAndNewlines(t *testing.T) {
	r := strings.NewReader("\n\n# this is a comment\n\nProperty1:<text>value1</text>\n#so is this\n\nProperty2:<text>value2</text>\nProperty3:<text>value3</text>\n#and this\n\n")

	doc, err := parse(r)

	if len(doc) != 3 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

	properties := []pair{
		{"Property1", "value1"},
		{"Property2", "value2"},
		{"Property3", "value3"},
	}

	if !sameDoc(properties, doc) {
		t.Errorf("Expected %s. Got %s", properties, doc)
	}
}

func TestMixedProperties(t *testing.T) {
	r := strings.NewReader("Property1:  <text>value1</text>\nProperty2:value2\nProperty3:<text>value3</text>\n")

	doc, err := parse(r)

	if len(doc) != 3 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
	}

	properties := []pair{
		{"Property1", "value1"},
		{"Property2", "value2"},
		{"Property3", "value3"},
	}

	if !sameDoc(properties, doc) {
		t.Errorf("Expected %s. Got %s", properties, doc)
	}
}

func TestInvalidTextValuePrefix(t *testing.T) {
	r := strings.NewReader("Property1: invalid <text>value1</text>\n")

	_, err := parse(r)

	if err == nil {
		t.Fail()
	}
	if err != ErrInvalidPrefix {
		t.Errorf("Another error: %s", err)
	}
}

func TestInvalidTextValueSuffix(t *testing.T) {
	r := strings.NewReader("Property1: <text>value1</text> invalid \n")

	_, err := parse(r)

	if err == nil {
		t.Fail()
	}
	if err != ErrInvalidSuffix {
		t.Errorf("Another error: %s", err)
	}
}

func TestInvalidTextValueSuffixComment(t *testing.T) {
	r := strings.NewReader("Property1: <text>value1</text># invalid \n")

	_, err := parse(r)

	if err == nil {
		t.Fail()
	}
	if err != ErrInvalidSuffix {
		t.Errorf("Another error: %s", err)
	}
}

func TestInvalidTextValueSuffixProperty(t *testing.T) {
	r := strings.NewReader("Property1: <text>value1</text>a:b\n")

	_, err := parse(r)
	t.Logf("Error: %s\n", err)

	if err == nil {
		t.Fail()
	}
	if err != ErrInvalidSuffix {
		t.Errorf("Another error: %s", err)
	}
}

func TestInvalidUnclosedText(t *testing.T) {
	r := strings.NewReader("Property1: <text>value1\n\n invalid \n")

	_, err := parse(r)
	t.Logf("Error: %s\n", err)

	if err == nil {
		t.Fail()
	}
	if err != ErrNoCloseTag {
		t.Errorf("Another error: %s", err)
	}
}

func TestCommentAsInlineValue(t *testing.T) {
	r := strings.NewReader("someKey:#someValue")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "#someValue"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestCommentAsTextValue(t *testing.T) {
	r := strings.NewReader("someKey:<text>#someValue</text>")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "#someValue"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestCommentAsMultilineTextValue(t *testing.T) {
	r := strings.NewReader("someKey:<text>#c\n#someValue\nd\n#a</text>")

	doc, err := parse(r)

	if len(doc) != 1 || err != nil {
		t.Errorf("Document: %s. Error: %s", doc, err)
		t.FailNow()
	}

	p := pair{"someKey", "#c\n#someValue\nd\n#a"}

	if len(doc) == 1 && doc[0] != p {
		t.Errorf("Expected %s. Got %s", p, doc[0])
	}
}

func TestSomeInvalidText(t *testing.T) {
	r := strings.NewReader("garbage")

	_, err := parse(r)
	if err == nil {
		t.Fail()
	}
	if err != ErrInvalidText {
		t.Errorf("Another error: %s", err)
	}

}

func TestPropertyWithNoValue(t *testing.T) {
	r := strings.NewReader("garbage:")
	doc, err := parse(r)
	t.Logf("Doc=%s, Err=%s", doc, err)
	p := pair{"garbage", ""}
	if err != nil || doc == nil || len(doc) != 1 {
		t.FailNow()
	}
	if doc[0] != p {
		t.Fail()
	}
}

func TestAllDataWhitespaceAtEOF(t *testing.T) {
	f := tokenize()
	data := []byte("  \n")
	advance, token, err := f(data, true)
	if advance != 0 || token != nil || err != nil {
		t.Errorf("Fail with: advance=%d, data=%s, err=%s\n", advance, token, err)
	}
}

func TestCommentEndingInNewlineAtEOF(t *testing.T) {
	f := tokenize()
	data := []byte("#comment\n")
	advance, token, err := f(data, true)
	if advance != 0 || token != nil || err != nil {
		t.Errorf("Fail with: advance=%d, data=%s, err=%s\n", advance, token, err)
	}
}