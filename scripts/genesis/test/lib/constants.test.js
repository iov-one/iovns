import { chainIds, conds, multisigs, names, source2multisig } from "../../lib/constants";

"use strict";


describe( "Tests ../../lib/constants.js.", () => {
   // conds
   const gov2 = "cond:gov/rule/0000000000000002";
   const multisig1 = "cond:multisig/usage/0000000000000001";
   const multisig2 = "cond:multisig/usage/0000000000000002";
   const multisig3 = "cond:multisig/usage/0000000000000003";
   const multisig4 = "cond:multisig/usage/0000000000000004";
   const multisig5 = "cond:multisig/usage/0000000000000005";
   const multisig6 = "cond:multisig/usage/0000000000000006";

   // names
   const reward = "reward fund";
   const iov = "IOV SAS";
   const employee = "IOV SAS employee bonus pool/colloboration appropriation pool";
   const pending = "IOV SAS pending deals pocket; close deal or burn";
   const bounty = "IOV SAS bounty fund";
   const cofounders = "Unconfirmed contributors/co-founders";
   const custodian = "Custodian of missing star1 accounts";


   it( `Should get multisigs keyed on iov1.`, () => {
      expect( Object.keys( multisigs ).length ).toEqual( 7 );

      expect( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n["//name"] ).toEqual( reward );
      expect( multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq["//name"] ).toEqual( iov );
      expect( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0["//name"] ).toEqual( employee );
      expect( multisigs.iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz["//name"] ).toEqual( pending );
      expect( multisigs.iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu["//name"] ).toEqual( bounty );
      expect( multisigs.iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw["//name"] ).toEqual( cofounders );
      expect( multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc["//name"] ).toEqual( custodian );

      expect( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.address ).toEqual( gov2 );
      expect( multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq.address ).toEqual( multisig1 );
      expect( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0.address ).toEqual( multisig2 );
      expect( multisigs.iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz.address ).toEqual( multisig3 );
      expect( multisigs.iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu.address ).toEqual( multisig4 );
      expect( multisigs.iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw.address ).toEqual( multisig5 );
      expect( multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc.address ).toEqual( multisig6 );

      expect( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.star1 ).toEqual( "star1scfumxscrm53s4dd3rl93py5ja2ypxmxlhs938" );
      expect( multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq.star1 ).toEqual( "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78" );
      expect( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0.star1 ).toEqual( "star16tm7scg0c2e04s0exk5rgpmws2wk4xkd84p5md" );
      expect( multisigs.iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz.star1 ).toEqual( "star1uyny88het6zaha4pmkwrkdyj9gnqkdfe4uqrwq" );
      expect( multisigs.iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu.star1 ).toEqual( "star1m7jkafh4gmds8r0w79y2wu2kvayqvrwt7cy7rf" );
      expect( multisigs.iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw.star1 ).toEqual( "star1p0d75y4vpftsx9z35s93eppkky7kdh220vrk8n" );
      expect( multisigs.iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc.star1 ).toEqual( "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y" );
   } );

   it( `Should get multisigs keyed on cond.`, () => {
      expect( Object.keys( conds ).length ).toEqual( 7 );

      expect( conds[gov2     ]["//name"] ).toEqual( reward );
      expect( conds[multisig1]["//name"] ).toEqual( iov );
      expect( conds[multisig2]["//name"] ).toEqual( employee );
      expect( conds[multisig3]["//name"] ).toEqual( pending );
      expect( conds[multisig4]["//name"] ).toEqual( bounty );
      expect( conds[multisig5]["//name"] ).toEqual( cofounders );

      expect( conds[gov2     ].iov1 ).toEqual( "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n" );
      expect( conds[multisig1].iov1 ).toEqual( "iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq" );
      expect( conds[multisig2].iov1 ).toEqual( "iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0" );
      expect( conds[multisig3].iov1 ).toEqual( "iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz" );
      expect( conds[multisig4].iov1 ).toEqual( "iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu" );
      expect( conds[multisig5].iov1 ).toEqual( "iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw" );

      expect( conds[gov2     ].star1 ).toEqual( "star1scfumxscrm53s4dd3rl93py5ja2ypxmxlhs938" );
      expect( conds[multisig1].star1 ).toEqual( "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78" );
      expect( conds[multisig2].star1 ).toEqual( "star16tm7scg0c2e04s0exk5rgpmws2wk4xkd84p5md" );
      expect( conds[multisig3].star1 ).toEqual( "star1uyny88het6zaha4pmkwrkdyj9gnqkdfe4uqrwq" );
      expect( conds[multisig4].star1 ).toEqual( "star1m7jkafh4gmds8r0w79y2wu2kvayqvrwt7cy7rf" );
      expect( conds[multisig5].star1 ).toEqual( "star1p0d75y4vpftsx9z35s93eppkky7kdh220vrk8n" );
   } );

   it( `Should get multisigs keyed on "//name".`, () => {
      expect( Object.keys( names ).length ).toEqual( 7 );

      expect( names[reward    ].cond ).toEqual( gov2 );
      expect( names[iov       ].cond ).toEqual( multisig1 );
      expect( names[employee  ].cond ).toEqual( multisig2 );
      expect( names[pending   ].cond ).toEqual( multisig3 );
      expect( names[bounty    ].cond ).toEqual( multisig4 );
      expect( names[cofounders].cond ).toEqual( multisig5 );
      expect( names[custodian ].cond ).toEqual( multisig6 );

      expect( names[reward    ].iov1 ).toEqual( "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n" );
      expect( names[iov       ].iov1 ).toEqual( "iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq" );
      expect( names[employee  ].iov1 ).toEqual( "iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0" );
      expect( names[pending   ].iov1 ).toEqual( "iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz" );
      expect( names[bounty    ].iov1 ).toEqual( "iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu" );
      expect( names[cofounders].iov1 ).toEqual( "iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw" );
      expect( names[custodian ].iov1 ).toEqual( "iov195cpqyk5sjh7qwfz8qlmlnz2vw4ylz394smqvc" );

      expect( names[reward    ].star1 ).toEqual( "star1scfumxscrm53s4dd3rl93py5ja2ypxmxlhs938" );
      expect( names[iov       ].star1 ).toEqual( "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78" );
      expect( names[employee  ].star1 ).toEqual( "star16tm7scg0c2e04s0exk5rgpmws2wk4xkd84p5md" );
      expect( names[pending   ].star1 ).toEqual( "star1uyny88het6zaha4pmkwrkdyj9gnqkdfe4uqrwq" );
      expect( names[bounty    ].star1 ).toEqual( "star1m7jkafh4gmds8r0w79y2wu2kvayqvrwt7cy7rf" );
      expect( names[cofounders].star1 ).toEqual( "star1p0d75y4vpftsx9z35s93eppkky7kdh220vrk8n" );
      expect( names[custodian ].star1 ).toEqual( "star12uv6k3c650kvm2wpa38wwlq8azayq6tlh75d3y" );
   } );

   it( `Should validate CAIP-based chain ids.`, () => {
      Object.values( chainIds ).forEach( id => {
         expect( id.indexOf( ":" ) ).not.toEqual( -1 );
      } );
   } );

   it( `Should get validate star1 addresses for escrows.`, () => {
      const isabella = source2multisig.iov1w2suyhrfcrv5h4wmq3rk3v4x95cxtu0a03gy6x;
      const kadima = source2multisig.iov1v9pzqxpywk05xn2paf3nnsjlefsyn5xu3nwgph;
      const joghurt = source2multisig.iov149cn0rauw2773lfdp34njyejg3cfz2d56c0m5t;
      const guaranteed = source2multisig.iov1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqvnwh0u;

      expect( isabella[  "//id"] ).toEqual( "escrow isabella*iov" );
      expect( kadima[    "//id"] ).toEqual( "escrow kadima*iov" );
      expect( joghurt[   "//id"] ).toEqual( "escrow joghurt*iov" );
      expect( guaranteed["//id"] ).toEqual( "vaildator guaranteed reward fund" );

      expect( isabella.star1   ).toEqual( "star1elad203jykd8la6wgfnvk43rzajyqpk0wsme9g" );
      expect( kadima.star1     ).toEqual( "star1hjf04872s9rlcdg2wqwvapwttvt3p4gjpp0xmc" );
      expect( joghurt.star1    ).toEqual( "star15u4kl3lalt8pm2g4m23erlqhylz76rfh50cuv8" );
      expect( guaranteed.star1 ).toEqual( "star17w7fjdkr9laphtyj4wxa32rf0evu94xgywxgl4" );
   } );
} );
