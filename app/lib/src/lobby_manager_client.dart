import 'dart:async';
import 'dart:convert';

import 'package:couchcampaign/src/api/api.pb.dart';
import 'package:couchcampaign/src/rpc_client.dart';

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
