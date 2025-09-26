package eml_test

import (
	"bytes"
	"maps"
	"strings"
	"testing"

	"github.com/fadend/eml"
)

func TestFileNameToAttachmentEmpty(t *testing.T) {
	r := strings.NewReader("")
	_, err := eml.ExtractFileNameToAttachment(r)
	if err != eml.ErrMissingBoundary {
		t.Errorf("Expect ErrMissingBoundary, got %q instead", err)
	}
}

func TestFileNameToAttachmentSuccess(t *testing.T) {
	testEml := `MIME-Version: 1.0
Date: Fri, 26 Sep 2025 11:41:26 -0700
Message-ID: <Fake-Fake-Fake@mail.revfad.com>
Subject: test
From: Fake <fake@revfad.com>
To: Fake <fake@revfad.com>
Content-Type: multipart/mixed; boundary="00000000000017e255063fb8a2c5"

--00000000000017e255063fb8a2c5
Content-Type: multipart/alternative; boundary="00000000000017e251063fb8a2c3"

--00000000000017e251063fb8a2c3
Content-Type: text/plain; charset="UTF-8"



--00000000000017e251063fb8a2c3
Content-Type: text/html; charset="UTF-8"

<div dir="ltr"><br></div>

--00000000000017e251063fb8a2c3--
--00000000000017e255063fb8a2c5
Content-Type: text/plain; charset="US-ASCII"; name="test.txt"
Content-Disposition: attachment; filename="test.txt"
Content-Transfer-Encoding: base64
X-Attachment-Id: f_mg16u0an0
Content-ID: <f_mg16u0an0>

YWJj
--00000000000017e255063fb8a2c5
Content-Type: text/plain; charset="US-ASCII"; name="test2.txt"
Content-Disposition: attachment; filename="test2.txt"
Content-Transfer-Encoding: base64
X-Attachment-Id: f_mg16u9la1
Content-ID: <f_mg16u9la1>

ZGVm
--00000000000017e255063fb8a2c5--`

	r := strings.NewReader(testEml)
	result, err := eml.ExtractFileNameToAttachment(r)
	if err != nil {
		t.Errorf("Expect nil error, got %q instead", err)
	}
	expected := map[string][]byte{"test.txt": []byte("abc"), "test2.txt": []byte("def")}
	if !maps.EqualFunc(result, expected, bytes.Equal) {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
