DELETE
FROM public.dosage
WHERE id = 1;
DELETE
FROM public.users
WHERE id = 1;
DELETE
FROM public.time
WHERE id = 1;
INSERT INTO public.users (id, username, auth_key, created_at, updated_at)
VALUES (1, 'GoRuGoo', 'test_auth', '2024-04-27 17:10:04.675314', '2024-04-27 17:10:04.675314');
INSERT INTO public.time (id, user_id, time, morning_flg, afternoon_flg, evening_flg, created_at, updated_at)
VALUES (1, 1, '2024-04-28 02:10:16.000000', true, false, false, '2024-04-27 17:10:28.524419',
        '2024-04-27 17:10:28.524419');
INSERT INTO public.dosage (id, user_id, time_id, name, amount, duration, created_at, updated_at)
VALUES (1, 1, 1, 'トラネキサム', 1, 3, '2024-04-27 17:10:49.023657', '2024-04-27 17:10:49.023657');
