package plancatalog

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

const (
	scalePathCount     = 100_000
	scaleMarkdownCount = 10_000
	scaleCandidateStep = 100
	maxRetainedHeap    = 64 << 20
)

func BenchmarkDiscoveryV2Classifier100kPaths10kMarkdown(b *testing.B) {
	blobs, candidateContent := scaleClassifierFixture()
	ordinaryContent := []byte("# Ordinary tracked documentation\n\n" + strings.Repeat("bounded inert content\n", 32))
	policy := DiscoveryPolicy{
		Version:               DiscoveryV2,
		ExcludedRoots:         []string{},
		AmbiguitySegments:     append([]string{}, conventionalAmbiguitySegments...),
		AuthoritativeUniverse: "git-regular-blobs",
	}
	expectedCandidates := scaleMarkdownCount / scaleCandidateStep
	expectedDigest := ""
	maxRetained := int64(0)

	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		readerCalls := 0
		runtime.GC()
		var before runtime.MemStats
		runtime.ReadMemStats(&before)
		classification, err := ClassifyGitBlobs(blobs, func(blob GitBlob) ([]byte, error) {
			readerCalls++
			if content, exists := candidateContent[blob.Path]; exists {
				return content, nil
			}
			return ordinaryContent, nil
		}, policy, ModeCanonical, "scale", nil)
		if err != nil {
			b.Fatal(err)
		}
		var after runtime.MemStats
		runtime.ReadMemStats(&after)
		retained := int64(after.HeapAlloc) - int64(before.HeapAlloc)
		if retained > maxRetained {
			maxRetained = retained
		}
		if retained > maxRetainedHeap {
			b.Fatalf("retained heap %d exceeds bounded %d-byte classifier budget", retained, maxRetainedHeap)
		}
		if readerCalls != scaleMarkdownCount {
			b.Fatalf("source reads=%d want=%d", readerCalls, scaleMarkdownCount)
		}
		if len(classification.Candidates) != expectedCandidates || len(classification.Discovery.Records) != expectedCandidates || !classification.Complete || HasErrors(classification.Discovery.Problems) {
			b.Fatalf("scale classification candidates=%d records=%d complete=%t problems=%#v", len(classification.Candidates), len(classification.Discovery.Records), classification.Complete, classification.Discovery.Problems)
		}
		if expectedDigest == "" {
			expectedDigest = classification.CandidateSetDigest
		} else if classification.CandidateSetDigest != expectedDigest {
			b.Fatalf("candidate digest changed: got %s want %s", classification.CandidateSetDigest, expectedDigest)
		}
	}
	b.StopTimer()
	b.ReportMetric(scalePathCount, "paths/op")
	b.ReportMetric(scaleMarkdownCount, "markdown/op")
	b.ReportMetric(0, "git-processes/op")
	b.ReportMetric(float64(maxRetained), "retained-heap-bytes/op")
	b.Logf("candidate_set_digest=%s", expectedDigest)
}

func scaleClassifierFixture() ([]GitBlob, map[string][]byte) {
	blobs := make([]GitBlob, 0, scalePathCount)
	candidateContent := make(map[string][]byte, scaleMarkdownCount/scaleCandidateStep)
	for index := 0; index < scalePathCount; index++ {
		objectID := fmt.Sprintf("%040x", index+1)
		if index < scaleMarkdownCount {
			path := fmt.Sprintf("products/scale/docs/reference/ordinary-%05d.md", index)
			if index%scaleCandidateStep == 0 {
				path = fmt.Sprintf("products/scale/docs/plans/scale-%05d_discovery_doc.md", index)
				candidateContent[path] = []byte(canonicalClassifierDocument(
					fmt.Sprintf("scale.discovery.plan-%05d", index),
					KindDiscovery,
					fmt.Sprintf("Scale Discovery %05d", index),
				))
			}
			blobs = append(blobs, GitBlob{Path: path, Mode: GitModeRegular, ObjectID: objectID})
			continue
		}
		blobs = append(blobs, GitBlob{Path: fmt.Sprintf("src/scale/path-%05d.txt", index), Mode: GitModeRegular, ObjectID: objectID})
	}
	return blobs, candidateContent
}
