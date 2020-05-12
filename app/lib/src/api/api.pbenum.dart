///
//  Generated code. Do not modify.
//  source: api.proto
//
// @dart = 2.3
// ignore_for_file: camel_case_types,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name,return_of_invalid_type

// ignore_for_file: UNDEFINED_SHOWN_NAME,UNUSED_SHOWN_NAME
import 'dart:core' as $core;
import 'package:protobuf/protobuf.dart' as $pb;

class SessionState extends $pb.ProtobufEnum {
  static const SessionState LOBBY = SessionState._(0, 'LOBBY');
  static const SessionState RUNNING = SessionState._(1, 'RUNNING');

  static const $core.List<SessionState> values = <SessionState> [
    LOBBY,
    RUNNING,
  ];

  static final $core.Map<$core.int, SessionState> _byValue = $pb.ProtobufEnum.initByValue(values);
  static SessionState valueOf($core.int value) => _byValue[value];

  const SessionState._($core.int v, $core.String n) : super(v, n);
}

