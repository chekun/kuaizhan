package kuaizhan

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	endPointURL = "https://cloud.kuaizhan.com/api"
)

// Response 公共返回结构体
type Response struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

// Error 转换错误为error对象
func (res *Response) Error() error {
	if res.Code == 200 {
		return nil
	}
	return fmt.Errorf("code:%d,%s", res.Code, res.Msg)
}

// Client client
type Client struct {
	appKey    string
	appSecret string
	client    *http.Client
	debug     bool
}

// NewClient 创建client
func NewClient(appKey, appSecret string, client *http.Client) *Client {
	if client == nil {
		client = defaultHTTPClient()
	}
	return &Client{
		appKey:    appKey,
		appSecret: appSecret,
		client:    client,
	}
}

func defaultHTTPClient() *http.Client {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxConnsPerHost:     200,
		MaxIdleConnsPerHost: 30,
		IdleConnTimeout:     30 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{Transport: tr}
}

func (c *Client) SetDebugMode(isDebugging bool) {
	c.debug = isDebugging
}

func (c *Client) Println(v ...interface{}) {
	if !c.debug {
		return
	}
	log.Println(v...)
}

func (c *Client) signParams(values url.Values) {
	values.Set("appKey", c.appKey)
	params := map[string]string{}
	for k := range values {
		if values.Get(k) == "" {
			continue
		}
		params[k] = values.Get(k)
	}
	keys := make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	rawData := make([]string, 0)
	for _, k := range keys {
		rawData = append(rawData, k+params[k])
	}
	h := md5.New()
	_, _ = h.Write([]byte(c.appSecret + strings.Join(rawData, "") + c.appSecret))
	values.Set("sign", fmt.Sprintf("%x", h.Sum(nil)))
}

func (c *Client) PostForm(api string, values url.Values) (json.RawMessage, error) {
	c.signParams(values)
	c.Println("Begin request", api, "with values", values)
	res, err := c.client.PostForm(endPointURL+api, values)
	if err != nil {
		c.Println("Request", api, "error", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.Println("Read response body", api, "error", err)
		return nil, err
	}
	c.Println("Got response body", api, string(body))
	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		c.Println("Unmarshal response body", api, "error", err)
		return nil, err
	}
	err = r.Error()
	if err != nil {
		c.Println("Api", api, "returned error", err)
		return nil, err
	}
	return r.Data, nil
}

func (c *Client) Get(api string, values url.Values) (json.RawMessage, error) {
	c.signParams(values)
	c.Println("Begin request", api, "with values", values)
	res, err := c.client.Get(endPointURL + api + "?" + values.Encode())
	if err != nil {
		c.Println("Request", api, "error", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.Println("Read response body", api, "error", err)
		return nil, err
	}
	c.Println("Got response body", api, string(body))
	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		c.Println("Unmarshal response body", api, "error", err)
		return nil, err
	}
	err = r.Error()
	if err != nil {
		c.Println("Api", api, "returned error", err)
		return nil, err
	}
	return r.Data, nil
}

// TbkDomain 根据站点获取站点域名, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/edk1gx
func (c *Client) TbkDomain(siteID string) (string, error) {
	body, err := c.PostForm("/v1/tbk/getDomain", url.Values{
		"siteId": []string{siteID},
	})
	if err != nil {
		return "", err
	}
	var f struct {
		Domain string `json:"domain"`
	}
	_ = json.Unmarshal(body, &f)
	return f.Domain, nil
}

type SiteTraffic struct {
	SiteID  string `json:"siteId"`
	PvCount string `json:"pvCount"`
	UvCount string `json:"uvCount"`
}

// TbkSiteTraffic 根据域名获取站点的流量信息, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/riha7l
func (c *Client) TbkSiteTraffic(domain string) (*SiteTraffic, error) {
	body, err := c.Get("/v1/tbk/getSitePvUv", url.Values{
		"domain": []string{domain},
	})
	if err != nil {
		return nil, err
	}
	var f SiteTraffic
	_ = json.Unmarshal(body, &f)
	return &f, nil
}

// TbkChangeDomain 给指定站点修改域名, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/imqwx3
func (c *Client) TbkChangeDomain(siteID, domain string, forceHttps bool) error {
	params := url.Values{
		"siteId": []string{siteID},
		"domain": []string{domain},
	}
	if forceHttps {
		params.Add("httpsForward", "true")
	} else {
		params.Add("httpsForward", "false")
	}
	_, err := c.PostForm("/v1/tbk/changeDomain", params)
	if err != nil {
		return err
	}
	return nil
}

