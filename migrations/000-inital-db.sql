CREATE TABLE IF NOT EXISTS customers (
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
	created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    scheduled_date timestamp without time zone NOT NULL,
    CONSTRAINT FK_LoanCustomer FOREIGN KEY (customer_id)
    REFERENCES customers(id)
);

CREATE TABLE repayments (
    id VARCHAR(16) PRIMARY KEY,
    loan_id INT,
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