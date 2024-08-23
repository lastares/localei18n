# localei18n

go的国际化还挺麻烦的，一些框架本身就不提供的国际化，为了不重复造轮子，针对 [go-i18n](https://github.com/nicksnyder/go-i18n)
进行了二次封装，该组件默认支持中文Bundle，调用翻译函数的时候，会实时从context中去获取locale,
如果context中未设置locale，则默认翻译信息为中文。

目前配置文件格式支持`json`和`toml`

功能函数

- Tran(ctx context.Context, msgId string) 翻译函数
- TranWithTemplate(ctx context.Context, msgId string, template map[string]any) 翻译函数，支持模板变量替换

使用步骤

1. 安装依赖

```go
go get github.com/lastares/localei18n
```

2. 配置翻译文件，文件目录会在初始化国际化组件( LocaleI18N )对象的时候当做参数传递，最好为绝对路径，个人自测的配置文件示例

<pre>
localei18n
├── README.md
├── go.mod
├── go.sum
├── i18n.go
├── i18n_test.go
└── resources
    └── lang
        ├── en.json
        └── zh.json
</pre>
那在我本机的`translateDir`(翻译资源文件目录地址)就是 `/Users/ares/GolandProjects/localei18n/resources/lang`

json 格式的配置文件内容示例

```json
{
   "app.version": "版本",
   "hello": "你好 {{.name}}"
}
```

带花括号内的字符串会被替换成对应的变量，比如 `{{.name}}` 会被替换成 `ares`
如果需要自定义变量占位符，可按照上面内容进行配置

3. 初始化国际化组件( LocaleI18N ) 对象
    ```
   type LocaleI18N struct {
         bd   *i18n.Bundle    // Bundle 用于管理消息文件
         i18n *i18n.Localizer // Localizer 用于本地化消息
   }
   ```
   - 直接调用 NewLocaleI18N(translateDir string) 方法即可，当然你自己要套一层，将项目内配置的语言翻译文件的目录地址配置传进来
    ```
      func NewLocaleI18N(translateDir string) *LocaleI18N {} 
    ```
   为了保证国际化组件项目内全局可用，有几种方式，可自行选择
   - 可以在项目入口文件 main.go 中初始化国际化组件，然后通过全局变量的方式获取国际化组件对象
   - 如果项目内使用wire依赖注入，可以在 wire.go 文件中初始化国际化组件，然后通过 wire.go 文件中定义的函数获取国际化组件对象，伪代码示例
   ```go
    func NewGlobalLocaleI18N(conf *data.conf) *localei18n.LocaleI18N {
        return localei18n.NewLocaleI18N(conf.TranslateDir)
    }
   ```
   依赖初始化
   ```go
   var DataProviderSet = wire.NewSet(
    自定义的包名.NewGlobalLocaleI18N,
   )
   ```
   - 或者使用单例模式自己封装一个全局翻译组件变量

4. 注入，可以在服务层，可以在控制器层等，在需要的地方都可以进行注入

```go
type UserService struct {
// 这里省略其他注入对象
i18n *localei18n.LocaleI18N,
}

func NewUserService(i18n *localei18n.LocaleI18N) *UserService {
return &UserService{
// 这里省略其他注入对象
i18n:       i18n,
}
}
```

4. 调用翻译函数,具体根据实际业务逻辑进行处理

```go
func (u *UserService) Store(ctx context.Context, media *model.user) (*service.User, error) {
// 其他逻辑，伪代码
// 翻译
msg := u.i18n.Tran(ctx, "user not found")
// 带变量模板的翻译
msg = u.i18n.TranWithTemplate(ctx, "user not found", map[string]any{"name": "ares"})
}
```


   

