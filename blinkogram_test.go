package blinkogram

import (
	"testing"

	"github.com/go-telegram/bot/models"
)

func TestFormatContent_Plain(t *testing.T) {
	in := "hello world"
	got := formatContent(in, nil)
	if got != in {
		t.Fatalf("plain: got %q want %q", got, in)
	}
}

func TestFormatContent_Bold(t *testing.T) {
	in := "this is bold text here"
	ents := []models.MessageEntity{{Type: models.MessageEntityTypeBold, Offset: 8, Length: 4}}
	got := formatContent(in, ents)
	want := "this is **bold** text here"
	if got != want {
		t.Fatalf("bold: got %q want %q", got, want)
	}
}

func TestFormatContent_Italic(t *testing.T) {
	in := "a italic b"
	// content: 'a ' offset0-1, 'italic' offset2-7, ' b' 8
	ents := []models.MessageEntity{{Type: models.MessageEntityTypeItalic, Offset: 2, Length: 6}}
	got := formatContent(in, ents)
	want := "a *italic* b"
	if got != want {
		t.Fatalf("italic: got %q want %q", got, want)
	}
}

func TestFormatContent_URL(t *testing.T) {
	in := "see https://example.com/x now"
	ents := []models.MessageEntity{{Type: models.MessageEntityTypeURL, Offset: 4, Length: 21}}
	got := formatContent(in, ents)
	want := "see [https://example.com/x](https://example.com/x) now"
	if got != want {
		t.Fatalf("url: got %q want %q", got, want)
	}
}

func TestFormatContent_TextLink(t *testing.T) {
	in := "check our site here"
	ents := []models.MessageEntity{{Type: models.MessageEntityTypeTextLink, Offset: 6, Length: 8, URL: "https://example.com"}}
	got := formatContent(in, ents)
	want := "check [our site](https://example.com) here"
	if got != want {
		t.Fatalf("textlink: got %q want %q", got, want)
	}
}

func TestFormatContent_Mixed(t *testing.T) {
	in := "bold then italic"
	ents := []models.MessageEntity{
		{Type: models.MessageEntityTypeBold, Offset: 0, Length: 4},
		{Type: models.MessageEntityTypeItalic, Offset: 10, Length: 6},
	}
	got := formatContent(in, ents)
	want := "**bold** then *italic*"
	if got != want {
		t.Fatalf("mixed: got %q want %q", got, want)
	}
}

func TestFormatContent_HashtagPreserved(t *testing.T) {
	// unhandled entity (hashtag) should not be dropped
	in := "#news and bold"
	ents := []models.MessageEntity{
		{Type: models.MessageEntityTypeHashtag, Offset: 0, Length: 5},
		{Type: models.MessageEntityTypeBold, Offset: 10, Length: 4},
	}
	got := formatContent(in, ents)
	want := "#news and **bold**"
	if got != want {
		t.Fatalf("hashtag: got %q want %q", got, want)
	}
}

func TestFormatContent_CJKUtf16Offset(t *testing.T) {
	// CJK chars are 1 UTF-16 unit each, so offsets are straightforward
	in := "你好世界 hello"
	// 你0 好1 世2 界3 ' '4 h5 e6 l7 l8 o9  -> "hello" at offset 5, length 5
	ents := []models.MessageEntity{{Type: models.MessageEntityTypeBold, Offset: 5, Length: 5}}
	got := formatContent(in, ents)
	want := "你好世界 **hello**"
	if got != want {
		t.Fatalf("cjk: got %q want %q", got, want)
	}
}

func TestFormatContent_ForwardedPrefixWithEntities(t *testing.T) {
	// simulates forwarded message with "Forwarded from X\n" prefix already added,
	// then a bold entity in the original text
	in := "Forwarded from [Channel](https://t.me/c)\nbold part here"
	ents := []models.MessageEntity{{Type: models.MessageEntityTypeBold, Offset: 41, Length: 4}}
	got := formatContent(in, ents)
	want := "Forwarded from [Channel](https://t.me/c)\n**bold** part here"
	if got != want {
		t.Fatalf("fwd: got %q want %q", got, want)
	}
}
