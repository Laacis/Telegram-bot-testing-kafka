CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    restaurant_id VARCHAR(3) NOT NULL,
    name VARCHAR(100) NOT NULL,
    contact_number VARCHAR(15) NOT NULL,
    tax_number VARCHAR(20) NOT NULL,
    address VARCHAR(200) NOT NULL
    );

CREATE TABLE IF NOT EXISTS destinations (
    id SERIAL PRIMARY KEY,
    restaurant_code VARCHAR(6) NOT NULL,
    restaurant_name VARCHAR(100) NOT NULL,
    address VARCHAR(200) NOT NULL,
    area_code VARCHAR(10) NOT NULL,
    customer_id INTEGER REFERENCES customers(id)
    );

INSERT INTO customers (restaurant_id, name, contact_number, tax_number, address) VALUES
    ('MCD', 'McDonalds LTD', '891237231', 'TX12345', '23 Main St, Springfield, 5500'),
    ('KFC', 'KFC foods', '90378923112', 'TX10015', '15 West St, Berry hill, 5209'),
    ('JLB', 'JolliBee', '885120134', 'TX24108', '92 San Jose St, Berry hill, 5200'),
    ('BNS', 'BONES Food services Inc', '9012316525', 'TX67890', '75 Elm St, Big city, 5502');

INSERT INTO destinations (restaurant_code, restaurant_name, address, area_code, customer_id) VALUES
    ('MCD001', 'McD Victory park', '20A Victory park, Springfield, 5501', '55', 1),
    ('MCD002', 'McD Walmart C', '11-200 Mist ave., Springfield, 5504', '55', 1),
    ('MCD003', 'McD SpringField W', '20A West st., Springfield, 5502', '55', 1),
    ('MCD031', 'McD Crown town', '10 Main st, Crown town, 5438', '54', 1),
    ('MCD032', 'McD Sand hill N', '55 Gate st., Long town, 5433', '54', 1),
    ('KFC101', 'KFC SpringField W', '20B West st., Springfield, 5502', '55', 2),
    ('KFC102', 'KFC Sand hills S', '502 South st.., Sand hills, 5404', '54', 2),
    ('JLB011', 'JolliBee SP. Trade Center', '5A Merchant gate, Springfield, 5505', '55', 3),
    ('JLB012', 'JolliBee TOYOTA C', '74 East st., Springfield, 5503', '55', 3),
    ('JLB023', 'JolliBee Paradise Golf club', '1 Green hill ave, Long town, 5428', '54', 3),
    ('BNS001', 'BONES APACHE PARK', '1-10C Apache park, Springfield, 5502', '55', 4),
    ('BNS002', 'BONES PARK ', '22 Victory park, Springfield, 5501', '55', 4);
