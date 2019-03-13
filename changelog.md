# Change Log

## `head`

## `v0.3.0`
- Add disposition batching
- Add NotFound errors for mgmt API
- Fix go routine leak when listening for messages upon context close
- Add batch sends for Topics

## `v0.2.0`
- Refactor disposition handler so that errors can be handled in handlers
- Add dead letter queues for entities
- Fix connection leaks when using multiple calls to Receive
- Ensure senders wait for message disposition before returning

## `v0.1.0`
- initial tag for Service Bus which includes Queues, Topics and Subscriptions using AMQP