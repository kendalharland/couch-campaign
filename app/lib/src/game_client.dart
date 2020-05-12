import 'dart:async';

import 'package:couchcampaign/src/api/api.pb.dart';
import 'package:couchcampaign/src/rpc_client.dart';
import 'package:couchcampaign/src/session.dart';

class GameClient {
  static const methodStart = 'start';
  static const methodSocket = 'socket';

  const GameClient(this.client);

  final RpcClient client;

  Future<GameSession> joinGame() async {
    final socketAddress = await client.call(methodSocket);
    return GameSession.connect(Uri.parse(socketAddress));
  }

  Future<void> startGame() => client.call(methodStart);
}

typedef StateHandler = void Function(GameState);
typedef LobbyHandler = void Function(GameState);
