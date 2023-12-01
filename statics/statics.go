package statics

import "embed"

//go:embed css/*css
//go:embed images/*
//go:embed fonts/*
//go:embed views/*
//go:embed javascript/*js

var StaticsFs embed.FS
