package search

// 搜索数据用的默认匹配器

// defaultMatcher 实现默认匹配器
type defaultMatcher struct {
}

// init 函数将默认匹配器注册到程序里面
func init() {
	var matcher defaultMatcher
	Register("default", matcher)
}

// Search 实现了默认匹配器的行为
func (m defaultMatcher) Search(feed *Feed, searchTerm string) ([]*Result, error) {
	return nil, nil
}
