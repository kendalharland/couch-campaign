import 'package:couchcampaign/src/api/api.pb.dart';
import 'package:couchcampaign/src/clients.dart';
import 'package:couchcampaign/src/session.dart';
import 'package:flutter/material.dart' hide Card;
import 'package:flutter/material.dart' as material;
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

class GameSessionView extends StatefulWidget {
  GameSessionView(this.client, this.session);

  final GameClient client;
  final GameSession session;

  @override
  State<StatefulWidget> createState() => GameSessionViewState();
}

class GameSessionViewState extends State<GameSessionView> {
  @override
  void initState() {
    super.initState();
    widget.session.addListener(() {
      setState(() {});
    });
  }

  @override
  void dispose() {
    super.dispose();
    widget.session.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final message = widget.session.value;
    if (message == null) {
      return Text('loading');
    }

    switch (message.whichContent()) {
      case Message_Content.sessionState:
        switch (message.sessionState) {
          case SessionState.LOBBY:
            return LobbyView(widget.client);
          case SessionState.RUNNING:
            // The next message will contain the player state.
            return Center(child: Text('running'));
          default:
            return Text('bad state ${message.sessionState}');
        }
        break;
      case Message_Content.playerState:
        return PlayerStateView(
          state: message.playerState,
          onInput: widget.session.send,
        );
      default:
        return Text('error: session has no state');
    }
  }
}

class LobbyView extends StatelessWidget {
  const LobbyView(this.client);

  final GameClient client;

  @override
  Widget build(BuildContext context) {
    return FlatButton(onPressed: _startGame, child: Text('START'));
  }

  void _startGame() {
    client.startGame();
  }
}

class PlayerStateView extends StatelessWidget {
  const PlayerStateView({@required this.state, @required this.onInput});

  final PlayerState state;
  final ValueSetter<String> onInput;

  @override
  Widget build(BuildContext context) {
    final shade = 50;
    final hue = Colors.indigo;
    final menuColor = hue[shade * 8];
    final bodyColor = hue[shade];

    final header = Container(
      decoration: BoxDecoration(color: menuColor),
      child: PlayerStats(
        health: state.health,
        wealth: state.wealth,
        stability: state.stability,
        color: bodyColor,
      ),
    );

    Widget body;
    switch (state.whichCard()) {
      case PlayerState_Card.actionCard:
        body = ActionCardBody(state.actionCard, onInput);
        break;
      case PlayerState_Card.infoCard:
        body = InfoCardBody(state.infoCard, onInput);
        break;
      case PlayerState_Card.votingCard:
        body = VotingCardBody(state.votingCard);
        onInput("got the voting card");
        break;
      default:
        body = Text('Unknown card type ${state.whichCard()}');
    }
    body = Container(
      decoration: BoxDecoration(color: bodyColor),
      child: Padding(padding: EdgeInsets.all(24), child: body),
    );

    final footer = Container(
      decoration: BoxDecoration(color: menuColor),
      child: Center(
        child: Padding(
          padding: EdgeInsets.all(16),
          child: LeaderStats(
            leader: state.leader,
            daysInOffice: 15,
            color: bodyColor,
          ),
        ),
      ),
    );

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Expanded(flex: 1, child: header),
        Expanded(flex: 4, child: body),
        Expanded(flex: 1, child: footer),
      ],
    );
  }
}

class LeaderStats extends StatelessWidget {
  const LeaderStats({
    @required this.leader,
    @required this.daysInOffice,
    this.color,
  });

  final String leader;
  final int daysInOffice;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Align(
          alignment: Alignment.centerLeft,
          child: Text(leader, style: TextStyle(fontSize: 22, color: color)),
        ),
        Align(
          alignment: Alignment.centerLeft,
          child: Text(
            '$daysInOffice days in office',
            style: TextStyle(fontSize: 36, color: color),
          ),
        ),
      ],
    );
  }
}

class PlayerStats extends StatelessWidget {
  const PlayerStats({
    @required this.health,
    @required this.wealth,
    @required this.stability,
    this.color,
  });

  final int health;
  final int wealth;
  final int stability;
  final Color color;

  @override
  Widget build(BuildContext context) {
    final style = TextStyle(
      fontSize: 24,
      fontWeight: FontWeight.bold,
      color: color,
    );

    return Center(
      child: Row(
        children: [
          Expanded(child: Center(child: Text('H($health)', style: style))),
          Expanded(child: Center(child: Text('W($wealth)', style: style))),
          Expanded(child: Center(child: Text('S($stability)', style: style))),
        ],
      ),
    );
  }
}

class ActionCardBody extends StatelessWidget {
  const ActionCardBody(this.card, this.onDismiss);

  final ActionCard card;
  final ValueSetter<String> onDismiss;

  @override
  Widget build(BuildContext context) {
    final style = TextStyle(fontSize: 18);
    final content = Center(child: Text(card.content, style: style));
    final advisor = Center(child: Text(card.advisor, style: style));
    final cardWidget = Center(child: material.Card(child: Text('picture')));
    return DismissibleCardView(
      header: content,
      card: cardWidget,
      footer: advisor,
      onDismiss: onDismiss,
    );
  }
}

class InfoCardBody extends StatelessWidget {
  const InfoCardBody(this.card, this.onDismiss);

  final InfoCard card;
  final ValueSetter<String> onDismiss;

  @override
  Widget build(BuildContext context) {
    return DismissibleCardView(
      header: Center(child: Text("")),
      card: Center(child: Text(card.text)),
      footer: Center(child: Text("")),
      onDismiss: onDismiss,
    );
  }
}

class VotingCardBody extends StatelessWidget {
  const VotingCardBody(this.card);

  final VotingCard card;

  @override
  Widget build(BuildContext context) {
    return Text('unimplemented: voting card');
  }
}

class DismissibleCardView extends StatelessWidget {
  const DismissibleCardView({
    @required this.header,
    @required this.card,
    @required this.footer,
    @required this.onDismiss,
  });

  final ValueChanged<String> onDismiss;
  final Widget header;
  final Widget card;
  final Widget footer;

  @override
  Widget build(BuildContext context) {
    final body = Dismissible(
      key: Key('${card.hashCode}'),
      direction: DismissDirection.horizontal,
      onDismissed: _onDismiss,
      child: card,
    );

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Expanded(flex: 1, child: header),
        Expanded(flex: 4, child: body),
        Expanded(flex: 1, child: footer),
      ],
    );
  }

  void _onDismiss(DismissDirection direction) {
    switch (direction) {
      case DismissDirection.endToStart:
        onDismiss("accept");
        break;
      case DismissDirection.startToEnd:
        onDismiss("reject");
        break;
      default:
        throw Exception('invalid dismiss direction: $direction');
    }
  }
}
