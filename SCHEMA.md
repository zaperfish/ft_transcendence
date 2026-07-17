# Database Schema Overview

We are generally using soft deletes. This means that deleted entries dont leave the database, a deleted_at column gets updated for the respective entry, indicating that it is considered to be deleted and consequently ignored by the application.

## users

This is the database representation of users.

| Column | Data Type | Nullable |
|---|---|---|
| id | bigint | NO |
| created_at | timestamp with time zone | YES |
| updated_at | timestamp with time zone | YES |
| deleted_at | timestamp with time zone | YES |
| name | text | YES |
| email | text | YES |
| password_hash | text | YES |

## events

This is the database representation of events.

| Column | Data Type | Nullable |
|---|---|---|
| id | bigint | NO |
| created_at | timestamp with time zone | YES |
| updated_at | timestamp with time zone | YES |
| deleted_at | timestamp with time zone | YES |
| title | character varying | NO |
| description | text | YES |
| start_time | timestamp with time zone | NO |
| duration | smallint | NO |
| location_name | character varying | YES |
| location_address | character varying | YES |
| max_capacity | bigint | NO |
| image_path | character varying | YES |

## event_users

This is a jointable describing the relationships between events and users. An entry means, that the user with the respective ID has joined the event with the respective id. The role column denotes the role of the user within the event. This can be "admin" or "member"

| Column | Data Type | Nullable |
|---|---|---|
| id | bigint | NO |
| created_at | timestamp with time zone | YES |
| updated_at | timestamp with time zone | YES |
| deleted_at | timestamp with time zone | YES |
| user_id | bigint | NO |
| event_id | bigint | NO |
| role | text | NO |

## messages

This is the database representation of chat messages.

| Column | Data Type | Nullable |
|---|---|---|
| id | bigint | NO |
| created_at | timestamp with time zone | YES |
| updated_at | timestamp with time zone | YES |
| deleted_at | timestamp with time zone | YES |
| event_id | bigint | NO |
| user_id | bigint | NO |
| content | text | NO |

## api_keys

This is the databaserepresentation of API keys used for the public API.

| Column | Data Type | Nullable |
|---|---|---|
| id | bigint | NO |
| created_at | timestamp with time zone | YES |
| updated_at | timestamp with time zone | YES |
| deleted_at | timestamp with time zone | YES |
| key_hash | text | YES |
| revoked | boolean | YES |
