create table localities(
    `id` varchar(50) not null primary key,
    local_name text not null,
    province_name text not null,
    country_name text not null
);

create table sellers(
    `id` int not null primary key auto_increment,
    cid int not null,
    company_name text not null,
    `address` text not null,
    telephone varchar(15) not null,
    locality_id varchar(50) not null,
    foreign key (locality_id ) references localities(id)
);
create table product_types(
    `id` int not null primary key auto_increment,
    `name` varchar(50) not null
);

create table products(
    `id` int not null primary key auto_increment,
    `description` text not null,
    expiration_rate float not null,
    freezing_rate float not null,
    height float not null,
    lenght float not null,
    netweight float not null,
    product_code text not null,
    recommended_freezing_temperature float not null,
    width float not null,
    id_product_type int not null,
    id_seller int not null,
    foreign key (id_seller) references sellers(id),
    foreign key (id_product_type) references product_types(id)
);
create table warehouses(
    `id` int not null primary key auto_increment,
    `address` text null,
    telephone text null,
    warehouse_code text null,
    minimum_capacity int null,
    minimum_temperature int null
);
create table employees(
    `id` int not null primary key auto_increment,
    card_number_id text not null,
    first_name text not null,
    last_name text not null,
    warehouse_id int not null,
    foreign key (warehouse_id) references warehouses(id)
);

create table sections(
    `id` int not null primary key auto_increment,
    section_number int not null unique,
    current_temperature int not null,
    minimum_temperature int not null,
    current_capacity int not null,
    minimum_capacity int not null,
    maximum_capacity int not null,
    warehouse_id int not null,
    id_product_type int not null,
    foreign key (warehouse_id) references warehouses(id),
    foreign key (id_product_type) references product_types(id)
);

create table buyers(
    `id` int not null primary key auto_increment,
    card_number_id text not null,
    first_name text not null,
    last_name text not null
);

/* tablas sprint 2 */


create table carries(
    `id` int not null primary key auto_increment,
    cid varchar(25) not null,
    company_name varchar(25) not null,
    `address` varchar(25) not null,
    telephone varchar(25) not null,
    locality_id varchar(50) not null,
    foreign key (locality_id) references localities(id)
);

create table products_batches(
    `id` int not null primary key auto_increment,
    batch_number int not null unique,
    current_quantity int not null,
    current_temperature int not null,
    due_date date not null,
    initial_quantity int not null,
    manufacturing_date date not null,
    manufacturing_hour time not null,
    minumum_temperature int not null,
    product_id int not null,
    section_id int not null,
    foreign key (product_id) references products(id),
    foreign key (section_id) references sections(id)
);

create table product_records(
	`id` int not null primary key auto_increment,
	last_update_date date,
	purchase_price float not null,
	sale_price float not null,
	product_id int not null,
	foreign key (product_id) references products(id)
);

create table inbound_orders(
    `id` int not null primary key auto_increment,
    order_date date not null,
    order_number varchar(25) not null unique,
    employee_id int not null,
    product_batch_id int not null,
    warehouse_id int not null,
    foreign key (employee_id) references employees(id),
    foreign key (product_batch_id) references products_batches(id),
    foreign key (warehouse_id) references warehouses(id)
);

CREATE TABLE purchase_orders (
	`id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    order_number VARCHAR(30) NOT NULL,
    order_date DATE NOT NULL,
    tracking_code VARCHAR(20) NOT NULL,
    buyer_id INT NOT NULL,
    product_record_id INT NOT NULL, 
    order_status_id INT NOT NULL, 
    FOREIGN KEY (`buyer_id`) references buyers(`id`),
    FOREIGN KEY (`product_record_id`) references product_records(`id`),
    FOREIGN KEY (`order_status_id`) references inbound_orders (`id`)
);
