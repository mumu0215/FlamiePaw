package crawler

type CrawlerJson struct {
	AllReqList []interface{} `json:"all_req_list"`
	ReqList []interface{} `json:"req_list"`
	AllDomainList []string `json:"all_domain_list"`
	SubDomainList []string `json:"sub_domain_list"`
}

//yaml解析
type CustomHeader struct {
	UserAgent string `json:"User-Agent"`
}
type Chrome struct {
	Path string `yaml:"Path"`
	MaxTapCount string `yaml:"MaxTapCount"`
}
type Crawler struct {
	CrawlerThread int `yaml:"CrawlerThread"`
	FilterMode string `yaml:"FilterMode"`
	Proxy string `yaml:"Proxy"`
	PushToPoolMax string `yaml:"PushToPoolMax"`
}
type CrawlerConfig struct {
	MyChrome Chrome `yaml:"Chrome"`
	MyCrawler Crawler `yaml:"Crawler"`
}
var UserAgent=[]string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) Gecko/20100101 Firefox/41.0",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:28.0) Gecko/20100101 Firefox/31.0",
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; WOW64; en-US; rv:2.0.4) Gecko/20120718 AskTbAVR-IDW/3.12.5.17700 Firefox/14.0.1",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.116 Safari/537.36 Mozilla/5.0 (iPad; U; CPU OS 3_2 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Version/4.0.4 Mobile/7B334b Safari/531.21.10",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/602.2.14 (KHTML, like Gecko) Version/10.0.1 Safari/602.2.14",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36",
}