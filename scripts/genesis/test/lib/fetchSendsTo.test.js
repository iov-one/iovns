import compareObjects from "../compareObjects";
import fetchSendsTo from "../../lib/fetchSendsTo";

"use strict";


describe( "Tests ../../lib/fetchSendsTo.js.", () => {
   it( `Should get at least 4 block ordered sends to syncnode.`, async () => {
      const sends = await fetchSendsTo( "iov1psjk38fdzxtx6ypsm5u5ujnt6dken4kfelja24" );

      expect( sends.length ).toBeGreaterThanOrEqual( 4 );
      expect( sends[0].block_height ).toBeLessThan( sends[1].block_height );
      expect( sends[1].block_height ).toBeLessThan( sends[2].block_height );
      expect( sends[2].block_height ).toBeLessThan( sends[3].block_height );
   } );


   it( `Should match dave's first 3 incoming payments.`, async () => {
      const payments = [
         { "hash": "842150818f53529e1aa40e8367300bbab5ff9a44530553c55679838ba57ef135", "block_height": 208, "message": { "path": "cash/send", "details": { "memo": "", "amount": { "whole": 454, "ticker": "IOV", "fractional": 500000000 }, "source": "iov15xzzgu5jkltm24hg9r2ykm6gm2d09tzrcqgrr9", "destination": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" } } },
         { "hash": "6104533d02587e870c12453ecfad1bd73289cfed7a8aaf2caadbed5685dbe209", "block_height": 9522, "message": { "path": "cash/send", "details": { "memo": "kick the chain", "amount": { "whole": 1, "ticker": "IOV" }, "source": "iov15xzzgu5jkltm24hg9r2ykm6gm2d09tzrcqgrr9", "destination": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" } } },
         { "hash": "c40d496564d867e436c8d3badde7fa34e35a17fbe4bb1d70982be868afd5565d", "block_height": 9607, "message": { "path": "cash/send", "details": { "memo": "testing name services", "amount": { "whole": 1, "ticker": "IOV" }, "source": "iov1psjk38fdzxtx6ypsm5u5ujnt6dken4kfelja24", "destination": "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" } } }
      ];
      const sends = await fetchSendsTo( "iov1qnpaklxv4n6cam7v99hl0tg0dkmu97sh6007un" );

      compareObjects( payments, sends.splice( 0, 3 ) );
   } );
} );
