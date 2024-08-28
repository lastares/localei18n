package localei18n

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// 常量定义区域，用于指定配置文件的类型
const (
	// Json 代表JSON格式的配置文件
	Json = "json"
	// Toml 代表TOML格式的配置文件
	Toml = "toml"
)

// LocaleKey 是一个结构体，用于标识上下文中存储的语言标签键。
type LocaleKey struct{}

// LocaleI18N 结构体用于存储国际化相关的数据。
type LocaleI18N struct {
	bd   *i18n.Bundle    // Bundle 用于管理消息文件
	i18n *i18n.Localizer // Localizer 用于本地化消息
}

// NewLocaleI18N 创建一个新的 LocaleI18N 实例。
// 参数 translateDir 指定了存放翻译文件的目录路径。
func NewLocaleI18N(translateDir string) *LocaleI18N {
	// 创建一个新的 Bundle，指定默认语言为中文
	bd := i18n.NewBundle(language.Chinese)
	localeI18N := &LocaleI18N{
		bd: bd,
	}
	// 注册配置文件解析函数
	localeI18N.bd.RegisterUnmarshalFunc(Json, json.Unmarshal)
	// 获取翻译文件列表
	translateFiles, err := GetDirFileList(translateDir)
	if err != nil {
		log.Fatalf("failed to load translate files: %v", err)
	}
	if len(translateFiles) == 0 {
		log.Fatalf("no translate file found in %s", translateDir)
	}
	// 加载所有翻译文件
	for _, translateFile := range translateFiles {
		// 加载单个翻译文件
		bd.LoadMessageFile(translateFile)
		fmt.Println("translate file loaded: ", translateFile)
	}

	// 返回新的 LocaleI18N 实例
	localeI18N.i18n = i18n.NewLocalizer(bd)
	return localeI18N
}

func (r *LocaleI18N) registerUnmarshalFunc() {
	r.bd.RegisterUnmarshalFunc(Json, json.Unmarshal)
}

type Localizer struct {
	Locale language.Tag
}

// switchLocalizer 方法用于切换语言环境
// 此方法在需要为 LocaleI18N 实例更改语言环境时调用
// 参数:
//   - locale: 一个指向 language.Tag 类型的指针，表示新的语言环境
//
// 通过这个方法，我们可以根据不同的语言环境需求动态地更改 LocaleI18N 实例的 Localizer
func (r *LocaleI18N) switchLocalizer(locale *language.Tag) {
	// 使用新的语言环境创建一个新的 Localizer 实例
	r.i18n = i18n.NewLocalizer(r.bd, locale.String())
}

// Tran 提供了一个简单的方法来翻译消息 ID。
// 如果需要，它会根据上下文中的语言标签选择合适的翻译。
// 参数:
// - ctx: 上下文，用于获取语言标签。
// - msgId: 要翻译的消息 ID。
// 返回:
// - 翻译后的字符串。
func (r *LocaleI18N) Tran(ctx context.Context, msgId string) string {
	return r.translate(ctx, msgId, nil, nil)
}

// TranWithTemplate 提供了一个方法来翻译带有模板的消息 ID。
// 这个方法允许传入一个模板参数，用于替换翻译字符串中的占位符。
// 参数:
// - ctx: 上下文，用于获取语言标签。
// - msgId: 要翻译的消息 ID。
// - template: 一个映射，用于替换翻译字符串中的占位符。
// 返回:
// - 翻译后的字符串。
func (r *LocaleI18N) TranWithTemplate(ctx context.Context, msgId string, template map[string]any) string {
	return r.translate(ctx, msgId, template, nil)
}

// Translate 方法用于将给定的消息 ID 翻译成指定语言的文本。
// 它接收一个语言标签和一个消息 ID，然后返回翻译后的文本。
// 参数:
//   - locale: 指定翻译所使用的语言标签。这个标签遵循 BCP 47 的语言标签格式，例如 "en-US" 表示美国英语。
//   - msgId: 需要翻译的消息的唯一标识符。消息 ID 通常是在翻译资源包中定义的消息模板的名称。
//
// 返回值:
//   - string: 翻译后的文本字符串。如果无法找到对应的消息 ID 或者翻译资源，将返回空字符串。
func (r *LocaleI18N) Translate(locale *language.Tag, msgId string) string {
	return r.translate(nil, msgId, nil, locale)
}

// TranslateWithTemplate 根据指定的语言标签和模板数据，翻译消息。
// 此方法用于需要模板渲染的翻译场景，允许传入自定义的语言环境和模板数据。
// 参数:
//   - locale  - 指定的语言标签，用于确定翻译的目标语言。
//   - msgId   - 要翻译的消息的唯一标识符，用于定位具体翻译内容。
//   - template- 模板数据，用于在翻译过程中动态替换变量。
//
// 返回值:
//   - 翻译并应用模板后的字符串。如果翻译失败或模板数据不匹配，可能返回原始消息ID或错误信息。
func (r *LocaleI18N) TranslateWithTemplate(locale *language.Tag, msgId string, template map[string]any) string {
	return r.translate(nil, msgId, template, locale)
}

// translate 翻译指定的消息。
//
// ctx: 上下文，用于传递请求范围的数据和截止时间信息。
// msgId: 消息的唯一标识符，用于定位要翻译的消息模板。
// template: 模板数据，用于填充消息模板中的变量。
//
// 返回值:
// string: 翻译后的消息字符串。如果发生错误，则返回空字符串。
func (r *LocaleI18N) translate(ctx context.Context, msgId string, template map[string]any, locale *language.Tag) string {
	if ctx != nil {
		// 初始化本地化器，确保 localize 方法可以正常调用。
		value := ctx.Value(LocaleKey{}) // 从上下文中获取 LocaleKey 的值
		if value != nil {
			newLocale, ok := value.(language.Tag)
			if ok {
				r.switchLocalizer(&newLocale)
			}
		}
	}
	if locale != nil {
		r.switchLocalizer(locale)
	}
	// 使用 i18n 包的 Localize 方法进行消息翻译。
	// MessageID 和 TemplateData 用于指定要翻译的消息和模板数据。
	msg, err := r.i18n.Localize(&i18n.LocalizeConfig{
		MessageID:    msgId,
		TemplateData: template,
	})
	if err != nil {
		// 如果翻译过程中发生错误，则记录错误日志。
		log.Printf("LocaleI18N Tran msgId %s error: %v", msgId, err)
	}
	return msg
}

// GetDirFileList 获取指定目录下的所有文件列表
// 参数 dir：需要扫描的目录路径
// 返回值 []string：目录下所有文件的完整路径列表
// 返回值 error：可能出现的错误，如果为nil，表示执行成功
func GetDirFileList(dir string) ([]string, error) {
	// 初始化文件路径切片
	var fullFileNames []string
	// 使用 filepath.Walk 遍历目录，err 用于接收 filepath.Walk 的执行结果
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		// 如果遍历过程中发生错误，直接返回错误
		if err != nil {
			return err
		}
		// 如果当前路径不是目录，则将其添加到文件路径切片中
		if !info.IsDir() {
			fullFileNames = append(fullFileNames, path)
		}
		// 遍历继续进行，无错误返回
		return nil
	})
	// 如果 filepath.Walk 执行过程中出现错误，记录日志并返回错误
	if err != nil {
		log.Printf("Error walking directory: %v", err)
		return nil, err
	}
	// 返回文件路径切片，表示目录遍历成功
	return fullFileNames, nil
}
