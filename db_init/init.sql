--
-- PostgreSQL database dump
--

-- Dumped from database version 15.7 (Debian 15.7-1.pgdg120+1)
-- Dumped by pg_dump version 16.3

-- Started on 2024-08-11 22:20:24

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- TOC entry 4 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: pg_database_owner
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO pg_database_owner;

--
-- TOC entry 3356 (class 0 OID 0)
-- Dependencies: 4
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: pg_database_owner
--

COMMENT ON SCHEMA public IS 'standard public schema';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 216 (class 1259 OID 16395)
-- Name: tokens; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.tokens (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    token character varying(300) NOT NULL
);


ALTER TABLE public.tokens OWNER TO baseuser;

--
-- TOC entry 215 (class 1259 OID 16394)
-- Name: tokens_id_seq; Type: SEQUENCE; Schema: public; Owner: baseuser
--

CREATE SEQUENCE public.tokens_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tokens_id_seq OWNER TO baseuser;

--
-- TOC entry 3357 (class 0 OID 0)
-- Dependencies: 215
-- Name: tokens_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: baseuser
--

ALTER SEQUENCE public.tokens_id_seq OWNED BY public.tokens.id;


--
-- TOC entry 214 (class 1259 OID 16389)
-- Name: users; Type: TABLE; Schema: public; Owner: baseuser
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    ip character varying(100) NOT NULL
);


ALTER TABLE public.users OWNER TO baseuser;

--
-- TOC entry 3203 (class 2604 OID 16398)
-- Name: tokens id; Type: DEFAULT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.tokens ALTER COLUMN id SET DEFAULT nextval('public.tokens_id_seq'::regclass);


--
-- TOC entry 3207 (class 2606 OID 16400)
-- Name: tokens tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_pkey PRIMARY KEY (id);


--
-- TOC entry 3205 (class 2606 OID 16393)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- TOC entry 3208 (class 2606 OID 16401)
-- Name: tokens tokens_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: baseuser
--

ALTER TABLE ONLY public.tokens
    ADD CONSTRAINT tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


-- Completed on 2024-08-11 22:20:25

--
-- PostgreSQL database dump complete
--

