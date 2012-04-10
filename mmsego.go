package mmsego

import(
    "unicode"
    "math"
    "darts"
    )

type Segmenter struct{
    dict darts.Darts
}
func max( a, b int ) int{
    if a < b {
	return b;
    }
    return a;
}
func min( a, b int ) int{
    if b < a {
	return b;
    }
    return a;
}
func average(in []darts.ResultPair) float64{
    numerator := 0
    denominator := 0
    for j := 0; j < len(in); j++{
	numerator += in[j].PrefixLen
	denominator++
    }
    return float64(numerator)/float64(denominator)
}
func variance(in []darts.ResultPair) float64{
    avg := average(in)
    cumulative := 0.
    denominator := 0.
    //in[j]0 means this item doesn't exist
    for j := 0; j < len(in); j++{
	v := float64(in[j].PrefixLen) - avg
	cumulative +=  v*v
	denominator++
    }
    return math.Sqrt(cumulative/denominator)
}
func morphemicFreedom(in []darts.ResultPair) (out float64) {
    for i := 0; i < len(in); i++ {
	if 1 == in[i].PrefixLen {
	    //add offset 3 to prevent negative log value
	    out += math.Log(float64(3+in[i].Freq))
	}
    }
    return out
}

//return value is the the chosen chunk
func filterChunksByRules(chunks [][]darts.ResultPair) ([]darts.ResultPair) {
    var candidates1,candidates2,candidates3,candidates4 [][]darts.ResultPair
    length := len(chunks)
    maxLength := 0
    for i :=0; i< length; i++{ //rule 1, Maximum matching
	var l int
	for j := 0 ; j < len(chunks[i]); j++{
	    l += chunks[i][j].PrefixLen
	}
	if l > maxLength {
	    maxLength = l
	    candidates1  = [][]darts.ResultPair{chunks[i]}
	}else if l == maxLength {
	    candidates1 = append(candidates1, chunks[i])
	}
    }
    if len(candidates1) == 1{
	return candidates1[0]
    }

    //else rule 2, Largest average word Rune length
    avgLen := 0.
    for i := 0; i < len(candidates1); i++{
	avg := average(candidates1[i])
	if avg > avgLen {
	    avgLen = avg
	    candidates2 = [][]darts.ResultPair{candidates1[i]}
	}else if avg == avgLen {
	    candidates2 = append(candidates2, candidates1[i])
	}
    }
    if len(candidates2) == 1{
	return candidates2[0]
    }

    //else rule 3, smallest variance
    smallestV := 65536. //large enough number
    for i := 0; i < len(candidates2); i++{
	v := variance(candidates2[i])
	if v < smallestV {
	    smallestV = v
	    candidates3 = [][]darts.ResultPair{candidates2[i]}
	}else if v == smallestV {
	    candidates3 = append(candidates3, candidates2[i])
	}
    }
    if len(candidates3) == 1{
	return candidates3[0]
    }

    //else rule 4, Largest sum of degree of morphemic freedom of one-character words
    smf := 0.
    for i := 0; i < len(candidates3); i++{
	v := morphemicFreedom(candidates3[i])
	if v > smf {
	    smf = v
	    candidates4 = [][]darts.ResultPair{candidates3[i]}
	}else if v == smf {
	    candidates4 = append(candidates4, candidates3[i])
	}
    }
    /*
    if len(candidates4) != 1{
	fmt.Println("exception!!", len(candidates4), candidates4)
	//exception 
    }
    */
    return candidates4[0]
}

type chunk struct{
    offSets []int
    values []int
}
func getChunks(inString []rune, d darts.Darts) (chunks [][]darts.ResultPair){
    results1 := d.CommonPrefixSearch(inString, 0)

    // no match or 1 match, 1 match assumes it's a 1 char match(or the dict is wrong)
    // can just return, according to the MMSEG algorithm
    if len(results1) == 0 {
	chunks = append(chunks, []darts.ResultPair{{PrefixLen:1, Value:darts.Value{Freq:1}}})
	return chunks
    }else if len(results1) == 1{
	chunks = append(chunks, results1)
	return chunks
    }
    //else
    for i := len(results1) - 1; i >= 0; i-- {
	word1 := results1[i].PrefixLen
	left1 := len(inString) - word1

	if left1 == 0{//meaning i == len(results) - 1, and there is only 1 word in this inString, done!
	    chunks = append(chunks, results1[i:i+1])
	    return chunks
	}
	//else
	results2 := d.CommonPrefixSearch(inString[word1:], 0)
	if len(results2) == 0{//no match, fake a result for convenience
	    results2 = []darts.ResultPair{{PrefixLen:1,Value:darts.Value{Freq:1}}}
	}
	for j := len(results2) - 1; j >= 0; j--{
	    word2 := word1 + results2[j].PrefixLen
	    left2 := len(inString) - word2

	    if left2 == 0{ //a 2 words chunk
		c := []darts.ResultPair{results1[i], results2[j]}
		chunks = append(chunks, c)
		continue
	    }
	    //else
	    results3 := d.CommonPrefixSearch(inString[word2:], 0)
	    if len(results3) == 0 { //fake a result for convenience
		results3 = []darts.ResultPair{{PrefixLen:1,Value:darts.Value{Freq:1}}}
	    }
	    for k := len(results3) - 1; k >= 0; k--{
		//word3 := word2 + results3[k].PrefixLen

		c := []darts.ResultPair{results1[i], results2[j], results3[k]}
		chunks = append(chunks, c)
	    }
	}
    }
    return chunks
}

func (s *Segmenter)Init(dictPath string){
    var err error
    s.dict, err = darts.Load(dictPath)
    if err != nil{
	panic(err)
    }
}
func (s *Segmenter)Mmseg(inString string, initOffset int, takeWord func(int, int), takeThesaurus func(int,int), lastPiece bool) (lstPos int){
    inRunes := []rune(inString)
    var pos = make([]int, len(inRunes)+1)
    pos[0] = 0
    for i, r := range inRunes{
	pos[i+1] = pos[i] + len(string(r))
    }
    offset := 0
    nextPunct := 0
    eol := false
    for ; offset < len(inRunes); {
	if offset == nextPunct && !eol{
	    offset++
	    //find the next non-Punct offset
	    for ; offset<len(inRunes) && !unicode.IsLetter(inRunes[offset]); {
		offset++
	    }
	    //only puncts left, return
	    if offset == len(inRunes){
		return pos[offset]
	    }
	    //find the next Punct after offset
	    for i, r := range inRunes[offset:]{
		if !unicode.IsLetter(r){
		    nextPunct = i + offset
		    break
		}
		if i+offset == len(inRunes) {
		    eol = true
		}
	    }
	}
	var chunks [][]darts.ResultPair
	if eol {
	    if lastPiece{
		chunks = getChunks(inRunes[offset:], s.dict);
	    }else{
		return pos[offset]
	    }
	}else{
	    chunks = getChunks(inRunes[offset:nextPunct], s.dict);
	}
	if 0 == len(chunks){
	    panic("error, length of chunks is zero")
	}
	chunk := filterChunksByRules(chunks)
	takeWord(initOffset+pos[offset], pos[offset+chunk[0].PrefixLen] - pos[offset])
	offset  += chunk[0].PrefixLen;
    }
    if offset > len(inRunes){
	panic("error MMseg")
    }
    return pos[offset]
}
