package stats

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type stat struct {
	word        string
	occurrences int
}

func (s stat) String() string {
	return fmt.Sprintf("%q: %v", s.word, s.occurrences)
}

type statMap map[string]stat

func (s statMap) toStatList() statList {
	out := make([]stat, 0, len(s))
	for _, v := range s {
		out = append(out, v)
	}
	return out
}

func (s statMap) add(a statMap) statMap {
	for k, v := range a {
		stat := s[k]
		stat.occurrences += v.occurrences
		stat.word = v.word
		s[k] = stat
	}
	return s
}

type statList []stat

func (s statList) sort() {
	slices.SortFunc(s, func(a, b stat) int {
		return b.occurrences - a.occurrences
	})
}

func (s statList) String() string {
	out := make([]string, 0, len(s))
	for _, v := range s {
		out = append(out, v.String())
	}
	return strings.Join(out, "\n")
}

func TopWords(dir string, positions int) (string, error) {
	stats, err := readStatsInDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read top words: %w", err)
	}

	sList := stats.toStatList()
	sList.sort()
	culled := sList[:positions]
	out := culled.String()
	return out, nil
}

func readStatsInDir(dir string) (statMap, error) {
	stats := statMap{}
	files, err := os.ReadDir(dir)
	if err != nil {
		return statMap{}, fmt.Errorf("failed to read stats in dir %q: %w", dir, err)
	}
	for _, v := range files {
		name := v.Name()
		path := filepath.Join(dir, name)
		if v.IsDir() {
			moreStats, err := readStatsInDir(path)
			if err != nil {
				return statMap{}, fmt.Errorf("failed to read subdir %q stats: %w", path, err)
			}
			stats.add(moreStats)
		}
		markdownFile := regexp.MustCompile(`\.md$`)
		// TODO: separate function, possibly concurrent
		if markdownFile.MatchString(name) {
			file, err := os.Open(path)
			if err != nil {
				return statMap{}, fmt.Errorf("failed to open file for stats: %w", err)
			}

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				w := scanner.Text()
				s := stats[w]
				s.occurrences++
				s.word = w
				stats[w] = s
			}

			err = scanner.Err()
			if err != nil {
				return statMap{}, fmt.Errorf("failed scanning for stats: %w", err)
			}
		}
	}
	return stats, nil
}
