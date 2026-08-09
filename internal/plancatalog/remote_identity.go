package plancatalog

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"path/filepath"
	"strings"
)

// uniqueTrackingRemoteForSource maps a fresh mirror source to one configured
// tracking remote without serializing the potentially sensitive remote URL.
// Missing or ambiguous mappings deliberately leave tracking and mirror
// observations separate.
func uniqueTrackingRemoteForSource(root, sourceIdentitySHA256 string) (string, bool) {
	if !validSHA256(sourceIdentitySHA256) {
		return "", false
	}
	remotes, err := gitText(root, "remote")
	if err != nil {
		return "", false
	}
	matches := []string{}
	for _, remote := range strings.Fields(remotes) {
		remoteURL, remoteErr := gitText(root, "remote", "get-url", remote)
		if remoteErr != nil {
			continue
		}
		if remoteIdentitySHA256(root, strings.TrimSpace(remoteURL)) == sourceIdentitySHA256 {
			matches = append(matches, remote)
		}
	}
	if len(matches) != 1 {
		return "", false
	}
	return matches[0], true
}

func remoteIdentitySHA256(root, value string) string {
	identity := canonicalRemoteIdentity(root, value)
	if identity == "" {
		return ""
	}
	digest := sha256.Sum256([]byte(identity))
	return hex.EncodeToString(digest[:])
}

func canonicalRemoteIdentity(root, value string) string {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return ""
	}
	if !strings.Contains(raw, "://") {
		if at := strings.LastIndex(raw, "@"); at >= 0 {
			if colon := strings.Index(raw[at+1:], ":"); colon >= 0 {
				hostStart := at + 1
				hostEnd := hostStart + colon
				return canonicalRemoteHostPath(raw[hostStart:hostEnd], raw[hostEnd+1:])
			}
		}
	}
	parsed, err := url.Parse(raw)
	if err == nil && parsed.Scheme != "" {
		if parsed.Host != "" {
			host := parsed.Hostname()
			if parsed.Port() != "" {
				host += ":" + parsed.Port()
			}
			return canonicalRemoteHostPath(host, parsed.EscapedPath())
		}
		if parsed.Scheme == "file" {
			return "file:" + filepath.Clean(parsed.Path)
		}
	}
	if !filepath.IsAbs(raw) && !strings.Contains(raw, "://") {
		if colon := strings.Index(raw, ":"); colon <= 0 || !strings.Contains(raw[:colon], "@") {
			raw = filepath.Clean(filepath.Join(root, raw))
		}
	}
	if filepath.IsAbs(raw) {
		return "file:" + filepath.Clean(raw)
	}
	return strings.TrimSuffix(raw, ".git")
}

func canonicalRemoteHostPath(host, repositoryPath string) string {
	host = strings.ToLower(strings.TrimSpace(host))
	repositoryPath = strings.Trim(strings.TrimSpace(repositoryPath), "/")
	repositoryPath = strings.TrimSuffix(repositoryPath, ".git")
	if host == "github.com" {
		repositoryPath = strings.ToLower(repositoryPath)
	}
	if host == "" || repositoryPath == "" {
		return ""
	}
	return host + "/" + repositoryPath
}
