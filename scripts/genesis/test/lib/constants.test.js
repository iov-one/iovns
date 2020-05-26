import { conds, multisigs, names } from "../../lib/constants";

"use strict";


describe( "Tests ../../lib/constants.js.", () => {
   // conds
   const gov2 = "cond:gov/rule/0000000000000002";
   const multisig1 = "cond:multisig/usage/0000000000000001";
   const multisig2 = "cond:multisig/usage/0000000000000002";
   const multisig3 = "cond:multisig/usage/0000000000000003";
   const multisig4 = "cond:multisig/usage/0000000000000004";
   const multisig5 = "cond:multisig/usage/0000000000000005";

   // names
   const reward = "reward fund";
   const iov = "IOV SAS";
   const employee = "IOV SAS employee bonus pool/colloboration appropriation pool";
   const pending = "IOV SAS pending deals pocket; close deal or burn";
   const bounty = "IOV SAS bounty fund";
   const cofounders = "Unconfirmed contributors/co-founders";


   it( `Should get multisigs keyed on iov1.`, () => {
      expect( Object.keys( multisigs ).length ).toEqual( 6 );

      expect( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n["//name"] ).toEqual( reward );
      expect( multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq["//name"] ).toEqual( iov );
      expect( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0["//name"] ).toEqual( employee );
      expect( multisigs.iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz["//name"] ).toEqual( pending );
      expect( multisigs.iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu["//name"] ).toEqual( bounty );
      expect( multisigs.iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw["//name"] ).toEqual( cofounders );

      expect( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.address ).toEqual( gov2 );
      expect( multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq.address ).toEqual( multisig1 );
      expect( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0.address ).toEqual( multisig2 );
      expect( multisigs.iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz.address ).toEqual( multisig3 );
      expect( multisigs.iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu.address ).toEqual( multisig4 );
      expect( multisigs.iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw.address ).toEqual( multisig5 );

      expect( multisigs.iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n.star1 ).toEqual( "TBD" ); // TODO
      expect( multisigs.iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq.star1 ).toEqual( "TBD" ); // TODO
      expect( multisigs.iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0.star1 ).toEqual( "TBD" ); // TODO
      expect( multisigs.iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz.star1 ).toEqual( "TBD" ); // TODO
      expect( multisigs.iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu.star1 ).toEqual( "TBD" ); // TODO
      expect( multisigs.iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw.star1 ).toEqual( "TBD" ); // TODO
   } );

   it( `Should get multisigs keyed on cond.`, () => {
      expect( Object.keys( conds ).length ).toEqual( 6 );

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

      expect( conds[gov2     ].star1 ).toEqual( "TBD" ); // TODO
      expect( conds[multisig1].star1 ).toEqual( "TBD" ); // TODO
      expect( conds[multisig2].star1 ).toEqual( "TBD" ); // TODO
      expect( conds[multisig3].star1 ).toEqual( "TBD" ); // TODO
      expect( conds[multisig4].star1 ).toEqual( "TBD" ); // TODO
      expect( conds[multisig5].star1 ).toEqual( "TBD" ); // TODO
   } );

   it( `Should get multisigs keyed on "//name".`, () => {
      expect( Object.keys( names ).length ).toEqual( 6 );

      expect( names[reward    ].cond ).toEqual( gov2 );
      expect( names[iov       ].cond ).toEqual( multisig1 );
      expect( names[employee  ].cond ).toEqual( multisig2 );
      expect( names[pending   ].cond ).toEqual( multisig3 );
      expect( names[bounty    ].cond ).toEqual( multisig4 );
      expect( names[cofounders].cond ).toEqual( multisig5 );

      expect( names[reward    ].iov1 ).toEqual( "iov1k0dp2fmdunscuwjjusqtk6mttx5ufk3zpwj90n" );
      expect( names[iov       ].iov1 ).toEqual( "iov1tt3vtpukkzk53ll8vqh2cv6nfzxgtx3t52qxwq" );
      expect( names[employee  ].iov1 ).toEqual( "iov1zd573wa38pxfvn9mxvpkjm6a8vteqvar2dwzs0" );
      expect( names[pending   ].iov1 ).toEqual( "iov1ppzrq5gwqlcsnwdvlz7x9mu98fntmp65m9a3mz" );
      expect( names[bounty    ].iov1 ).toEqual( "iov1ym3uxcfv9zar2md0xd3hq2vah02u3fm6zn8mnu" );
      expect( names[cofounders].iov1 ).toEqual( "iov1myq53ry9pa6awl88m0xgp224q0dgwjdvz2dcsw" );

      expect( names[reward    ].star1 ).toEqual( "TBD" ); // TODO
      expect( names[iov       ].star1 ).toEqual( "TBD" ); // TODO
      expect( names[employee  ].star1 ).toEqual( "TBD" ); // TODO
      expect( names[pending   ].star1 ).toEqual( "TBD" ); // TODO
      expect( names[bounty    ].star1 ).toEqual( "TBD" ); // TODO
      expect( names[cofounders].star1 ).toEqual( "TBD" ); // TODO
   } );
} );
