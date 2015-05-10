package goexample

import "sort"
import "strings"

//returns a copy of str with all occurences of char removed
func StripChars(str, char string) string{
    return strings.Map(func(r rune) rune {
        if strings.IndexRune(char, r) < 0 {
            return r
        }
        return -1
    }, str)
}

type WordCounter struct{
    Counts map[string]int
    Words []string
}

func (wc *WordCounter) Len() int {
    return len(wc.Counts)
}

func (wc *WordCounter) Less(i, j int) bool {
    return wc.Counts[wc.Words[i]] > wc.Counts[wc.Words[j]]
}

func (wc *WordCounter) Swap(i, j int) {
    wc.Words[i], wc.Words[j] = wc.Words[j], wc.Words[i]
}

func SortWords(c map[string]int) []string {
    wc := new(WordCounter)
    wc.Counts = c
    wc.Words = make([]string, len(c))
    i := 0
    for word, _ := range c {
        wc.Words[i] = word
        i++
    }
    sort.Sort(wc)
    return wc.Words
}