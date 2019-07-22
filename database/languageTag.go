package database

import (
	"github.com/jinzhu/gorm"
)

// LanguageTag are the tags of some language
type LanguageTag struct {
	gorm.Model

	Language string
	TagName  string
}

// InsertLanguageTagsValue inserts in the database a new language tag
func (db *Gorm) InsertLanguageTagsValue(value LanguageTag) {
	db.Conn.Create(&value)
}

// GetRecommendationTagByLanguage returns an array of string ordered by usage based on the language
func (db *Gorm) GetRecommendationTagByLanguage(language string) []string {
	var tags []LanguageTag
	db.Conn.Select("tag_name").Where("language = ?", language).Group("tag_name").Limit(10).Find(&tags)
	var mostUsedTags []string
	for _, tag := range tags {
		mostUsedTags = append(mostUsedTags, tag.TagName)
	}
	return mostUsedTags
}
