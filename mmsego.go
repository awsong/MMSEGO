package mmsego

import(
//    "strings"
    "exp/utf8string"
    "fmt"
    "math"
    )

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
    for j := 0; j < len(in) && 0 != in[j]; j++{
	v := float64(in[j]) - avg
	cumulative +=  v*v
	denominator++
    }
    return math.Sqrt(cumulative/denominator)
}
func morphemicFreedom(in *MatchItem) (out float64) {
    for i := 0; i < len(in.word); i++ {
	if 1 == in.wordLen[i] {
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
    wordLen [3]int
    wordItem [3]*DictItem
}

//return value is the byte length of the first word
func filterChunksByRules(chunks []MatchItem) (index int) {
    var retVecRule1 []int
    var retVecRule2 []int
    var retVecRule3 []int
    var retVecRule4 []int
    length := len(chunks)
    maxLength := 0
    for i :=0; i< length; i++{ //rule 1, Maximum matching
	l := chunks[i].wordLen[0] + chunks[i].wordLen[1] + chunks[i].wordLen[2]
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
	//rule 2, Largest average word length
	avgLen := 0.
	for i := 0; i < len(retVecRule1); i++{
	    avg := average(chunks[retVecRule1[i]].wordLen[:])
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
		v := variance(chunks[retVecRule2[i]].wordLen[:])
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

func getChunks(inString string, dict Dict) (chunks []MatchItem){
    utf8InString := utf8string.NewString(inString);
    var wordLen [3]int;
    var word [3]string;
    var wordItem [3]*DictItem;
    var chunkLen = utf8InString.RuneCount();
    var present bool;
    for wordLen[0] = min(8, chunkLen); wordLen[0]>0; wordLen[0]-- {
	word[0] = utf8InString.Slice(0, wordLen[0])
	wordItem[0], present = dict[word[0]];
	if present || wordLen[0] == 1 {
	    left := chunkLen-wordLen[0];
	    if left == 0 {
		//left==0 means this is THE best match, accroding to all 3 rules
		//we should return from here
		/*word[1] = ""
		word[2] = ""
		wordLen[1] = 0
		wordLen[2] = 0
		wordItem[1] = nil
		wordItem[2] = nil
		chunks = append(chunks, MatchItem{word,wordLen,wordItem});*/
		return []MatchItem{MatchItem{word,wordLen,wordItem}}[:]
	    }
w1:	    for wordLen[1] = min(8, left); wordLen[1]>0; wordLen[1]-- {
		word[1] = utf8InString.Slice(wordLen[0], wordLen[0]+wordLen[1])
		wordItem[1], present = dict[word[1]];
		if present || wordLen[1] == 1 {
		    left = chunkLen-wordLen[0]-wordLen[1]
		    if left == 0 {
			word[2] = ""
			wordLen[2] = 0
			wordItem[2] = nil
			chunks = append(chunks, MatchItem{word,wordLen,wordItem});
			break w1;
		    }
		    for wordLen[2] = min(8, left); wordLen[2]>0; wordLen[2]-- {
			word[2] = utf8InString.Slice(wordLen[0]+wordLen[1], wordLen[0]+wordLen[1]+wordLen[2])
			wordItem[2], present = dict[word[2]];
			if present || wordLen[2] == 1 {
			    chunks = append(chunks, MatchItem{word,wordLen,wordItem});
			}
		    }
		}
	    }
	}
    }
    for _, v := range chunks{
	fmt.Printf("%v, %v\n", v.word, v.wordLen);
    }
    return chunks;
}
