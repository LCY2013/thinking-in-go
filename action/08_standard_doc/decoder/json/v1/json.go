package v1

// json 编解码, 基础库 - Json 和 NewDecoder 函数,来解析JSON响应
type (
	// GgResult 映射到从搜索拿到的结果文档
	GgResult struct {
		GsearchResultClass string `json:"GsearchResultClass"`
		UnescapedURL       string `json:"unescapedUrl"`
		URL                string `json:"url"`
		VisibleURL         string `json:"visibleUrl"`
		CacheURL           string `json:"cacheUrl"`
		Title              string `json:"title"`
		TitleNoFormatting  string `json:"titleNoFormatting"`
		Content            string `json:"content"`
	}

	// GgResponse 包含顶级的文档
	GgResponse struct {
		ResponseData struct {
			Results []GgResult `json:"results"`
		} `json:"responseData"`
	}
)
