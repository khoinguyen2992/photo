package text

import (
	"regexp"
	"strings"
)

var pattern = map[string]*regexp.Regexp{
	"a": regexp.MustCompile(`à|á|ạ|ả|ã|â|ầ|ấ|ậ|ẩ|ẫ|ă|ằ|ắ|ặ|ẳ|ẵ`),
	"e": regexp.MustCompile(`è|é|ẹ|ẻ|ẽ|ê|ề|ế|ệ|ể|ễ`),
	"i": regexp.MustCompile(`ì|í|ị|ỉ|ĩ`),
	"o": regexp.MustCompile(`ò|ó|ọ|ỏ|õ|ô|ồ|ố|ộ|ổ|ỗ|ơ|ờ|ớ|ợ|ở|ỡ`),
	"u": regexp.MustCompile(`ù|ú|ụ|ủ|ũ|ư|ừ|ứ|ự|ử|ữ`),
	"y": regexp.MustCompile(`ỳ|ý|ỵ|ỷ|ỹ`),
	"d": regexp.MustCompile(`đ`),
	"A": regexp.MustCompile(`À|Á|Ạ|Ả|Ã|Â|Ầ|Ấ|Ậ|Ẩ|Ẫ|Ă|Ằ|Ắ|Ặ|Ẳ|Ẵ`),
	"E": regexp.MustCompile(`È|É|Ẹ|Ẻ|Ẽ|Ê|Ề|Ế|Ệ|Ể|Ễ`),
	"I": regexp.MustCompile(`Ì|Í|Ị|Ỉ|Ĩ`),
	"O": regexp.MustCompile(`Ò|Ó|Ọ|Ỏ|Õ|Ô|Ồ|Ố|Ộ|Ổ|Ỗ|Ơ|Ờ|Ớ|Ợ|Ở|Ỡ`),
	"U": regexp.MustCompile(`Ù|Ú|Ụ|Ủ|Ũ|Ư|Ừ|Ứ|Ự|Ử|Ữ`),
	"Y": regexp.MustCompile(`Ỳ|Ý|Ỵ|Ỷ|Ỹ`),
	"D": regexp.MustCompile("Đ"),
}

var tagPattern = regexp.MustCompile(`#\w+`)

func Normalize(s string) string {
	var res = s
	for k, v := range pattern {
		res = v.ReplaceAllString(res, k)
	}

	res = strings.ToLower(res)
	return res
}

func ToArray(s string, sep string) []string {
	str := strings.Split(s, sep)
	res := make([]string, 0)
	for i := range str {
		res = append(res, strings.TrimSpace(str[i]))
	}

	return res
}

func ToTagArray(s string) []string {
	rs := Normalize(s)
	return tagPattern.FindAllString(rs, -1)
}

func Deduplicate(texts []string) []string {
	appended := make(map[string]bool)
	res := make([]string, 0)
	for i := range texts {
		if !appended[texts[i]] {
			res = append(res, texts[i])
			appended[texts[i]] = true
		}
	}

	return res
}

func DeduplicateTag(texts []string) []string {
	appended := make(map[string]bool)
	res := make([]string, 0)
	for i := range texts {
		if !appended[texts[i][1:]] {
			res = append(res, texts[i][1:])
			appended[texts[i][1:]] = true
		}
	}

	return res
}

func StringToInterfaceArray(texts []string) []interface{} {
	res := make([]interface{}, 0)

	for _, word := range texts {
		word = strings.TrimSpace(strings.ToLower(word))
		if len(word) > 0 {
			res = append(res, word)
		}
	}

	return res
}
