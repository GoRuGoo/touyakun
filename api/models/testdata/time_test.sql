DELETE
FROM public.users
WHERE id = 1;
INSERT INTO public.users (id, line_user_id, created_at, updated_at)
VALUES (1, 'test_id', '2024-04-27 17:10:04.675314', '2024-04-27 17:10:04.675314');