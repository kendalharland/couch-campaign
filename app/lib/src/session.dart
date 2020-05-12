import 'dart:async';

import 'package:couchcampaign/src/api/api.pb.dart';
import 'package:couchcampaign/src/message_decoder.dart';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

//
// Problems:
// * In-order: buffered delivery of game messages
// * Messy code: the protocol should be entirely specified in terms of protos.

typedef MessageHandler = void Function(Message);

class GameSession extends ValueNotifier<Message> {
  static Future<GameSession> connect(Uri address) async =>
      GameSession(WebSocketChannel.connect(address));

  GameSession(this._ws) : super(null) {
    _wsSub = _ws.stream.transform<Message>(messageDecoder()).listen((message) {
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
