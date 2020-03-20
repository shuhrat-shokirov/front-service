Create table if not exists users
(
    id       Bigserial primary key,
    login    Text UNIQUE not null,
    password text        not null
);

INSERT INTO users (login, password)
VALUES ('admin', 'pass'),
       ('qwe', 'pass');

INSERT INTO users (login, password)
VALUES (?, ?);

Create table if not exists rooms
(
    id       Bigserial primary key,
    status   bool,
    timeInHour int,
    timeInMinutes int,
    timeOutHour int,
    timeOutMinutes int,
    fileName text
);

INSERT INTO rooms (status)
VALUES (false);

SELECT id, timeinhour, timeinminutes, timeouthour, timeoutminutes, filename FROM mitings;

SELECT id, status, timeinhour, timeinminutes, timeouthour, timeoutminutes, filename FROM rooms where status = true;