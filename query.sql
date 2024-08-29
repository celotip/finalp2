-- Table: Users
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    birth_date DATE,
    address TEXT,
    contact_no VARCHAR(20),
    deposit INT DEFAULT 0,
    jwt_token TEXT
);

-- Table: Authors
CREATE TABLE Authors (
    author_id SERIAL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    nationality VARCHAR(255),
    birth_date DATE
);

-- Table: Categories
CREATE TABLE Categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Table: Books
CREATE TABLE Books (
    book_id SERIAL PRIMARY KEY,
    author_id INT REFERENCES Authors(author_id) ON DELETE SET NULL,
    category_id INT REFERENCES Categories(category_id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    ISBN VARCHAR(13) UNIQUE,
    stock INT DEFAULT 0,
    price INT NOT NULL,
    reading_days INT DEFAULT 0
);

-- Table: Rentals
CREATE TABLE Rentals (
    rental_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    rental_date DATE,
    rental_status VARCHAR(50) DEFAULT 'created',
    total_price INT DEFAULT 0
);

-- Table: Cart
CREATE TABLE Carts (
    cart_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    book_id INT REFERENCES Books(book_id) ON DELETE CASCADE
);

-- Table: Rental_Details
CREATE TABLE Rental_Details (
    rental_detail_id SERIAL PRIMARY KEY,
    rental_id INT REFERENCES Rentals(rental_id) ON DELETE CASCADE,
    book_id INT REFERENCES Books(book_id) ON DELETE CASCADE,
    returned BOOLEAN DEFAULT FALSE
);

-- Table: Payments
CREATE TABLE Payments (
    payment_id SERIAL PRIMARY KEY,
    rental_id INT REFERENCES Rentals(rental_id) ON DELETE CASCADE,
    payment_date DATE,
    payment_amount DECIMAL(10, 2) NOT NULL
);

-- Insert data into Authors
INSERT INTO Authors (first_name, last_name, nationality, birth_date)
VALUES
('George', 'Orwell', 'British', '1903-06-25'),
('J.K.', 'Rowling', 'British', '1965-07-31');

-- Insert data into Categories
INSERT INTO Categories (name)
VALUES
('Fiction'),
('Science Fiction'),
('Fantasy'),
('Non-fiction');

-- Insert data into Books
INSERT INTO Books (author_id, category_id, title, ISBN, stock, price, reading_days)
VALUES
(1, 2, '1984', '9780451524935', 10, 20000, 14),
(2, 3, 'Harry Potter and the Philosopher''s Stone', '9780747532743', 15, 30000, 21);

-- Insert data into Users
INSERT INTO Users (email, password_hash, first_name, last_name, birth_date, address, contact_no, deposit, jwt_token)
VALUES
('user1@example.com', 'hashed_password_1', 'John', 'Doe', '1990-01-01', '123 Main St', '081234567890', 50000, 'jwt_token_example_1'),
('user2@example.com', 'hashed_password_2', 'Jane', 'Doe', '1992-02-02', '456 Main St', '081234567891', 100000, 'jwt_token_example_2');

-- Insert data into Rentals
INSERT INTO Rentals (user_id, rental_date, rental_status, total_price)
VALUES
(1, '2024-08-01', 'completed', 20000),
(2, '2024-08-05', 'created', 30000);

-- Insert data into Cart
INSERT INTO Carts (user_id, book_id)
VALUES
(1, 1),
(2, 2);

-- Insert data into Rental_Details
INSERT INTO Rental_Details (rental_id, book_id, returned)
VALUES
(1, 1, TRUE),
(2, 2, FALSE);

-- Insert data into Payments
INSERT INTO Payments (rental_id, payment_date, payment_amount)
VALUES
(1, '2024-08-02', 20000.00),
(2, '2024-08-06', 30000.00);
