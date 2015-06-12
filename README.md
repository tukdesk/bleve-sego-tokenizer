#### sego-tokenizer for bleve

使用了 [sego](https://github.com/huichen/sego) 作为分词器.  
为了适配 `bleve` 的 `analysis.Tokenizer` 接口, 我 fork 了[一个版本](https://github.com/tukdesk/sego), 有一些修改和修正.  

借用了 [jiebago](https://github.com/wangbin/jiebago) 的测试用例.
测试过程中, 发现 `sego` 和 `jiebago` 的分词结果会有少量的差异, 主要集中在一个分词结果包含多个子分词的情况下. 
比如 `中华人民`:  
`jieba` 的分词结果是 `中华`, `华人`, `人民`,  
`sego` 的分词结果是 `中华` `人民`.

全部测试样本中大约有10来例.

粗略看了一下, 这大概和 sego [从尾部开始逆序分词, 并且有可选结果时, 跳过中间的位置](https://github.com/tukdesk/sego/blob/master/segmenter.go#L194-L207) 有关

`sego` 看上去已经挺久没有维护.  
我会尽量排除 fork 过来的版本在使用过程中出现的 bug

非常感谢 [huichen](https://github.com/huichen) 和 [wangbin](https://github.com/wangbin) 提供了 golang 的中文分词实现.

使用方法参考:

```
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve"
	sego "github.com/tukdesk/bleve-sego-tokenizer"
)

func main() {
	// open a new index
	indexMapping := bleve.NewIndexMapping()

	err := indexMapping.AddCustomTokenizer("sego",
		map[string]interface{}{
			"files": "../dict.txt",
			"type":  sego.Name,
		})
	if err != nil {
		log.Fatal(err)
	}

	// create a custom analyzer
	err = indexMapping.AddCustomAnalyzer("sego",
		map[string]interface{}{
			"type":      "custom",
			"tokenizer": "sego",
			"token_filters": []string{
				"possessive_en",
				"to_lower",
				"stop_en",
			},
		})

	if err != nil {
		log.Fatal(err)
	}

	indexMapping.DefaultAnalyzer = "sego"
	cacheDir := "sego.beleve"
	os.RemoveAll(cacheDir)
	index, err := bleve.New(cacheDir, indexMapping)

	if err != nil {
		log.Fatal(err)
	}

	docs := []struct {
		Title string
		Name  string
	}{
		{
			Title: "Doc 1",
			Name:  "This is the first document we’ve added",
		},
		{
			Title: "Doc 2",
			Name:  "The second one 你 中文测试中文 is even more interesting! 吃水果",
		},
		{
			Title: "Doc 3",
			Name:  "买水果然后来世博园。",
		},
		{
			Title: "Doc 4",
			Name:  "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作",
		},
		{
			Title: "Doc 5",
			Name:  "咱俩交换一下吧。",
		},
	}
	// index docs
	for _, doc := range docs {
		index.Index(doc.Title, doc)
	}

	// search for some text
	for _, keyword := range []string{"水果世博园", "你", "first", "中文", "交换机", "交换"} {
		query := bleve.NewQueryStringQuery(keyword)
		search := bleve.NewSearchRequest(query)
		search.Highlight = bleve.NewHighlight()
		searchResults, err := index.Search(search)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Result of \"%s\": %d matches:\n", keyword, searchResults.Total)
		for i, hit := range searchResults.Hits {
			rv := fmt.Sprintf("%d. %s, (%f)\n", i+searchResults.Request.From+1, hit.ID, hit.Score)
			for fragmentField, fragments := range hit.Fragments {
				rv += fmt.Sprintf("%s: ", fragmentField)
				for _, fragment := range fragments {
					rv += fmt.Sprintf("%s", fragment)
				}
			}
			fmt.Printf("%s\n", rv)
		}
	}

	// Output:
	// Result of "水果世博园": 2 matches:
	// 1. Doc 3, (0.981169)
	// Name: 买<span class="highlight">水果</span>然后来<span class="highlight">世博</span>园。
	// 2. Doc 2, (0.051133)
	// Name: The second one 你 中文测试中文 is even more interesting! 吃<span class="highlight">水果</span>
	// Result of "你": 1 matches:
	// 1. Doc 2, (0.399574)
	// Name: The second one <span class="highlight">你</span> 中文测试中文 is even more interesting! 吃水果
	// Result of "first": 1 matches:
	// 1. Doc 1, (0.512150)
	// Name: This is the <span class="highlight">first</span> document we’ve added
	// Result of "中文": 1 matches:
	// 1. Doc 2, (0.565083)
	// Name: The second one 你 <span class="highlight">中文</span>测试<span class="highlight">中文</span> is even more interesting! 吃水果
	// Result of "交换机": 2 matches:
	// 1. Doc 4, (0.488048)
	// Name: 工信处女干事每月经过下属科室都要亲口交代24口<span class="highlight">交换机</span>等技术性器件的安装工作
	// 2. Doc 5, (0.165357)
	// Name: 咱俩<span class="highlight">交换</span>一下吧。
	// Result of "交换": 2 matches:
	// 1. Doc 5, (0.534158)
	// Name: 咱俩<span class="highlight">交换</span>一下吧。
	// 2. Doc 4, (0.302165)
	// Name: 工信处女干事每月经过下属科室都要亲口交代24口<span class="highlight">交换</span>机等技术性器件的安装工作
}

```