import os
import glob
import ptflow.pflow as pflow_loader

EVENTSTORE = {}
""" loaded eventstore objects indexed by schema """

def initialize(provider, **kwargs):
    """ reload eventstorage instances """

    pflow_loader.set_provider(provider)
    provider.reconnect(**kwargs)
    provider.migrate()

    if not 'dirname' in kwargs:
        kwargs['dirname'] = os.environ.get(
            'PTFLOW_DIR',
            os.path.dirname(os.path.abspath(__file__)) + "/examples/"
        )

    for pf in  glob.glob(kwargs['dirname'] + "*.pflow"):
        es, _ = pflow_loader.load_file(pf)
        EVENTSTORE[es.name] = es.to_module().Machine

def eventstore(schema=None, oid=None):
    """ get statemachine eventstore by schema name """
    return EVENTSTORE[schema](oid, schema)

def schemata():
    """ list names of state machines """
    return [ k for k in EVENTSTORE ]
