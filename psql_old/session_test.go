package psql

import (
	"fmt"
	"log"
	"testing"
)

func TestCreateSession(t *testing.T) {
	tx, err := Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	sessionID := "sess123"
	username := "user1"
	country := "US"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	err = CreateSession(tx, sessionID, username, country, ipAddress, userAgent)
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = CreateSession(tx, sessionID+"2", username, country, ipAddress, userAgent)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}
}

func TestGetSession(t *testing.T) {
	tx, err := Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	sessionID := "sess123"
	session, err := GetSession(tx, sessionID)
	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Printf("%v", session)

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}
}

func TestUpdateSession(t *testing.T) {
	tx, err := Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	// sessionID := "sess123"
	sessions, err := GetSessions(tx, "user1")
	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Printf("%v", sessions)
	// fmt.Printf("%v", sessions[0])
	// fmt.Printf("%v", sessions[1])

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}
}

func TestDeleteSession(t *testing.T) {
	tx, err := Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	sessionID := "sess123"
	err = DeleteSession(tx, sessionID)
	// err = DeleteSession(tx, sessionID+"2")
	if err != nil {
		t.Fatalf("%v", err)
	}

	// fmt.Printf("%v", sessions)
	// fmt.Printf("%v", sessions[0])
	// fmt.Printf("%v", sessions[1])

	if err := tx.Commit(); err != nil {
		log.Printf("error on commit: %v", err.Error())
		tx.Rollback()
	}
}
