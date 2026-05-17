--
-- PostgreSQL database dump
--

\restrict jpkuQs1ngiJ6Adbg16s3cymMbqvqqV80yGQJMegGyEdKsOBtkKJc9Es4b0q2qAF

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

--
-- Data for Name: events; Type: TABLE DATA; Schema: public; Owner: ft_transcendence_user
--

COPY public.events (id, created_at, updated_at, deleted_at, title, description, start_time, duration, location_name, location_address, max_capacity, num_registered) FROM stdin;
1	2026-05-17 20:21:24.555699+00	2026-05-17 21:19:44.714214+00	\N	Go Meetup Berlin	A monthly meetup for Go developers	2026-06-15 18:00:00+00	120	Betahaus	Prinzessinnenstraße 19, 10969 Berlin	100	0
2	2026-05-17 20:22:26.428689+00	2026-05-17 21:19:44.736159+00	\N	Rust Workshop for Beginners	A hands-on introduction to Rust, covering ownership, borrowing and basic data structures	2026-06-22 10:00:00+00	240	Factory Berlin	Rheinsberger Str. 76, 10115 Berlin	30	0
3	2026-05-17 20:22:34.943521+00	2026-05-17 21:19:44.763289+00	\N	TypeScript & Node.js Hackathon	A full day hackathon building real world APIs with TypeScript and Node.js, teams of up to 4	2026-07-05 09:00:00+00	480	Mindspace Moritzplatz	Oranienstraße 185, 10999 Berlin	80	0
4	2026-05-17 20:22:44.322143+00	2026-05-17 21:19:44.783056+00	\N	Linux Kernel Development Talk	An evening talk on contributing to the Linux kernel, covering the patch submission process and common pitfalls	2026-07-10 19:00:00+00	90	c-base	Rungestraße 20, 10179 Berlin	50	0
\.


--
-- Data for Name: labels; Type: TABLE DATA; Schema: public; Owner: ft_transcendence_user
--

