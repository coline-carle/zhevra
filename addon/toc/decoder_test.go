package toc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stripTests = []struct {
	Name   string
	Input  string
	Output string
}{
	{
		Name:   "Simple String with no tag",
		Input:  "Simple String with no tag",
		Output: "Simple String with no tag",
	},
	{
		Name:   "String with a tag",
		Input:  "non tagged.|cFFFFFFFFTagged|r.non tagged",
		Output: "non tagged.Tagged.non tagged",
	},
	{
		Name:   "Two tags",
		Input:  "notag.|cFF000000tag1|cFFFFFFFFTag2|rTag1F|r.notag",
		Output: "notag.tag1Tag2Tag1F.notag",
	},
	{
		Name:   "non matching tags",
		Input:  "non tagged.|cFFFFFFFFTagged|x.non tagged",
		Output: "non tagged.|cFFFFFFFFTagged|x.non tagged",
	},
}

func TestStrip(t *testing.T) {
	for _, tt := range stripTests {
		out := StripColorTags(tt.Input)
		if out != tt.Output {
			t.Errorf("%s: out=%s want %s", tt.Name, out, tt.Output)
		}
	}
}

var bigwigsStr = `
## Interface: 70300

## Title: BigWigs [|cffeda55fCore|r]
## Title-zhCN: BigWigs [|cffeda55f核心|r]
## Title-zhTW: BigWigs [|cffeda55f核心|r]
## Notes: BigWigs Core
## Notes-zhCN: BigWigs 核心模块。
## Notes-zhTW: BigWigs 核心模組。

## LoadOnDemand: 1
## Dependencies: BigWigs

libs.xml
modules.xml
`

type BigWigsInfo struct {
	Interface int64
	Title     string
	TitleZhCN string `toc:"Title-zhCN"`

	Notes        string
	NotesZhCN    string `toc:"Notes-zhCN"`
	LoadOnDemand int
	WhatElse     string
}

var expected = BigWigsInfo{
	Title:        "BigWigs [|cffeda55fCore|r]",
	TitleZhCN:    "BigWigs [|cffeda55f核心|r]",
	NotesZhCN:    "BigWigs 核心模块。",
	Notes:        "BigWigs Core",
	LoadOnDemand: 1,
	Interface:    70300,
}

func TestDecode(t *testing.T) {
	r := strings.NewReader(bigwigsStr)
	var bwInfo BigWigsInfo
	decoder := NewDecoder(r)
	decoder.Decode(&bwInfo)
	assert.Equal(t, expected, bwInfo)
}
