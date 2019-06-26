#!/usr/bin/env python3

import copy

import lib

def process_genesis(genesis, parsed_args):
    genesis['app_state']['account'] = {
        'params': {
            'mint_denom': 'trusteak',
            'inflation_rate_change': '0.150000000000000000',
            'inflation_max': '0.250000000000000000',
            'inflation_min': '0.100000000000000000',
            'goal_bonded': '0.670000000000000000',
            'blocks_per_year': '6311520',
        }
    }

    # Set new chain ID and genesis start time
    genesis['chain_id'] = parsed_args.chain_id.strip()
    genesis['genesis_time'] = parsed_args.start_time

    return genesis

if __name__ == '__main__':
    parser = lib.init_default_argument_parser(
        prog_desc='Migrate genesis.json to add account module',
        default_chain_id='devnet-n',
        default_start_time='2019-02-11T12:00:00Z',
    )
    lib.main(parser, process_genesis)
