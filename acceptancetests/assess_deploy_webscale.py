#!/usr/bin/env python3
""" Test webscale deployment

    1. deploying kubernetes core and asserting it is `healthy`
    2. inspect the logs to parse timings from trace logs
"""

from __future__ import print_function

import argparse
import logging
import sys
import os
import subprocess
import re
import requests
import functools
import time

from deploy_stack import (
    BootstrapManager,
    deploy_caas_stack,
)
from utility import (
    add_basic_testing_arguments,
    configure_logging,
    JujuAssertionError,
    get_current_model,
)

from jujucharm import (
    local_charm_path,
)
from jujupy.utility import until_timeout

__metaclass__ = type

log = logging.getLogger("assess_deploy_webscale")

def deploy_bundle(client, charm_bundle):
    """Deploy the given charm bundle

    :param client: Jujupy ModelClient object
    :param charm_bundle: Optional charm bundle string
    """
    model_name = "webscale"

    bundle = None
    if not charm_bundle:
        bundle = local_charm_path(
            charm='bundles-kubernetes-core-lxd.yaml',
            repository=os.environ['JUJU_REPOSITORY'],
            juju_ver=client.version,
        )
    else:
        bundle = charm_bundle

    caas_client = deploy_caas_stack(
        path=bundle,
        client=client,
        charm=(not not charm_bundle),
    )

    if not caas_client.is_cluster_healthy:
        raise JujuAssertionError('k8s cluster is not healthy because kubectl is not accessible')

    current_model = caas_client.add_model(model_name)
    current_model.juju(current_model._show_status, ('--format', 'tabular'))

    current_model.destroy_model()

def extract_module_logs(client, module):
    """Extract the logs from destination module.

    :param module: string containing the information to extract from the destination module.
    """
    deploy_logs = client.get_juju_output(
        'debug-log', '-m', 'controller',
        '--no-tail', '--replay', '-l', 'TRACE',
        '--include-module', module,
    )
    return deploy_logs

def extract_txn_timings(logs, module):
    """Extract the transaction timings (txn) from the deploy logs.

    It's expected that the timings are in seconds to 3 decimal places ("0.042s")

    :param logs: string containing all the logs from the module
    :param module: string containing the destination module.
    """
    exp = re.compile(r'{} ran transaction in (?P<seconds>\d+\.\d+)s'.format(module), re.IGNORECASE)
    timings = []
    for timing in exp.finditer(logs):
        timings.append(timing.group("seconds"))
    return list(map(float, timings))

def calculate_total_time(timings):
    """Accumulate transaction timings (txn) from the timings.

    :param timings: expects timings to be floats
    """
    return functools.reduce(lambda x, y: x + y, timings)

def calculate_max_time(timings):
    """Calculate maximum transaction timing from (txn).

    :param timings: expects timings to be floats
    """
    return functools.reduce(lambda x, y: x if x > y else y, timings)

def parse_args(argv):
    """Parse all arguments."""
    parser = argparse.ArgumentParser(description="Webscale charm deployment CI test")
    parser.add_argument(
        '--charm-bundle',
        help="Override the charm bundle to deploy",
    )
    parser.add_argument(
        '--logging-config',
        help="Override logging configuration for a deploy",
        default="juju.state.txn=TRACE;<root>=INFO;unit=INFO",
    )
    parser.add_argument(
        '--logging-module',
        help="Override default module to extract",
        default="juju.state.txn",
    )
    add_basic_testing_arguments(parser, existing=False)
    return parser.parse_args(argv)

def main(argv=None):
    args = parse_args(argv)
    configure_logging(args.verbose)
    begin = time.time()
    bs_manager = BootstrapManager.from_args(args)
    with bs_manager.booted_context(args.upload_tools):
        client = bs_manager.client
        deploy_bundle(client, charm_bundle=args.charm_bundle)
        raw_logs = extract_module_logs(client, module=args.logging_module)
        timings = extract_txn_timings(raw_logs, module=args.logging_module)
        # Calculate the timings to forward to the datastore
        total_time = calculate_total_time(timings)
        total_txns = len(timings)
        max_time = calculate_max_time(timings)
        since = (time.time() - begin)

        log.info("The timings for deployment: total txn time: {}, total nums txn: {}, max txn time: {}, max time (seconds): {}".format(total_time, total_txns, max_time, since))
    return 0

if __name__ == '__main__':
    sys.exit(main())