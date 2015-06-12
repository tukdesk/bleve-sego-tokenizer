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
