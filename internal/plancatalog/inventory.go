package plancatalog

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/state"
)

type InventoryRecord struct {
	Path                    string                    `json:"path" yaml:"path"`
	Kind                    Kind                      `json:"kind" yaml:"kind"`
	Title                   string                    `json:"title" yaml:"title"`
	Lifecycle               Lifecycle                 `json:"lifecycle" yaml:"lifecycle"`
	SourceRevision          string                    `json:"source_revision" yaml:"source_revision"`
	SourceSHA256            string                    `json:"source_sha256" yaml:"source_sha256"`
	MetadataCoverage        string                    `json:"metadata_coverage" yaml:"metadata_coverage"`
	ParseStatus             string                    `json:"parse_status" yaml:"parse_status"`
	ProblemCodes            []string                  `json:"problem_codes" yaml:"problem_codes"`
	PlanID                  string                    `json:"plan_id,omitempty" yaml:"plan_id,omitempty"`
	LegacyEvidence          []LegacyEntry             `json:"legacy_evidence" yaml:"legacy_evidence"`
	LegacyAmbiguous         bool                      `json:"legacy_ambiguous" yaml:"legacy_ambiguous"`
	ActiveBranchTouch       []string                  `json:"-" yaml:"-"`
	LegacyActiveBranchTouch []string                  `json:"active_branch_touches,omitempty" yaml:"active_branch_touches,omitempty"`
	BranchTouchCandidates   []BranchCandidateEvidence `json:"branch_touch_candidates" yaml:"branch_touch_candidates"`
	DirtyOverlap            bool                      `json:"dirty_overlap" yaml:"dirty_overlap"`
	Signal                  CandidateSignal           `json:"signal,omitempty" yaml:"signal,omitempty"`
	Ownership               OwnershipClass            `json:"ownership,omitempty" yaml:"ownership,omitempty"`
	Provenance              *CandidateProvenance      `json:"provenance,omitempty" yaml:"provenance,omitempty"`
	Authoritative           *bool                     `json:"authoritative,omitempty" yaml:"authoritative,omitempty"`
	Preview                 bool                      `json:"preview,omitempty" yaml:"preview,omitempty"`
}

type InventoryCoverage struct {
	FormalRecords             int `json:"formal_records" yaml:"formal_records"`
	CanonicalMetadata         int `json:"canonical_metadata" yaml:"canonical_metadata"`
	LegacyRecords             int `json:"legacy_records" yaml:"legacy_records"`
	InvalidRecords            int `json:"invalid_records" yaml:"invalid_records"`
	UnreadableRecords         int `json:"unreadable_records" yaml:"unreadable_records"`
	UnsafeRecords             int `json:"unsafe_records" yaml:"unsafe_records"`
	UnpairedLegacyEvidence    int `json:"unpaired_legacy_evidence" yaml:"unpaired_legacy_evidence"`
	ActiveBranchTouches       int `json:"-" yaml:"-"`
	LegacyActiveBranchTouches int `json:"active_branch_touches,omitempty" yaml:"active_branch_touches,omitempty"`
	BranchTouchCandidates     int `json:"branch_touch_candidates" yaml:"branch_touch_candidates"`
	BlockingBranchTouches     int `json:"blocking_branch_touches" yaml:"blocking_branch_touches"`
	DirtyOverlaps             int `json:"dirty_overlaps" yaml:"dirty_overlaps"`
	Candidates                int `json:"candidates" yaml:"candidates"`
	IncludedCandidates        int `json:"included_candidates" yaml:"included_candidates"`
	ExcludedCandidates        int `json:"excluded_candidates" yaml:"excluded_candidates"`
	BlockedCandidates         int `json:"blocked_candidates" yaml:"blocked_candidates"`
	UnownedCandidates         int `json:"unowned_candidates" yaml:"unowned_candidates"`
	PreviewCandidates         int `json:"preview_candidates,omitempty" yaml:"preview_candidates,omitempty"`
}

type Inventory struct {
	SchemaVersion                  int                     `json:"schema_version" yaml:"schema_version"`
	RepositoryID                   string                  `json:"repository_id,omitempty" yaml:"repository_id,omitempty"`
	CatalogMode                    CatalogMode             `json:"catalog_mode" yaml:"catalog_mode"`
	SourceRevision                 string                  `json:"source_revision" yaml:"source_revision"`
	EvidenceRevision               string                  `json:"evidence_revision,omitempty" yaml:"evidence_revision,omitempty"`
	GeneratedAt                    string                  `json:"generated_at" yaml:"generated_at"`
	Records                        []InventoryRecord       `json:"records" yaml:"records"`
	UnpairedLegacy                 []LegacyEntry           `json:"unpaired_legacy_evidence" yaml:"unpaired_legacy_evidence"`
	Coverage                       InventoryCoverage       `json:"coverage" yaml:"coverage"`
	Problems                       []Problem               `json:"problems" yaml:"problems"`
	RemoteOverlays                 []RemoteOverlayEvidence `json:"remote_overlays,omitempty" yaml:"remote_overlays,omitempty"`
	DiscoveryVersion               DiscoveryVersion        `json:"discovery_version,omitempty" yaml:"discovery_version,omitempty"`
	ConfiguredDiscoveryVersion     DiscoveryVersion        `json:"configured_discovery_version,omitempty" yaml:"configured_discovery_version,omitempty"`
	ConfiguredCatalogMode          CatalogMode             `json:"configured_catalog_mode,omitempty" yaml:"configured_catalog_mode,omitempty"`
	TargetCatalogMode              CatalogMode             `json:"target_catalog_mode,omitempty" yaml:"target_catalog_mode,omitempty"`
	PolicyDigest                   string                  `json:"policy_digest,omitempty" yaml:"policy_digest,omitempty"`
	CandidateSetDigest             string                  `json:"candidate_set_digest,omitempty" yaml:"candidate_set_digest,omitempty"`
	TargetConfigPreconditionSHA256 string                  `json:"target_config_precondition_sha256,omitempty" yaml:"target_config_precondition_sha256,omitempty"`
	Complete                       *bool                   `json:"complete,omitempty" yaml:"complete,omitempty"`
	MixedCoverageComplete          *bool                   `json:"mixed_coverage_complete,omitempty" yaml:"mixed_coverage_complete,omitempty"`
	CanonicalReady                 *bool                   `json:"canonical_ready,omitempty" yaml:"canonical_ready,omitempty"`
	BranchEvidence                 *BranchEvidenceSummary  `json:"branch_evidence,omitempty" yaml:"branch_evidence,omitempty"`
	Candidates                     []Candidate             `json:"candidates" yaml:"candidates"`
	PreviewCandidates              []Candidate             `json:"preview_candidates,omitempty" yaml:"preview_candidates,omitempty"`
	PreviewProblems                []Problem               `json:"preview_problems,omitempty" yaml:"preview_problems,omitempty"`
	PreviewValid                   *bool                   `json:"preview_valid,omitempty" yaml:"preview_valid,omitempty"`
}

