import 'package:couchcampaign/src/lobby_manager_client.dart';
import 'package:couchcampaign/src/main_menu_screen.dart';
import 'package:couchcampaign/src/rpc_client.dart';
import 'package:flutter/material.dart' hide Card;
import 'package:xhttp/xhttp.dart' as http;

class App extends StatelessWidget {
  static const gameServerAddress = 'http://localhost:8080';

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: Scaffold(body: CouchCampaign(gameServerAddress)),
    );
  }
}

class CouchCampaign extends StatefulWidget {
  const CouchCampaign(this.address);

  final String address;

  @override
  State<StatefulWidget> createState() => CouchCampaignState();
}

class CouchCampaignState extends State<CouchCampaign> {
  Widget _screen;

  @override
  Widget build(BuildContext context) {
    if (_screen == null) {
      final client = http.Client();
      final lmc = LobbyManagerClient(RpcClient(client, widget.address));
      _screen = MainMenuScreen(lmc, onGameCreated: _setScreen);
    }
    return _screen;
  }

  void _setScreen(Widget screen) {
    setState(() {
      _screen = screen;
    });
  }
}
