package chat

import (
	// Std
	"testing"
	"time"

	// Extern
	"gorm.io/gorm"
)

// These backend-only tests cover constructors, room lifecycle, DTO conversion,
// message input validation, repository limit normalization, and basic WebSocket
// rejection paths. Full end-to-end verification still needs a running backend
// and an actual WebSocket client or browser.

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

func TestMessageToDTO(t *testing.T) {
	createdAt := time.Date(2026, 5, 31, 10, 30, 0, 0, time.UTC)

	message := Message{
		Model: gorm.Model{
			ID:        7,
			CreatedAt: createdAt,
		},
		EventID: 42,
		UserID:  3,
		Content: "hello chat",
	}

	dto := message.toDTO()

	if dto.ID != message.ID {
		t.Fatalf("expected ID %d, got %d", message.ID, dto.ID)
	}
	if dto.EventID != message.EventID {
		t.Fatalf("expected eventID %d, got %d", message.EventID, dto.EventID)
	}
	if dto.UserID != message.UserID {
		t.Fatalf("expected userID %d, got %d", message.UserID, dto.UserID)
	}
	if dto.Content != message.Content {
		t.Fatalf("expected content %q, got %q", message.Content, dto.Content)
	}
	if !dto.CreatedAt.Equal(message.CreatedAt) {
		t.Fatalf("expected createdAt %s, got %s", message.CreatedAt, dto.CreatedAt)
	}
}

func TestMessagesToDTOsOldestFirst(t *testing.T) {
	messages := []Message{
		{Model: gorm.Model{ID: 3}, Content: "newest"},
		{Model: gorm.Model{ID: 2}, Content: "middle"},
		{Model: gorm.Model{ID: 1}, Content: "oldest"},
	}

	dtos := messagesToDTOsOldestFirst(messages)

	if len(dtos) != len(messages) {
		t.Fatalf("expected %d DTOs, got %d", len(messages), len(dtos))
	}

	expectedIDs := []uint{1, 2, 3}
	for i, expectedID := range expectedIDs {
		if dtos[i].ID != expectedID {
			t.Fatalf("expected DTO at index %d to have ID %d, got %d", i, expectedID, dtos[i].ID)
		}
	}
}

func TestNormalizeMessageLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{name: "uses default for zero", limit: 0, want: messageHistoryLimit},
		{name: "uses default for negative", limit: -1, want: messageHistoryLimit},
		{name: "keeps positive value under max", limit: 25, want: 25},
		{name: "keeps max value", limit: messageHistoryLimit, want: messageHistoryLimit},
		{name: "uses default above max", limit: messageHistoryLimit + 1, want: messageHistoryLimit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeMessageLimit(tt.limit)
			if got != tt.want {
				t.Fatalf("expected limit %d, got %d", tt.want, got)
			}
		})
	}
}