type RemoteOverlayRef struct {
	LogicalRef string `json:"logical_ref" yaml:"logical_ref"`
	Tip        string `json:"tip" yaml:"tip"`
}

type RemoteOverlayEvidence struct {
	Remote                string                    `json:"remote" yaml:"remote"`
	SourceIdentitySHA256  string                    `json:"source_identity_sha256,omitempty" yaml:"source_identity_sha256,omitempty"`
	Status                string                    `json:"status" yaml:"status"`
	ObservedAt            string                    `json:"observed_at" yaml:"observed_at"`
	Digest                string                    `json:"digest,omitempty" yaml:"digest,omitempty"`
	Refs                  []RemoteOverlayRef        `json:"refs" yaml:"refs"`
	BranchTouchCandidates []BranchCandidateEvidence `json:"branch_touch_candidates" yaml:"branch_touch_candidates"`
}

// CanonicalRemoteOverlayDigest binds the selected public-safe source identity,
// completeness state, exact ref tips, and portable branch evidence. Diagnostic
// timestamps and Git version strings are deliberately excluded.
func CanonicalRemoteOverlayDigest(overlay RemoteOverlayEvidence) string {
	refs := append([]RemoteOverlayRef{}, overlay.Refs...)
	sort.SliceStable(refs, func(i, j int) bool {
		if refs[i].LogicalRef != refs[j].LogicalRef {
			return refs[i].LogicalRef < refs[j].LogicalRef
		}
		return refs[i].Tip < refs[j].Tip
	})
	uniqueRefs := refs[:0]
	for _, ref := range refs {
		if len(uniqueRefs) == 0 || uniqueRefs[len(uniqueRefs)-1] != ref {
			uniqueRefs = append(uniqueRefs, ref)
		}
	}
	evidence := append([]BranchCandidateEvidence{}, overlay.BranchTouchCandidates...)
	for index := range evidence {
		evidence[index].GitVersion = ""
		evidence[index].Digest = CanonicalBranchEvidenceDigest(evidence[index])
	}
	sort.SliceStable(evidence, func(i, j int) bool {
		left := evidence[i].Identity.LogicalRef + "\x00" + evidence[i].Identity.CandidatePath
		right := evidence[j].Identity.LogicalRef + "\x00" + evidence[j].Identity.CandidatePath
		return left < right
	})
	canonical := struct {
		Remote                string                    `json:"remote"`
		SourceIdentitySHA256  string                    `json:"source_identity_sha256"`
		Status                string                    `json:"status"`
		Refs                  []RemoteOverlayRef        `json:"refs"`
		BranchTouchCandidates []BranchCandidateEvidence `json:"branch_touch_candidates"`
	}{overlay.Remote, overlay.SourceIdentitySHA256, overlay.Status, uniqueRefs, evidence}
	return canonicalJSONDigest(canonical)
}

