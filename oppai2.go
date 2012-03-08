package main

import (
	"bufio";
	"crypto/md5";
	"encoding/hex";
	"fmt";
	"http";
	"io/ioutil";
	"json";
	"os";
	"regexp";
	"runtime";
	"strconv";
	"url";
	"utf8";
)

const (
	appid = "CAB81CE7FCEFE5385E59E6E7680D707D0E3B7C48"	// Bing API App ID
	query_uri = "http://api.bing.net/json.aspx"
	dir = "./data2/"
)


type thumbnail struct{
	Url string;
	ContentType string;
	Width int;
	Height int;
	FileSize int;
}

type results struct{
	Title string;
	MediaUrl string;
	Url string;
	Width int;
	Height int;
	FileSize int;
	ContentType int;
	Thumbnail thumbnail;
}

type image struct{
	Total int;
	Offset int;
	Results []results;
}

type query struct {
	SearchTerms string;
}

type search_response struct {
	Version string;
	Query query;
	Image image;
}

type json_root struct{
	SearchResponse search_response;
}

func encode_utf8(str string) string {
	return utf8.NewString(str).String();
}


func md5hex(str string) string {
	h := md5.New();
	h.Write([]byte(str));
	return hex.EncodeToString([]byte(h.Sum()));;
}


func get_request_uri(param map[string]string) string {
	post_body := "";
	cnt := 0;
	
	for k, v := range param {
		post_body += k + "=" + url.QueryEscape(v);
		if (cnt < len(param) - 1) {
			post_body += "&";
		}
		cnt++;
	}
	return string(query_uri + "?" + post_body);
}


func main() {
	page_count := 1;
	download_count := 0;
	
	for {
		offset := page_count * 50;
		param := map[string]string {
				"AppId": appid,
				"Version": "2.2",
				"Market": "ja-JP",
				"Sources": "Image",
				"Image.Count": strconv.Itoa(50),
				"Image.Offset": strconv.Itoa(offset),
				"Adult": "off",
				"Query": "おっぱい",
			};
		
		var sr *json_root;
		res, err := http.Get(get_request_uri(param));
		if err != nil { break; };
		reader := bufio.NewReader(res.Body);
		line, err := reader.ReadBytes('\n');
		if err == nil { break; }
		json.Unmarshal(line, &sr);
		
		q := make(chan int);
		for i := 0; i < len(sr.SearchResponse.Image.Results); i++ {
			result := sr.SearchResponse.Image.Results[i];
			if regexp.MustCompile(".jpg$").FindString(result.MediaUrl) == "" {
				continue;
			}
			download_count++;
			
			filename := md5hex(encode_utf8(result.MediaUrl)) + ".jpg";
			filepath := dir + filename;
			
			if _, err := os.Stat(filepath); err == nil { continue; }
			fmt.Printf("%d : Download... %s\n", download_count, result.MediaUrl);
			
			go func() {
				q <- 1;
				res, err := http.Get(result.MediaUrl);
				if err != nil { runtime.Goexit(); };
				data, err := ioutil.ReadAll(res.Body);
				if err != nil { runtime.Goexit(); };
				if regexp.MustCompile("^image").FindString(http.DetectContentType(data)) != "" {
					ioutil.WriteFile(filepath, data, 0666);
				}
			}();
			<- q;
		}
		
		close(q);
		page_count++;
	}
}

