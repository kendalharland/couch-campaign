import 'package:couchcampaign/src/api/api.pb.dart';
import 'package:couchcampaign/src/game_client.dart';
import 'package:couchcampaign/src/session.dart';
import 'package:flutter/material.dart' hide Card;
import 'package:flutter/material.dart' as material;

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
            return Center(child: Text('running'));
          default:
            return Text('bad state ${message.sessionState}');
        }
        break;
      case Message_Content.gameState:
        return GameStateView(
          state: message.gameState,
          onInput: widget.session.send,
        );
      default:
        return Text('session has no state');
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

class GameStateView extends StatelessWidget {
  const GameStateView({@required this.state, @required this.onInput});

  final GameState state;
  final ValueSetter<String> onInput;

  @override
  Widget build(BuildContext context) {
    final header = PlayerStats(
      health: state.health,
      wealth: state.wealth,
      stability: state.stability,
    );

    final footer = Center(child: Text(state.leader));

    Widget body;
    switch (state.whichCard()) {
      case GameState_Card.actionCard:
        body = ActionCardBody(state.actionCard);
        break;
      case GameState_Card.infoCard:
        body = InfoCardBody(state.infoCard);
        break;
      case GameState_Card.votingCard:
        body = VotingCardBody(state.votingCard);
        onInput("got the voting card");
        break;
      default:
        body = Text('Unknown card type ${state.whichCard()}');
    }

    return DismissibleCardView(
      header: header,
      body: body,
      footer: footer,
      onDismiss: onInput,
    );
  }
}

class PlayerStats extends StatelessWidget {
  const PlayerStats({
    @required this.health,
    @required this.wealth,
    @required this.stability,
  });

  final int health;
  final int wealth;
  final int stability;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Row(
        children: [
          Expanded(child: Center(child: Text('H($health)'))),
          Expanded(child: Center(child: Text('W($wealth)'))),
          Expanded(child: Center(child: Text('S($stability)'))),
        ],
      ),
    );
  }
}

class ActionCardBody extends StatelessWidget {
  const ActionCardBody(this.card);

  final ActionCard card;

  @override
  Widget build(BuildContext context) {
    final content = material.Card(child: Center(child: Text(card.content)));
    final advisor = Text(card.advisor);
    return Column(children: [
      Expanded(flex: 3, child: content),
      Expanded(child: advisor),
    ]);
  }
}

class InfoCardBody extends StatelessWidget {
  const InfoCardBody(this.card);

  final InfoCard card;

  @override
  Widget build(BuildContext context) {
    return Text(card.text);
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
    @required this.body,
    @required this.footer,
    @required this.onDismiss,
  });

  final ValueChanged<String> onDismiss;
  final Widget header;
  final Widget body;
  final Widget footer;

  @override
  Widget build(BuildContext context) {
    final card = Dismissible(
      key: Key('${body.hashCode}'),
      direction: DismissDirection.horizontal,
      onDismissed: _onDismiss,
      child: body,
    );

    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Expanded(flex: 1, child: header),
        Expanded(flex: 4, child: card),
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
