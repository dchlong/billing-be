CREATE TABLE call_histories
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_name  varchar(255),
    duration   INT,
    created_at datetime,
    updated_at datetime,
    deleted_at bigint
);