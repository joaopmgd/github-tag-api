package database

import (
	"github.com/jinzhu/gorm"
)

// RepoTag are the tags of some repo
type RepoTag struct {
	gorm.Model

	UserID  string
	RepoID  int64
	TagName string
}

// InsertRepoTagsValue inserts in the database a new repo tag
func (db *Gorm) InsertRepoTagsValue(value RepoTag) {
	db.Conn.Create(&value)
}

// GetAllRepoTagsMap recovers all repo tags for and user id
func (db *Gorm) GetAllRepoTagsMap(userID string) map[int64][]string {
	var repoTags []RepoTag
	db.Conn.Where("user_id = ?", userID).Find(&repoTags)
	repoTagsMap := make(map[int64][]string)
	for _, repoTag := range repoTags {
		repoTagsMap[repoTag.RepoID] = append(repoTagsMap[repoTag.RepoID], repoTag.TagName)
	}
	return repoTagsMap
}

// GetAllRepoTagsByRepoID recovers all repo tags for an repo id and user id
func (db *Gorm) GetAllRepoTagsByRepoID(userID string, repoID int64) []RepoTag {
	var repoTags []RepoTag
	db.Conn.Where("user_id = ? AND repo_id = ?", userID, repoID).Find(&repoTags)
	return repoTags
}
