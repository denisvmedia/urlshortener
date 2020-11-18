package model

import (
	"encoding/base64"
	"github.com/google/uuid"
	"strings"
)

func generateShortName() string {
	hash, _ := uuid.New().MarshalBinary()
	var b strings.Builder
	encoder := base64.NewEncoder(base64.URLEncoding, &b)
	_, _ = encoder.Write(hash)
	_ = encoder.Close()

	// youtube-like name
	res := strings.ReplaceAll(strings.Trim(b.String(), "_-=\n"), "_", "")
	if len(res) > 8 {
		res = res[:8]
	}

	return res
}
