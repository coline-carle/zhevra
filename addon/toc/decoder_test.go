package toc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
