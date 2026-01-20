--
-- PostgreSQL database dump
--

\restrict QLJ6eXabXaoWCBrSbBqRbSbJhkXkd9Z8WKWlghItBJR5YNv8EwO3mIMnC5pKYHS

-- Dumped from database version 18.1
-- Dumped by pg_dump version 18.1

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
-- Name: gift_ratings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.gift_ratings (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    gift_id bigint NOT NULL,
    rate numeric,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.gift_ratings OWNER TO postgres;

--
-- Name: gift_ratings_gift_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.gift_ratings_gift_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.gift_ratings_gift_id_seq OWNER TO postgres;

--
-- Name: gift_ratings_gift_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.gift_ratings_gift_id_seq OWNED BY public.gift_ratings.gift_id;


--
-- Name: gift_ratings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.gift_ratings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.gift_ratings_id_seq OWNER TO postgres;

--
-- Name: gift_ratings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.gift_ratings_id_seq OWNED BY public.gift_ratings.id;


--
-- Name: gift_ratings_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.gift_ratings_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.gift_ratings_user_id_seq OWNER TO postgres;

--
-- Name: gift_ratings_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.gift_ratings_user_id_seq OWNED BY public.gift_ratings.user_id;


--
-- Name: gifts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.gifts (
    id bigint NOT NULL,
    title text,
    description text,
    stock integer,
    redeem_point integer NOT NULL,
    image_url text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.gifts OWNER TO postgres;

--
-- Name: gifts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.gifts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.gifts_id_seq OWNER TO postgres;

--
-- Name: gifts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.gifts_id_seq OWNED BY public.gifts.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username character varying(100) NOT NULL,
    password text NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: gift_ratings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gift_ratings ALTER COLUMN id SET DEFAULT nextval('public.gift_ratings_id_seq'::regclass);


--
-- Name: gift_ratings user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gift_ratings ALTER COLUMN user_id SET DEFAULT nextval('public.gift_ratings_user_id_seq'::regclass);


--
-- Name: gift_ratings gift_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gift_ratings ALTER COLUMN gift_id SET DEFAULT nextval('public.gift_ratings_gift_id_seq'::regclass);


--
-- Name: gifts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gifts ALTER COLUMN id SET DEFAULT nextval('public.gifts_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: gift_ratings; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.gift_ratings (id, user_id, gift_id, rate, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: gifts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.gifts (id, title, description, stock, redeem_point, image_url, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version, dirty) FROM stdin;
3	f
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, password, created_at, updated_at, deleted_at) FROM stdin;
\.


--
-- Name: gift_ratings_gift_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.gift_ratings_gift_id_seq', 1, false);


--
-- Name: gift_ratings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.gift_ratings_id_seq', 1, false);


--
-- Name: gift_ratings_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.gift_ratings_user_id_seq', 1, false);


--
-- Name: gifts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.gifts_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: gift_ratings pk_gift_ratings_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gift_ratings
    ADD CONSTRAINT pk_gift_ratings_id PRIMARY KEY (id);


--
-- Name: gifts pk_gifts_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gifts
    ADD CONSTRAINT pk_gifts_id PRIMARY KEY (id);


--
-- Name: users pk_users_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT pk_users_id PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: gift_ratings_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX gift_ratings_active ON public.gift_ratings USING btree (user_id, gift_id) WHERE (deleted_at IS NULL);


--
-- Name: idx_gift_ratings_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_gift_ratings_active ON public.gift_ratings USING btree (id) WHERE (deleted_at IS NULL);


--
-- Name: idx_gifts_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_gifts_active ON public.gifts USING btree (id) WHERE (deleted_at IS NULL);


--
-- Name: idx_users_active; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_active ON public.users USING btree (id) WHERE (deleted_at IS NULL);


--
-- Name: gift_ratings fk_gift_ratings_gift_id_gifts_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gift_ratings
    ADD CONSTRAINT fk_gift_ratings_gift_id_gifts_id FOREIGN KEY (gift_id) REFERENCES public.gifts(id);


--
-- Name: gift_ratings fk_gift_ratings_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gift_ratings
    ADD CONSTRAINT fk_gift_ratings_user_id_users_id FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict QLJ6eXabXaoWCBrSbBqRbSbJhkXkd9Z8WKWlghItBJR5YNv8EwO3mIMnC5pKYHS

