import 'dart:async';

import 'package:couchcampaign/src/api/api.pb.dart';

StreamTransformer<dynamic, Message> messageDecoder() =>
    StreamTransformer<dynamic, Message>.fromHandlers(
      handleData: (data, EventSink<Message> sink) =>
          sink.add(Message.fromBuffer(data)),
    );
