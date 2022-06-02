--
-- PostgreSQL database dump
--

-- Dumped from database version 14.3 (Ubuntu 14.3-1.pgdg21.10+1)
-- Dumped by pg_dump version 14.3 (Ubuntu 14.3-1.pgdg21.10+1)

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
-- Name: livros; Type: TABLE; Schema: public; Owner: xx_stefan_xx
--

CREATE TABLE public.livros (
    id_livro integer NOT NULL,
    nome_livro character varying(255) DEFAULT ''::character varying NOT NULL,
    num_emprestimos integer DEFAULT 0 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.livros OWNER TO xx_stefan_xx;

--
-- Name: livros_id_livro_seq; Type: SEQUENCE; Schema: public; Owner: xx_stefan_xx
--

CREATE SEQUENCE public.livros_id_livro_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.livros_id_livro_seq OWNER TO xx_stefan_xx;

--
-- Name: livros_id_livro_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: xx_stefan_xx
--

ALTER SEQUENCE public.livros_id_livro_seq OWNED BY public.livros.id_livro;


--
-- Name: livros_restricoes; Type: TABLE; Schema: public; Owner: xx_stefan_xx
--

CREATE TABLE public.livros_restricoes (
    id integer NOT NULL,
    data_inicio date NOT NULL,
    data_final date NOT NULL,
    id_livro integer NOT NULL,
    id_reserva integer NOT NULL,
    id_restricao integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.livros_restricoes OWNER TO xx_stefan_xx;

--
-- Name: livros_restricoes_id_seq; Type: SEQUENCE; Schema: public; Owner: xx_stefan_xx
--

CREATE SEQUENCE public.livros_restricoes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.livros_restricoes_id_seq OWNER TO xx_stefan_xx;

--
-- Name: livros_restricoes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: xx_stefan_xx
--

ALTER SEQUENCE public.livros_restricoes_id_seq OWNED BY public.livros_restricoes.id;


--
-- Name: reservas; Type: TABLE; Schema: public; Owner: xx_stefan_xx
--

CREATE TABLE public.reservas (
    id integer NOT NULL,
    nome character varying(255) DEFAULT ''::character varying NOT NULL,
    sobrenome character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) NOT NULL,
    phone character varying(255) DEFAULT ''::character varying NOT NULL,
    data_inicio date NOT NULL,
    data_final date NOT NULL,
    livro_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.reservas OWNER TO xx_stefan_xx;

--
-- Name: reservas_id_seq; Type: SEQUENCE; Schema: public; Owner: xx_stefan_xx
--

CREATE SEQUENCE public.reservas_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reservas_id_seq OWNER TO xx_stefan_xx;

--
-- Name: reservas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: xx_stefan_xx
--

ALTER SEQUENCE public.reservas_id_seq OWNED BY public.reservas.id;


--
-- Name: restricoes; Type: TABLE; Schema: public; Owner: xx_stefan_xx
--

CREATE TABLE public.restricoes (
    id_restricao integer NOT NULL,
    nome_restricao character varying(255) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.restricoes OWNER TO xx_stefan_xx;

--
-- Name: restricoes_id_restricao_seq; Type: SEQUENCE; Schema: public; Owner: xx_stefan_xx
--

CREATE SEQUENCE public.restricoes_id_restricao_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.restricoes_id_restricao_seq OWNER TO xx_stefan_xx;

--
-- Name: restricoes_id_restricao_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: xx_stefan_xx
--

ALTER SEQUENCE public.restricoes_id_restricao_seq OWNED BY public.restricoes.id_restricao;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: xx_stefan_xx
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO xx_stefan_xx;

--
-- Name: users; Type: TABLE; Schema: public; Owner: xx_stefan_xx
--

CREATE TABLE public.users (
    id integer NOT NULL,
    nome character varying(255) DEFAULT ''::character varying NOT NULL,
    sobrenome character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    acces_level integer DEFAULT 1 NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.users OWNER TO xx_stefan_xx;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: xx_stefan_xx
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO xx_stefan_xx;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: xx_stefan_xx
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: livros id_livro; Type: DEFAULT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.livros ALTER COLUMN id_livro SET DEFAULT nextval('public.livros_id_livro_seq'::regclass);


--
-- Name: livros_restricoes id; Type: DEFAULT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.livros_restricoes ALTER COLUMN id SET DEFAULT nextval('public.livros_restricoes_id_seq'::regclass);


--
-- Name: reservas id; Type: DEFAULT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.reservas ALTER COLUMN id SET DEFAULT nextval('public.reservas_id_seq'::regclass);


--
-- Name: restricoes id_restricao; Type: DEFAULT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.restricoes ALTER COLUMN id_restricao SET DEFAULT nextval('public.restricoes_id_restricao_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: livros livros_pkey; Type: CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.livros
    ADD CONSTRAINT livros_pkey PRIMARY KEY (id_livro);


--
-- Name: livros_restricoes livros_restricoes_pkey; Type: CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.livros_restricoes
    ADD CONSTRAINT livros_restricoes_pkey PRIMARY KEY (id);


--
-- Name: reservas reservas_pkey; Type: CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.reservas
    ADD CONSTRAINT reservas_pkey PRIMARY KEY (id);


--
-- Name: restricoes restricoes_pkey; Type: CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.restricoes
    ADD CONSTRAINT restricoes_pkey PRIMARY KEY (id_restricao);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: livros_restricoes_data_inicio_data_final_idx; Type: INDEX; Schema: public; Owner: xx_stefan_xx
--

CREATE INDEX livros_restricoes_data_inicio_data_final_idx ON public.livros_restricoes USING btree (data_inicio, data_final);


--
-- Name: livros_restricoes_id_livro_idx; Type: INDEX; Schema: public; Owner: xx_stefan_xx
--

CREATE INDEX livros_restricoes_id_livro_idx ON public.livros_restricoes USING btree (id_livro);


--
-- Name: livros_restricoes_id_reserva_idx; Type: INDEX; Schema: public; Owner: xx_stefan_xx
--

CREATE INDEX livros_restricoes_id_reserva_idx ON public.livros_restricoes USING btree (id_reserva);


--
-- Name: reservas_email_idx; Type: INDEX; Schema: public; Owner: xx_stefan_xx
--

CREATE INDEX reservas_email_idx ON public.reservas USING btree (email);


--
-- Name: reservas_sobrenome_idx; Type: INDEX; Schema: public; Owner: xx_stefan_xx
--

CREATE INDEX reservas_sobrenome_idx ON public.reservas USING btree (sobrenome);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: xx_stefan_xx
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: users_email_idx; Type: INDEX; Schema: public; Owner: xx_stefan_xx
--

CREATE UNIQUE INDEX users_email_idx ON public.users USING btree (email);


--
-- Name: livros_restricoes livros_restricoes_livros_id_livro_fk; Type: FK CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.livros_restricoes
    ADD CONSTRAINT livros_restricoes_livros_id_livro_fk FOREIGN KEY (id_livro) REFERENCES public.livros(id_livro) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: livros_restricoes livros_restricoes_reservas_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.livros_restricoes
    ADD CONSTRAINT livros_restricoes_reservas_id_fk FOREIGN KEY (id_reserva) REFERENCES public.reservas(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: livros_restricoes livros_restricoes_restricoes_id_restricao_fk; Type: FK CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.livros_restricoes
    ADD CONSTRAINT livros_restricoes_restricoes_id_restricao_fk FOREIGN KEY (id_restricao) REFERENCES public.restricoes(id_restricao) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: reservas reservas_livros_id_livro_fk; Type: FK CONSTRAINT; Schema: public; Owner: xx_stefan_xx
--

ALTER TABLE ONLY public.reservas
    ADD CONSTRAINT reservas_livros_id_livro_fk FOREIGN KEY (livro_id) REFERENCES public.livros(id_livro) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

