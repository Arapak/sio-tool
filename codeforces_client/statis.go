package codeforces_client

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Arapak/sio-tool/util"
)

type StatisInfo struct {
	ID     string
	Name   string
	IO     string
	Limit  string
	Passed string
	State  string
}

func findStatisBlock(body []byte) ([]byte, error) {
	reg := regexp.MustCompile(`class="problems"[\s\S]+?</tr>([\s\S]+?)</table>`)
	tmp := reg.FindSubmatch(body)
	if tmp == nil {
		return nil, errors.New("cannot find any problem statis")
	}
	return tmp[1], nil
}

func findProblems(body []byte) ([]StatisInfo, error) {
	reg := regexp.MustCompile(`<tr[\s\S]*?>`)
	tmp := reg.FindAllIndex(body, -1)
	if tmp == nil {
		return nil, errors.New("cannot find any problem")
	}
	var ret []StatisInfo
	scr := regexp.MustCompile(`<script[\s\S]*?>[\s\S]*?</script>`)
	cls := regexp.MustCompile(`class="(.+?)"`)
	rep := regexp.MustCompile(`<[\s\S]+?>`)
	ton := regexp.MustCompile(`<\s+`)
	rmv := regexp.MustCompile(`<+`)
	tmp = append(tmp, []int{len(body), 0})
	for i := 1; i < len(tmp); i++ {
		state := ""
		if x := cls.FindSubmatch(body[tmp[i-1][0]:tmp[i-1][1]]); x != nil {
			state = string(x[1])
		}
		b := scr.ReplaceAll(body[tmp[i-1][0]:tmp[i][0]], []byte{})
		b = rep.ReplaceAll(b, []byte("<"))
		b = ton.ReplaceAll(b, []byte("<"))
		b = rmv.ReplaceAll(b, []byte("<"))
		data := strings.Split(string(b), "<")
		var tot []string
		for j := 0; j < len(data); j++ {
			s := strings.TrimSpace(data[j])
			if s != "" {
				tot = append(tot, s)
			}
		}
		if len(tot) >= 5 {
			tot[4] = strings.ReplaceAll(tot[4], "x", "")
			tot[4] = strings.ReplaceAll(tot[4], "&nbsp;", "")
			if tot[4] == "" {
				tot[4] = "0"
			}
			ret = append(ret, StatisInfo{
				tot[0], tot[1], tot[2], tot[3],
				tot[4], state,
			})
		}
	}
	return ret, nil
}

func (c *CodeforcesClient) Statis(info Info) (problems []StatisInfo, perf util.Performance, err error) {
	URL, err := info.ProblemSetURL(c.host)
	if err != nil {
		return
	}
	if info.ProblemType == "acmsguru" {
		return nil, perf, errors.New(ErrorNotSupportAcmsguru)
	}

	perf.StartFetching()

	body, err := util.GetBody(c.client, URL)
	if err != nil {
		return
	}

	perf.StopFetching()
	perf.StartParsing()

	if _, err = findHandle(body); err != nil {
		return
	}

	block, err := findStatisBlock(body)
	if err != nil {
		return
	}

	problems, err = findProblems(block)

	perf.StopParsing()
	return
}
