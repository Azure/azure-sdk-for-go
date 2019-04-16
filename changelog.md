# Change Log

## `head`

## `v0.4.1`
- fix issue with sender when SB returns a different receiver disposition [#119](https://github.com/Azure/azure-service-bus-go/issues/119)

## `v0.4.0`
- Update to AMQP 0.11.0 which introduces strict settlement mode
  ([#111](https://github.com/Azure/azure-service-bus-go/issues/111))

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