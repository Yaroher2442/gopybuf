from __future__ import annotations

import collections
from typing import Dict, Protocol, Mapping

from gopybuf.stream import CustomStream

Handler = collections.namedtuple(
    'Handler', 'func, cardinality, request_type, reply_type',
)


class IServable(Protocol):
    def __mapping__(self) -> Mapping[str, Handler]: ...


IncomingBytes = bytes
OutgoingBytes = bytes
ErrorBytes = bytes


class CustomServer:
    _instance: CustomServer = None
    _mappings: Dict[str, Handler]

    def __new__(cls, *args, **kwargs):
        if not cls._instance:
            cls._instance = super(CustomServer, cls).__new__(cls, *args, **kwargs)
        return cls._instance

    def add_mapping(self, mappings: IServable):
        self._mappings.update(mappings.__mapping__())

    def __call__(self, method_name: str, request: IncomingBytes) -> OutgoingBytes:
        method = self._mappings[method_name]
        stream = CustomStream(method.request_type.parse(request), method.func)
        return OutgoingBytes(stream())


_global_server = CustomServer()
