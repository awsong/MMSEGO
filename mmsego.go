package mmsego

import(
    "unicode"
//    "exp/utf8string"
    "fmt"
    "math"
    )

type Segmenter struct{
    dict Dict
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
func average(in []int) float64{
    numerator := 0
    denominator := 0
    for j := 0; j < len(in); j++{
	numerator += in[j]
	//in[j]==0 means this item doesn't exist
	if 0 != in[j] {
	    denominator++
	}
    }
    return float64(numerator)/float64(denominator)
}
func variance(in []int) float64{
    avg := average(in)
    cumulative := 0.
    denominator := 0.
    //in[j]0 means this item doesn't exist
    for j := 0; j < len(in) && 0 != in[j]; j++{
	v := float64(in[j]) - avg
	cumulative +=  v*v
	denominator++
    }
    return math.Sqrt(cumulative/denominator)
}
func morphemicFreedom(in *MatchItem) (out float64) {
    for i := 0; i < len(in.word); i++ {
	if 1 == in.wordRuneLen[i] {
	    if nil != in.wordItem[i]{
		//add offset 3 to prevent negative log value
		out += math.Log(float64(3+in.wordItem[i].Freq))
	    }else{
		//wordItem==nil means we did NOT find the char in dictionary
		//we assume the frequency as 1
		//add offset 3 to prevent negative log value
		out += math.Log(3+1)
	    }
	}
    }
    return out
}
type MatchItem struct{
    word [3]string
    wordRuneLen [3]int
    wordItem [3]*DictItem
}

//return value is the index of the chunk
func filterChunksByRules(chunks []MatchItem) (index int) {
    var retVecRule1 []int
    var retVecRule2 []int
    var retVecRule3 []int
    var retVecRule4 []int
    length := len(chunks)
    maxLength := 0
    for i :=0; i< length; i++{ //rule 1, Maximum matching
	l := chunks[i].wordRuneLen[0] + chunks[i].wordRuneLen[1] + chunks[i].wordRuneLen[2]
	if l > maxLength {
	    maxLength = l
	    retVecRule1 = []int{i}
	}else if l == maxLength {
	    retVecRule1 = append(retVecRule1, i)
	}
    }
    if len(retVecRule1) == 1{
	return retVecRule1[0]
    }else{
	//rule 2, Largest average word Rune length
	avgLen := 0.
	for i := 0; i < len(retVecRule1); i++{
	    avg := average(chunks[retVecRule1[i]].wordRuneLen[:])
	    if avg > avgLen {
		avgLen = avg
		retVecRule2 = []int{retVecRule1[i]}
	    }else if avg == avgLen {
		retVecRule2 = append(retVecRule2, retVecRule1[i])
	    }
	}
	if len(retVecRule2) == 1{
	    return retVecRule2[0]
	}else{
	    //rule 3, smallest variance
	    smallestV := 65536. //large enough number
	    for i := 0; i < len(retVecRule2); i++{
		v := variance(chunks[retVecRule2[i]].wordRuneLen[:])
		if v < smallestV {
		    smallestV = v
		    retVecRule3 = []int{retVecRule2[i]}
		}else if v == smallestV {
		    retVecRule3 = append(retVecRule3, retVecRule2[i])
		}
	    }
	    if len(retVecRule3) == 1{
		return retVecRule3[0]
	    }else{
		//rule 4, Largest sum of degree of morphemic freedom of one-character words
		smf := 0.
		for i := 0; i < len(retVecRule3); i++{
		    v := morphemicFreedom(&chunks[retVecRule3[i]])
		    if v > smf {
			smf = v
			retVecRule4 = []int{retVecRule3[i]}
		    }else if v == smf {
			retVecRule4 = append(retVecRule4, retVecRule3[i])
		    }
		}
		if len(retVecRule4) != 1{
		    fmt.Println("exception!!")
		    //exception 
		}
		return retVecRule4[0]
	    }
	}
    }
    fmt.Println("error")
    return 0
}

func getChunks(inString string, inOffset []int,dict Dict) (chunks []MatchItem){
    var wordRuneLen [3]int;
    var word [3]string;
    var wordItem [3]*DictItem;
    var chunkLen = len(inOffset)
    var present bool;
    for wordRuneLen[0] = min(4, chunkLen); wordRuneLen[0]>0; wordRuneLen[0]-- {
	if wordRuneLen[0] == len(inOffset) {
	    word[0] = inString[inOffset[0] : ]
	}else{
	    word[0] = inString[inOffset[0] : inOffset[wordRuneLen[0]]]
	}
	wordItem[0], present = dict[word[0]];
	if present || wordRuneLen[0] == 1 {
	    left := chunkLen-wordRuneLen[0];
	    if left == 0 {
		return []MatchItem{MatchItem{word,wordRuneLen,wordItem}}[:]
	    }
w1:	    for wordRuneLen[1] = min(4, left); wordRuneLen[1]>0; wordRuneLen[1]-- {
		if wordRuneLen[0] + wordRuneLen[1] == len(inOffset) {
		    word[1] = inString[inOffset[wordRuneLen[0]] : ]
		}else{
		    word[1] = inString[inOffset[wordRuneLen[0]] : inOffset[wordRuneLen[0]+ wordRuneLen[1]]]
		}
		wordItem[1], present = dict[word[1]];
		if present || wordRuneLen[1] == 1 {
		    left = chunkLen-wordRuneLen[0]-wordRuneLen[1]
		    if left == 0 {
			word[2] = ""
			wordRuneLen[2] = 0
			wordItem[2] = nil
			chunks = append(chunks, MatchItem{word,wordRuneLen,wordItem});
			break w1;
		    }
		    for wordRuneLen[2] = min(4, left); wordRuneLen[2]>0; wordRuneLen[2]-- {
			if wordRuneLen[0] + wordRuneLen[1] + wordRuneLen[2]== len(inOffset) {
			    word[1] = inString[inOffset[wordRuneLen[0]+wordRuneLen[1]] : ]
			}else{
			    word[1] = inString[inOffset[wordRuneLen[0]+wordRuneLen[1]] : inOffset[wordRuneLen[0]+wordRuneLen[1]+wordRuneLen[2]]]
			}
			wordItem[2], present = dict[word[2]];
			if present || wordRuneLen[2] == 1 {
			    chunks = append(chunks, MatchItem{word,wordRuneLen,wordItem});
			}
		    }
		}
	    }
	}
    }
    /*
    for _, v := range chunks{
	fmt.Printf("%v, %v\n", v.word, v.wordRuneLen);
    }*/
    return chunks;
}

func (s *Segmenter)Init(dictPath string){
    s.dict = LoadDictionary(dictPath)
}
func (s *Segmenter)Mmseg(inString string, out chan [2]int) bool{
    var pos = make([]int,len(inString))
    runeLen := 0
    for i,r := range inString{
	if unicode.IsPunct(r){
	    pos[runeLen] = -1
	}else{
	    pos[runeLen] = i
	}
	runeLen++
    }
    offset := 0
    nextPunct := 0
    for i, v := range pos{
	if v == -1 {
	    nextPunct = i
	    break
	}
    }
    eol := false
f0: for ; offset < runeLen; {
	if offset == nextPunct && !eol{
	    offset++
	    //find the next none Punct offset
	    for ;pos[offset] == -1 && offset<runeLen;{
		offset++
	    }
	    if offset == runeLen{
		break f0
	    }
	    //find the next Punct after offset
	    for i, v := range pos[offset:]{
		if v == -1 {
		    nextPunct = i+offset
		    break
		}
		if i+offset == runeLen {
		    eol = true
		}
	    }
	}
	var chunks []MatchItem
	if eol {
	    chunks = getChunks(inString[:], pos[offset:runeLen], s.dict);
	}else{
	    chunks = getChunks(inString[:], pos[offset:nextPunct], s.dict);
	}
	if 0 == len(chunks){
	    fmt.Println("chunks is 0",offset, nextPunct, runeLen,inString[:])
	}
	index := filterChunksByRules(chunks);
	if offset < 0 ||offset > len(pos) || index > len(chunks){
	    fmt.Println("oops:",offset, pos[offset], index)
	}
	//fmt.Printf("%v, %v\n", offset, chunks[index].wordRuneLen[0]);
	out <- [2]int{pos[offset], pos[offset]+len(chunks[index].word[0])}
	offset  += chunks[index].wordRuneLen[0];
    }
    close(out);
    return true;
}
