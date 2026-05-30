package chat

import "testing"

// These tests cover the initial chat package skeleton.
// They intentionally verify constructor initialization rather than chat behavior:
// nil maps would panic on writes, and nil channels would block forever once the
// room/client loops are implemented.

func TestNewHubInitializesRooms(t *testing.T) {
	hub := NewHub()

	if hub == nil {
		t.Fatal("expected hub, got nil")
	}
	if hub.rooms == nil {
		t.Fatal("expected rooms map to be initialized")
	}
	if len(hub.rooms) != 0 {
		t.Fatalf("expected no rooms, got %d", len(hub.rooms))
	}
}

func TestNewRoomInitializesState(t *testing.T) {
	room := NewRoom(42)

	if room == nil {
		t.Fatal("expected room, got nil")
	}
	if room.eventID != 42 {
		t.Fatalf("expected eventID 42, got %d", room.eventID)
	}
	if room.clients == nil {
		t.Fatal("expected clients map to be initialized")
	}
	if room.join == nil {
		t.Fatal("expected join channel to be initialized")
	}
	if room.leave == nil {
		t.Fatal("expected leave channel to be initialized")
	}
	if room.broadcast == nil {
		t.Fatal("expected broadcast channel to be initialized")
	}
}

func TestNewHandlerInitializesHubAndDB(t *testing.T) {
	handler := NewHandler(nil)

	if handler.Hub == nil {
		t.Fatal("expected handler hub to be initialized")
	}
	if handler.Hub.rooms == nil {
		t.Fatal("expected handler hub rooms map to be initialized")
	}
	if handler.DB != nil {
		t.Fatal("expected nil DB to be preserved")
	}
}