func AttachRemoteOverlay(root string, inventory *Inventory, overlay RemoteOverlayEvidence) {
	trackingRemote, trackingRemoteMatched := uniqueTrackingRemoteForSource(root, overlay.SourceIdentitySHA256)
	remoteEvidenceBlocked := false
	for recordIndex := range inventory.Records {
		for evidenceIndex := range inventory.Records[recordIndex].BranchTouchCandidates {
			evidence := &inventory.Records[recordIndex].BranchTouchCandidates[evidenceIndex]
			evidence.Identity.EvidenceScope = "remote-aware"
			evidence.Digest = CanonicalBranchEvidenceDigest(*evidence)
		}
	}
	for _, evidence := range overlay.BranchTouchCandidates {
		if blocker := remoteOverlayEvidenceIdentityBlocker(overlay, evidence); blocker != nil {
			evidence.ProposedDisposition = BranchDispositionBlocking
			evidence.Proof = nil
			evidence.OwnerTipCandidate = nil
			evidence.Blockers = uniqueBranchBlockers(append(evidence.Blockers, *blocker))
			evidence.Digest = CanonicalBranchEvidenceDigest(evidence)
			remoteEvidenceBlocked = true
		}
		matchedRecord := false
		for recordIndex := range inventory.Records {
			if inventory.Records[recordIndex].Path != evidence.Identity.CandidatePath {
				continue
			}
			matchedRecord = true
			coalesced := false
			trackingLogical := trackingLogicalRefForRemote(evidence.Identity.LogicalRef, overlay.Remote, trackingRemote, trackingRemoteMatched)
			if trackingLogical != "" {
				for existingIndex := range inventory.Records[recordIndex].BranchTouchCandidates {
					existing := &inventory.Records[recordIndex].BranchTouchCandidates[existingIndex]
					if existing.Identity.LogicalRef != trackingLogical {
						continue
					}
					if EquivalentBranchCandidateEvidence(*existing, evidence) {
						inventory.Records[recordIndex].BranchTouchCandidates[existingIndex] = evidence
						coalesced = true
						break
					}
					if existing.ProposedDisposition != BranchDispositionBlocking {
						inventory.Coverage.BlockingBranchTouches++
					}
					existing.ProposedDisposition = BranchDispositionBlocking
					existing.Proof = nil
					existing.OwnerTipCandidate = nil
					existing.Blockers = append(existing.Blockers, BranchEvidenceBlocker{Code: "branch_tracking_overlay_mismatch", Message: "tracking ref and fresh mirror evidence differ", Path: evidence.Identity.CandidatePath, Ref: existing.Identity.LogicalRef})
					existing.Digest = CanonicalBranchEvidenceDigest(*existing)
					evidence.ProposedDisposition = BranchDispositionBlocking
					evidence.Proof = nil
					evidence.OwnerTipCandidate = nil
					evidence.Blockers = append(evidence.Blockers, BranchEvidenceBlocker{Code: "branch_tracking_overlay_mismatch", Message: "fresh mirror and tracking ref evidence differ", Path: evidence.Identity.CandidatePath, Ref: evidence.Identity.LogicalRef})
					evidence.Digest = CanonicalBranchEvidenceDigest(evidence)
					break
				}
			}
			if coalesced {
				break
			}
			inventory.Records[recordIndex].BranchTouchCandidates = append(inventory.Records[recordIndex].BranchTouchCandidates, evidence)
			inventory.Coverage.BranchTouchCandidates++
			if evidence.ProposedDisposition == BranchDispositionBlocking {
				inventory.Coverage.BlockingBranchTouches++
			}
			break
		}
		if !matchedRecord {
			remoteEvidenceBlocked = true
			inventory.Problems = append(inventory.Problems, Problem{Code: "remote_overlay_candidate_unmatched", Message: "remote branch evidence names a candidate outside the reviewed inventory", Path: evidence.Identity.CandidatePath, Severity: SeverityError})
		}
	}
	inventory.RemoteOverlays = []RemoteOverlayEvidence{overlay}
	if inventory.BranchEvidence == nil {
		return
	}
	inventory.BranchEvidence.EvidenceScope = "remote-aware"
	inventory.BranchEvidence.RemoteOverlayStatus = overlay.Status
	inventory.BranchEvidence.RemoteOverlayDigest = overlay.Digest
	inventory.BranchEvidence.RemoteSourceIdentitySHA256 = overlay.SourceIdentitySHA256
	evidence := []BranchCandidateEvidence{}
	for _, record := range inventory.Records {
		evidence = append(evidence, record.BranchTouchCandidates...)
	}
	inventory.BranchEvidence.Digest = branchEvidenceSetDigest(evidence, overlay.Digest)
	if overlay.Status != "complete" || remoteEvidenceBlocked {
		if inventory.Complete != nil {
			*inventory.Complete = false
		}
		if inventory.MixedCoverageComplete != nil {
			*inventory.MixedCoverageComplete = false
		}
		if inventory.CanonicalReady != nil {
			*inventory.CanonicalReady = false
		}
	}
}

func remoteOverlayEvidenceIdentityBlocker(overlay RemoteOverlayEvidence, evidence BranchCandidateEvidence) *BranchEvidenceBlocker {
	identity := evidence.Identity
	expectedPrefix := "remote:" + overlay.Remote + ":refs/heads/"
	if overlay.Remote == "" || identity.RepositoryID != overlay.Remote || identity.EvidenceScope != "remote-aware" ||
		identity.CandidatePath == "" || !strings.HasPrefix(identity.LogicalRef, expectedPrefix) || strings.TrimPrefix(identity.LogicalRef, expectedPrefix) == "" {
		return &BranchEvidenceBlocker{Code: "remote_overlay_identity_mismatch", Message: "remote branch evidence does not match the selected source, scope, ref, and candidate binding", Path: identity.CandidatePath, Ref: identity.LogicalRef}
	}
	if !validSHA256(evidence.Digest) || CanonicalBranchEvidenceDigest(evidence) != evidence.Digest {
		return &BranchEvidenceBlocker{Code: "remote_overlay_evidence_tampered", Message: "remote branch evidence digest does not match its structured content", Path: identity.CandidatePath, Ref: identity.LogicalRef}
	}
	return nil
}

func trackingLogicalRefForRemote(remoteLogical, repositoryID, remoteName string, matched bool) string {
	if !matched || remoteName == "" {
		return ""
	}
	prefix := "remote:" + repositoryID + ":refs/heads/"
	if !strings.HasPrefix(remoteLogical, prefix) {
		return ""
	}
	branch := strings.TrimPrefix(remoteLogical, prefix)
	if branch == "" {
		return ""
	}
	return "tracking:" + remoteName + ":refs/heads/" + branch
}

func BuildInventory(root string, now time.Time) (Inventory, error) {
	return BuildInventoryWithOptions(root, now, SnapshotOptions{})
}

