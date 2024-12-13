// gpt4o @ 240801

package psql

import (
	"reflect"
	"testing"

	orderedmap "github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
)

func TestCreateUser(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	username := "user1"
	displayName := "User One"
	avatarURL := "http://example.com/avatar.png"
	settings := map[string]interface{}{
		"theme":         "dark",
		"notifications": true,
	}
	settingsMap := orderedmap.NewFromMap(settings)

	if err := CreateUser(tx, username, "", displayName, avatarURL, "", "", settingsMap); err != nil {
		t.Fatal("failed to create user:", err)
		tx.Rollback()
	}

	// Verify the user was inserted correctly
	user, err := GetUser(tx, username)
	if err != nil {
		t.Fatal("failed to get user:", err)
		tx.Rollback()
	}

	if user.Username != username {
		t.Fatalf("expected username %s, got %s", username, user.Username)
	}
	if user.DisplayName != displayName {
		t.Fatalf("expected displayName %s, got %s", displayName, user.DisplayName)
	}
	if user.AvatarURL != avatarURL {
		t.Fatalf("expected avatarURL %s, got %s", avatarURL, user.AvatarURL)
	}
	if user.Flag != "active" {
		t.Fatalf("expected flag 'active', got %s", user.Flag)
	}
	if !reflect.DeepEqual(user.Settings, settingsMap) {
		t.Fatalf("expected settings %v, got %v", settingsMap, user.Settings)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

func TestGetUser(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	// Assume a user with username "user1" exists
	username := "user1"
	user, err := GetUser(tx, username)
	if err != nil {
		t.Fatal("failed to get user:", err)
		tx.Rollback()
	}

	if user.Username != username {
		t.Fatalf("expected username %s, got %s", username, user.Username)
	}

	// Try to get a non-existing user
	_, err = GetUser(tx, "non_existing_user")
	if err == nil {
		t.Fatal("expected error for non-existing user, got nil")
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

func TestUpdateUser(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	username := "user1"
	newDisplayName := "Updated User One"
	newAvatarURL := "http://example.com/new_avatar.png"
	newSettings := map[string]interface{}{
		"theme":         "light",
		"notifications": false,
	}
	newSettingsMap := orderedmap.NewFromMap(newSettings)
	newFlag := "inactive"

	// Update the user
	if err := UpdateUser(tx, username, newDisplayName, newAvatarURL, "newbio", "{}", newSettingsMap, newFlag); err != nil {
		t.Fatal("failed to update user:", err)
		tx.Rollback()
	}

	// Verify the user was updated correctly
	user, err := GetUser(tx, username)
	if err != nil {
		t.Fatal("failed to get user:", err)
		tx.Rollback()
	}

	if user.DisplayName != newDisplayName {
		t.Fatalf("expected displayName %s, got %s", newDisplayName, user.DisplayName)
	}
	if user.AvatarURL != newAvatarURL {
		t.Fatalf("expected avatarURL %s, got %s", newAvatarURL, user.AvatarURL)
	}
	if !reflect.DeepEqual(user.Settings, newSettingsMap) {
		t.Fatalf("expected settings %v, got %v", newSettingsMap, user.Settings)
	}
	if user.Flag != newFlag {
		t.Fatalf("expected flag %s, got %s", newFlag, user.Flag)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}

func TestSoftDeleteUser(t *testing.T) {
	tx, err := DB.Begin()
	if err != nil {
		t.Fatal("failed to begin transaction:", err)
	}

	username := "user1"

	// Soft delete the user
	if err := SoftDeleteUser(tx, username); err != nil {
		t.Fatal("failed to soft delete user:", err)
		tx.Rollback()
	}

	// Verify the user was marked as deleted
	user, err := GetUser(tx, username)
	if err != nil {
		t.Fatal("failed to get user:", err)
		tx.Rollback()
	}

	if user.Flag != "deleted" {
		t.Fatalf("expected flag 'deleted', got %s", user.Flag)
	}

	if err := tx.Commit(); err != nil {
		t.Logf("error on commit: %v", err)
		tx.Rollback()
	}
}
