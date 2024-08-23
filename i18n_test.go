package localei18n

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

var localeI18n *LocaleI18N

func init() {
	localeI18n = NewLocaleI18N("/Users/ares/GolandProjects/localei18n/resources/lang")
}

//func TestNewLocaleI18N(t *testing.T) {
//	NewLocaleI18N("/Users/ares/GolandProjects/localei18n/resources/lang")
//}

func TestTran(t *testing.T) {
	TestTranByZh(t)
	TestTranByEn(t)
}

func TestTranByZh(t *testing.T) {
	msg := localeI18n.Tran(context.Background(), "app.version")
	assert.Equal(t, "版本", msg)
}

func TestTranByEn(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, LocaleKey{}, language.English)
	msg := localeI18n.Tran(ctx, "app.version")
	assert.Equal(t, "version", msg)
}

func TestTranWithTemplate(t *testing.T) {
	TestTranWithTemplateByZh(t)
	TestTranWithTemplateByEn(t)
}

func TestTranWithTemplateByZh(t *testing.T) {
	msg := localeI18n.TranWithTemplate(context.Background(), "hello", map[string]any{"name": "阿瑞斯"})
	assert.Equal(t, "你好 阿瑞斯", msg)
}

func TestTranWithTemplateByEn(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, LocaleKey{}, language.English)
	msg := localeI18n.TranWithTemplate(ctx, "hello", map[string]any{"name": "ares"})
	assert.Equal(t, "hello ares", msg)
}

func TestListDir(t *testing.T) {
	_, err := GetDirFileList("/Users/ares/GolandProjects/localei18n/resources/lang")
	if err != nil {
		t.Errorf("ListDir() error = %v", err)
		return
	}
}
