from __future__ import annotations

from typing import Generic, Callable, TypeVar

SendType = TypeVar('SendType')
RecvType = TypeVar('RecvType')


class CustomStream(Generic[RecvType, SendType]):
    __target_call = None
    __resp: SendType

    def __init__(self, req: RecvType, func: TargetCall):
        self._req = req
        self.__resp = None
        self.__target_call = func

    async def recv_message(self) -> RecvType:
        return self._req

    async def send_message(self, message: SendType) -> None:
        self.__resp = message

    async def __call__(self, *args, **kwargs) -> SendType:
        await self.__target_call(self)
        return self.__resp


TargetCall = Callable[[CustomStream[RecvType, SendType]], None]
