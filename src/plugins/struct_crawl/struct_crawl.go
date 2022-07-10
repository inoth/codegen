package struct_crawl

import (
	plugin "codegenfornet/src/plugins"
	"errors"
)

var crawlMap = make(map[string]plugin.DbStructureGet)

func AddCrawl(key string, crawl plugin.DbStructureGet) {
	crawlMap[key] = crawl
}

func GetCrawl(db string) (plugin.DbStructureGet, error) {
	r, ok := crawlMap[db]
	if !ok {
		return nil, errors.New("尚未支持的数据库")
	}
	return r, nil
}
