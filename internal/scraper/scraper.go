package scraper

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/eust-w/urlreader/internal/logger"
)

// Scraper 定义网页抓取器
type Scraper struct {
	collector *colly.Collector
}

// NewScraper 创建一个新的网页抓取器
func NewScraper() *Scraper {
	log := logger.GetLogger()
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.MaxDepth(1),
		colly.AllowURLRevisit(), // 允许重复访问同一URL
	)

	c.SetRequestTimeout(30 * time.Second)
	log.Infow("Scraper 初始化完成")

	return &Scraper{
		collector: c,
	}
}

// ScrapedContent 存储抓取的网页内容
type ScrapedContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

// ScrapeURL 抓取指定URL的内容
func (s *Scraper) ScrapeURL(url string) (*ScrapedContent, error) {
	log := logger.GetLogger()
	log.Infow("开始抓取URL", "url", url)
	if url == "" {
		log.Errorw("URL不能为空")
		return nil, errors.New("URL不能为空")
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
		log.Infow("自动补全URL为https", "url", url)
	}

	content := &ScrapedContent{
		URL: url,
	}

	// 提取标题
	s.collector.OnHTML("title", func(e *colly.HTMLElement) {
		content.Title = e.Text
		log.Infow("抓取到网页标题", "title", e.Text)
	})

	// 提取正文内容
	var textParts []string

	// 提取主要文本内容
	s.collector.OnHTML("body", func(e *colly.HTMLElement) {
		log.Infow("抓取到网页body")
		// 提取段落
		e.ForEach("p", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text != "" {
				textParts = append(textParts, text)
			}
		})

		// 提取标题
		e.ForEach("h1, h2, h3, h4, h5, h6", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text != "" {
				textParts = append(textParts, fmt.Sprintf("[%s] %s", el.Name, text))
			}
		})

		// 提取列表
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text != "" {
				textParts = append(textParts, "- "+text)
			}
		})

		// 提取表格内容
		e.ForEach("table", func(_ int, el *colly.HTMLElement) {
			textParts = append(textParts, "[表格]")
			el.ForEach("tr", func(_ int, row *colly.HTMLElement) {
				var rowTexts []string
				row.ForEach("td, th", func(_ int, cell *colly.HTMLElement) {
					text := strings.TrimSpace(cell.Text)
					if text != "" {
						rowTexts = append(rowTexts, text)
					}
				})
				if len(rowTexts) > 0 {
					textParts = append(textParts, strings.Join(rowTexts, " | "))
				}
			})
		})

		// 提取文章内容
		e.ForEach("article", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text != "" {
				textParts = append(textParts, text)
			}
		})

		// 提取div内容（可能包含主要内容）
		e.ForEach("div.content, div.main, div.article, div#content, div#main, div#article", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text != "" {
				textParts = append(textParts, text)
			}
		})
	})

	// 错误处理
	var scrapeErr error
	s.collector.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("抓取错误 %s: %w", r.Request.URL, err)
	})

	// 访问URL
	err := s.collector.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("访问URL失败: %w", err)
	}

	// 等待抓取完成
	s.collector.Wait()

	if scrapeErr != nil {
		return nil, scrapeErr
	}

	// 合并所有文本部分
	content.Content = strings.Join(textParts, "\n\n")

	// 如果内容为空，返回错误
	if content.Content == "" {
		return nil, errors.New("无法提取网页内容")
	}

	return content, nil
}
