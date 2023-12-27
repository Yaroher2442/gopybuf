from __future__ import annotations

import asyncio
import collections
import json
import traceback
from typing import Dict, Protocol, Mapping, Tuple

from gopybuf.stream import CustomStream

Handler = collections.namedtuple(
    'Handler', 'func, cardinality, request_type, reply_type',
)


class IServable(Protocol):
    def __mapping__(self) -> Mapping[str, Handler]: ...


IncomingBytes = bytes
OutgoingBytes = bytes
ErrorBytes = bytes


class CustomErr(Dict):

    def __init__(self, e: Exception, traceback_str: str = None):
        super().__init__()
        self["error"] = str(e)
        self["traceback"] = traceback_str

    def dump(self) -> bytes:
        return bytes(json.dumps(self), 'utf-8')


class CustomServer:
    _instance: CustomServer = None
    _mappings: Dict[str, Handler]

    def __new__(cls, *args, **kwargs):
        if not cls._instance:
            cls._instance = super(CustomServer, cls).__new__(cls, *args, **kwargs)
        return cls._instance

    def add_mapping(self, mappings: IServable):
        self._mappings.update(mappings.__mapping__())

    async def __call__(self, method_name: str, request: IncomingBytes) -> OutgoingBytes:
        method = self._mappings[method_name]
        stream = CustomStream(method.request_type.parse(request), method.func)
        return OutgoingBytes(await stream())


_global_server = CustomServer()


def register_service(mapping: IServable):
    _global_server.add_mapping(mapping)


async def call_async(method_name: str, arg: IncomingBytes) -> Tuple[OutgoingBytes, ErrorBytes]:
    try:
        return await _global_server(method_name, arg), ErrorBytes()
    except Exception as e:
        return OutgoingBytes(), CustomErr(e, traceback.format_exc()).dump()


def go_py_buf(method_name: str, arg: IncomingBytes) -> Tuple[OutgoingBytes, ErrorBytes]:
    return asyncio.run(call_async(method_name, arg))

# async def __call(method_name: str, arg: bytes) -> Tuple[bytes, bytes]:
#     try:
#         server = CustomServer([EchoService()])
#         return await server(method_name, arg), bytes()
#     except Exception as e:
#         return bytes(), CustomErr(e, traceback.format_exc()).dump()
#
#
# def call_go_py(method_name: str, arg: bytes) -> Tuple[bytes, bytes]:
#     return asyncio.run(__call(method_name, arg))
#
