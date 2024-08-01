// gpt4o mini @ 240801
// passed.

package psql

import (
	"fmt"
	"reflect"
	"testing"
	"time"
	// "github.com/stretchr/testify/assert"
	// 需要导入你的 orderedmap 库
)

// TestCreateStatus 测试插入状态记录
func TestCreateStatus(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	id := int64(1)
	username := "user1"
	warning := "This is a warning"
	content := "This is the content of the status"
	visibility := "public"

	if err := CreateStatus(tx, id, username, warning, content, visibility); err != nil {
		t.Fatal("failed to create status:", err)
		tx.Rollback()
	}

	// 验证状态记录是否已插入
	status, err := GetStatus(tx, id) // 这里假设 ID 是 1，需要根据实际情况调整
	if err != nil {
		t.Fatal("failed to get status:", err)
		tx.Rollback()
	}

	if status.Username != username {
		t.Fatalf("expected username %s, got %s", username, status.Username)
	}
	if status.Warning != warning {
		t.Fatalf("expected warning %s, got %s", warning, status.Warning)
	}
	if status.Content != content {
		t.Fatalf("expected content %s, got %s", content, status.Content)
	}
	if status.Visibility != visibility {
		t.Fatalf("expected visibility %s, got %s", visibility, status.Visibility)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestGetStatus 测试根据 ID 获取状态记录
func TestGetStatus(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	id := int64(1)
	// 假设一个状态记录的 ID 为 1 存在
	status, err := GetStatus(tx, id) // 这里假设 ID 是 1，需要根据实际情况调整
	if err != nil {
		t.Fatal("failed to get status:", err)
		tx.Rollback()
	}

	if status.ID != 1 { // 需要确保 ID 是 1，或者根据实际情况调整
		t.Fatalf("expected ID 1, got %d", status.ID)
	}

	// 尝试获取一个不存在的状态记录
	_, err = GetStatus(tx, int64(999)) // 这里假设 9999 是一个不存在的 ID
	if err == nil {
		t.Fatal("expected error for non-existing status, got nil")
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestUpdateStatus 测试更新状态记录
func TestUpdateStatus(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	id := int64(1)
	newWarning := "Updated warning"
	newContent := "Updated content"
	newVisibility := "private"

	// 更新状态记录
	if err := UpdateStatus(tx, id, newWarning, newContent, newVisibility); err != nil {
		t.Fatal("failed to update status:", err)
		tx.Rollback()
	}

	// 验证状态记录是否已更新
	status, err := GetStatus(tx, id)
	if err != nil {
		t.Fatal("failed to get status:", err)
		tx.Rollback()
	}

	if status.Warning != newWarning {
		t.Fatalf("expected warning %s, got %s", newWarning, status.Warning)
	}
	if status.Content != newContent {
		t.Fatalf("expected content %s, got %s", newContent, status.Content)
	}
	if status.Visibility != newVisibility {
		t.Fatalf("expected visibility %s, got %s", newVisibility, status.Visibility)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// TestSoftDeleteStatus 测试软删除状态记录
func TestSoftDeleteStatus(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	id := int64(1)

	// 软删除状态记录
	if err := SoftDeleteStatus(tx, id); err != nil {
		t.Fatal("failed to soft delete status:", err)
		tx.Rollback()
	}

	// 验证状态记录是否已标记为删除
	status, err := GetStatus(tx, id)
	if err != nil {
		t.Fatal("failed to get status:", err)
		tx.Rollback()
	}

	if status.Visibility != "deleted" {
		t.Fatalf("expected visibility 'deleted', got %s", status.Visibility)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

// it will not work. by gpt4o
// TestGetStatuses tests the GetStatuses function
func TestGetStatuses(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	testStatuses := []struct {
		ID         int64
		Username   string
		Warning    string
		Content    string
		Visibility string
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}{
		{1, "user1", "warning1", "content1", "public", time.Now(), time.Now()},
		{2, "user1", "warning2", "content2", "private", time.Now(), time.Now()},
		{3, "user1", "warning3", "content3", "unlisted", time.Now(), time.Now()},
	}

	for _, s := range testStatuses {
		err := CreateStatus(tx, s.ID, s.Username, s.Warning, s.Content, s.Visibility)
		if err != nil {
			// t.Fatal("failed to insert test data:", err)
			fmt.Println("failed to insert test data:", err)
		}
	}

	// Commit transaction to persist the changes
	if err := tx.Commit(); err != nil {
		t.Fatal("failed to commit transaction:", err)
	}

	// Test cases
	tests := []struct {
		name      string
		ids       []int64
		wantCount int
		wantIDs   []int64
		wantError bool
	}{
		{
			name:      "Valid IDs",
			ids:       []int64{1, 2, 3},
			wantCount: 3,
			wantIDs:   []int64{1, 2, 3},
		},
		{
			name:      "Empty IDs",
			ids:       []int64{},
			wantCount: 0,
			wantIDs:   []int64{},
		},
		{
			name:      "Partial Match",
			ids:       []int64{1, 4},
			wantCount: 1,
			wantIDs:   []int64{1},
		},
		{
			name:      "Invalid Query",
			ids:       []int64{100}, // assuming this ID does not exist
			wantCount: 0,
			wantIDs:   []int64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := Begin()
			if err != nil {
				t.Fatal("failed to begin transaction:", err)
			}
			defer tx.Rollback()

			statuses, err := GetStatuses(tx, tt.ids)
			if (err != nil) != tt.wantError {
				t.Errorf("GetStatuses() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if len(statuses) != tt.wantCount {
				t.Errorf("GetStatuses() returned %d statuses, want %d", len(statuses), tt.wantCount)
			}

			gotIDs := make([]int64, len(statuses))
			for i, s := range statuses {
				gotIDs[i] = s.ID
			}

			if !reflect.DeepEqual(gotIDs, tt.wantIDs) {
				t.Errorf("GetStatuses() returned IDs = %v, want %v", gotIDs, tt.wantIDs)
			}

		})
	}
}
