package kuaizhan_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/chekun/kuaizhan"
)

func setUpClient(t *testing.T) *kuaizhan.Client {
	appKey := os.Getenv("APP_KEY")
	appSecret := os.Getenv("APP_SECRET")
	if appKey == "" || appSecret == "" {
		t.Fatalf("app_id, app_secret must be provided\n")
	}
	client := kuaizhan.NewClient(appKey, appSecret, nil)
	client.SetDebugMode(false)
	return client
}

func TestTbkDomain(t *testing.T) {
	client := setUpClient(t)
	domain, err := client.TbkDomain(os.Getenv("SITE_ID"))
	if err != nil {
		t.Logf("failed to get tbk domain %s\n", err)
		return
	}
	t.Log("domain is", domain)
}

func TestTbkSiteTraffic(t *testing.T) {
	client := setUpClient(t)
	trafficData, err := client.TbkSiteTraffic(os.Getenv("SITE_DOMAIN"))
	if err != nil {
		t.Logf("failed to get tbk traffic data %s\n", err)
		return
	}
	t.Log("traffic data is", "pv", trafficData.PvCount, "uv", trafficData.UvCount)
}

func TestTbkChangeDomain(t *testing.T) {
	client := setUpClient(t)
	err := client.TbkChangeDomain(os.Getenv("SITE_ID"), os.Getenv("SITE_DOMAIN"), false)
	if err != nil {
		t.Logf("failed to change domain %s\n", err)
		return
	}
	t.Log("change domain ok")
}

func TestTbkChangeDomainHttpsForward(t *testing.T) {
	client := setUpClient(t)
	err := client.TbkChangeDomainHttpsForward(os.Getenv("SITE_ID"), os.Getenv("SITE_DOMAIN"), false)
	if err != nil {
		t.Logf("failed to change domain https forward %s\n", err)
		return
	}
	t.Log("change domain https forward ok")
}

func TestTbkGenKzShortURL(t *testing.T) {
	client := setUpClient(t)
	longURL := fmt.Sprintf("https://%s.kuaizhan.com/?a=b", os.Getenv("SITE_DOMAIN"))
	shortURL, err := client.TbkGenKzShortURL(longURL)
	if err != nil {
		t.Logf("failed to gen kuaizhan short url, %s\n", err)
		return
	}
	t.Log("short url is", shortURL)
}

func TestTbkGenShortURL(t *testing.T) {
	client := setUpClient(t)
	longURL := fmt.Sprintf("https://%s.kuaizhan.com/?a=b", os.Getenv("SITE_DOMAIN"))
	shortURL, err := client.TbkGenShortURL(longURL, "default")
	if err != nil {
		t.Logf("failed to gen short url, %s\n", err)
		return
	}
	t.Log("short url is", shortURL)
}

func TestTbkRevertShortURL(t *testing.T) {
	client := setUpClient(t)
	r, err := client.TbkRevertShortURL("http://kzurl08.cn/wGQS")
	if err != nil {
		t.Logf("failed to revert short url, %s\n", err)
		return
	}
	t.Log("short url is kuaizhan link ?", r.IsKzLink, "original url is ", r.OriginLink)
}

func TestTbkModifyPageJs(t *testing.T) {
	client := setUpClient(t)
	err := client.TbkModifyPageJs(os.Getenv("SITE_ID"), "", `alert('it works!')`, false)
	if err != nil {
		t.Logf("failed to modify page js, %s\n", err)
		return
	}
	t.Log("modify page js ok")
}

func TestTbkPublishPage(t *testing.T) {
	client := setUpClient(t)
	pageURL, err := client.TbkPublishPage(os.Getenv("SITE_ID"), "")
	if err != nil {
		t.Logf("failed to modify page js, %s\n", err)
		return
	}
	t.Log("page url is", pageURL)
}

func TestTbkCheckDomainBan(t *testing.T) {
	client := setUpClient(t)
	status, err := client.TbkCheckDomainBan("https://xxx.kuaizhan.com")
	if err != nil {
		t.Logf("failed to check domain status, %s\n", err)
		return
	}
	t.Log("domain banned status", status)
}

