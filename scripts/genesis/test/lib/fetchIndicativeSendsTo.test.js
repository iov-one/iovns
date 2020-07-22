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

   it( `Should match star1's first 2 incoming payments.`, async () => {
      const indicative = [
         { "hash": "e0d65bc5377e0806de18f76e07c3234632fad570a799c1063df1f69809bf4337", "block_height": 65609, "message": { "path": "cash/send", "details": { "memo": "star1cnywewxct2p4d5j2fapgkse6yxgh7ecnj4uwpu", "amount": { "whole": 1, "ticker": "IOV" }, "source": "iov1yhk8qqp3wsdg7tefd8u457n9zqsny4nqzp6960", "destination": "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6" } } },
         { "hash": "20894f0429901e402bb0520d117da9b64dacce2a97b647c66645bf6436af17d7", "block_height": 67029, "message": { "path": "cash/send", "details": { "memo": "star19m9ufykj5ur67l822fpxvz49p535wp3j0m5v3h", "amount": { "ticker": "IOV", "fractional": 1 }, "source": "iov1a9duw7yyxdfh8mrjxmuc0slu8a48muvxkcxvg8", "destination": "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6" } } }
      ];
      const sends = await fetchIndicativeSendsTo( "iov10v69k57z2v0pr3yvtr60pp8g2jx8tdd7f55sv6", /star1[qpzry9x8gf2tvdw0s3jn54khce6mua7l]{38}/ );

      compareObjects( indicative, sends );
   } );
} );
