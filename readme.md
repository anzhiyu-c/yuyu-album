# 这是一个极简的相册展示项目

![首页](https://upload-bbs.miyoushe.com/upload/2025/06/12/125766904/edc34204be7ed35ecf2cf928095f7501_5939306686125312503.png)

预览地址：[http://album.anheyu.com/](http://album.anheyu.com/)

## 技术栈

前端: vue3 + vite + element-plus

后端: go + gin

UI设计：[张洪](https://plog.zhheo.com/)

得益于社区有良好的生态后台直接使用了 [Pure Admin](https://pure-admin.cn/) 构建

## 配置文件参数

| 配置分类       | 参数名称              | 示例值                                                                                     | 说明                                                                   |
| :------------- | :-------------------- | :----------------------------------------------------------------------------------------- | :--------------------------------------------------------------------- |
| **管理员账户** | `ADMIN_USERNAME`      | `anzhiyu`                                                                                  | 管理员用户的登录名                                                     |
|                | `ADMIN_PASSWORD`      | `anzhiyu`                                                                                  | 管理员用户的登录密码                                                   |
| **数据库连接** | `DB_USER`             | `root`                                                                                     | 数据库连接用户名                                                       |
|                | `DB_PASS`             | `root`                                                                                     | 数据库连接密码                                                         |
|                | `DB_NAME`             | `album`                                                                                    | 要连接的数据库名称                                                     |
|                | `DB_HOST`             | `127.0.0.1`                                                                                | 数据库服务器的 IP 地址或主机名                                         |
|                | `DB_PORT`             | `3306`                                                                                     | 数据库服务器的端口号                                                   |
| **Redis 连接** | `REDIS_ADDR`          | `localhost:6379`                                                                           | Redis 服务器的地址和端口                                               |
|                | `REDIS_PASSWORD`      | (空字符串)                                                                                 | Redis 连接密码（如果无密码则留空）                                     |
|                | `REDIS_DB`            | `10`                                                                                       | 要连接的 Redis 数据库索引                                              |
| **安全及应用** | `JWT_SECRET`          | `yuyu_album`                                                                               | 用于 JWT 签名和验证的密钥，请务必使用强密钥                            |
|                | `ABOUT_LINK`          | `https://github.com/anzhiyu-c/yuyu-album`                                                  | 关于页面或项目相关信息的外部链接                                       |
|                | `APP_NAME`            | `鱼鱼相册`                                                                                 | 应用程序的名称                                                         |
|                | `APP_VERSION`         | `1.0.0`                                                                                    | 应用程序的版本号                                                       |
|                | `ICP_NUMBER`          | `湘ICP备2023015794号-2`                                                                    | 网站的 ICP 备案号（用于国内网站）                                      |
| **资源链接**   | `USER_AVATAR`         | `https://npm.elemecdn.com/anzhiyu-blog-static@1.0.4/img/avatar.jpg`                        | 默认用户头像的 URL                                                     |
|                | `API_URL`             | `https://wallpaper.anheyu.com/`                                                            | 应用程序的 API 基础 URL                                                |
|                | `LOGO_URL`            | `https://wallpaper.anheyu.com/logo.svg`                                                    | 应用程序 Logo 图片的 URL                                               |
|                | `ICON_URL`            | `https://wallpaper.anheyu.com/logo.svg`                                                    | 应用程序 Icon 图片的 URL                                               |
| **图片处理**   | `DEFAULT_THUMB_PARAM` | `"x-oss-process=image//resize,h_600/quality,q_100/auto-orient,0/interlace,1/format,avif"`  | 图片缩略图处理的默认参数字符串，通常用于 OSS（对象存储服务）的图片处理 |
|                | `DEFAULT_BIG_PARAM`   | `"x-oss-process=image//resize,s_2000/quality,q_100/auto-orient,0/interlace,1/format,avif"` | 大图处理的默认参数字符串，通常用于 OSS 的图片处理                      |

## 项目运行

1. 从[https://github.com/anzhiyu-c/yuyu-album/releases](https://github.com/anzhiyu-c/yuyu-album/releases)下载对应服务器最新的发布版本

2. 目前只支持服务器部署

