-- optigit schema v1

-- Track a single Git Repository, inside of a given org
CREATE TABLE repos (
	id       integer      not null primary key,

	org      varchar(100) not null,
	name     varchar(200) not null,
	included integer(1)   not null
);

-- Track a branch, on a single Git Repository
CREATE TABLE branches (
	id      integer      not null,
	repo_id integer      not null,
	name    varchar(100) not null,
	sha1    varchar(40)  not null,

	primary key (id, repo_id)
);

-- Track a Pull Request, against a single Git Repository
CREATE TABLE pulls (
	id          integer not null, -- GH-assigned PR number
	repo_id     integer not null, -- repo identifier

	created_at  integer NOT NULL,
	updated_at  integer NOT NULL,
	assignees   text NOT NULL, -- comma-separated usernames

	title       text NOT NULL,

	primary key (id, repo_id)
);

-- Track an Issue, against a single Git Repository
CREATE TABLE issues (
	id      integer not null, -- GH-assigned issue number
	repo_id integer not null, -- repo identifier

	created_at  integer NOT NULL,
	updated_at  integer NOT NULL,
	assignees   text NOT NULL, -- comma-separated usernames

	title       text NOT NULL,

	primary key (id, repo_id)
);
