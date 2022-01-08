package hugo

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseFrontMatterFromReader(t *testing.T) {
	testCases := []struct {
		desc      string
		input     string
		expectFM  *FrontMatter
		expectErr error
	}{
		{
			desc: "Parse YAML front matter",
			input: `---
title: "HugoでもTwitterCardを自動生成したい"
author: ["@Ladicle"]
date: 2020-06-21T03:56:24+09:00
tags: ["hugo", "go", "OGP"]
categories: ["program"]
---
content`,
			expectFM: &FrontMatter{
				Title:    "HugoでもTwitterCardを自動生成したい",
				Author:   "@Ladicle",
				Category: "program",
				Tags:     []string{"hugo", "go", "OGP"},
				Date:     mustParseRFC3339(t, "2020-06-21T03:56:24+09:00"),
			},
		},
		{
			desc: "Parse TOML front matter",
			input: `+++
title = "HugoでもTwitterCardを自動生成したい"
author = ["@Ladicle"]
date = "2020-06-21T03:56:24+09:00"
tags = ["hugo", "go", "OGP"]
categories = ["program"]
+++
content`,
			expectFM: &FrontMatter{
				Title:    "HugoでもTwitterCardを自動生成したい",
				Author:   "@Ladicle",
				Category: "program",
				Tags:     []string{"hugo", "go", "OGP"},
				Date:     mustParseRFC3339(t, "2020-06-21T03:56:24+09:00"),
			},
		},
		{
			desc:      "Failed to parse empty file",
			expectErr: NewFMNotExistError(fmTitle),
		},
		{
			desc: "Failed to parse invalid front matter",
			input: `---
title = "invalid format'
---`,
			expectErr: errors.New("failed to unmarshal YAML: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `title =...` into map[string]interface {}"),
		},
		{
			desc: "Title is missing",
			input: `+++
author = ["@Ladicle"]
+++`,
			expectErr: NewFMNotExistError(fmTitle),
		},
		{
			desc: "Title is empty",
			input: `+++
title = ""
+++`,
			expectErr: NewFMNotExistError(fmTitle),
		},
		{
			desc: "Author is missing",
			input: `+++
title = "Title"
+++`,
			expectErr: NewFMNotExistError(fmAuthor),
		},
		{
			desc: "Category is empty",
			input: `+++
title = "Title"
author = ["@Ladicle"]
categories = [""]
+++`,
			expectErr: NewFMNotExistError(fmCategories),
		},
		{
			desc: "Tag is missing",
			input: `+++
title = "Title"
author = ["@Ladicle"]
categories = ["Program"]
+++`,
			expectErr: NewFMNotExistError(fmTags),
		},
		{
			desc: "Time is missing",
			input: `+++
title = "Title"
author = ["@Ladicle"]
categories = ["cat11"]
tags = ["tag1"]
+++`,
			expectErr: errors.New("\"date, lastmod, publishDate\" is not defined or empty"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			fm, err := parseFrontMatter(r)
			if err != nil {
				if tc.expectErr != nil {
					if tc.expectErr.Error() == err.Error() {
						return
					}
					t.Fatalf("parseFrontMatter() returns unexpected error: got=%#+v, want=%#+v",
						err, tc.expectErr)
				}
				t.Fatalf("failed to parse front matter: %v", err)
			}
			if tc.expectErr != nil {
				t.Fatalf("expect to occur %+v error but it didn't", tc.expectErr)
			}
			if !reflect.DeepEqual(fm, tc.expectFM) {
				t.Fatalf("parseFrontMatter() returns unexpected value: got=%#+v, want=%#+v",
					*fm, *tc.expectFM)
			}
		})
	}
}

func mustParseRFC3339(t *testing.T, timeStr string) time.Time {
	tt, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Fatal(err)
	}
	return tt
}
