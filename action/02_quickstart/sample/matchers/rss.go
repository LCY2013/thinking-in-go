package matchers

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"fufeng.org/sample/search"
)

// 搜索rss源的匹配器
type (
	// item根据 item 字段的标签，将定义的字段与rss文档的字段关联起来
	item struct {
		XMLName     xml.Name `xml:"item"`
		PubDate     string   `xml:"pubDate"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Link        string   `xml:"link"`
		GUID        string   `xml:"guid"`
		GeoRssPoint string   `xml:"georss:point"`
	}

	// image根据 image 字段的标签，将定义的字段与rss文档的字段关联起来
	image struct {
		XMLName xml.Name `xml:"image"`
		URL     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}

	// channel根据 channel 字段的标签，将定义的字段与rss文档的字段关联起来
	channel struct {
		XMLName        xml.Name `xml:"channel"`
		Title          string   `xml:"title"`
		Description    string   `xml:"description"`
		Link           string   `xml:"link"`
		PubDate        string   `xml:"pubDate"`
		LastBuildDate  string   `xml:"lastBuildDate"`
		TTL            string   `xml:"ttl"`
		Language       string   `xml:"language"`
		ManagingEditor string   `xml:"managingEditor"`
		WebMaster      string   `xml:"webMaster"`
		Image          image    `xml:"image"`
		Item           []item   `xml:"item"`
	}

	// rssDocument 定义了与rss文档关联的字段
	rssDocument struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)

// rssMatcher 实现 Matcher 接口
type rssMatcher struct {
}

// retrieve 发送 HTTP Get 请求获取 rss 数据源并解析
func (m rssMatcher) retrieve(feed *search.Feed) (*rssDocument, error) {
	if feed.URI == "" {
		return nil, errors.New("No rss feed URI provided")
	}

	// 从网络获取 rss 数据源文档
	resp, err := http.Get(feed.URI)
	if err != nil {
		return nil, err
	}

	// 函数返回后，关闭返回的响应链
	defer resp.Body.Close()

	// 检测网络状态是不是 200 ，以此来判断网络返回的数据是否正常
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP Response Error %d\n", resp.StatusCode)
	}

	// 将 rss 数据源文档解析到我们定义的结构体中，不需要检测错误，调用者会做这件事
	var document rssDocument
	err = xml.NewDecoder(resp.Body).Decode(&document)
	return &document, err
}

// Search 在文档中查找特定的搜索项
func (m rssMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
	var results []*search.Result
	log.Printf("Search Feed Type[%s] Site[%s] For Uri[%s]\n",
		feed.Type, feed.Name, feed.URI)

	// 获取想要搜索的项
	document, err := m.retrieve(feed)
	if err != err || document == nil {
		return nil, err
	}

	for _, channelItem := range document.Channel.Item {
		// 检查标题部分是否包含搜索项
		matched, err := regexp.MatchString(searchTerm, channelItem.Title)
		if err != nil {
			return nil, err
		}

		// 如果找到匹配项，将其作为结果保存
		if matched {
			results = append(results, &search.Result{
				Field:   "Title",
				Content: channelItem.Title,
			})
		}

		// 检查描述部分是否包含搜索项
		matched, err = regexp.MatchString(searchTerm, channelItem.Description)
		if err != nil {
			return nil, err
		}

		// 如果找到匹配项，将其作为结果保存
		if matched {
			results = append(results, &search.Result{
				Field:   "Description",
				Content: channelItem.Description,
			})
		}
	}
	return results, nil
}

// init 将匹配器注册到程序里
func init() {
	var matcher rssMatcher
	search.Register("rss", matcher)
}
