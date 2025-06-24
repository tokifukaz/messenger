CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL,
    content TEXT NOT NULL CHECK (char_length(content) <= 160),
    sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP
);

-- Insert mock data (unsent messages)
INSERT INTO messages (phone_number, content) VALUES
('1234567890', 'Hello, this is a test message 1'),
('1234567891', 'Hello, this is a test message 2'),
('1234567892', 'Hello, this is a test message 3'),
('1234567893', 'Hello, this is a test message 4'),
('1234567894', 'Hello, this is a test message 5')
ON CONFLICT DO NOTHING;
