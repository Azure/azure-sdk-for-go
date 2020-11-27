package report

type Status interface {
	IsEmpty() bool
	HasBreakingChanges() bool
	HasAdditiveChanges() bool
}

// represents a collection of per-package reports, one for each commit hash
type CommitPkgReport struct {
	BreakingChanges []string           `json:"breakingChanges,omitempty"`
	CommitsReports  map[string]Package `json:"deltas"`
}

// returns true if the report contains no data
func (c CommitPkgReport) IsEmpty() bool {
	for _, rpt := range c.CommitsReports {
		if !rpt.IsEmpty() {
			return false
		}
	}
	return true
}

// returns true if the report contains breaking changes
func (c CommitPkgReport) HasBreakingChanges() bool {
	for _, r := range c.CommitsReports {
		if r.HasBreakingChanges() {
			return true
		}
	}
	return false
}

// returns true if the report contains additive changes
func (c CommitPkgReport) HasAdditiveChanges() bool {
	for _, r := range c.CommitsReports {
		if r.HasAdditiveChanges() {
			return true
		}
	}
	return false
}
