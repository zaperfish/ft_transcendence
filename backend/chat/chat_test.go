package chat

import (
	// Std
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	// Extern
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

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

func TestHubJoinRoomReusesRoom(t *testing.T) {
	hub := NewHub()
	firstClient := &Client{send: make(chan MessageDTO)}
	secondClient := &Client{send: make(chan MessageDTO)}
	otherClient := &Client{send: make(chan MessageDTO)}

	firstRoom := hub.JoinRoom(42, firstClient)
	secondRoom := hub.JoinRoom(42, secondClient)
	otherRoom := hub.JoinRoom(43, otherClient)

	if firstRoom != secondRoom {
		t.Fatal("expected same event ID to reuse the existing room")
	}
	if firstRoom == otherRoom {
		t.Fatal("expected different event IDs to use different rooms")
	}

	firstRoom.Leave(firstClient)
	secondRoom.Leave(secondClient)
	otherRoom.Leave(otherClient)
}

func TestHubJoinRoomStartsRoomRunLoop(t *testing.T) {
	hub := NewHub()
	client := &Client{send: make(chan MessageDTO)}

	room := hub.JoinRoom(42, client)
	room.Leave(client)

	select {
	case _, ok := <-client.send:
		if ok {
			t.Fatal("expected client send channel to be closed after leaving")
		}
	case <-time.After(time.Second):
		t.Fatal("expected leaving client send channel to be closed")
	}
}

func TestHubRemovesRoomWhenLastClientLeaves(t *testing.T) {
	hub := NewHub()
	client := &Client{send: make(chan MessageDTO)}

	room := hub.JoinRoom(42, client)
	room.Leave(client)
	waitForRoomClosed(t, room)

	hub.mu.Lock()
	defer hub.mu.Unlock()

	if _, ok := hub.rooms[42]; ok {
		t.Fatal("expected room to be removed after last client leaves")
	}
}

func TestHubRecreatesRoomAfterPreviousRoomClosed(t *testing.T) {
	hub := NewHub()
	firstClient := &Client{send: make(chan MessageDTO)}

	firstRoom := hub.JoinRoom(42, firstClient)
	firstRoom.Leave(firstClient)
	waitForRoomClosed(t, firstRoom)

	secondClient := &Client{send: make(chan MessageDTO)}
	secondRoom := hub.JoinRoom(42, secondClient)
	defer secondRoom.Leave(secondClient)

	if firstRoom == secondRoom {
		t.Fatal("expected a new room after the previous room closed")
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
	if room.removeUser == nil {
		t.Fatal("expected removeUser channel to be initialized")
	}
	if room.broadcast == nil {
		t.Fatal("expected broadcast channel to be initialized")
	}
	if room.done == nil {
		t.Fatal("expected done channel to be initialized")
	}
}

func TestRoomRunBroadcastsToJoinedClients(t *testing.T) {
	room := NewRoom(42)
	client := &Client{send: make(chan MessageDTO, 1)}
	message := MessageDTO{
		EventID:    42,
		UserID:     3,
		SenderName: "sender",
		Content:    "hello",
	}

	go room.run()

	if !room.Join(client) {
		t.Fatal("expected room run loop to receive joined client")
	}

	if !room.Broadcast(message) {
		t.Fatal("expected room run loop to receive broadcast message")
	}

	select {
	case got := <-client.send:
		if got.Content != message.Content {
			t.Fatalf("expected message content %q, got %q", message.Content, got.Content)
		}
	case <-time.After(time.Second):
		t.Fatal("expected joined client to receive broadcast message")
	}

	room.Leave(client)
	waitForRoomClosed(t, room)
}

func TestRoomRunRemovesAllConnectionsForUser(t *testing.T) {
	room := NewRoom(42)
	firstRemoved := &Client{userID: 3, send: make(chan MessageDTO, 1)}
	secondRemoved := &Client{userID: 3, send: make(chan MessageDTO, 1)}
	remaining := &Client{userID: 4, send: make(chan MessageDTO, 1)}

	go room.run()

	if !room.Join(firstRemoved) || !room.Join(secondRemoved) || !room.Join(remaining) {
		t.Fatal("expected clients to join the room")
	}
	if !room.RemoveUser(3) {
		t.Fatal("expected room to receive user removal event")
	}

	message := MessageDTO{Content: "still connected"}
	if !room.Broadcast(message) {
		t.Fatal("expected room to remain active for the other user")
	}

	assertClientSendClosed(t, firstRemoved)
	assertClientSendClosed(t, secondRemoved)

	select {
	case got := <-remaining.send:
		if got.Content != message.Content {
			t.Fatalf("expected message content %q, got %q", message.Content, got.Content)
		}
	case <-time.After(time.Second):
		t.Fatal("expected remaining client to receive broadcast")
	}

	room.Leave(remaining)
	waitForRoomClosed(t, room)
}

func TestHubDisconnectParticipantDoesNotCreateRoom(t *testing.T) {
	hub := NewHub()

	hub.DisconnectParticipant(42, 3)

	hub.mu.Lock()
	defer hub.mu.Unlock()
	if len(hub.rooms) != 0 {
		t.Fatalf("expected no rooms, got %d", len(hub.rooms))
	}
}

func TestHubDisconnectParticipantRemovesConnectedUser(t *testing.T) {
	hub := NewHub()
	client := &Client{userID: 3, send: make(chan MessageDTO)}

	room := hub.JoinRoom(42, client)
	hub.DisconnectParticipant(42, 3)

	waitForRoomClosed(t, room)
	assertClientSendClosed(t, client)

	hub.mu.Lock()
	defer hub.mu.Unlock()
	if len(hub.rooms) != 0 {
		t.Fatalf("expected no rooms, got %d", len(hub.rooms))
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

func TestNewMessageFromInput(t *testing.T) {
	message, err := newMessageFromInput(42, 3, createMessageInput{
		Content: " hello chat ",
	})
	if err != nil {
		t.Fatalf("expected message, got error: %v", err)
	}

	if message.EventID != 42 {
		t.Fatalf("expected eventID 42, got %d", message.EventID)
	}
	if message.UserID != 3 {
		t.Fatalf("expected userID 3, got %d", message.UserID)
	}
	if message.Content != "hello chat" {
		t.Fatalf("expected trimmed content %q, got %q", "hello chat", message.Content)
	}
}

func TestNewMessageFromInputRejectsEmptyContent(t *testing.T) {
	_, err := newMessageFromInput(42, 3, createMessageInput{
		Content: "   ",
	})
	if err == nil {
		t.Fatal("expected empty message content error")
	}
}

func TestNewMessageFromInputAcceptsContentAtLimits(t *testing.T) {
	content := strings.Repeat("🙂", maxMessageCharacters)

	message, err := newMessageFromInput(42, 3, createMessageInput{Content: content})
	if err != nil {
		t.Fatalf("expected message at limits to be accepted, got error: %v", err)
	}
	if message.Content != content {
		t.Fatal("expected message content to be preserved")
	}
}

func TestNewMessageFromInputRejectsTooManyCharacters(t *testing.T) {
	content := strings.Repeat("a", maxMessageCharacters+1)

	_, err := newMessageFromInput(42, 3, createMessageInput{Content: content})
	if err == nil {
		t.Fatal("expected overlong message content error")
	}
	if !strings.Contains(err.Error(), "characters") {
		t.Fatalf("expected character limit error, got: %v", err)
	}
}

func TestNewMessageFromInputRejectsTooManyBytes(t *testing.T) {
	content := strings.Repeat("é", maxMessageBytes/2+1)

	_, err := newMessageFromInput(42, 3, createMessageInput{Content: content})
	if err == nil {
		t.Fatal("expected oversized message content error")
	}
	if !strings.Contains(err.Error(), "bytes") {
		t.Fatalf("expected byte limit error, got: %v", err)
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

	dto := message.toDTO("sender")

	if dto.ID != message.ID {
		t.Fatalf("expected ID %d, got %d", message.ID, dto.ID)
	}
	if dto.EventID != message.EventID {
		t.Fatalf("expected eventID %d, got %d", message.EventID, dto.EventID)
	}
	if dto.UserID != message.UserID {
		t.Fatalf("expected userID %d, got %d", message.UserID, dto.UserID)
	}
	if dto.SenderName != "sender" {
		t.Fatalf("expected sender name %q, got %q", "sender", dto.SenderName)
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

	senderNames := map[uint]string{
		0: "sender",
	}
	dtos := messagesToDTOsOldestFirst(messages, senderNames)

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

func TestHandleEventChatWebSocketRejectsInvalidEventID(t *testing.T) {
	handler := NewHandler(nil)
	req := newChatWebSocketRequest("invalid")
	recorder := httptest.NewRecorder()

	handler.handleEventChatWebSocket(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func TestHandleEventChatWebSocketRejectsMissingAuthCookie(t *testing.T) {
	handler := NewHandler(nil)
	req := newChatWebSocketRequest("42")
	recorder := httptest.NewRecorder()

	handler.handleEventChatWebSocket(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func newChatWebSocketRequest(eventID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/api/events/"+eventID+"/chat/ws", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", eventID)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}

func waitForRoomClosed(t *testing.T, room *Room) {
	t.Helper()

	select {
	case <-room.done:
	case <-time.After(time.Second):
		t.Fatal("expected room to close")
	}
}

func assertClientSendClosed(t *testing.T, client *Client) {
	t.Helper()

	select {
	case _, ok := <-client.send:
		if ok {
			t.Fatal("expected client send channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("expected client send channel to be closed")
	}
}
