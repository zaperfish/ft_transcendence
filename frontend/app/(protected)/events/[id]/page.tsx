'use client';
import { useParams, useRouter } from "next/navigation";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getEventById, getEventParticipants, updateEvent, deleteEvent, removeParticipant, leaveEvent } from "@/lib/api/events";
import { useAuth } from "@/lib/hooks/useAuth";
import EditEventModal from "@/components/features/events/EditEventModal";
import { useState } from "react";
// Import necessary UI later

export default function EventDetailPage() {
  return <div>Event detail page Coming Soon</div>;
}