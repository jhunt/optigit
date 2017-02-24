INSERT INTO repos (id, org, name, included) VALUES (1, "starkandwayne", "shield", 1);
INSERT INTO repos (id, org, name, included) VALUES (2, "starkandwayne", "shield-boshrelease", 1);
INSERT INTO repos (id, org, name, included) VALUES (3, "cloudfoundry-community", "vault-boshrelease", 1);
INSERT INTO repos (id, org, name, included) VALUES (4, "cloudfoundry-community", "ignored-boshrelease", 0);

-- INSERT INTO branches (id, repo_id, name, sha1) VALUES (1, 1, "master", "deadbeef");
-- INSERT INTO branches (id, repo_id, name, sha1) VALUES (2, 2, "master", "deadbeef");
-- INSERT INTO branches (id, repo_id, name, sha1) VALUES (3, 3, "master", "deadbeef");
-- INSERT INTO branches (id, repo_id, name, sha1) VALUES (4, 4, "master", "deadbeef");
-- INSERT INTO branches (id, repo_id, name, sha1) VALUES (5, 1, "devs", "deadbeef");

INSERT INTO pulls (id, repo_id, created_at, updated_at, assignees, title) VALUES (1, 2, 1487878540, -1, "jhunt", "Fix this BOSH release");
INSERT INTO pulls (id, repo_id, created_at, updated_at, assignees, title) VALUES (2, 2, 1487878540, -1, "jhunt", "Fix this BOSH release fo realz");

INSERT INTO issues (id, repo_id, created_at, updated_at, assignees, title) VALUES (1, 3, 1487811111, 1487855555, "quintessence,jhunt", "Vault is borked");
INSERT INTO issues (id, repo_id, created_at, updated_at, assignees, title) VALUES (2, 3, 1487811111, 1487855555, "", "Vault is super borked");
