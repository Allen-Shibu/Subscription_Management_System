package email

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"testing"
)

func TestBuildMessageIncludesAttachment(t *testing.T) {
	raw, err := buildMessage("Brightloop <billing@example.com>", Message{
		To:             "ada@example.com",
		Subject:        "Invoice INV-000001",
		Text:           "Your invoice is attached.",
		HTML:           "<p>Your invoice is attached.</p>",
		AttachmentName: "INV-000001.pdf",
		Attachment:     []byte("%PDF-test"),
	})
	if err != nil {
		t.Fatalf("buildMessage() error = %v", err)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("read message: %v", err)
	}
	if got := msg.Header.Get("To"); got != "ada@example.com" {
		t.Fatalf("To = %q", got)
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/mixed" {
		t.Fatalf("Content-Type = %q, err = %v", msg.Header.Get("Content-Type"), err)
	}

	found := false
	reader := multipart.NewReader(msg.Body, params["boundary"])
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		if part.FileName() != "INV-000001.pdf" {
			continue
		}
		found = true
		data, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, part))
		if err != nil {
			t.Fatalf("decode attachment: %v", err)
		}
		if string(data) != "%PDF-test" {
			t.Fatalf("attachment = %q", data)
		}
	}
	if !found {
		t.Fatal("attachment part not found")
	}
}
