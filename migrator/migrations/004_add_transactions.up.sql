ALTER TABLE transactions
ADD CONSTRAINT sender_balance_check 
CHECK (amount > 0);
