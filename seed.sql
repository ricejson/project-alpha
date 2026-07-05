BEGIN;

TRUNCATE TABLE ticket_tags, tickets, tags RESTART IDENTITY CASCADE;

WITH inserted_tags AS (
  INSERT INTO tags (name, color, created_at, updated_at, deleted_at)
  VALUES
    ('ios', '#007AFF', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('android', '#3DDC84', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('web', '#2563EB', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('backend', '#0F766E', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('project-alpha', '#7C3AED', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('bugfix', '#DC2626', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('rendering-fix', '#F97316', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('performance', '#0891B2', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('ui-polish', '#DB2777', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL),
    ('release-blocker', '#B91C1C', '2026-05-23 09:00:00+08', '2026-05-23 09:00:00+08', NULL)
  RETURNING id, name
),
inserted_tickets AS (
  INSERT INTO tickets (title, description, completed, completed_at, created_at, updated_at, deleted_at)
  VALUES
    ('Fix iOS status badge overlap', 'Status badge overlaps the title on narrow iPhone screens. Rework the ticket row layout and preserve wrapping.', false, NULL, '2026-05-23 09:05:00+08', '2026-05-23 09:05:00+08', NULL),
    ('Add project onboarding checklist', 'Create a lightweight checklist for new projects so common setup steps are not forgotten.', false, NULL, '2026-05-23 09:10:00+08', '2026-05-23 09:10:00+08', NULL),
    ('Investigate slow ticket list pagination', 'The second page takes noticeably longer to render when there are many tags. Profile query and frontend rendering.', true, '2026-05-22 18:30:00+08', '2026-05-23 09:15:00+08', '2026-05-23 09:15:00+08', NULL),
    ('Fix rendering glitch in tag picker', 'Selected tags sometimes flicker when the dialog reopens. Keep selection state stable across renders.', false, NULL, '2026-05-23 09:20:00+08', '2026-05-23 09:20:00+08', NULL),
    ('Ship backend health check improvements', 'Return clearer error output when the database cannot be reached during startup.', true, '2026-05-21 13:00:00+08', '2026-05-23 09:25:00+08', '2026-05-23 09:25:00+08', NULL),
    ('Reduce search debounce latency', 'Make title search feel more responsive without spamming API requests.', false, NULL, '2026-05-23 09:30:00+08', '2026-05-23 09:30:00+08', NULL),
    ('Polish empty state for filtered views', 'When filters return no tickets, show a clearer empty state with reset guidance.', false, NULL, '2026-05-23 09:35:00+08', '2026-05-23 09:35:00+08', NULL),
    ('Resolve duplicate tag color edge case', 'Ensure repeated edits do not leave stale color values in tag records.', true, '2026-05-20 10:45:00+08', '2026-05-23 09:40:00+08', '2026-05-23 09:40:00+08', NULL),
    ('Prepare release blocker triage', 'Review critical issues before the next release candidate and prioritize the blocked items.', false, NULL, '2026-05-23 09:45:00+08', '2026-05-23 09:45:00+08', NULL),
    ('Improve mobile ticket form spacing', 'Tighten the mobile form spacing so controls do not crowd each other on smaller screens.', false, NULL, '2026-05-23 09:50:00+08', '2026-05-23 09:50:00+08', NULL)
  RETURNING id, title
)
INSERT INTO ticket_tags (ticket_id, tag_id, created_at)
SELECT t.id, tag.id, '2026-05-23 10:00:00+08'
FROM inserted_tickets t
JOIN inserted_tags tag ON
  (t.title = 'Fix iOS status badge overlap' AND tag.name IN ('ios', 'ui-polish', 'rendering-fix')) OR
  (t.title = 'Add project onboarding checklist' AND tag.name IN ('project-alpha', 'web')) OR
  (t.title = 'Investigate slow ticket list pagination' AND tag.name IN ('backend', 'performance')) OR
  (t.title = 'Fix rendering glitch in tag picker' AND tag.name IN ('rendering-fix', 'ui-polish', 'bugfix')) OR
  (t.title = 'Ship backend health check improvements' AND tag.name IN ('backend', 'release-blocker')) OR
  (t.title = 'Reduce search debounce latency' AND tag.name IN ('web', 'performance')) OR
  (t.title = 'Polish empty state for filtered views' AND tag.name IN ('ui-polish', 'project-alpha')) OR
  (t.title = 'Resolve duplicate tag color edge case' AND tag.name IN ('bugfix', 'backend')) OR
  (t.title = 'Prepare release blocker triage' AND tag.name IN ('release-blocker', 'project-alpha')) OR
  (t.title = 'Improve mobile ticket form spacing' AND tag.name IN ('ios', 'android', 'ui-polish'));

COMMIT;
