-- MySQL 8.0 初始化校验
-- 根用户密码与远程访问主机由 docker-compose 中的
-- MYSQL_ROOT_PASSWORD / MYSQL_ROOT_HOST 控制，不在 SQL 中硬编码。

USE mysql;

SELECT user, host, plugin FROM mysql.user WHERE user = 'root';
