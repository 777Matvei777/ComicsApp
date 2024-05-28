package xkcd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	ctx := context.Background()
	exist := make(map[int]bool)
	for i := 2; i < 2920; i++ {
		exist[i] = true
	}
	results := Parse("https://xkcd.com", 1, ctx, 1, exist)
	assert.Equal(t, 1, results[0].Id)
	assert.Equal(t, "Don't we all.", results[0].Alt)
	assert.Equal(t, "[[A boy sits in a barrel which is floating in an ocean.]]\nBoy: I wonder where I'll float next?\n[[The barrel drifts into the distance. Nothing else can be seen.]]\n{{Alt: Don't we all.}}", results[0].Transcript)
	assert.Equal(t, "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg", results[0].Url)
}
