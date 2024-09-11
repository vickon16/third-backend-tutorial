-- +goose Up

CREATE TABLE IF NOT EXISTS projects (
  `id` CHAR(36) PRIMARY KEY NOT NULL /* chat(36) for UUID */,
  `userId` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `repoURL` TEXT,
  `siteURL` TEXT,
  `status`  ENUM('backlog', 'developing', 'done') NOT NULL DEFAULT 'backlog',
  `dependencies` TEXT,
  `createdAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updatedAt` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  FOREIGN KEY (`userId`) REFERENCES users(`id`) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS projects;