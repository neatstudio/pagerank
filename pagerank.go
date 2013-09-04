package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

/**
//https://github.com/phurix/pagerank/blob/master/pagerank2.php
function GetPageRank($q,$host='toolbarqueries.google.com',$context=NULL) {
	$seed = "Mining PageRank is AGAINST GOOGLE'S TERMS OF SERVICE. Yes, I'm talking to you, scammer.";
	$result = 0x01020345;
	$len = strlen($q);
	for ($i=0; $i<$len; $i++) {
		$result ^= ord($seed{$i%strlen($seed)}) ^ ord($q{$i});
		$result = (($result >> 23) & 0x1ff) | $result << 9;
	}
    if (PHP_INT_MAX != 2147483647) { $result = -(~($result & 0xFFFFFFFF) + 1); }
	$ch=sprintf('8%x', $result);
	$url='http://%s/tbr?client=navclient-auto&ch=%s&features=Rank&q=info:%s';
	$url=sprintf($url,$host,$ch,$q);
	@$pr=file_get_contents($url,false,$context);
	return $pr?substr(strrchr($pr, ':'), 1):false;
}
*/

func main() {
	fmt.Println(getPageRank("http://neaststudio.com"))
}

func getPageRank(urlstring string) string {
	seed := []byte("Mining PageRank is AGAINST GOOGLE'S TERMS OF SERVICE. Yes, I'm talking to you, scammer.")
	seedlen := len(seed)
	result := 0x01020345
	q := []byte(urlstring)
	qlen := len(q)
	for i := 0; i < qlen; i++ {
		seedmod := i % seedlen
		seed_ascii, _ := utf8.DecodeLastRuneInString(string(seed[seedmod]))
		q_ascii, _ := utf8.DecodeLastRuneInString(string(q[i]))
		result ^= int(seed_ascii ^ q_ascii)
		result = ((result >> 23) & 0x1ff) | result<<9
	}
	result = -(^(result & 0xffffffff) + 1)
	ch := fmt.Sprintf("8%x", result)
	u, _ := url.Parse("http://toolbarqueries.google.com/tbr")
	uq := u.Query()
	uq.Set("client", "navclient-auto")
	uq.Set("ch", ch)
	uq.Set("features", "Rank")
	uq.Set("q", "info:"+string(q))
	u.RawQuery = uq.Encode()
	resp, err := http.Get(u.String())
	if nil != err {
		return err.Error()
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return err.Error()
	}
	//
	if len(body) <= 0 {
		return "error"
	}
	pagerankbody := strings.TrimSpace(string(body[:]))
	return pagerankbody[len(pagerankbody)-1:]
}