// TbkChangeDomainHttpsForward 开启/关闭域名Https跳转, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/mlvv6r
func (c *Client) TbkChangeDomainHttpsForward(siteID, domain string, forceHttps bool) error {
	params := url.Values{
		"siteId": []string{siteID},
		"domain": []string{domain},
	}
	if forceHttps {
		params.Add("httpsForward", "true")
	} else {
		params.Add("httpsForward", "false")
	}
	_, err := c.PostForm("/v1/tbk/changeDomainHttpsForward", params)
	if err != nil {
		return err
	}
	return nil
}

// TbkGenKzShortUrl 生成快站短链接, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/naqeg6
func (c *Client) TbkGenKzShortURL(longURL string) (string, error) {
	body, err := c.PostForm("/v1/tbk/genKzShortUrl", url.Values{
		"url": []string{longURL},
	})
	if err != nil {
		return "", err
	}
	var f struct {
		ShortURL string `json:"shortUrl"`
	}
	_ = json.Unmarshal(body, &f)
	return f.ShortURL, nil
}

// TbkGenKzShortUrl 第三方短链接生成, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/qeepyp
func (c *Client) TbkGenShortURL(longURL, urlType string) (string, error) {
	if urlType == "" {
		urlType = "default"
	}
	body, err := c.PostForm("/v1/tbk/genShortUrl", url.Values{
		"url":     []string{longURL},
		"urlType": []string{urlType},
	})
	if err != nil {
		return "", err
	}
	var f struct {
		ShortURL string `json:"shortUrl"`
	}
	_ = json.Unmarshal(body, &f)
	return f.ShortURL, nil
}

type RevertedShortURLInfo struct {
	IsKzLink   bool   `json:"isKzLink"`
	OriginLink string `json:"originLink"`
}

// TbkRevertShortURL 快码短链还原长链, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/gywhwq
func (c *Client) TbkRevertShortURL(shortURL string) (*RevertedShortURLInfo, error) {
	body, err := c.Get("/v1/tbk/shortUrlRevert", url.Values{
		"url": []string{shortURL},
	})
	if err != nil {
		return nil, err
	}
	var f RevertedShortURLInfo
	_ = json.Unmarshal(body, &f)
	return &f, nil
}

// TbkModifyPageJs 更新页面js, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/kahcrz
func (c *Client) TbkModifyPageJs(siteID, pageID, content string, isEncrypted bool) error {
	isEncryptedContent := "false"
	if isEncrypted {
		isEncryptedContent = "true"
	}
	_, err := c.PostForm("/v1/tbk/modifyPageJs", url.Values{
		"siteId":           []string{siteID},
		"pageId":           []string{pageID},
		"content":          []string{content},
		"isEncryptContent": []string{isEncryptedContent},
	})
	if err != nil {
		return err
	}
	return nil
}

// TbkPublishPage 发布页面, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/tbrdea
func (c *Client) TbkPublishPage(siteID, pageID string) (string, error) {
	body, err := c.PostForm("/v1/tbk/publishPage", url.Values{
		"siteId": []string{siteID},
		"pageId": []string{pageID},
	})
	if err != nil {
		return "", err
	}
	var f struct {
		URL string `json:"url"`
	}
	_ = json.Unmarshal(body, &f)
	return f.URL, nil
}

// TbkCheckDomainBan 检测域名是否被封禁, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/bbn7u7
func (c *Client) TbkCheckDomainBan(domain string) (bool, error) {
	body, err := c.PostForm("/v1/tbk/checkDomainBan", url.Values{
		"domain": []string{domain},
	})
	if err != nil {
		return false, err
	}
	var f struct {
		IsBannedWX bool `json:"isBannedWX"`
	}
	_ = json.Unmarshal(body, &f)
	return f.IsBannedWX, nil
}

// TbkGenPromoteLink 生成淘口令推广链接, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/fk2vtr
func (c *Client) TbkGenPromoteLink(siteID, code, imageURL string) (string, error) {
	body, err := c.PostForm("/v1/tbk/genPromoteLink", url.Values{
		"siteId": []string{siteID},
		"tkl":    []string{code},
		"image":  []string{imageURL},
	})
	if err != nil {
		return "", err
	}
	var f struct {
		Link string `json:"link"`
	}
	_ = json.Unmarshal(body, &f)
	return f.Link, nil
}

// TbkGetSiteIds 获取站点ID列表, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/gvltmw
func (c *Client) TbkGetSiteIds() ([]uint, error) {
	body, err := c.PostForm("/v1/tbk/getSiteIds", url.Values{})
	if err != nil {
		return nil, err
	}
	var f struct {
		SiteIDs []uint `json:"siteIds"`
	}
	_ = json.Unmarshal(body, &f)
	return f.SiteIDs, nil
}

