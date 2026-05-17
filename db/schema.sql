--
-- PostgreSQL database dump
--

\restrict jTjCj6kJVjY7VbK6FZ3jMMopxGIui0v1xumG3URbjP0OvQrgK1ubwRhVTuHshdp

-- Dumped from database version 18.3 (Debian 18.3-1.pgdg13+1)
-- Dumped by pg_dump version 18.3 (Debian 18.3-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: event_labels; Type: TABLE; Schema: public; Owner: ft_transcendence_user
--

CREATE TABLE public.event_labels (
    event_id bigint NOT NULL,
    label_id bigint NOT NULL
);


ALTER TABLE public.event_labels OWNER TO ft_transcendence_user;

--
-- Name: events; Type: TABLE; Schema: public; Owner: ft_transcendence_user
--

CREATE TABLE public.events (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text NOT NULL,
    description text,
    start_time timestamp with time zone,
    duration bigint,
    location_name text,
    location_address text,
    max_capacity bigint,
    num_registered bigint,
    CONSTRAINT chk_events_max_capacity CHECK ((max_capacity >= 0)),
    CONSTRAINT chk_events_num_registered CHECK ((max_capacity >= 0)),
    CONSTRAINT chk_events_title CHECK ((length(title) >= 3))
);


ALTER TABLE public.events OWNER TO ft_transcendence_user;

--
-- Name: events_id_seq; Type: SEQUENCE; Schema: public; Owner: ft_transcendence_user
--

CREATE SEQUENCE public.events_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.events_id_seq OWNER TO ft_transcendence_user;

--
-- Name: events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ft_transcendence_user
--

ALTER SEQUENCE public.events_id_seq OWNED BY public.events.id;


--
-- Name: labels; Type: TABLE; Schema: public; Owner: ft_transcendence_user
--

CREATE TABLE public.labels (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text NOT NULL,
    CONSTRAINT chk_labels_name CHECK ((length(name) >= 1))
);


ALTER TABLE public.labels OWNER TO ft_transcendence_user;

--
-- Name: labels_id_seq; Type: SEQUENCE; Schema: public; Owner: ft_transcendence_user
--

CREATE SEQUENCE public.labels_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.labels_id_seq OWNER TO ft_transcendence_user;

--
-- Name: labels_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ft_transcendence_user
--

ALTER SEQUENCE public.labels_id_seq OWNED BY public.labels.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: ft_transcendence_user
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text
);


ALTER TABLE public.users OWNER TO ft_transcendence_user;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: ft_transcendence_user
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO ft_transcendence_user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ft_transcendence_user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: events id; Type: DEFAULT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.events ALTER COLUMN id SET DEFAULT nextval('public.events_id_seq'::regclass);


--
-- Name: labels id; Type: DEFAULT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.labels ALTER COLUMN id SET DEFAULT nextval('public.labels_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: event_labels event_labels_pkey; Type: CONSTRAINT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.event_labels
    ADD CONSTRAINT event_labels_pkey PRIMARY KEY (event_id, label_id);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: labels labels_pkey; Type: CONSTRAINT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.labels
    ADD CONSTRAINT labels_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_events_deleted_at; Type: INDEX; Schema: public; Owner: ft_transcendence_user
--

CREATE INDEX idx_events_deleted_at ON public.events USING btree (deleted_at);


--
-- Name: idx_labels_deleted_at; Type: INDEX; Schema: public; Owner: ft_transcendence_user
--

CREATE INDEX idx_labels_deleted_at ON public.labels USING btree (deleted_at);


--
-- Name: idx_labels_name; Type: INDEX; Schema: public; Owner: ft_transcendence_user
--

CREATE UNIQUE INDEX idx_labels_name ON public.labels USING btree (name) WHERE (deleted_at IS NULL);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: ft_transcendence_user
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: event_labels fk_event_labels_event; Type: FK CONSTRAINT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.event_labels
    ADD CONSTRAINT fk_event_labels_event FOREIGN KEY (event_id) REFERENCES public.events(id);


--
-- Name: event_labels fk_event_labels_label; Type: FK CONSTRAINT; Schema: public; Owner: ft_transcendence_user
--

ALTER TABLE ONLY public.event_labels
    ADD CONSTRAINT fk_event_labels_label FOREIGN KEY (label_id) REFERENCES public.labels(id);


--
-- PostgreSQL database dump complete
--

\unrestrict jTjCj6kJVjY7VbK6FZ3jMMopxGIui0v1xumG3URbjP0OvQrgK1ubwRhVTuHshdp

