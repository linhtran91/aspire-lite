CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255),
    password VARCHAR(255)
);

CREATE TABLE loans (
    id SERIAL PRIMARY KEY,
    amount     FLOAT,
	term       INT,
	customer_id INT,
	status     INT,
	created_at  TIMESTAMP NOT NULL,
	updated_at  TIMESTAMP,
    CONSTRAINT FK_LoanCustomer FOREIGN KEY (customer_id)
    REFERENCES customers(id)
);

CREATE TABLE repayments (
    id VARCHAR(255),
    loan_id INT,
    scheduled_amount FLOAT,
    actual_amount FLOAT,
    status INT,
    created_at TIMESTAMP,
    paid_at TIMESTAMP,
    scheduled_pay_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT FK_RepaymentLoan FOREIGN KEY (loan_id)
    REFERENCES loans(id)
);