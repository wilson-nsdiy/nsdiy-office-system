-- initial_schema.sql
-- Created by Ent auto-migration on 2026-06-13
-- This file documents the initial database schema.

CREATE TABLE IF NOT EXISTS `api_tokens` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `user_id` INTEGER NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `token_hash` TEXT NOT NULL UNIQUE,
    `token_prefix` VARCHAR(20) NOT NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'active',
    `expires_at` DATETIME NULL,
    `last_used_at` DATETIME NULL,
    `usage_count` INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `articles` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `article_no` VARCHAR(50) NOT NULL UNIQUE,
    `title` VARCHAR(200) NOT NULL,
    `content` TEXT NULL,
    `summary` VARCHAR(1000) NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    `author_id` INTEGER NOT NULL,
    `cover_description` VARCHAR(500) NULL,
    `cover_url` VARCHAR(500) NULL,
    `first_published_at` DATETIME NULL,
    FOREIGN KEY (`author_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `article_versions` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `article_id` INTEGER NOT NULL,
    `version_no` INTEGER NOT NULL,
    `title` VARCHAR(200) NOT NULL,
    `content` TEXT NULL,
    `cover_description` VARCHAR(500) NULL,
    `summary` VARCHAR(1000) NULL,
    `status` VARCHAR(20) NOT NULL,
    `editor_id` INTEGER NULL,
    `edit_reason` VARCHAR(500) NULL,
    FOREIGN KEY (`article_id`) REFERENCES `articles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`editor_id`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    UNIQUE(`article_id`, `version_no`)
);

CREATE TABLE IF NOT EXISTS `media_accounts` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `name` VARCHAR(100) NOT NULL,
    `platform` VARCHAR(50) NOT NULL,
    `account_id` VARCHAR(100) NOT NULL,
    `avatar` VARCHAR(500) NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'active',
    `access_token` VARCHAR(500) NULL,
    `refresh_token` VARCHAR(500) NULL,
    `token_expires_at` DATETIME NULL
);

CREATE TABLE IF NOT EXISTS `media_account_fans_snapshots` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `account_id` INTEGER NOT NULL,
    `fans_count` INTEGER NOT NULL,
    `snapshot_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`account_id`) REFERENCES `media_accounts`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `media_contents` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `title` VARCHAR(200) NOT NULL,
    `content` TEXT NULL,
    `cover_image` VARCHAR(500) NULL,
    `platform` VARCHAR(50) NOT NULL,
    `account_id` INTEGER NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'draft',
    `views` INTEGER NOT NULL DEFAULT 0,
    `likes` INTEGER NOT NULL DEFAULT 0,
    `comments` INTEGER NOT NULL DEFAULT 0,
    `shares` INTEGER NOT NULL DEFAULT 0,
    `publish_time` DATETIME NULL,
    FOREIGN KEY (`account_id`) REFERENCES `media_accounts`(`id`) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS `media_content_versions` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `content_id` INTEGER NOT NULL,
    `version_no` INTEGER NOT NULL,
    `title` VARCHAR(200) NOT NULL,
    `content` TEXT NULL,
    `cover_image` VARCHAR(500) NULL,
    `status` VARCHAR(20) NOT NULL,
    `editor_id` INTEGER NULL,
    `edit_reason` VARCHAR(500) NULL,
    FOREIGN KEY (`content_id`) REFERENCES `media_contents`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`editor_id`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    UNIQUE(`content_id`, `version_no`)
);

CREATE TABLE IF NOT EXISTS `news` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `group_id` INTEGER NOT NULL,
    `title` VARCHAR(200) NOT NULL,
    `content` TEXT NULL,
    `creator_id` INTEGER NOT NULL,
    FOREIGN KEY (`group_id`) REFERENCES `news_groups`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`creator_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `news_groups` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `name` VARCHAR(100) NOT NULL UNIQUE,
    `description` VARCHAR(500) NULL,
    `sort_order` INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `user_id` INTEGER NULL,
    `username` VARCHAR(100) NULL,
    `operation_type` VARCHAR(50) NOT NULL,
    `module` VARCHAR(50) NOT NULL,
    `action` VARCHAR(50) NOT NULL,
    `resource_type` VARCHAR(50) NULL,
    `resource_id` INTEGER NULL,
    `resource_name` VARCHAR(200) NULL,
    `project_id` INTEGER NULL,
    `detail` VARCHAR(2000) NULL,
    `ip_address` VARCHAR(50) NULL,
    `user_agent` VARCHAR(500) NULL,
    `status` VARCHAR(20) NOT NULL,
    `error_message` VARCHAR(2000) NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS `permissions` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `pid` INTEGER NULL,
    `name` VARCHAR(100) NOT NULL UNIQUE,
    `resource_type` VARCHAR(50) NOT NULL,
    `resource_path` VARCHAR(200) NOT NULL,
    `http_method` VARCHAR(10) NULL,
    `description` VARCHAR(500) NULL,
    `is_active` INTEGER NOT NULL DEFAULT 1,
    `is_builtin` INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS `projects` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `name` VARCHAR(200) NOT NULL,
    `project_no` VARCHAR(50) NOT NULL UNIQUE,
    `description` VARCHAR(2000) NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'TODO',
    `priority` VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    `expected_start_date` DATETIME NULL,
    `expected_end_date` DATETIME NULL,
    `start_date` DATETIME NULL,
    `end_date` DATETIME NULL,
    `owner_id` INTEGER NOT NULL,
    FOREIGN KEY (`owner_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `project_members` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `project_id` INTEGER NOT NULL,
    `user_id` INTEGER NOT NULL,
    `role` VARCHAR(20) NOT NULL DEFAULT 'MEMBER',
    `joined_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    UNIQUE(`project_id`, `user_id`)
);

