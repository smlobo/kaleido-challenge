-- Distributed tracing: the W3C traceparent of the request that queued a
-- transaction, so the submitter and receipt loops, which run later and on any
-- instance, put their work (the send, the receipt, the final-status hook) in
-- that request's trace. NULL when the request was not traced.
ALTER TABLE transactions ADD COLUMN traceparent text;
