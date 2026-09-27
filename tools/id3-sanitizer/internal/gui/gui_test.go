package gui

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/lburfeindt/trackops/tools/id3-sanitizer/internal/batch"
	"github.com/ncruces/zenity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_SanitizesSelectionAndReportsResult(t *testing.T) {
	// Given
	fixture, err := os.ReadFile("../../testdata/Silence.mp3")
	require.NoError(t, err)
	track := filepath.Join(t.TempDir(), "track.mp3")
	require.NoError(t, os.WriteFile(track, fixture, 0o600))
	ui := newTestDialogs([]string{track})

	// When
	err = Run(ui)

	// Then
	require.NoError(t, err)
	require.NotNil(t, ui.result)
	assert.Equal(t, 1, ui.result.Updated)
	assert.Empty(t, ui.result.Failures)
	assert.True(t, ui.progress.closed)
	assert.True(t, ui.closedBeforeResult)
	updated, err := os.ReadFile(track)
	require.NoError(t, err)
	assert.NotEqual(t, fixture, updated)
}

func TestRun_CancelOrDialogFailureDoesNotModifyFiles(t *testing.T) {
	for _, scenario := range []string{"selection canceled", "empty selection", "selection failed", "progress failed", "progress canceled", "progress update failed"} {
		t.Run(scenario, func(t *testing.T) {
			// Given
			fixture, err := os.ReadFile("../../testdata/Silence.mp3")
			require.NoError(t, err)
			track := filepath.Join(t.TempDir(), "track.mp3")
			require.NoError(t, os.WriteFile(track, fixture, 0o600))
			ui := newTestDialogs([]string{track})
			dialogErr := errors.New("dialog unavailable")
			wantError := false
			switch scenario {
			case "selection canceled":
				ui.chooseErr = zenity.ErrCanceled
			case "empty selection":
				ui.paths = nil
			case "selection failed":
				ui.chooseErr, wantError = dialogErr, true
			case "progress failed":
				ui.progressErr, wantError = dialogErr, true
			case "progress canceled":
				close(ui.progress.done)
			case "progress update failed":
				ui.progress.textErr, wantError = dialogErr, true
			}

			// When
			err = Run(ui)

			// Then
			if wantError {
				assert.ErrorIs(t, err, dialogErr)
			} else {
				assert.NoError(t, err)
			}
			contents, readErr := os.ReadFile(track)
			require.NoError(t, readErr)
			assert.Equal(t, fixture, contents)
			if scenario == "progress canceled" {
				require.NotNil(t, ui.result)
				assert.True(t, ui.result.Canceled)
			}
		})
	}
}

func TestSummary_ReportsPartialResultsAndLimitsErrorDetails(t *testing.T) {
	// Given
	result := batch.Result{Updated: 2, Unchanged: 3, Canceled: true}
	for i := 0; i < 12; i++ {
		result.Failures = append(result.Failures, batch.Failure{Path: "broken.mp3", Err: errors.New("invalid tag")})
	}

	// When
	text := summary(result)

	// Then
	assert.Contains(t, text, "Canceled")
	assert.Contains(t, text, "Updated: 2")
	assert.Contains(t, text, "Unchanged: 3")
	assert.Contains(t, text, "Failed: 12")
	assert.Contains(t, text, "broken.mp3: invalid tag")
	assert.Contains(t, text, "2 more")
}

type testDialogs struct {
	paths              []string
	chooseErr          error
	progressErr        error
	progress           *testProgress
	result             *batch.Result
	closedBeforeResult bool
}

func newTestDialogs(paths []string) *testDialogs {
	return &testDialogs{paths: paths, progress: &testProgress{done: make(chan struct{})}}
}

func (d *testDialogs) Choose() ([]string, error)   { return d.paths, d.chooseErr }
func (d *testDialogs) Progress() (Progress, error) { return d.progress, d.progressErr }
func (d *testDialogs) Result(result batch.Result) error {
	d.result = &result
	d.closedBeforeResult = d.progress.closed
	return nil
}

type testProgress struct {
	done    chan struct{}
	closed  bool
	textErr error
}

func (p *testProgress) Text(string) error     { return p.textErr }
func (p *testProgress) Done() <-chan struct{} { return p.done }
func (p *testProgress) Close() error          { p.closed = true; return nil }