COPY public.labels (id, created_at, updated_at, deleted_at, name) FROM stdin;
70	2026-05-17 21:16:31.682461+00	2026-05-17 21:16:31.682461+00	2026-05-17 21:17:29.269247+00	Go
73	2026-05-17 21:17:40.335485+00	2026-05-17 21:17:40.335485+00	2026-05-17 21:17:59.341535+00	Go
74	2026-05-17 21:18:49.26122+00	2026-05-17 21:18:49.26122+00	\N	go
75	2026-05-17 21:18:49.273296+00	2026-05-17 21:18:49.273296+00	\N	rust
76	2026-05-17 21:18:49.282232+00	2026-05-17 21:18:49.282232+00	\N	python
77	2026-05-17 21:18:49.290611+00	2026-05-17 21:18:49.290611+00	\N	typescript
78	2026-05-17 21:18:49.29873+00	2026-05-17 21:18:49.29873+00	\N	javascript
80	2026-05-17 21:18:49.317983+00	2026-05-17 21:18:49.317983+00	\N	zig
81	2026-05-17 21:18:49.327757+00	2026-05-17 21:18:49.327757+00	\N	linux
82	2026-05-17 21:18:49.336186+00	2026-05-17 21:18:49.336186+00	\N	web
83	2026-05-17 21:18:49.346942+00	2026-05-17 21:18:49.346942+00	\N	backend
84	2026-05-17 21:18:49.355349+00	2026-05-17 21:18:49.355349+00	\N	frontend
85	2026-05-17 21:18:49.362262+00	2026-05-17 21:18:49.362262+00	\N	devops
86	2026-05-17 21:18:49.36861+00	2026-05-17 21:18:49.36861+00	\N	security
87	2026-05-17 21:18:49.375048+00	2026-05-17 21:18:49.375048+00	\N	embedded
88	2026-05-17 21:18:49.381202+00	2026-05-17 21:18:49.381202+00	\N	gamedev
1	2026-05-17 20:27:34.18999+00	2026-05-17 20:27:34.18999+00	2026-05-17 20:33:37.468466+00	Go
2	2026-05-17 20:29:04.499287+00	2026-05-17 20:29:04.499287+00	2026-05-17 20:47:51.328574+00	go
3	2026-05-17 20:29:04.508476+00	2026-05-17 20:29:04.508476+00	2026-05-17 20:51:30.774496+00	rust
4	2026-05-17 20:29:04.518704+00	2026-05-17 20:29:04.518704+00	2026-05-17 20:51:30.783975+00	python
5	2026-05-17 20:29:04.528284+00	2026-05-17 20:29:04.528284+00	2026-05-17 20:51:30.792059+00	typescript
6	2026-05-17 20:29:04.537202+00	2026-05-17 20:29:04.537202+00	2026-05-17 20:51:30.799299+00	javascript
8	2026-05-17 20:29:04.552535+00	2026-05-17 20:29:04.552535+00	2026-05-17 20:51:30.811623+00	zig
9	2026-05-17 20:29:04.560109+00	2026-05-17 20:29:04.560109+00	2026-05-17 20:51:30.817948+00	linux
10	2026-05-17 20:29:04.568165+00	2026-05-17 20:29:04.568165+00	2026-05-17 20:51:30.824944+00	web
11	2026-05-17 20:29:04.57612+00	2026-05-17 20:29:04.57612+00	2026-05-17 20:51:30.831294+00	backend
12	2026-05-17 20:29:04.582226+00	2026-05-17 20:29:04.582226+00	2026-05-17 20:51:30.837591+00	frontend
13	2026-05-17 20:29:04.588609+00	2026-05-17 20:29:04.588609+00	2026-05-17 20:51:30.844579+00	devops
14	2026-05-17 20:29:04.59524+00	2026-05-17 20:29:04.59524+00	2026-05-17 20:51:30.851923+00	security
15	2026-05-17 20:29:04.60228+00	2026-05-17 20:29:04.60228+00	2026-05-17 20:51:30.859686+00	embedded
16	2026-05-17 20:29:04.608508+00	2026-05-17 20:29:04.608508+00	2026-05-17 20:51:30.867471+00	gamedev
17	2026-05-17 20:29:04.614785+00	2026-05-17 20:29:04.614785+00	2026-05-17 20:51:30.875541+00	ai
18	2026-05-17 20:29:04.620552+00	2026-05-17 20:29:04.620552+00	2026-05-17 20:51:30.885592+00	workshop
19	2026-05-17 20:29:04.62752+00	2026-05-17 20:29:04.62752+00	2026-05-17 20:51:30.894033+00	talk
20	2026-05-17 20:29:04.633952+00	2026-05-17 20:29:04.633952+00	2026-05-17 20:51:30.904593+00	hackathon
21	2026-05-17 20:29:04.640801+00	2026-05-17 20:29:04.640801+00	2026-05-17 20:51:30.91433+00	beginner
22	2026-05-17 20:29:04.649014+00	2026-05-17 20:29:04.649014+00	2026-05-17 20:51:30.921904+00	advanced
23	2026-05-17 20:29:04.660516+00	2026-05-17 20:29:04.660516+00	2026-05-17 20:51:30.92871+00	networking
89	2026-05-17 21:18:49.38734+00	2026-05-17 21:18:49.38734+00	\N	ai
90	2026-05-17 21:18:49.393278+00	2026-05-17 21:18:49.393278+00	\N	workshop
91	2026-05-17 21:18:49.399554+00	2026-05-17 21:18:49.399554+00	\N	talk
92	2026-05-17 21:18:49.405518+00	2026-05-17 21:18:49.405518+00	\N	hackathon
93	2026-05-17 21:18:49.412225+00	2026-05-17 21:18:49.412225+00	\N	beginner
94	2026-05-17 21:18:49.418264+00	2026-05-17 21:18:49.418264+00	\N	advanced
95	2026-05-17 21:18:49.424197+00	2026-05-17 21:18:49.424197+00	\N	networking
\.


--
-- Data for Name: event_labels; Type: TABLE DATA; Schema: public; Owner: ft_transcendence_user
--

COPY public.event_labels (event_id, label_id) FROM stdin;
1	1
1	18
2	2
2	20
2	17
3	4
3	5
3	10
3	19
4	8
4	6
4	21
4	18
1	2
1	11
1	19
2	21
2	18
3	6
3	11
3	20
4	9
4	22
4	19
1	74
1	83
1	91
2	75
2	93
2	90
3	77
3	78
3	83
3	92
4	81
4	94
4	91
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: ft_transcendence_user
--

COPY public.users (id, created_at, updated_at, deleted_at, name) FROM stdin;
\.


--
-- Name: events_id_seq; Type: SEQUENCE SET; Schema: public; Owner: ft_transcendence_user
--

SELECT pg_catalog.setval('public.events_id_seq', 4, true);


--
-- Name: labels_id_seq; Type: SEQUENCE SET; Schema: public; Owner: ft_transcendence_user
--

SELECT pg_catalog.setval('public.labels_id_seq', 95, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: ft_transcendence_user
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- PostgreSQL database dump complete
--

\unrestrict jpkuQs1ngiJ6Adbg16s3cymMbqvqqV80yGQJMegGyEdKsOBtkKJc9Es4b0q2qAF

