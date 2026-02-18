// Package postgres
package postgres

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func BuildPaginationURL(baseURL, cursor, direction string, limit int) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	q := u.Query()
	q.Set("limit", fmt.Sprintf("%d", limit))
	switch direction {
	case "next":
		q.Set("after", cursor)
		q.Del("before")
	case "prev":
		q.Set("before", cursor)
		q.Del("after")
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func GetBaseULR(fullURL string) string {
	u, err := url.Parse(fullURL)
	if err != nil {
		return fullURL
	}
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

type CursorData struct {
	Timestamp time.Time
	ID        string
}

func EncodeCursor(timestamp time.Time, id string) string {
	cursorStr := fmt.Sprintf("%d_%s", timestamp.Unix(), id)
	return base64.StdEncoding.EncodeToString([]byte(cursorStr))
}

func DecodeCursor(cursor string) (*CursorData, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor format %w", err)
	}
	parts := string(decoded)
	split := strings.Split(parts, "_")
	if len(split) != 2 {
		return nil, fmt.Errorf("invalid cursor structure")
	}
	timestamp, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp ir cursor %w", err)
	}
	return &CursorData{
		Timestamp: time.Unix(timestamp, 0),
		ID:        split[1],
	}, nil
}
