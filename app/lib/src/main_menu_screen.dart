import 'package:couchcampaign/src/api/api.pb.dart';
import 'package:couchcampaign/src/game_client.dart';
import 'package:couchcampaign/src/game.dart';
import 'package:couchcampaign/src/lobby_manager_client.dart';
import 'package:couchcampaign/src/rpc_client.dart';
import 'package:flutter/material.dart';
import 'package:xhttp/xhttp.dart' as http;

class MainMenuScreen extends StatelessWidget {
  const MainMenuScreen(this.lmc, {@required this.onGameCreated});

  final LobbyManagerClient lmc;
  final ValueSetter<Widget> onGameCreated;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: EdgeInsets.all(24),
        child: ListView(
          shrinkWrap: true,
          children: [
            RaisedButton(onPressed: _createGame, child: Text('Create game')),
          ],
        ),
      ),
    );
  }

  Future<void> _createGame() async {
    final response = await lmc.createGame(CreateGameRequest());
    await Future.delayed(const Duration(milliseconds: 200));
    final client = GameClient(RpcClient(http.Client(), response.gameUrl));
    final session = await client.joinGame();
    onGameCreated(GameSessionView(client, session));
  }
}
