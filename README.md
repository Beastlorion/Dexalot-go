# Dexalot-go
Open source bots for dexalot.

# Table of Contents
## Known Areas of Improvement
### Receiving Headers
   * The HeaderReceiver in `src/geth` will query all logs for the block
   * This can be improved by using the `types.Bloom` on the `types.Header`
   * For now, filtering is done downstream 
### Websocket Client should be refactored 
   * The design of the websocket client should be improved
   * Passing around the reference to the connection within the Client methods is not ideal
### Adding Request for Open Orders on Dexalot
   * This should be fairly straightforward to add 
   * It requires expanding the `RESTHelper` struct to support signing requests
### Order Books
   * Currently, we support `maker.OrderBook` which is specific to order placement
   * A better way to do this would be to have our own order book implementation
   * This should not rely on the `RedBlackTree` wrapper, as we should have our own custom implementation
### Curve Models
   * Currently, the curve models are very rudimentary 
   * They were primarily added for testing purposes and to provide a starting point
   * They should be expanded, but should work for simple Market Makers that simply want to place orders
### Sending Orders 
   * The `maker.Dispatcher` interface is also very rudimentary
   * Depending on the concurrency design, it is expected to change as will the API surrounding the `maker.Engine`
