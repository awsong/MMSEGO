package mmsego

import(
    "testing"
    "unicode/utf8"
    )

type minMaxTest struct{
    in [2]int
    out int
}

var maxTests = []minMaxTest{
    {[2]int{10,20}, 20},
    {[2]int{30,20}, 30},
    {[2]int{3,2}, 3},
    {[2]int{-3,2}, 2},
    {[2]int{3,-2}, 3},
    {[2]int{-3,-2}, -2},
    {[2]int{0,-2}, 0},
    {[2]int{-3,0}, 0},
    {[2]int{0,2}, 2},
    {[2]int{3,0}, 3},
}

var minTests = []minMaxTest{
    {[2]int{10,20}, 10},
    {[2]int{30,20}, 20},
    {[2]int{3,2}, 2},
    {[2]int{-3,2}, -3},
    {[2]int{3,-2}, -2},
    {[2]int{-3,-2}, -3},
    {[2]int{0,-2}, -2},
    {[2]int{-3,0}, -3},
    {[2]int{0,2}, 0},
    {[2]int{3,0}, 0},
}

func TestMax(t *testing.T) {
    for _, dt := range maxTests {
	v := max(dt.in[0], dt.in[1]);
	if v != dt.out {
	    t.Errorf("In = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}
    }
}
func TestMin(t *testing.T) {
    for _, dt := range minTests {
	v := min(dt.in[0], dt.in[1]);
	if v != dt.out {
	    t.Errorf("In = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}
    }
}

type getChunksTest struct{
    in string
    out []MatchItem
}

var getChunksTests = []getChunksTest{
    {"研究生命起源",
	[]MatchItem{MatchItem{[3]string{"中国人民解放军","第二","炮兵部队"}, [3]int{0,0,0}, [3]*DictItem{nil,nil,nil}}, MatchItem{}} },
    {"中国人民解放军第二炮兵部队",
	[]MatchItem{MatchItem{[3]string{"中国人民解放军","第二","炮兵部队"}, [3]int{0,0,0}, [3]*DictItem{nil,nil,nil}}, MatchItem{}} },
    {"中国人民",
	[]MatchItem{MatchItem{[3]string{"中国人民解放军","第二","炮兵部队"}, [3]int{0,0,0}, [3]*DictItem{nil,nil,nil}}, MatchItem{}} },
    {"中国人",
	[]MatchItem{MatchItem{[3]string{"中国人民解放军","第二","炮兵部队"}, [3]int{0,0,0}, [3]*DictItem{nil,nil,nil}}, MatchItem{}} },
}

func testGetChunks(t *testing.T) {
    dict := LoadDictionary("/public/development/go_projects/src/mmsego/uni.lib");
    for _, dt := range getChunksTests {
	var pos = make([]int,len(dt.in))
        pos[0] = 0
        runeLen := 0
        for i,s := range dt.in{
            pos[i+1] = pos[i] + utf8.RuneLen(s)
            runeLen++
        }
	v := getChunks(dt.in, pos[:runeLen], dict);
	t.Errorf("dt.in:%v, dict:%v\n", dt.in, dict["人民政府"].TT[0]);
	return;
	t.Errorf("number of chunks: %v\n", len(v));
	for i:=0; i<len(v); i++ {
	    if v[0].word != dt.out[0].word {
		    t.Errorf("Key = %v, Real = %v, want %v.", dt.in, v, dt.out)
	    }
	}
    }
}

type filterChunksByRulesTest struct{
    in []MatchItem
    out int
}

var filterChunksByRulesTests = []filterChunksByRulesTest{
    {[]MatchItem{MatchItem{[3]string{"中国人","第二","炮兵部队"}, [3]int{3,2,4}, [3]*DictItem{nil,nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,2,2}, [3]*DictItem{nil,nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,2,2}, [3]*DictItem{nil,nil,nil}}},
    0},
    {[]MatchItem{MatchItem{[3]string{"中国人","第二","炮兵部队"}, [3]int{3,2,4}, [3]*DictItem{nil,nil,nil}}},
    0},
    {[]MatchItem{MatchItem{[3]string{"中国人","第二","炮兵部队"}, [3]int{3,2,4}, [3]*DictItem{nil,nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,5,4}, [3]*DictItem{nil,nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,4,4}, [3]*DictItem{nil,nil,nil}}},
    1},
    {[]MatchItem{MatchItem{[3]string{"中国人","第二","炮兵部队"}, [3]int{3,2,1}, [3]*DictItem{nil,nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,1,4}, [3]*DictItem{nil,nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,2,3}, [3]*DictItem{nil,nil,nil}}},
    2},
    {[]MatchItem{MatchItem{[3]string{"中国人","第二","炮兵部队"}, [3]int{3,5,0}, [3]*DictItem{nil,nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,1,4}, [3]*DictItem{nil,nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,2,3}, [3]*DictItem{nil,nil,nil}}},
    0},
    {[]MatchItem{MatchItem{[3]string{"中国人","第二","炮兵部队"}, [3]int{3,3,1}, [3]*DictItem{nil,nil,&DictItem{Freq:38}}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{1,3,3}, [3]*DictItem{&DictItem{Freq:28},nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,1,3}, [3]*DictItem{nil,&DictItem{Freq:18},nil}}},
    0},
    {[]MatchItem{MatchItem{[3]string{"中国人","第二","炮兵部队"}, [3]int{3,3,1}, [3]*DictItem{nil,nil,&DictItem{Freq:38}}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{1,3,3}, [3]*DictItem{&DictItem{Freq:98},nil,nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{3,1,3}, [3]*DictItem{nil,&DictItem{Freq:18},nil}}},
    1},
    {[]MatchItem{MatchItem{[3]string{"中国人","第二","炮兵部队"}, [3]int{5,1,1}, [3]*DictItem{nil,&DictItem{Freq:48},&DictItem{Freq:38}}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{1,1,5}, [3]*DictItem{&DictItem{Freq:98},&DictItem{Freq:18},nil}},
		MatchItem{[3]string{"解放军","第二","部队"}, [3]int{1,5,1}, [3]*DictItem{&DictItem{Freq:1},nil,&DictItem{Freq:108}}}},
    1},
}

func TestFilterChunksByRules(t *testing.T) {
    for _, dt := range filterChunksByRulesTests{
	v := filterChunksByRules(dt.in)
	if v != dt.out {
	    t.Errorf("In = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}
    }
}

type averageTest struct{
    in []int
    out float64
}

var averageTests = []averageTest{
    {[]int{3,4,5}, 4.},
    {[]int{22,33,10}, 21.666666666666668},
    {[]int{30,40,50}, 40.},
}

func TestAverage(t *testing.T){
    for _, dt := range averageTests{
	v := average(dt.in)
	if v != dt.out {
	    t.Errorf("In = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}
    }
}
type varianceTest struct{
    in []int
    out float64
}

var varianceTests = []varianceTest{
    {[]int{3,4,5}, 0.816496580927726},
    {[]int{3,4,6}, 1.247219128924647},
    {[]int{22,33,10}, 9.392668535736915},
}

func TestVariance(t *testing.T){
    for _, dt := range varianceTests{
	v := variance(dt.in)
	if v != dt.out {
	    t.Errorf("In = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}
    }
}

type mmsegTest struct{
    in string
    out int
}

var mmsegTests = []mmsegTest{
    {"中国人民解放军是，\"世界上最强大的军队，\"每个人都说共产党好", 0},
    {"中国人民解放军是世界上最强大的军队，每个人都说共产党好", 0},
    {"我爱北京天安门，天安门上太阳升，伟大领袖毛主席，领导人民闹革命！", 0},
}

func TestMmseg(t *testing.T){
    s := new(Segmenter)
    outChan := make(chan [2]int, 50);

    s.init("/public/development/go_projects/src/mmsego/uni.lib")
    for _, dt := range mmsegTests{
	go s.Mmseg(dt.in[:], outChan)
	for m := range outChan {
	    t.Logf("%v ",m);
	}
	/*
	v := mmseg(dt.in[:], outChan)
	if v != dt.out {
	    t.Errorf("In = %v, Real = %v, want %v.", dt.in, v, dt.out)
	}*/
    }
    t.Errorf("fail\n");
}
