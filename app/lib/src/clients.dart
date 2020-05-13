import 'dart:async';
import 'dart:convert';

import 'package:couchcampaign/src/api/api.pb.dart';
import 'package:couchcampaign/src/session.dart';
import 'package:xhttp/xhttp.dart' as http;

typedef StateHandler = void Function(PlayerState);

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

class LobbyManagerClient {
  static const methodCreateLobby = 'lobby/create';

  const LobbyManagerClient(this.client);

  final RpcClient client;

  Future<CreateGameResponse> createGame(CreateGameRequest request) async {
    final request = json.encode(CreateGameRequest().toProto3Json());
    final response = await client.call(methodCreateLobby, request);
    final responseJs = json.decode(response);
    return CreateGameResponse()..mergeFromProto3Json(responseJs);
  }
}

class RpcClient {
  const RpcClient(this.client, this.address);

  final http.Client client;
  final String address;

  Future<String> call(String method, [String message = '']) async {
    ArgumentError.checkNotNull(message, 'message');
    final url = Uri.parse('$address/$method');
    final response = await client.post(url, bodyBytes: message.codeUnits);
    if (response.statusCode != 200) {
      throw Exception('${response.statusCode}: ${await response.body}');
    }
    return response.body;
  }
}