func BuildInventoryWithOptions(root string, now time.Time, options SnapshotOptions) (Inventory, error) {
	snapshot, err := LoadRepositorySnapshotWithOptions(root, options)
	if err != nil {
		return Inventory{}, err
	}
	revision, err := gitText(root, "rev-parse", "--verify", "HEAD")
	if err != nil {
		return Inventory{}, fmt.Errorf("git_revision_unavailable: %w", err)
	}
	touches, touchProblems := activeBranchTouchesWithSettings(root, snapshot.Settings)
	candidates := append([]Candidate{}, snapshot.Candidates...)
	if snapshot.Settings.DiscoveryVersion != DiscoveryV2 {
		candidates, err = Enumerate(root)
		if err != nil {
			return Inventory{}, err
		}
	}
	allCandidates := append([]Candidate{}, candidates...)
	allCandidates = append(allCandidates, snapshot.PreviewCandidates...)
	if now.IsZero() {
		now = time.Now()
	}
	inventory := Inventory{
		SchemaVersion:  1,
		RepositoryID:   snapshot.Settings.RepositoryID,
		CatalogMode:    snapshot.Settings.Mode,
		SourceRevision: revision,
		GeneratedAt:    now.UTC().Truncate(time.Second).Format(time.RFC3339),
		Records:        []InventoryRecord{},
		UnpairedLegacy: append([]LegacyEntry{}, snapshot.Reconciliation.Unpaired...),
		Problems:       append([]Problem{}, snapshot.Problems...),
	}
	if snapshot.Settings.DiscoveryVersion == DiscoveryV2 || snapshot.Targeted {
		complete := snapshot.Complete
		configSHA, configErr := regularSourceDigest(root, state.ConfigPath)
		if configErr != nil || len(configSHA) != 64 {
			complete = false
			message := "target config precondition is unavailable"
			if configErr != nil {
				message = configErr.Error()
			}
			inventory.Problems = append(inventory.Problems, Problem{Code: "target_config_precondition_unavailable", Message: message, Path: state.ConfigPath, Severity: SeverityError})
		}
		inventory.SchemaVersion = 3
		inventory.EvidenceRevision = revision
		inventory.DiscoveryVersion = snapshot.Settings.DiscoveryVersion
		inventory.ConfiguredDiscoveryVersion = snapshot.ConfiguredSettings.DiscoveryVersion
		inventory.ConfiguredCatalogMode = snapshot.ConfiguredSettings.Mode
		inventory.TargetCatalogMode = snapshot.Settings.Mode
		inventory.PolicyDigest = snapshot.PolicyDigest
		inventory.CandidateSetDigest = snapshot.CandidateSetDigest
		inventory.TargetConfigPreconditionSHA256 = configSHA
		inventory.Complete = &complete
		mixedCoverageComplete := complete
		canonicalReady := complete
		inventory.MixedCoverageComplete = &mixedCoverageComplete
		inventory.CanonicalReady = &canonicalReady
		inventory.Candidates = append([]Candidate{}, candidates...)
		inventory.PreviewCandidates = append([]Candidate{}, snapshot.PreviewCandidates...)
		inventory.PreviewProblems = append([]Problem{}, snapshot.PreviewProblems...)
		if options.IncludeUntracked {
			previewValid := !HasErrors(snapshot.PreviewProblems)
			inventory.PreviewValid = &previewValid
		}
	}
	if inventory.SchemaVersion == 1 {
		inventory.UnpairedLegacy = legacyEntriesForInventoryV1(inventory.UnpairedLegacy)
	}
	inventory.Problems = append(inventory.Problems, touchProblems...)
	if inventory.Complete != nil && HasErrors(touchProblems) {
		complete := false
		inventory.Complete = &complete
	}
	recordsByPath := map[string]Record{}
	for _, record := range snapshot.Records {
		recordsByPath[record.Path] = record
	}
	for _, record := range snapshot.PreviewRecords {
		recordsByPath[record.Path] = record
	}
	problemCodesByPath := map[string][]string{}
	for _, problem := range inventory.Problems {
		if problem.Path != "" {
			problemCodesByPath[problem.Path] = append(problemCodesByPath[problem.Path], problem.Code)
		}
	}
	for _, problem := range inventory.PreviewProblems {
		if problem.Path != "" {
			problemCodesByPath[problem.Path] = append(problemCodesByPath[problem.Path], problem.Code)
		}
	}
	allBranchEvidence := []BranchCandidateEvidence{}
	for _, candidate := range allCandidates {
		record, parsed := recordsByPath[candidate.Path]
		matches := append([]LegacyEntry{}, snapshot.Reconciliation.ByCanonicalPath[candidate.Path]...)
		if inventory.SchemaVersion == 1 {
			matches = legacyEntriesForInventoryV1(matches)
		}
		coverage := "invalid"
		parseStatus := "invalid"
		planID := ""
		title := ""
		lifecycle := Lifecycle("")
		sourceSHA := ""
		preview := candidate.Provenance.Source.Preview
		authoritative := !preview && (snapshot.Settings.DiscoveryVersion != DiscoveryV2 || candidate.Ownership == OwnershipOwned)
		if preview {
			coverage = "preview"
			parseStatus = "preview"
			inventory.Coverage.PreviewCandidates++
		} else if snapshot.Settings.DiscoveryVersion == DiscoveryV2 && candidate.Ownership != OwnershipOwned {
			coverage = string(candidate.Ownership)
			parseStatus = string(candidate.Ownership)
			switch candidate.Ownership {
			case OwnershipExcluded:
				inventory.Coverage.ExcludedCandidates++
			case OwnershipProspectiveBlocked:
				inventory.Coverage.BlockedCandidates++
			case OwnershipHardUnowned:
				inventory.Coverage.UnownedCandidates++
			}
		} else if parsed && record.Metadata != nil {
			coverage = "canonical"
			parseStatus = "canonical"
			planID = record.Metadata.ID
			title = DisplayTitle(record, matches)
			lifecycle = record.Header.Lifecycle
			sourceSHA = record.ContentSHA256
			if authoritative {
				inventory.Coverage.CanonicalMetadata++
			}
		} else if parsed {
			coverage = "legacy"
			parseStatus = "legacy"
			title = DisplayTitle(record, matches)
			lifecycle = record.Header.Lifecycle
			sourceSHA = record.ContentSHA256
			if authoritative {
				inventory.Coverage.LegacyRecords++
			}
		} else if snapshot.Settings.DiscoveryVersion != DiscoveryV2 {
			data, readErr := readRegularSource(root, candidate.Path)
			if readErr != nil {
				if ErrorCode(readErr) == "source_unsafe" {
					coverage = "unsafe"
					parseStatus = "unsafe"
					inventory.Coverage.UnsafeRecords++
				} else {
					coverage = "unreadable"
					parseStatus = "unreadable"
					inventory.Coverage.UnreadableRecords++
				}
			} else {
				sourceSHA = sha256Text(data)
				inventory.Coverage.InvalidRecords++
			}
		} else {
			sourceSHA = candidate.Provenance.Source.ContentSHA256
			if authoritative {
				inventory.Coverage.InvalidRecords++
			}
		}
		if sourceSHA == "" {
			sourceSHA = candidate.Provenance.Source.ContentSHA256
		}
		if parsed {
			title = DisplayTitle(record, matches)
			lifecycle = record.Header.Lifecycle
			if record.Metadata != nil {
				planID = record.Metadata.ID
			}
		}
		dirty, dirtyErr := gitPathDirty(root, candidate.Path)
		if dirtyErr != nil {
			inventory.Problems = append(inventory.Problems, Problem{Code: "git_dirty_check_failed", Message: dirtyErr.Error(), Path: candidate.Path, Severity: SeverityError})
		} else if authoritative && !dirty && sourceSHA != "" {
			sourceBytes, sourceErr := gitBytes(root, "show", revision+":"+candidate.Path)
			if sourceErr != nil {
				inventory.Problems = append(inventory.Problems, Problem{Code: "source_revision_unavailable", Message: sourceErr.Error(), Path: candidate.Path, Severity: SeverityError, Remediation: "commit the reviewed plan or regenerate inventory from a coherent revision"})
			} else {
				sourceSHA = sha256Text(sourceBytes)
			}
		}
		branchTouches := append([]string{}, touches[candidate.Path]...)
		candidateBranchEvidence := []BranchCandidateEvidence{}
		if snapshot.Settings.DiscoveryVersion == DiscoveryV2 && authoritative && !preview {
			for _, ref := range branchTouches {
				paths, scopeBlocker := branchCandidateProofPaths(root, revision, ref, candidate.Path)
				var evidence BranchCandidateEvidence
				if scopeBlocker != nil {
					evidence = BranchCandidateEvidence{
						Algorithm:        BranchEvidenceAlgorithm,
						Identity:         BranchReviewIdentity{EvidenceScope: "local", RepositoryID: inventory.RepositoryID, LogicalRef: logicalRefForLocalTouch(ref), CandidatePath: candidate.Path},
						EvidenceRevision: revision, PolicyDigest: inventory.PolicyDigest, CandidateSetDigest: inventory.CandidateSetDigest,
						ProposedDisposition: BranchDispositionBlocking, Blockers: []BranchEvidenceBlocker{*scopeBlocker},
					}
					evidence.Digest = CanonicalBranchEvidenceDigest(evidence)
					inventory.Problems = append(inventory.Problems, Problem{Code: scopeBlocker.Code, Message: scopeBlocker.Message, Path: candidate.Path, Severity: SeverityError, Remediation: "regenerate inventory after repairing the exact ref/path evidence"})
					if inventory.Complete != nil {
						*inventory.Complete = false
					}
				} else {
					evidence = proposeBranchCandidateEvidence(BranchEvidenceRequest{
						Root: root, EvidenceScope: "local", RepositoryID: inventory.RepositoryID,
						LogicalRef: logicalRefForLocalTouch(ref), Ref: ref, EvidenceRevision: revision,
						PolicyDigest: inventory.PolicyDigest, CandidateSetDigest: inventory.CandidateSetDigest,
						CandidatePath: candidate.Path, Paths: paths, ExpectedKind: candidate.ExpectedKind,
					})
				}
				candidateBranchEvidence = append(candidateBranchEvidence, evidence)
				allBranchEvidence = append(allBranchEvidence, evidence)
			}
		}
		if dirty && authoritative {
			inventory.Coverage.DirtyOverlaps++
		}
		if len(branchTouches) > 0 && authoritative {
			inventory.Coverage.ActiveBranchTouches++
		}
		inventory.Coverage.BranchTouchCandidates += len(candidateBranchEvidence)
		for _, evidence := range candidateBranchEvidence {
			if evidence.ProposedDisposition == BranchDispositionBlocking {
				inventory.Coverage.BlockingBranchTouches++
			}
		}
		problemCodes := uniqueSortedStrings(problemCodesByPath[candidate.Path])
		// Filename/family qualification owns the human-readable record kind.
		// Metadata kind remains identity evidence and mismatches stay visible as
		// validation problems rather than rewriting this field.
		kind := candidate.ExpectedKind
		authority := authoritative
		recordRevision := revision
		if preview {
			recordRevision = ""
		}
		var provenance *CandidateProvenance
		if snapshot.Settings.DiscoveryVersion == DiscoveryV2 {
			copy := candidate.Provenance
			provenance = &copy
		}
		legacyBranchTouches := []string(nil)
		if inventory.SchemaVersion == 1 {
			legacyBranchTouches = append([]string{}, branchTouches...)
		}
		inventory.Records = append(inventory.Records, InventoryRecord{
			Path:                    candidate.Path,
			Kind:                    kind,
			Title:                   title,
			Lifecycle:               lifecycle,
			SourceRevision:          recordRevision,
			SourceSHA256:            sourceSHA,
			MetadataCoverage:        coverage,
			ParseStatus:             parseStatus,
			ProblemCodes:            problemCodes,
			PlanID:                  planID,
			LegacyEvidence:          matches,
			LegacyAmbiguous:         len(matches) > 1,
			ActiveBranchTouch:       branchTouches,
			LegacyActiveBranchTouch: legacyBranchTouches,
			BranchTouchCandidates:   candidateBranchEvidence,
			DirtyOverlap:            dirty && authoritative,
			Signal:                  candidate.Signal,
			Ownership:               candidate.Ownership,
			Provenance:              provenance,
			Authoritative: func() *bool {
				if snapshot.Settings.DiscoveryVersion == DiscoveryV2 {
					return &authority
				}
				return nil
			}(),
			Preview: preview,
		})
		if snapshot.Settings.DiscoveryVersion == DiscoveryV2 && !preview {
			inventory.Coverage.Candidates++
			if authoritative {
				inventory.Coverage.IncludedCandidates++
			}
		}
	}
	if snapshot.Settings.DiscoveryVersion == DiscoveryV2 {
		inventory.Coverage.FormalRecords = inventory.Coverage.IncludedCandidates
	} else {
		inventory.Coverage.FormalRecords = len(candidates)
	}
	inventory.Coverage.UnpairedLegacyEvidence = len(inventory.UnpairedLegacy)
	sort.SliceStable(inventory.Records, func(i, j int) bool { return inventory.Records[i].Path < inventory.Records[j].Path })
	SortProblems(inventory.Problems)
	if snapshot.Settings.DiscoveryVersion == DiscoveryV2 {
		gitVersion, blocker := CheckBranchEvidenceGitVersion(root)
		if blocker != nil {
			inventory.Problems = append(inventory.Problems, Problem{Code: blocker.Code, Message: blocker.Message, Severity: SeverityError})
			if inventory.Complete != nil {
				*inventory.Complete = false
			}
		}
		inventory.BranchEvidence = &BranchEvidenceSummary{Algorithm: BranchEvidenceAlgorithm, GitVersion: gitVersion, EvidenceScope: "local", RemoteOverlayStatus: "not-requested", Digest: branchEvidenceSetDigest(allBranchEvidence, "")}
		if inventory.MixedCoverageComplete != nil {
			*inventory.MixedCoverageComplete = inventory.Complete != nil && *inventory.Complete && len(allBranchEvidence) == 0
		}
		if inventory.CanonicalReady != nil {
			*inventory.CanonicalReady = inventory.MixedCoverageComplete != nil && *inventory.MixedCoverageComplete && inventory.Coverage.LegacyRecords == 0
		}
		SortProblems(inventory.Problems)
	}
	if inventory.SchemaVersion == 1 {
		inventory.Coverage.LegacyActiveBranchTouches = inventory.Coverage.ActiveBranchTouches
	}
	return inventory, nil
}

