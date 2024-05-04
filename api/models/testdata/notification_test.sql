DELETE
FROM users
WHERE id IN (1, 2);
INSERT INTO users (id, line_user_id, morning_medication_time, afternoon_medication_time, evening_medication_time)
VALUES (1, 'test_user_1', '08:00', '12:00', '18:00'),
       (2, 'test_user_2', '09:00', '13:00', '19:00');

INSERT INTO dosage (user_id, name, amount, duration, morning_flg, afternoon_flg, evening_flg)
VALUES (1, 'Test Drug 1', 2, 7, true, false, true),
       (2, 'Test Drug 2', 3, 14, false, true, false);