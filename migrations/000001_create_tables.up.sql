--
-- PostgreSQL database dump
--

-- Dumped from database version 14.4
-- Dumped by pg_dump version 14.4

-- Started on 2022-07-26 16:23:57

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
-- TOC entry 3379 (class 1262 OID 24791)
-- Name: notifier_server; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE notifier_server WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'Russian_Russia.1251';


ALTER DATABASE notifier_server OWNER TO postgres;

\connect notifier_server

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
-- TOC entry 843 (class 1247 OID 24846)
-- Name: command_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.command_type AS ENUM (
    'post',
    'put',
    'delete'
);


ALTER TYPE public.command_type OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 214 (class 1259 OID 24822)
-- Name: ns_journal; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ns_journal (
                                   id integer NOT NULL,
                                   message_id integer NOT NULL,
                                   user_id integer NOT NULL,
                                   resource_id integer NOT NULL,
                                   response_status integer,
                                   response_content character varying(300),
                                   description character varying(100) DEFAULT ''::character varying NOT NULL,
                                   created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.ns_journal OWNER TO postgres;

--
-- TOC entry 213 (class 1259 OID 24821)
-- Name: ns_journal_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ns_journal_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ns_journal_id_seq OWNER TO postgres;

--
-- TOC entry 3380 (class 0 OID 0)
-- Dependencies: 213
-- Name: ns_journal_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ns_journal_id_seq OWNED BY public.ns_journal.id;


--
-- TOC entry 210 (class 1259 OID 24793)
-- Name: ns_messages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ns_messages (
                                    id integer NOT NULL,
                                    code character varying(32) DEFAULT ''::character varying NOT NULL,
                                    user_id integer NOT NULL,
                                    resource_id integer NOT NULL,
                                    command character varying(10) DEFAULT ''::character varying NOT NULL,
                                    priority character varying(10) DEFAULT 'normal'::character varying NOT NULL,
                                    is_send_callback boolean DEFAULT false NOT NULL,
                                    content text DEFAULT ''::text NOT NULL,
                                    is_sent boolean DEFAULT false NOT NULL,
                                    attempt_count integer DEFAULT 3 NOT NULL,
                                    description character varying(100),
                                    send_at timestamp without time zone,
                                    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                    updated_at timestamp without time zone,
                                    deleted_at timestamp without time zone
);


ALTER TABLE public.ns_messages OWNER TO postgres;

--
-- TOC entry 209 (class 1259 OID 24792)
-- Name: ns_messages_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ns_messages_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ns_messages_id_seq OWNER TO postgres;

--
-- TOC entry 3381 (class 0 OID 0)
-- Dependencies: 209
-- Name: ns_messages_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ns_messages_id_seq OWNED BY public.ns_messages.id;


--
-- TOC entry 212 (class 1259 OID 24810)
-- Name: ns_resources; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ns_resources (
                                     id integer NOT NULL,
                                     user_id integer NOT NULL,
                                     url character varying(1000) DEFAULT ''::character varying NOT NULL,
                                     description character varying(100) DEFAULT ''::character varying NOT NULL,
                                     created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                     updated_at timestamp without time zone NOT NULL,
                                     deleted_at timestamp without time zone
);


ALTER TABLE public.ns_resources OWNER TO postgres;

--
-- TOC entry 211 (class 1259 OID 24809)
-- Name: ns_resources_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ns_resources_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ns_resources_id_seq OWNER TO postgres;

--
-- TOC entry 3382 (class 0 OID 0)
-- Dependencies: 211
-- Name: ns_resources_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ns_resources_id_seq OWNED BY public.ns_resources.id;


--
-- TOC entry 216 (class 1259 OID 24831)
-- Name: ns_settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ns_settings (
                                    id integer NOT NULL,
                                    user_id integer NOT NULL,
                                    resource_id integer NOT NULL,
                                    code character varying(32) DEFAULT ''::character varying NOT NULL,
                                    title character varying(100) DEFAULT ''::character varying NOT NULL,
                                    count integer DEFAULT 3 NOT NULL,
                                    intervals integer[] NOT NULL,
                                    timeout integer DEFAULT 3 NOT NULL,
                                    description character varying(100) DEFAULT ''::character varying NOT NULL,
                                    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                    updated_at timestamp without time zone,
                                    deleted_at timestamp without time zone
);


ALTER TABLE public.ns_settings OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 24830)
-- Name: ns_settings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ns_settings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ns_settings_id_seq OWNER TO postgres;

--
-- TOC entry 3383 (class 0 OID 0)
-- Dependencies: 215
-- Name: ns_settings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ns_settings_id_seq OWNED BY public.ns_settings.id;


