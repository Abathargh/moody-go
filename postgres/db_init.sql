--
-- PostgreSQL database dump
--

-- Dumped from database version 12.3 (Debian 12.3-1.pgdg100+1)
-- Dumped by pg_dump version 12.3 (Debian 12.3-1.pgdg100+1)

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
-- Name: moody; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE moody WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE moody OWNER TO postgres;

\connect moody

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

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: datatype; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.datatype (
    id integer NOT NULL,
    type character varying(6)
);


ALTER TABLE public.datatype OWNER TO postgres;

--
-- Name: datatype_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.datatype_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.datatype_id_seq OWNER TO postgres;

--
-- Name: datatype_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.datatype_id_seq OWNED BY public.datatype.id;


--
-- Name: service; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service (
    id integer NOT NULL,
    name character varying(50),
    datatype integer,
    state integer
);


ALTER TABLE public.service OWNER TO postgres;

--
-- Name: service_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.service_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.service_id_seq OWNER TO postgres;

--
-- Name: service_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.service_id_seq OWNED BY public.service.id;


--
-- Name: service_state; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.service_state (
    id integer NOT NULL,
    state character varying(3)
);


ALTER TABLE public.service_state OWNER TO postgres;

--
-- Name: situation; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.situation (
    id integer NOT NULL,
    name character varying(100)
);


ALTER TABLE public.situation OWNER TO postgres;

--
-- Name: situation_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.situation_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.situation_id_seq OWNER TO postgres;

--
-- Name: situation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.situation_id_seq OWNED BY public.situation.id;


--
-- Name: datatype id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.datatype ALTER COLUMN id SET DEFAULT nextval('public.datatype_id_seq'::regclass);


--
-- Name: service id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service ALTER COLUMN id SET DEFAULT nextval('public.service_id_seq'::regclass);


--
-- Name: situation id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.situation ALTER COLUMN id SET DEFAULT nextval('public.situation_id_seq'::regclass);


--
-- Name: datatype datatype_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.datatype
    ADD CONSTRAINT datatype_pk PRIMARY KEY (id);


--
-- Name: service service_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service
    ADD CONSTRAINT service_pk PRIMARY KEY (id);


--
-- Name: service_state service_state_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service_state
    ADD CONSTRAINT service_state_pk PRIMARY KEY (id);


--
-- Name: situation situation_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.situation
    ADD CONSTRAINT situation_pk PRIMARY KEY (id);


--
-- Name: datatype_type_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX datatype_type_uindex ON public.datatype USING btree (type);


--
-- Name: service_name_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX service_name_uindex ON public.service USING btree (name);


--
-- Name: service_state_state_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX service_state_state_uindex ON public.service_state USING btree (state);


--
-- Name: situation_name_uindex; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX situation_name_uindex ON public.situation USING btree (name);


--
-- Name: service service_datatype_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service
    ADD CONSTRAINT service_datatype_id_fk FOREIGN KEY (datatype) REFERENCES public.datatype(id);


--
-- Name: service service_service_state_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.service
    ADD CONSTRAINT service_service_state_id_fk FOREIGN KEY (state) REFERENCES public.service_state(id);


--
-- PostgreSQL database dump complete
--


--
-- Data for Name: datatype; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.datatype (id, type) FROM stdin;
0	int
1	float
2	string
\.


--
-- Data for Name: service_state; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.service_state (id, state) FROM stdin;
0	off
1	on
\.

