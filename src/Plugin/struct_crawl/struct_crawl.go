package struct_crawl

import (
	plugin "codegenfornet/src/Plugin"
	"errors"
)

var crawlMap = make(map[string]plugin.DbStructureGet)

func AddCrawl(key string, crawl plugin.DbStructureGet) {
	crawlMap[key] = crawl
}

func GetCrawl(db string) (plugin.DbStructureGet, error) {
	r, ok := crawlMap[db]
	if !ok {
		return nil, errors.New("该类型数据库尚未支持")
	}
	return r, nil
}
