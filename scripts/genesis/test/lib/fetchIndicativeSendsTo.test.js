import compareObjects from "../compareObjects";
import fetchIndicativeSendsTo from "../../lib/fetchIndicativeSendsTo";

"use strict";


describe( "Tests ../../lib/fetchIndicativeSendsTo.js.", () => {
   it( `Should match dave's first 3 incoming payments.`, async () => {
      const indicative = [
         { "hash": "6104533d02587e870c12453ecfad1bd73289cfed7a8aaf2caadbed5685dbe209", "block_height": 9522, "message": { "path": "cash/send", "details": { "memo": "kick the chain", "amount": { "whole": 1, "ticker": "IOV" }, "source": "iov15xzzgu5jkltm24hg9r2ykm6gm2d09tzrcqgrr9", "destination": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" } } },
         { "hash": "acfd7092f62e8e6c119cbd66257aa11c043ca8adf2f60111518370c7de4e537f", "block_height": 17625, "message": { "path": "cash/send", "details": { "memo": "kick", "amount": { "whole": 1, "ticker": "IOV" }, "source": "iov1c5vttpanpp6wxq7naee5am7yw6l0vu8u3d0qug", "destination": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" } } },
         { "hash": "be68add27fd834eaef3ae1e39bbf714d0e8a5c9b67febf8e0c1daebee5828936", "block_height": 17629, "message": { "path": "cash/send", "details": { "memo": "kick", "amount": { "ticker": "IOV", "fractional": 100000000 }, "source": "iov1c5vttpanpp6wxq7naee5am7yw6l0vu8u3d0qug", "destination": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" } } }
      ];
      const sends = await fetchIndicativeSendsTo( "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un", /kick/ );

      compareObjects( indicative, sends );
   } );
} );
