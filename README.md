# bilibili

Go 1.18 Bilibili SDK，参考 `bilibili-api` 的 `Api + client + module` 设计思路，重建为适合 Go 的同步类型化客户端。

## 架构

- `Client`
  - 管理配置、认证信息、模块实例与共享状态。
- `RequestBuilder`
  - 负责参数组装、WBI 签名、CSRF 注入、响应解码。
- `endpoint`
  - 用声明式元数据描述接口路径、鉴权要求、返回载荷位置。
- `*Service`
  - 面向业务域提供稳定 API，目前包含 `Video`、`User`、`Search`、`Live`、`Login`。

## 示例

```go
package main

import (
	"context"
	"log"

	"github.com/bilibili-go/bilibili"
)

func main() {
	client := bilibili.NewClient()

	info, err := client.Video().InfoByBVID(context.Background(), "BV1xx411c7mD")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(info.Title, info.Owner.Name)
}
```

## 迁移原则

- `biligo` 中的旧式扁平请求函数不直接照搬。
- 先稳定基础设施，再按模块逐步迁移旧能力和类型。
- 新增接口时，优先补 `endpoint` 和 `Service`，避免继续扩散全局函数。
