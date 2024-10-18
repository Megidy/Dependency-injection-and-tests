CREATE DATABASE IF NOT EXISTS api ;
USE api;

CREATE TABLE IF NOT EXISTS users(
    id INT AUTO_INCREMENT PRIMARY KEY,
    firstname VARCHAR(255),
    lastname VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255)
);


CREATE TABLE IF NOT EXISTS products(
    id INT PRIMARY KEY,
    name VARCHAR(255),
    quantity INT ,
    price INT 
);

CREATE TABLE IF NOT EXISTS orders(
    id VARCHAR(255),
    user_id INT ,
    product_id INT ,
    quantity INT ,
    status VARCHAR(255)
);
 