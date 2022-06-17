/*
Package rpc implements bridge to Orion full node API interface.

We recommend using local IPC for fast and the most efficient inter-process communication between the API server
and an Orion node. Any remote RPC connection will work, but the performance may be significantly degraded
by extra networking overhead of remote RPC calls.

You should also consider security implications of opening Orion RPC interface for a remote access.
If you considering it as your deployment strategy, you should establish encrypted channel between the API server
and NEXT RPC interface with connection limited to specified endpoints.

We strongly discourage opening Orion RPC interface for unrestricted Internet access.
*/

package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"time"
)

// nextHeadsObserverSubscribeTick represents the time between subscription attempts.
const nextHeadsObserverSubscribeTick = 30 * time.Second

// observeBlocks collects new blocks from the blockchain network
// and posts them into the proxy channel for processing.
func (next *NextBridge) observeBlocks() {
	var sub ethereum.Subscription
	defer func() {
		if sub != nil {
			sub.Unsubscribe()
		}
		next.log.Noticef("block observer done")
		next.wg.Done()
	}()

	sub = next.blockSubscription()
	for {
		// re-subscribe if the subscription ref is not valid
		if sub == nil {
			tm := time.NewTimer(nextHeadsObserverSubscribeTick)
			select {
			case <-next.sigClose:
				return
			case <-tm.C:
				sub = next.blockSubscription()
				continue
			}
		}

		// use the subscriptions
		select {
		case <-next.sigClose:
			return
		case err := <-sub.Err():
			next.log.Errorf("block subscription failed; %s", err.Error())
			sub = nil
		}
	}
}

// blockSubscription provides a subscription for new blocks received
// by the connected blockchain node.
func (next *NextBridge) blockSubscription() ethereum.Subscription {
	sub, err := next.rpc.EthSubscribe(context.Background(), next.headers, "newHeads")
	if err != nil {
		next.log.Criticalf("can not observe new blocks; %s", err.Error())
		return nil
	}
	return sub
}