// TbkGetPageIds 获取页面ID列表, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/bkm43g
func (c *Client) TbkGetPageIds(siteID string) ([]uint, error) {
	body, err := c.PostForm("/v1/tbk/getPageIds", url.Values{
		"siteId": []string{siteID},
	})
	if err != nil {
		return nil, err
	}
	var f struct {
		PageIDs []uint `json:"pageIds"`
	}
	_ = json.Unmarshal(body, &f)
	return f.PageIDs, nil
}

type Page struct {
	PageID uint   `json:"pageId"`
	Title  string `json:"title"`
}

// TbkGetPageName 获取所有站点页面名称, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/mehpsf
func (c *Client) TbkGetPageName(siteID string) ([]*Page, error) {
	body, err := c.Get("/v1/tbk/getPageName", url.Values{
		"siteId": []string{siteID},
	})
	if err != nil {
		return nil, err
	}
	pages := make([]*Page, 0)
	_ = json.Unmarshal(body, &pages)
	return pages, nil
}

// TbkCreateSitePage 新建极速版站点页面, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/hexyr9
func (c *Client) TbkCreateSitePage(siteID, template string) (uint, error) {
	if template == "" {
		template = "WHITE"
	}
	body, err := c.PostForm("/v1/tbk/createSitePage", url.Values{
		"siteId": []string{siteID},
		"tpl":    []string{template},
	})
	if err != nil {
		return 0, err
	}
	var f struct {
		PageID uint `json:"pageId"`
	}
	_ = json.Unmarshal(body, &f)
	return f.PageID, nil
}

// TbkDeleteSitePage 新建极速版站点页面, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/dyi82m
func (c *Client) TbkDeleteSitePage(pageId string) error {
	_, err := c.PostForm("/v1/tbk/deleteSitePage", url.Values{
		"pageId": []string{pageId},
	})
	if err != nil {
		return err
	}
	return nil
}

type Site struct {
	ID                   string `json:"siteId"`
	PageID               string `json:"pageId"`
	Domain               string `json:"siteDomain"`
	Status               string `json:"siteStatus"`
	PackageName          string `json:"packageName"`
	PackageRemainingDays uint   `json:"packageRemainingDays"`
}

// TbkCreateSite 创建快站站点, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/fv1x9b
func (c *Client) TbkCreateSite(name, domain, siteType string) (*Site, error) {
	if siteType == "" {
		siteType = "FAST"
	}
	body, err := c.PostForm("/v1/tbk/createSite", url.Values{
		"siteName": []string{name},
		"domain":   []string{domain},
		"siteType": []string{siteType},
	})
	if err != nil {
		return nil, err
	}
	var f Site
	_ = json.Unmarshal(body, &f)
	return &f, nil
}

// TbkPublishSite 发布站点, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/max5xw
func (c *Client) TbkPublishSite(siteID string) (string, error) {
	body, err := c.PostForm("/v1/tbk/publishSite", url.Values{
		"siteId": []string{siteID},
	})
	if err != nil {
		return "", err
	}
	var f struct {
		URL string `json:"url"`
	}
	_ = json.Unmarshal(body, &f)
	return f.URL, nil
}

// TbkUpdateSiteSetting 修改站点基本信息, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/motq6t
func (c *Client) TbkUpdateSiteSetting(siteID, siteName string) error {
	_, err := c.PostForm("/v1/tbk/updateSiteSetting", url.Values{
		"siteId":   []string{siteID},
		"siteName": []string{siteName},
	})
	if err != nil {
		return err
	}
	return nil
}

// TbkGetSiteInfo 获取站点基本信息, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/acgiwn
func (c *Client) TbkGetSiteInfo(siteID, siteDomain string) (*Site, error) {
	params := url.Values{}
	if siteID != "" {
		params.Set("siteId", siteID)
	} else if siteDomain != "" {
		params.Set("siteDomain", siteDomain)
	}
	body, err := c.PostForm("/v1/tbk/getSiteInfo", params)
	if err != nil {
		return nil, err
	}
	var f Site
	_ = json.Unmarshal(body, &f)
	return &f, nil
}

// TbkGetSiteBansCount 获取站点当月被封禁的次数, 详细查看 https://www.yuque.com/kuaizhan_help/ndcqmp/yigagy
func (c *Client) TbkGetSiteBansCount(siteID string) (uint, error) {
	body, err := c.PostForm("/v1/tbk/getSiteBanCount", url.Values{
		"siteId": []string{siteID},
	})
	if err != nil {
		return 0, err
	}
	var f struct {
		Count uint `json:"count"`
	}
	_ = json.Unmarshal(body, &f)
	return f.Count, nil
}

// AgentChangeDomain 修改客户域名
func (c *Client) AgentChangeDomain(siteID, domain string) error {
	_, err := c.PostForm("/v1/agent/changeDomain", url.Values{
		"siteId": []string{siteID},
		"domain": []string{domain},
	})
	if err != nil {
		return err
	}
	return nil
}
