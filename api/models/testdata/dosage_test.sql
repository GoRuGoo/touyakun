DELETE
FROM public.dosage
WHERE id = 1;
DELETE
FROM public.users
WHERE id = 1;
INSERT INTO public.users (id, line_user_id, created_at, updated_at)
VALUES (1, 'test_id', '2024-04-27 17:10:04.675314', '2024-04-27 17:10:04.675314');
INSERT INTO public.dosage (id, user_id, name, amount, duration, morning_flg, afternoon_flg, evening_flg)
VALUES (1, 1, 'トラネキサム', 2, 7, true, false, true);

DELETE
FROM public.users
WHERE id = 333;

INSERT INTO public.users (id, line_user_id, created_at, updated_at)
VALUES (333, 'test_id_for_register_medications', '2024-04-27 17:10:04.675314', '2024-04-27 17:10:04.675314');
