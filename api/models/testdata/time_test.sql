DELETE
FROM public.users
WHERE id = 1;

INSERT INTO public.users (line_user_id, morning_medication_time, afternoon_medication_time, evening_medication_time)
VALUES ('test_id', '08:00:00', '12:00:00', '18:00:00');