func uniqueSortedStrings(values []string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func branchCandidateProofPaths(root, evidenceRevision, ref, candidatePath string) ([]string, *BranchEvidenceBlocker) {
	data, err := gitBytes(root, "diff", "--name-status", "-z", "--find-renames", "--diff-filter=ACDMRT", evidenceRevision+"..."+ref)
	if err != nil {
		return nil, &BranchEvidenceBlocker{Code: "branch_evidence_scope_unavailable", Message: "candidate proof path scope could not be recomputed", Path: candidatePath, Ref: ref}
	}
	parts := bytes.Split(data, []byte{0})
	paths := []string{candidatePath}
	for index := 0; index < len(parts); {
		if len(parts[index]) == 0 {
			index++
			continue
		}
		status := string(parts[index])
		index++
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			if index+1 >= len(parts) {
				return nil, &BranchEvidenceBlocker{Code: "branch_evidence_scope_malformed", Message: "candidate proof rename/copy scope is malformed", Path: candidatePath, Ref: ref}
			}
			left, right := string(parts[index]), string(parts[index+1])
			if left == candidatePath || right == candidatePath {
				paths = append(paths, left, right)
			}
			index += 2
			continue
		}
		if index >= len(parts) {
			return nil, &BranchEvidenceBlocker{Code: "branch_evidence_scope_malformed", Message: "candidate proof path scope is truncated", Path: candidatePath, Ref: ref}
		}
		index++
	}
	return uniqueSortedStrings(paths), nil
}

