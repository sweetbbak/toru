package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/sweetbbak/toru/pkg/nyaa"
)

var LeechColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#E75B6B"))
var SeedColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#0BEB8D"))

func FormatMedia(m nyaa.Media) string {
	return fmt.Sprintf("%s\n%s\nDownloads: %d\n%s\nSize: %v\n",
		m.Name,
		m.Date.Format(time.DateTime),
		m.Downloads,
		formatPeers(m),
		humanize.Bytes(m.Size),
	)
}

func formatPeers(m nyaa.Media) string {
	s := SeedColor.Render(fmt.Sprintf("%d", m.Seeders))
	l := LeechColor.Render(fmt.Sprintf("%d", m.Leechers))
	return fmt.Sprintf("[%s|%s]", s, l)
}

// accepts any int format and returns a human readable string of the size in bytes
func formatSize(num interface{}) string {
	i := uint64(num.(uint64))
	return humanize.Bytes(i)
}