CREATE TABLE IF NOT EXISTS `role_perms` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `role_id` INTEGER NOT NULL,
    `permission_id` INTEGER NOT NULL,
    FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`permission_id`) REFERENCES `permissions`(`id`) ON DELETE CASCADE,
    UNIQUE(`role_id`, `permission_id`)
);

CREATE TABLE IF NOT EXISTS `roles` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `name` VARCHAR(100) NOT NULL UNIQUE,
    `code` VARCHAR(50) NOT NULL UNIQUE,
    `description` VARCHAR(500) NULL,
    `is_active` INTEGER NOT NULL DEFAULT 1,
    `is_default` INTEGER NOT NULL DEFAULT 0,
    `is_builtin` INTEGER NOT NULL DEFAULT 0,
    `role_type` VARCHAR(50) NULL
);

CREATE TABLE IF NOT EXISTS `tasks` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `project_id` INTEGER NOT NULL,
    `parent_id` INTEGER NULL,
    `title` VARCHAR(200) NOT NULL,
    `description` VARCHAR(2000) NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'TODO',
    `priority` VARCHAR(20) NOT NULL DEFAULT 'MEDIUM',
    `assignee_id` INTEGER NULL,
    `creator_id` INTEGER NOT NULL,
    `planned_start_date` DATETIME NULL,
    `planned_end_date` DATETIME NULL,
    `actual_start_time` DATETIME NULL,
    `actual_end_time` DATETIME NULL,
    `estimated_hours` REAL NULL,
    FOREIGN KEY (`project_id`) REFERENCES `projects`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`assignee_id`) REFERENCES `users`(`id`) ON DELETE SET NULL,
    FOREIGN KEY (`creator_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `upload_files` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `filename` VARCHAR(200) NOT NULL UNIQUE,
    `original_filename` VARCHAR(200) NOT NULL,
    `file_path` VARCHAR(500) NOT NULL,
    `file_size` INTEGER NOT NULL,
    `mime_type` VARCHAR(100) NOT NULL,
    `file_type` VARCHAR(50) NOT NULL,
    `extension` VARCHAR(20) NOT NULL,
    `uploader_id` INTEGER NOT NULL,
    `purpose` VARCHAR(100) NULL,
    `md5` VARCHAR(32) NULL,
    `reference_count` INTEGER NOT NULL DEFAULT 1,
    FOREIGN KEY (`uploader_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `users` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `username` VARCHAR(100) NOT NULL UNIQUE,
    `email` VARCHAR(200) NOT NULL UNIQUE,
    `nickname` VARCHAR(100) NULL,
    `salt` TEXT NOT NULL,
    `hashed_password` TEXT NOT NULL,
    `role_id` INTEGER NULL,
    `user_type` VARCHAR(20) NOT NULL DEFAULT 'HUMAN',
    `is_active` INTEGER NOT NULL DEFAULT 1,
    `token_version` INTEGER NOT NULL DEFAULT 1,
    `verification_code` VARCHAR(10) NULL,
    `verification_code_expires_at` DATETIME NULL,
    FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS `verification_codes` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `code` VARCHAR(10) NOT NULL,
    `target` VARCHAR(200) NOT NULL,
    `channel` VARCHAR(50) NOT NULL,
    `scene` VARCHAR(50) NOT NULL,
    `is_used` INTEGER NOT NULL DEFAULT 0,
    `attempts` INTEGER NOT NULL DEFAULT 0,
    `max_attempts` INTEGER NOT NULL DEFAULT 5,
    `expires_at` DATETIME NOT NULL,
    `used_at` DATETIME NULL,
    `ip_address` VARCHAR(50) NULL,
    `user_agent` VARCHAR(500) NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS `users_username_email` ON `users`(`username`, `email`);
CREATE UNIQUE INDEX IF NOT EXISTS `project_members_project_id_user_id` ON `project_members`(`project_id`, `user_id`);
CREATE UNIQUE INDEX IF NOT EXISTS `role_perms_role_id_permission_id` ON `role_perms`(`role_id`, `permission_id`);
CREATE UNIQUE INDEX IF NOT EXISTS `article_versions_article_id_version_no` ON `article_versions`(`article_id`, `version_no`);
CREATE UNIQUE INDEX IF NOT EXISTS `media_content_versions_content_id_version_no` ON `media_content_versions`(`content_id`, `version_no`);
