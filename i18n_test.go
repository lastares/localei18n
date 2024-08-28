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

func TestTran(t *testing.T) {
	TestTranByZh(t)
	TestTranByEn(t)
}

func TestTranslate(t *testing.T) {
	cases := []struct {
		locale   language.Tag
		msgId    string
		template map[string]any
		want     string
	}{
		{
			locale: language.Chinese,
			msgId:  "app.version",
			want:   "版本",
		},
		{
			locale: language.English,
			msgId:  "app.version",
			want:   "version",
		},
		{
			locale: language.Chinese,
			msgId:  "hello",
			template: map[string]any{
				"name": "阿瑞斯",
			},
			want: "你好 阿瑞斯",
		},
		{
			locale: language.English,
			msgId:  "hello",
			template: map[string]any{
				"name": "ares",
			},
			want: "hello ares",
		},
	}
	t.Run("Translate", func(t *testing.T) {
		for i, item := range cases {
			got := localeI18n.Translate(&item.locale, item.msgId)
			if len(item.template) > 0 {
				got = localeI18n.TranslateWithTemplate(&item.locale, item.msgId, item.template)
			}
			if !assert.Equal(t, item.want, got) {
				t.Errorf("Tran() case %d failed", i)
			}
		}
	})
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
