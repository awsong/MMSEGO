package mmsego

import(
    "testing"
    )

type dictTest struct{
    in string
    out DictItem
}

var dictTests = []dictTest{
    dictTest{
	"信仰主义",
	DictItem{1, 2, [8]Thesaurus{{0, "信仰"}, {6, "主义"}}}},
    dictTest{
	"令人心悸",
	DictItem{1, 3, [8]Thesaurus{{6, "心悸"}, {0, "令人"}, {3, "人心"}}}},
    dictTest{
	"河",
	DictItem{187, 0, [8]Thesaurus{}}},
    dictTest{
	"哲",
	DictItem{11, 0, [8]Thesaurus{}}},
    dictTest{
	"电信业者",
	DictItem{1, 3, [8]Thesaurus{{6, "业者"}, {0, "电信业"}, {0, "电信"}}}},
}

func TestLoadDictionary(t *testing.T) {
    dict := LoadDictionary("/public/development/go_projects/src/mmsego/uni.lib");
    for _, dt := range dictTests {
	v := dict[dt.in]
	if *v != dt.out {
		t.Errorf("Key = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}
    }
}
