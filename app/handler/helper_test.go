package handler

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/joaopmgd/github-tag-api/app/model"
)

func TestAddLanguage(t *testing.T) {
	tt := map[string]struct {
		language      string
		tags          []string
		expectedSlice []string
	}{
		"empty_language":    {"", []string{"golang", "java", "docker"}, []string{"golang", "java", "docker"}},
		"new_language":      {"python", []string{"golang", "java", "docker"}, []string{"golang", "java", "docker", "python"}},
		"repeated_language": {"golang", []string{"golang", "java", "docker"}, []string{"golang", "java", "docker"}},
		"empty_slice":       {"golang", []string{}, []string{"golang"}},
		"nil_slice":         {"golang", nil, []string{"golang"}},
	}
	for testName, tc := range tt {

		finalSlice := addLanguage(tc.language, tc.tags)

		if len(finalSlice) != len(tc.expectedSlice) ||
			!reflect.DeepEqual(finalSlice, tc.expectedSlice) {
			t.Errorf("\nTest %s\nLanguage '%s' and Initial Tags %s\nGot Slice %s with Length %v\nWant Slice %s with Lenght %v",
				testName, tc.language, tc.tags, finalSlice, len(finalSlice), tc.expectedSlice, len(tc.expectedSlice))
		}
	}
}

func TestRepoHasTags(t *testing.T) {
	tt := map[string]struct {
		selectedTag string
		tags        []string
		isPresent   bool
	}{
		"empty_tag_list":     {"golang", []string{}, false},
		"empty_selected_tag": {"", []string{"golang", "java", "docker"}, false},
		"tag_not_found":      {"golang", []string{"java", "docker"}, false},
		"tag_found":          {"golang", []string{"golang", "java", "docker"}, true},
		"nil_slice":          {"golang", nil, false},
	}
	for testName, tc := range tt {

		decision := repoHasTag(tc.tags, tc.selectedTag)

		if decision != tc.isPresent {
			t.Errorf("\nTest %s\nSelected tag '%s' and Initial Tags %s\nGot %s\nWant %s",
				testName, tc.selectedTag, tc.tags, strconv.FormatBool(decision), strconv.FormatBool(tc.isPresent))
		}
	}
}

func TestCreateMessageStarredReposSelectedTag(t *testing.T) {
	tt := map[string]struct {
		repos       []model.StarredRepoRequest
		tags        map[int64][]string
		selectedTag string
		response    []model.StarredRepoTags
	}{
		"empty_tag_list":         {"golang", []string{}, false},
		"empty_repos":            {"", []string{"golang", "java", "docker"}, false},
		"selected_tag_found":     {"golang", []string{"java", "docker"}, false},
		"selected_tag_not_found": {"golang", []string{"golang", "java", "docker"}, true},
		"nil_repos":              {"golang", nil, false},
		"nil_tags":               {"golang", nil, false},
	}
	for testName, tc := range tt {

		decision := repoHasTag(tc.tags, tc.selectedTag)

		if decision != tc.isPresent {
			t.Errorf("\nTest %s\nSelected tag '%s' and Initial Tags %s\nGot %s\nWant %s",
				testName, tc.selectedTag, tc.tags, strconv.FormatBool(decision), strconv.FormatBool(tc.isPresent))
		}
	}
}
