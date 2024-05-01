DELETE
FROM public.dosage
WHERE id = 1;
DELETE
FROM public.users
WHERE id = 1;
INSERT INTO public.users (line_user_id, morning_medication_time, afternoon_medication_time, evening_medication_time)
VALUES ('test_user', '08:00:00', '12:00:00', '18:00:00');
INSERT INTO public.dosage (user_id, name, amount, duration, morning_flg, afternoon_flg, evening_flg)
VALUES (1, 'トラネキサム', 2, 7, true, false, true);