-- +goose Up
INSERT INTO todos (title, completed, created_at, completed_at)
VALUES (
    'Bake a cake',
    TRUE,
    STR_TO_DATE('2025-02-18 15:44:04', '%Y-%m-%d %H:%i:%s'),
    STR_TO_DATE('2025-02-18 16:44:04', '%Y-%m-%d %H:%i:%s')
  ),
  (
    'Feed the cat',
    FALSE,
    STR_TO_DATE('2025-02-18 15:55:04', '%Y-%m-%d %H:%i:%s'),
    NULL
  ),
  (
    'Take out the trash',
    FALSE,
    STR_TO_DATE('2025-02-18 15:57:10', '%Y-%m-%d %H:%i:%s'),
    NULL
  );
-- +goose Down
DELETE FROM todos;