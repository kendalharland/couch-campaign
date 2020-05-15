import 'package:couchcampaign/src/clients.dart';
import 'package:couchcampaign/src/screens.dart';
import 'package:flutter/material.dart' hide Card;
import 'package:google_fonts/google_fonts.dart';
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
        textTheme: GoogleFonts.patrickHandScTextTheme(),
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
      final rpcClient = RpcClient(http.Client(), widget.address);
      final lmc = LobbyManagerClient(rpcClient);
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
