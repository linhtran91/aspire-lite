CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);
CREATE UNIQUE INDEX idx_username ON customers(username);

CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    amount FLOAT NOT NULL,
    term INT NOT NULL,
	customer_id INT NOT NULL,
	status INT,
	created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    scheduled_date timestamp without time zone NOT NULL,
    CONSTRAINT FK_LoanCustomer FOREIGN KEY (customer_id)
    REFERENCES customers(id)
);
CREATE INDEX idx_customer_id ON loans(customer_id);

CREATE TABLE IF NOT EXISTS repayments (
    id VARCHAR(16) PRIMARY KEY,
    loan_id INT NOT NULL,
    scheduled_amount FLOAT,
    actual_amount FLOAT,
    status INT,
    created_at timestamp without time zone NOT NULL,
    paid_at timestamp without time zone,
    scheduled_pay_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    CONSTRAINT FK_RepaymentLoan FOREIGN KEY (loan_id)
    REFERENCES loans(id)
);
CREATE INDEX idx_loan_id ON repayments(loan_id);
