# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: marble.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0cmarble.proto\x12\x03rpc\"(\n\x11StartOperationReq\x12\x13\n\x0bnOperations\x18\x01 \x01(\x05\" \n\x12StartOperationResp\x12\n\n\x02ok\x18\x01 \x01(\x08\x32U\n\x10OperationTracker\x12\x41\n\x0eStartOperation\x12\x16.rpc.StartOperationReq\x1a\x17.rpc.StartOperationRespB-Z+github.com/edgelesssys/marblerun/marble/rpcb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'marble_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z+github.com/edgelesssys/marblerun/marble/rpc'
  _STARTOPERATIONREQ._serialized_start=21
  _STARTOPERATIONREQ._serialized_end=61
  _STARTOPERATIONRESP._serialized_start=63
  _STARTOPERATIONRESP._serialized_end=95
  _OPERATIONTRACKER._serialized_start=97
  _OPERATIONTRACKER._serialized_end=182
# @@protoc_insertion_point(module_scope)