--
-- TOC entry 217 (class 1259 OID 24853)
-- Name: ns_users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ns_users (
                                 id integer NOT NULL,
                                 code character varying(64) DEFAULT ''::character varying NOT NULL,
                                 role character varying(20) DEFAULT 'service'::character varying NOT NULL,
                                 title character varying(100) DEFAULT ''::character varying NOT NULL,
                                 auth_token character varying(256) DEFAULT ''::character varying NOT NULL,
                                 description character varying(256) DEFAULT ''::character varying NOT NULL,
                                 created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
                                 updated_at timestamp without time zone,
                                 deleted_at timestamp without time zone
);


ALTER TABLE public.ns_users OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 24864)
-- Name: ns_users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ns_users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ns_users_id_seq OWNER TO postgres;

--
-- TOC entry 3384 (class 0 OID 0)
-- Dependencies: 218
-- Name: ns_users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ns_users_id_seq OWNED BY public.ns_users.id;


--
-- TOC entry 3200 (class 2604 OID 24825)
-- Name: ns_journal id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_journal ALTER COLUMN id SET DEFAULT nextval('public.ns_journal_id_seq'::regclass);


--
-- TOC entry 3187 (class 2604 OID 24796)
-- Name: ns_messages id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_messages ALTER COLUMN id SET DEFAULT nextval('public.ns_messages_id_seq'::regclass);


--
-- TOC entry 3196 (class 2604 OID 24813)
-- Name: ns_resources id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_resources ALTER COLUMN id SET DEFAULT nextval('public.ns_resources_id_seq'::regclass);


--
-- TOC entry 3203 (class 2604 OID 24834)
-- Name: ns_settings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_settings ALTER COLUMN id SET DEFAULT nextval('public.ns_settings_id_seq'::regclass);


--
-- TOC entry 3216 (class 2604 OID 24865)
-- Name: ns_users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_users ALTER COLUMN id SET DEFAULT nextval('public.ns_users_id_seq'::regclass);


--
-- TOC entry 3222 (class 2606 OID 24829)
-- Name: ns_journal ns_journal_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_journal
    ADD CONSTRAINT ns_journal_pk PRIMARY KEY (id);


--
-- TOC entry 3218 (class 2606 OID 24808)
-- Name: ns_messages ns_messages_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_messages
    ADD CONSTRAINT ns_messages_pk PRIMARY KEY (id);


--
-- TOC entry 3220 (class 2606 OID 24820)
-- Name: ns_resources ns_resources_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_resources
    ADD CONSTRAINT ns_resources_pk PRIMARY KEY (id);


--
-- TOC entry 3224 (class 2606 OID 24844)
-- Name: ns_settings ns_settings_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_settings
    ADD CONSTRAINT ns_settings_pk PRIMARY KEY (id);


--
-- TOC entry 3226 (class 2606 OID 24867)
-- Name: ns_users ns_users_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_users
    ADD CONSTRAINT ns_users_pk PRIMARY KEY (id);


--
-- TOC entry 3232 (class 2606 OID 24903)
-- Name: ns_journal ns_journal_ns_messages_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_journal
    ADD CONSTRAINT ns_journal_ns_messages_id_fk FOREIGN KEY (message_id) REFERENCES public.ns_messages(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3231 (class 2606 OID 24893)
-- Name: ns_journal ns_journal_ns_resources_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_journal
    ADD CONSTRAINT ns_journal_ns_resources_id_fk FOREIGN KEY (resource_id) REFERENCES public.ns_resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3230 (class 2606 OID 24883)
-- Name: ns_journal ns_journal_ns_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_journal
    ADD CONSTRAINT ns_journal_ns_users_id_fk FOREIGN KEY (user_id) REFERENCES public.ns_users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3228 (class 2606 OID 24888)
-- Name: ns_messages ns_messages_ns_resources_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_messages
    ADD CONSTRAINT ns_messages_ns_resources_id_fk FOREIGN KEY (resource_id) REFERENCES public.ns_resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3227 (class 2606 OID 24868)
-- Name: ns_messages ns_messages_ns_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_messages
    ADD CONSTRAINT ns_messages_ns_users_id_fk FOREIGN KEY (user_id) REFERENCES public.ns_users(id);


--
-- TOC entry 3229 (class 2606 OID 24873)
-- Name: ns_resources ns_resources_ns_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_resources
    ADD CONSTRAINT ns_resources_ns_users_id_fk FOREIGN KEY (user_id) REFERENCES public.ns_users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3234 (class 2606 OID 24898)
-- Name: ns_settings ns_settings_ns_resources_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_settings
    ADD CONSTRAINT ns_settings_ns_resources_id_fk FOREIGN KEY (resource_id) REFERENCES public.ns_resources(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- TOC entry 3233 (class 2606 OID 24878)
-- Name: ns_settings ns_settings_ns_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ns_settings
    ADD CONSTRAINT ns_settings_ns_users_id_fk FOREIGN KEY (user_id) REFERENCES public.ns_users(id) ON UPDATE CASCADE ON DELETE CASCADE;


-- Completed on 2022-07-26 16:23:57

--
-- PostgreSQL database dump complete
--