func proposeBranchCandidateEvidence(request BranchEvidenceRequest) BranchCandidateEvidence {
	initial := RecomputeBranchCandidateEvidence(request)
	if initial.ProposedDisposition == BranchDispositionSameContentNonOwner || len(initial.Blockers) != 1 || initial.Blockers[0].Code != "branch_evidence_non_owner_unproven" {
		return initial
	}
	args := []string{"rev-list", "--max-count=513", request.EvidenceRevision, "--"}
	args = append(args, request.Paths...)
	runner, runnerBlocker := newBranchGitRunner(request.Root)
	if runnerBlocker != nil {
		initial.Blockers = append(initial.Blockers, *runnerBlocker)
		initial.Digest = CanonicalBranchEvidenceDigest(initial)
		return initial
	}
	defer runner.close()
	if blocker := runner.checkRepositoryIsolation(); blocker != nil {
		initial.Blockers = append(initial.Blockers, *blocker)
		initial.Digest = CanonicalBranchEvidenceDigest(initial)
		return initial
	}
	commitsData, searchBlocker := runner.run(513*80, nil, args...)
	if searchBlocker != nil {
		searchBlocker.Code = "branch_evidence_incorporated_search_unavailable"
		initial.Blockers = append(initial.Blockers, *searchBlocker)
		initial.Digest = CanonicalBranchEvidenceDigest(initial)
		return initial
	}
	commits := strings.Fields(string(commitsData))
	if len(commits) > 512 {
		initial.Blockers = append(initial.Blockers, BranchEvidenceBlocker{Code: "branch_evidence_incorporated_search_limit", Message: "reachable incorporated-commit search exceeded 512 candidates", Path: request.CandidatePath})
		initial.Digest = CanonicalBranchEvidenceDigest(initial)
		return initial
	}
	matches := []BranchCandidateEvidence{}
	for _, commit := range commits {
		candidateRequest := request
		candidateRequest.IncorporatedCommit = commit
		candidate := RecomputeBranchCandidateEvidence(candidateRequest)
		if candidate.ProposedDisposition == BranchDispositionIncorporatedHistoryNonOwner && len(candidate.Blockers) == 0 {
			matches = append(matches, candidate)
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}
	if len(matches) > 1 {
		initial.Blockers = []BranchEvidenceBlocker{{Code: "branch_evidence_incorporated_commit_ambiguous", Message: "more than one reachable commit matches the complete candidate transition", Path: request.CandidatePath, Ref: request.Ref}}
		initial.Digest = CanonicalBranchEvidenceDigest(initial)
		return initial
	}
	state, blocker := ReadBranchPathState(request.Root, initial.RefTip, request.CandidatePath)
	if blocker != nil || state.State != "present" || (state.Mode != GitModeRegular && state.Mode != GitModeExecutable) {
		initial.ProposedDisposition = BranchDispositionBlocking
		initial.Blockers = []BranchEvidenceBlocker{{Code: "branch_active_owner_candidate_invalid", Message: "divergent ref tip does not retain a regular candidate at the reviewed path", Path: request.CandidatePath, Ref: request.Ref}}
		initial.Digest = CanonicalBranchEvidenceDigest(initial)
		return initial
	}
	kind := request.ExpectedKind
	if kind == "" {
		kind, _ = kindForFilename(path.Base(request.CandidatePath))
	}
	data, readErr := gitBytes(request.Root, "show", initial.RefTip+":"+request.CandidatePath)
	record, parseErr := ParseDocument(request.CandidatePath, data, kind)
	if readErr != nil || parseErr != nil || record.Header.Lifecycle != LifecycleActive {
		initial.ProposedDisposition = BranchDispositionBlocking
		initial.Blockers = []BranchEvidenceBlocker{{Code: "branch_active_owner_lifecycle_invalid", Message: "divergent ref-tip candidate must parse with active lifecycle before it can be reviewed as an owner", Path: request.CandidatePath, Ref: request.Ref}}
		initial.Digest = CanonicalBranchEvidenceDigest(initial)
		return initial
	}
	initial.ProposedDisposition = BranchDispositionActiveOwner
	initial.OwnerTipCandidate = &OwnerTipCandidate{Path: request.CandidatePath, Mode: state.Mode, ObjectID: state.ObjectID, SHA256: state.SHA256, Lifecycle: LifecycleActive}
	initial.Blockers = nil
	initial.Digest = CanonicalBranchEvidenceDigest(initial)
	return initial
}

// ProposeBranchCandidateEvidence recomputes deterministic evidence and, when
// the candidate differs, searches the bounded reachable history for one exact
// incorporated transition before conservatively proposing active ownership.
func ProposeBranchCandidateEvidence(request BranchEvidenceRequest) BranchCandidateEvidence {
	return proposeBranchCandidateEvidence(request)
}

func activeBranchTouches(root string) (map[string][]string, []Problem) {
	settings, problems := LoadRepositorySettings(root)
	if HasErrors(problems) {
		return map[string][]string{}, problems
	}
	return activeBranchTouchesWithSettings(root, settings)
}

func activeBranchTouchesWithSettings(root string, settings RepositorySettings) (map[string][]string, []Problem) {
	if settings.DiscoveryVersion != DiscoveryV2 {
		return legacyActiveBranchTouches(root)
	}
	result := map[string][]string{}
	seen := map[string]map[string]bool{}
	problems := []Problem{}
	current, _ := gitText(root, "symbolic-ref", "-q", "HEAD")
	refsText, err := gitText(root, "for-each-ref", "--format=%(refname)", "refs/heads", "refs/remotes")
	if err != nil {
		return result, []Problem{{Code: "branch_enumeration_failed", Message: err.Error(), Severity: SeverityError}}
	}
	refs := strings.Fields(refsText)
	sort.Strings(refs)
	for _, ref := range refs {
		if ref == current || strings.HasSuffix(ref, "/HEAD") {
			continue
		}
		merged := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", ref, "HEAD")
		merged.Env = append(merged.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1")
		if err := merged.Run(); err == nil {
			continue
		}
		headTip, headErr := gitText(root, "rev-parse", "--verify", "HEAD^{commit}")
		refTip, refErr := gitText(root, "rev-parse", "--verify", ref+"^{commit}")
		if headErr != nil || refErr != nil {
			problems = append(problems, Problem{Code: "branch_evidence_ref_missing", Message: "target or touching ref tip could not be resolved exactly", Path: ref, Severity: SeverityError, Remediation: "repair unavailable ref objects before migration"})
			continue
		}
		if _, blocker := ResolveBranchUniqueMergeBase(root, headTip, refTip); blocker != nil {
			problems = append(problems, Problem{Code: blocker.Code, Message: blocker.Message, Path: ref, Severity: SeverityError, Remediation: "repair ambiguous or unavailable merge-base evidence before migration"})
			continue
		}
		diff, diffErr := gitBytes(root, "diff", "--name-status", "-z", "--find-renames", "--diff-filter=ACDMRT", "HEAD..."+ref)
		if diffErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: diffErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair the ref or merge-base evidence before migration"})
			continue
		}
		affected, parseErr := parseNameStatusZ(diff)
		if parseErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: parseErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair the ref diff evidence before migration"})
			continue
		}
		headCandidates, classifyErr := candidatePathsAtRevision(root, "HEAD", affected, settings)
		if classifyErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: classifyErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair local Git object evidence before migration"})
			continue
		}
		refCandidates, classifyErr := candidatePathsAtRevision(root, ref, affected, settings)
		if classifyErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: classifyErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair branch Git object evidence before migration"})
			continue
		}
		for _, candidatePath := range affected {
			if !headCandidates[candidatePath] && !refCandidates[candidatePath] {
				continue
			}
			if seen[candidatePath] == nil {
				seen[candidatePath] = map[string]bool{}
			}
			if seen[candidatePath][ref] {
				continue
			}
			seen[candidatePath][ref] = true
			result[candidatePath] = append(result[candidatePath], ref)
		}
	}
	for candidatePath := range result {
		sort.Strings(result[candidatePath])
	}
	SortProblems(problems)
	return result, problems
}

