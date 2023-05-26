CREATE TABLE couriers
(
    id serial not null unique,
    courier_type varchar(4) not null,
    regions integer[] not null,
    working_hours varchar(11)[] not null
);

CREATE TABLE orders
(
    order_id serial not null unique,
    weight float not null,
    regions varchar not null,
    delivery_hours varchar[] not null,
    order_cost int not null,
    completed_time	varchar --string($date-time)
);