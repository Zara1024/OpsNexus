DROP TABLE IF EXISTS sys_departments;

ALTER TABLE sys_menus
    DROP COLUMN IF EXISTS permission,
    DROP COLUMN IF EXISTS menu_type;

ALTER TABLE sys_users
    DROP COLUMN IF EXISTS department_id;