func legacyActiveBranchTouches(root string) (map[string][]string, []Problem) {
	result := map[string][]string{}
	seen := map[string]map[string]bool{}
	problems := []Problem{}
	current, _ := gitText(root, "symbolic-ref", "-q", "HEAD")
	refsText, err := gitText(root, "for-each-ref", "--format=%(refname)", "refs/heads", "refs/remotes")
	if err != nil {
		return result, []Problem{{Code: "branch_enumeration_failed", Message: err.Error(), Severity: SeverityError}}
	}
	refs := strings.Fields(refsText)
	sort.Strings(refs)
	for _, ref := range refs {
		if ref == current || strings.HasSuffix(ref, "/HEAD") {
			continue
		}
		merged := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", ref, "HEAD")
		if err := merged.Run(); err == nil {
			continue
		}
		paths, diffErr := gitText(root, "diff", "--name-status", "--find-renames", "--diff-filter=ACDMRT", "HEAD..."+ref, "--", "docs/repo/plans")
		if diffErr != nil {
			problems = append(problems, Problem{Code: "branch_touch_unavailable", Message: diffErr.Error(), Path: ref, Severity: SeverityError, Remediation: "repair the ref or merge-base evidence before migration"})
			continue
		}
		for _, line := range strings.Split(strings.TrimSpace(paths), "\n") {
			fields := strings.Split(line, "\t")
			if len(fields) < 2 {
				continue
			}
			affected := fields[1:2]
			if (strings.HasPrefix(fields[0], "R") || strings.HasPrefix(fields[0], "C")) && len(fields) >= 3 {
				affected = fields[1:3]
			}
			for _, path := range affected {
				path = strings.TrimSpace(path)
				if path == "" {
					continue
				}
				if seen[path] == nil {
					seen[path] = map[string]bool{}
				}
				if seen[path][ref] {
					continue
				}
				seen[path][ref] = true
				result[path] = append(result[path], ref)
			}
		}
	}
	for path := range result {
		sort.Strings(result[path])
	}
	SortProblems(problems)
	return result, problems
}

