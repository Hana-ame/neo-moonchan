// gpt4o mini @ 240802

package psql

import (
	"reflect"
	"testing"
)

// TestCreateLink 测试插入链接记录
func TestCreateLink(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	linkID := int64(1)
	username := "user1"
	statusID := int64(100)
	visibility := "public"

	if err := CreateLink(tx, linkID, username, statusID, visibility); err != nil {
		t.Fatal("failed to create link:", err)
		tx.Rollback()
	}

	// 验证链接记录是否已插入
	link, err := GetLink(tx, linkID)
	if err != nil {
		t.Fatal("failed to get link:", err)
		tx.Rollback()
	}

	if link.Username != username {
		t.Fatalf("expected username %s, got %s", username, link.Username)
	}
	if link.StatusID != statusID {
		t.Fatalf("expected statusID %d, got %d", statusID, link.StatusID)
	}
	if link.Visibility != visibility {
		t.Fatalf("expected visibility %s, got %s", visibility, link.Visibility)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestGetLink 测试根据 LinkID 获取链接记录
func TestGetLink(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	linkID := int64(1)
	link, err := GetLink(tx, linkID)
	if err != nil {
		t.Fatal("failed to get link:", err)
		tx.Rollback()
	}

	if link.LinkID != linkID {
		t.Fatalf("expected LinkID %d, got %d", linkID, link.LinkID)
	}

	// 尝试获取一个不存在的链接记录
	_, err = GetLink(tx, int64(999))
	if err == nil {
		t.Fatal("expected error for non-existing link, got nil")
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestUpdateLink 测试更新链接记录
func TestUpdateLink(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	linkID := int64(1)
	newUsername := "user2"
	newStatusID := int64(200)
	newVisibility := "private"

	if err := UpdateLink(tx, linkID, newUsername, newStatusID, newVisibility); err != nil {
		t.Fatal("failed to update link:", err)
		tx.Rollback()
	}

	// 验证链接记录是否已更新
	link, err := GetLink(tx, linkID)
	if err != nil {
		t.Fatal("failed to get link:", err)
		tx.Rollback()
	}

	if link.Username != newUsername {
		t.Fatalf("expected username %s, got %s", newUsername, link.Username)
	}
	if link.StatusID != newStatusID {
		t.Fatalf("expected statusID %d, got %d", newStatusID, link.StatusID)
	}
	if link.Visibility != newVisibility {
		t.Fatalf("expected visibility %s, got %s", newVisibility, link.Visibility)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestSoftDeleteLink 测试软删除链接记录
func TestSoftDeleteLink(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	linkID := int64(1)

	if err := SoftDeleteLink(tx, linkID); err != nil {
		t.Fatal("failed to soft delete link:", err)
		tx.Rollback()
	}

	// 验证链接记录是否已标记为删除
	link, err := GetLink(tx, linkID)
	if err != nil {
		t.Fatal("failed to get link:", err)
		tx.Rollback()
	}

	if link.Visibility != "deleted" {
		t.Fatalf("expected visibility 'deleted', got %s", link.Visibility)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestGetLinks 测试根据多个 LinkID 获取链接记录
func TestGetLinks(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	linkIDs := []int64{1, 2, 3}
	links, err := GetLinks(tx, linkIDs)
	if err != nil {
		t.Fatal("failed to get links:", err)
		tx.Rollback()
	}

	if len(links) != len(linkIDs) {
		t.Fatalf("expected %d links, got %d", len(linkIDs), len(links))
	}

	gotIDs := make([]int64, len(links))
	for i, link := range links {
		gotIDs[i] = link.LinkID
	}

	if !reflect.DeepEqual(gotIDs, linkIDs) {
		t.Errorf("GetLinks() returned IDs = %v, want %v", gotIDs, linkIDs)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestGetStatusesFromLinksMaxID 测试根据小于某个 ID 和用户名获取链接记录的状态信息
func TestGetStatusesFromLinksMaxID(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	maxID := int64(10)
	username := "user1"
	limit := 5

	statuses, err := GetStatusesByUsernameFromLinksMaxID(tx, maxID, username, limit)
	if err != nil {
		t.Fatal("failed to get statuses:", err)
		tx.Rollback()
	}

	// 验证返回的状态记录数量
	if len(statuses) > limit {
		t.Fatalf("expected at most %d statuses, got %d", limit, len(statuses))
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestGetStatusesFromLinksMinID 测试根据大于某个 ID 和用户名获取链接记录的状态信息
func TestGetStatusesFromLinksMinID(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	minID := int64(5)
	username := "user1"
	limit := 5

	statuses, err := GetStatusesByUsernameFromLinksMinID(tx, minID, username, limit)
	if err != nil {
		t.Fatal("failed to get statuses:", err)
		tx.Rollback()
	}

	// 验证返回的状态记录数量
	if len(statuses) > limit {
		t.Fatalf("expected at most %d statuses, got %d", limit, len(statuses))
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}
