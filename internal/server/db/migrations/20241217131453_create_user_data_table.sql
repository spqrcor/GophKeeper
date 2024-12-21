-- +goose Up
-- +goose StatementBegin
CREATE extension IF NOT EXISTS pgcrypto;

CREATE TABLE public.user_data (
                                  id uuid DEFAULT gen_random_uuid() NOT NULL,
                                  user_id uuid NOT NULL,
                                  data bytea NOT NULL,
                                  file bytea,
                                  created_at timestamptz DEFAULT now() NOT NULL,
                                  CONSTRAINT user_data_pkey PRIMARY KEY (id),
                                  CONSTRAINT user_data_user_id_fkey
                                      FOREIGN KEY (user_id)
                                          REFERENCES public.users(id)
                                          ON DELETE CASCADE
);
CREATE INDEX user_data_user_id_idx ON public.user_data (user_id);
CREATE INDEX user_data_created_at_idx ON public.user_data (created_at DESC);

COMMENT ON COLUMN public.user_data.id IS 'UUID';
COMMENT ON COLUMN public.user_data.user_id IS 'ID пользователя';
COMMENT ON COLUMN public.user_data.data IS 'Набор шифрованных данных';
COMMENT ON COLUMN public.user_data.file IS 'Зашифрованный файл';
COMMENT ON COLUMN public.user_data.created_at IS 'Дата создания';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.user_data;
-- +goose StatementEnd
