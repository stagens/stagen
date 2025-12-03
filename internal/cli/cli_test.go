package cli

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"stagen/internal/config"
	"stagen/pkg/git"
)

type fakeClock struct {
	time time.Time
}

func newFakeClock(time time.Time) *fakeClock {
	return &fakeClock{
		time: time,
	}
}

func (c *fakeClock) Now() time.Time {
	return c.time
}

func (c *fakeClock) Sleep(duration time.Duration) {}

func (c *fakeClock) Since(value time.Time) time.Duration {
	return c.Now().Sub(value)
}

func rootDir() string {
	currentRootDir := config.RootDir()

	return filepath.Join(currentRootDir, "../../")
}

func TestInit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	clocks := newFakeClock(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	gitTool := git.New()

	cliTool := New(clocks, gitTool)

	workDir := filepath.Join(rootDir(), "tests/new_project")

	err := os.RemoveAll(workDir)
	require.NoError(t, err)

	t.Cleanup(func() {
		if !t.Failed() {
			err := os.RemoveAll(workDir)
			require.NoError(t, err)
		}
	})

	err = cliTool.Init(ctx, workDir, "My Cool Website", false)
	require.NoError(t, err)

	checkDir := filepath.Join(rootDir(), "tests/_new_project_check")

	diffs, err := DiffDirs(workDir, checkDir)
	require.NoError(t, err)

	if len(diffs) > 0 {
		t.Log(strings.Join(diffs, "\n"))
	}

	require.Empty(t, diffs)
}

func TestBuild(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	clocks := newFakeClock(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

	tests := []struct {
		name    string
		testDir string
	}{
		{
			name:    "base",
			testDir: filepath.Join(rootDir(), "tests/01-base"),
		},
		{
			name:    "vars",
			testDir: filepath.Join(rootDir(), "tests/02-vars"),
		},
		{
			name:    "imports",
			testDir: filepath.Join(rootDir(), "tests/03-imports"),
		},
		{
			name:    "macros",
			testDir: filepath.Join(rootDir(), "tests/04-macros"),
		},
		{
			name:    "generators",
			testDir: filepath.Join(rootDir(), "tests/05-generators"),
		},
		// @todo includes
		// @todo extras
		// @todo theme changing
		// @todo layout changing
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			workDir := testCase.testDir

			buildDir := filepath.Join(workDir, "build")
			checkDir := filepath.Join(workDir, "_check")

			err := os.RemoveAll(buildDir)
			require.NoError(t, err)

			t.Cleanup(func() {
				if !t.Failed() {
					err = os.RemoveAll(buildDir)
					require.NoError(t, err)
				}
			})

			gitTool := git.New()

			cliTool := New(clocks, gitTool)

			err = cliTool.Build(ctx, workDir)
			require.NoError(t, err)

			diffs, err := DiffDirs(buildDir, checkDir)
			require.NoError(t, err)

			if len(diffs) > 0 {
				t.Log(strings.Join(diffs, "\n"))
			}

			require.Empty(t, diffs)
		})
	}
}

func DiffDirs(buildDir, checkDir string) ([]string, error) {
	buildDir = filepath.Clean(buildDir)
	checkDir = filepath.Clean(checkDir)

	checkMap := make(map[string]fs.DirEntry)
	buildMap := make(map[string]fs.DirEntry)

	var diffs []string

	err := filepath.WalkDir(checkDir, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == checkDir {
			return nil // пропускаем корень
		}

		rel, err := filepath.Rel(checkDir, path)
		if err != nil {
			return err
		}

		checkMap[rel] = dirEntry

		return nil
	})
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(buildDir, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == buildDir {
			return nil // пропускаем корень
		}

		rel, err := filepath.Rel(buildDir, path)
		if err != nil {
			return err
		}

		buildMap[rel] = dirEntry

		return nil
	})
	if err != nil {
		return nil, err
	}

	for rel, checkEntry := range checkMap {
		buildEntry, ok := buildMap[rel]
		if !ok {
			diffs = append(diffs, "MISSING in buildDir: "+rel)

			continue
		}

		if checkEntry.IsDir() != buildEntry.IsDir() {
			diffs = append(diffs, fmt.Sprintf("TYPE MISMATCH: %s (checkDir is %s, buildDir is %s)",
				rel,
				dirOrFile(checkEntry),
				dirOrFile(buildEntry),
			))

			continue
		}

		if !checkEntry.IsDir() {
			checkPath := filepath.Join(checkDir, rel)
			buildPath := filepath.Join(buildDir, rel)

			equal, err := filesEqual(checkPath, buildPath)
			if err != nil {
				return nil, err
			}

			if !equal {
				diffs = append(diffs, "CONTENT MISMATCH: "+rel)
			}
		}
	}

	for rel := range buildMap {
		if _, ok := checkMap[rel]; !ok {
			diffs = append(diffs, "EXTRA in buildDir: "+rel)
		}
	}

	return diffs, nil
}

func dirOrFile(d fs.DirEntry) string {
	if d.IsDir() {
		return "dir"
	}

	return "file"
}

func filesEqual(path1, path2 string) (bool, error) {
	info1, err := os.Stat(path1)
	if err != nil {
		return false, err
	}

	info2, err := os.Stat(path2)
	if err != nil {
		return false, err
	}

	if info1.Size() != info2.Size() {
		return false, nil
	}

	file1, err := os.Open(path1)
	if err != nil {
		return false, err
	}

	defer func() {
		//nolint:staticcheck
		if fErr := file1.Close(); fErr != nil {
			// @todo log
		}
	}()

	file2, err := os.Open(path2)
	if err != nil {
		return false, err
	}

	defer func() {
		//nolint:staticcheck
		if fErr := file2.Close(); fErr != nil {
			// @todo log
		}
	}()

	const bufSize = 32 * 1024

	buf1 := make([]byte, bufSize)
	buf2 := make([]byte, bufSize)

	for {
		len1, err1 := file1.Read(buf1)
		len2, err2 := file2.Read(buf2)

		if len1 != len2 {
			return false, nil
		}

		if len1 == 0 { // оба EOF
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil
			}

			return false, fmt.Errorf("unexpected read errors: %w, %w", err1, err2)
		}

		if !equalBytes(buf1[:len1], buf2[:len2]) {
			return false, nil
		}

		if err1 == io.EOF || err2 == io.EOF {
			break
		}
	}

	return true, nil
}

func equalBytes(buf1 []byte, buf2 []byte) bool {
	if len(buf1) != len(buf2) {
		return false
	}

	for i := range buf1 {
		if buf1[i] != buf2[i] {
			return false
		}
	}

	return true
}