func parseNameStatusZ(data []byte) ([]string, error) {
	parts := bytes.Split(data, []byte{0})
	paths := []string{}
	for index := 0; index < len(parts); {
		if len(parts[index]) == 0 {
			index++
			continue
		}
		status := string(parts[index])
		index++
		count := 1
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			count = 2
		}
		if index+count > len(parts) {
			return nil, fmt.Errorf("invalid NUL-delimited name-status record for %q", status)
		}
		for offset := 0; offset < count; offset++ {
			if len(parts[index+offset]) == 0 {
				return nil, fmt.Errorf("empty path in NUL-delimited name-status record for %q", status)
			}
			paths = append(paths, string(parts[index+offset]))
		}
		index += count
	}
	return uniqueSortedStrings(paths), nil
}

func candidatePathsAtRevision(root, revision string, paths []string, settings RepositorySettings) (map[string]bool, error) {
	result := map[string]bool{}
	if len(paths) == 0 {
		return result, nil
	}
	blobs := []GitBlob{}
	const pathChunkSize = 256
	for start := 0; start < len(paths); start += pathChunkSize {
		end := start + pathChunkSize
		if end > len(paths) {
			end = len(paths)
		}
		args := []string{"-C", root, "ls-tree", "-z", revision, "--"}
		args = append(args, paths[start:end]...)
		command := exec.Command("git", args...)
		command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_LITERAL_PATHSPECS=1", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1", "LC_ALL=C")
		output, err := command.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("git ls-tree %s: %w: %s", revision, err, strings.TrimSpace(string(output)))
		}
		for _, item := range bytes.Split(output, []byte{0}) {
			if len(item) == 0 {
				continue
			}
			header, rawPath, ok := bytes.Cut(item, []byte{'\t'})
			fields := strings.Fields(string(header))
			if !ok || len(fields) != 3 {
				return nil, fmt.Errorf("invalid ls-tree record at %s", revision)
			}
			blobs = append(blobs, GitBlob{Path: string(rawPath), Mode: GitMode(fields[0]), ObjectID: fields[2], Revision: revision})
		}
	}
	batch, err := newIndexBatchReader(root)
	if err != nil {
		return nil, err
	}
	classification, classifyErr := ClassifyGitBlobs(blobs, batch.Read, settings.DiscoveryPolicy(false), settings.Mode, settings.RepositoryID, nil)
	closeErr := batch.Close()
	if classifyErr != nil {
		return nil, classifyErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	for _, candidate := range classification.Candidates {
		result[candidate.Path] = true
	}
	return result, nil
}

func gitPathDirty(root, path string) (bool, error) {
	output, err := gitBytes(root, "status", "--porcelain=v1", "--untracked-files=all", "--", path)
	if err != nil {
		return false, err
	}
	return len(bytes.TrimSpace(output)) > 0, nil
}

func gitText(root string, args ...string) (string, error) {
	output, err := gitBytes(root, args...)
	return strings.TrimSpace(string(output)), err
}

func gitBytes(root string, args ...string) ([]byte, error) {
	commandArgs := append([]string{"-C", root}, args...)
	command := exec.Command("git", commandArgs...)
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_NO_LAZY_FETCH=1", "GIT_NO_REPLACE_OBJECTS=1")
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return output, nil
}
