import 'dart:async';

import 'package:couchcampaign/src/api/api.pb.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

typedef MessageHandler = void Function(Message);

class GameSession extends ValueNotifier<Message> {
  static Future<GameSession> connect(Uri address) async =>
      GameSession(WebSocketChannel.connect(address));

  static StreamTransformer<dynamic, Message> _decoder() =>
      StreamTransformer<dynamic, Message>.fromHandlers(
        handleData: (data, EventSink<Message> sink) =>
            sink.add(Message.fromBuffer(data)),
      );

  GameSession(this._ws) : super(null) {
    _wsSub = _ws.stream.transform<Message>(_decoder()).listen((message) {
      value = message;
    });
  }

  final WebSocketChannel _ws;
  StreamSubscription<dynamic> _wsSub;

  @override
  void dispose() {
    super.dispose();
    _wsSub.cancel();
  }

  Future<void> send(String input) async {
    _ws.sink.add(input);
  }
}