func TestTbkGenPromoteLink(t *testing.T) {
	client := setUpClient(t)
	pURL, err := client.TbkGenPromoteLink(os.Getenv("SITE_ID"), "￥Do39YvZOJEq￥", "//img.alicdn.com/i4/2939922051/O1CN01i12Oqd1R1OtyrrCIS_!!0-item_pic.jpg")
	if err != nil {
		t.Logf("failed to generate promotion link, %s\n", err)
		return
	}
	t.Log("promotion link is ", pURL)
}

func TestTbkGetSiteIds(t *testing.T) {
	client := setUpClient(t)
	siteIDs, err := client.TbkGetSiteIds()
	if err != nil {
		t.Logf("failed to get site ids, %s\n", err)
		return
	}
	t.Log("site ids are ", siteIDs)
}

func TestTbkGetPageIds(t *testing.T) {
	client := setUpClient(t)
	pageIDs, err := client.TbkGetPageIds(os.Getenv("SITE_ID"))
	if err != nil {
		t.Logf("failed to get page ids, %s\n", err)
		return
	}
	t.Log("page ids are ", pageIDs)
}

func TestTbkGetPageName(t *testing.T) {
	client := setUpClient(t)
	pages, err := client.TbkGetPageName(os.Getenv("SITE_ID"))
	if err != nil {
		t.Logf("failed to get pages, %s\n", err)
		return
	}
	t.Logf("pages are %+v \n", pages)
}

func TestTbkCreateSitePage(t *testing.T) {
	client := setUpClient(t)
	pageID, err := client.TbkCreateSitePage(os.Getenv("SITE_ID"), "")
	if err != nil {
		t.Logf("failed to create page, %s\n", err)
		return
	}
	t.Log("created page id is ", pageID)
}

func TestTbkDeleteSitePage(t *testing.T) {
	client := setUpClient(t)
	err := client.TbkDeleteSitePage("111")
	if err != nil {
		t.Logf("failed to delete page, %s\n", err)
		return
	}
	t.Log("page deleted")
}

func TestTbkCreateSite(t *testing.T) {
	client := setUpClient(t)
	site, err := client.TbkCreateSite("测试", "test", "FAST")
	if err != nil {
		t.Logf("failed to create site, %s\n", err)
		return
	}
	t.Logf("site created %+v\n", *site)
}

func TestTbkPublishSite(t *testing.T) {
	client := setUpClient(t)
	siteURL, err := client.TbkPublishSite(os.Getenv("SITE_ID"))
	if err != nil {
		t.Logf("failed to publish site, %s\n", err)
		return
	}
	t.Log("site published, url is ", siteURL)
}

func TestTbkUpdateSiteSetting(t *testing.T) {
	client := setUpClient(t)
	err := client.TbkUpdateSiteSetting(os.Getenv("SITE_ID"), "优惠购买")
	if err != nil {
		t.Logf("failed to update site settings, %s\n", err)
		return
	}
	t.Log("site settings updated")
}

func TestTbkGetSiteInfo(t *testing.T) {
	client := setUpClient(t)
	site, err := client.TbkGetSiteInfo(os.Getenv("SITE_ID"), "")
	if err != nil {
		t.Logf("failed to get site info, %s\n", err)
		return
	}
	t.Logf("site info %+v\n", *site)
}

func TestTbkGetSiteBansCount(t *testing.T) {
	client := setUpClient(t)
	count, err := client.TbkGetSiteBansCount(os.Getenv("SITE_ID"))
	if err != nil {
		t.Logf("failed to get site bans count, %s\n", err)
		return
	}
	t.Logf("site bans %d times", count)
}

func TestTbkSitePvUvBySiteId(t *testing.T) {
	client := setUpClient(t)
	traffic, err := client.TbkSitePvUvBySiteId(os.Getenv("SITE_ID"), "20210617", "20210617")
	if err != nil {
		t.Logf("failed to get site traffic, %s\n", err)
		return
	}
	t.Logf("site traffic %+v", *traffic)
}
