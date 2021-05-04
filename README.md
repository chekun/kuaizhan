# 快站淘客开放接口API

[![Build Status](https://travis-ci.com/chekun/kuaizhan.svg?branch=main)](https://travis-ci.com/chekun/kuaizhan)
[![Go Report Card](https://goreportcard.com/badge/github.com/chekun/kuaizhan)](https://goreportcard.com/report/github.com/chekun/kuaizhan)
[![Go Doc](https://godoc.org/github.com/chekun/kuaizhan?status.svg)](http://godoc.org/github.com/chekun/kuaizhan)

## 快速使用

```go
appKey := "your app key"
appSecret := "your app secret"
// 创建client
client := kuaizhan.NewClient(appKey, appSecret, nil)
// 获取站点信息
siteID := "1234567"
siteInfo, err := kuaizhan.TbkGetSiteInfo(siteID, "")
if err != nil {
  panic(err)
}
// 查看结果
fmt.Printf("%+v\n", *siteInfo)
```

> 应用的申请参见官方文档：[官方文档](https://www.yuque.com/kuaizhan_help/ndcqmp)

## 支持快站所有接口

- [TbkDomain 根据站点获取站点域名](https://www.yuque.com/kuaizhan_help/ndcqmp/edk1gx)
- [TbkSiteTraffic 根据域名获取站点的流量信息](https://www.yuque.com/kuaizhan_help/ndcqmp/riha7l)
- [TbkChangeDomain 给指定站点修改域名](https://www.yuque.com/kuaizhan_help/ndcqmp/imqwx3)
- [TbkChangeDomainHttpsForward 开启/关闭域名Https跳转](https://www.yuque.com/kuaizhan_help/ndcqmp/mlvv6r)
- [TbkGenKzShortUrl 生成快站短链接](https://www.yuque.com/kuaizhan_help/ndcqmp/naqeg6)
- [TbkGenKzShortUrl 第三方短链接生成](https://www.yuque.com/kuaizhan_help/ndcqmp/qeepyp)
- [TbkRevertShortURL 快码短链还原长链](https://www.yuque.com/kuaizhan_help/ndcqmp/gywhwq)
- [TbkModifyPageJs 更新页面js](https://www.yuque.com/kuaizhan_help/ndcqmp/kahcrz)
- [TbkPublishPage 发布页面](https://www.yuque.com/kuaizhan_help/ndcqmp/tbrdea)
- [TbkCheckDomainBan 检测域名是否被封禁](https://www.yuque.com/kuaizhan_help/ndcqmp/bbn7u7)
- [TbkGenPromoteLink 生成淘口令推广链接](https://www.yuque.com/kuaizhan_help/ndcqmp/fk2vtr)
- [TbkGetSiteIds 获取站点ID列表](https://www.yuque.com/kuaizhan_help/ndcqmp/gvltmw)
- [TbkGetPageIds 获取页面ID列表](https://www.yuque.com/kuaizhan_help/ndcqmp/bkm43g)
- [TbkCreateSitePage 新建极速版站点页面](https://www.yuque.com/kuaizhan_help/ndcqmp/hexyr9)
- [TbkDeleteSitePage 新建极速版站点页面](https://www.yuque.com/kuaizhan_help/ndcqmp/dyi82m)
- [TbkCreateSite 创建快站站点](https://www.yuque.com/kuaizhan_help/ndcqmp/fv1x9b)
- [TbkPublishSite 发布站点](https://www.yuque.com/kuaizhan_help/ndcqmp/max5xw)
- [TbkUpdateSiteSetting 修改站点基本信息](https://www.yuque.com/kuaizhan_help/ndcqmp/motq6t)
- [TbkGetSiteInfo 获取站点基本信息](https://www.yuque.com/kuaizhan_help/ndcqmp/acgiwn)
- [TbkGetSiteBansCount 获取站点当月被封禁的次数](https://www.yuque.com/kuaizhan_help/ndcqmp/yigagy)

