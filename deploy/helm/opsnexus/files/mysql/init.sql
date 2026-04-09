-- Keep root credentials aligned with MYSQL_ROOT_PASSWORD from the container env.
-- The image already initializes root with caching_sha2_password, and
-- MYSQL_ROOT_HOST=% enables cross-pod access for the Helm deployment.

USE mysql;

FLUSH PRIVILEGES;

SELECT user, host, plugin FROM mysql.user WHERE user = 'root';
