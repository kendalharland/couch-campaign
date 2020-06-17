///
//  Generated code. Do not modify.
//  source: api.proto
//
// @dart = 2.3
// ignore_for_file: camel_case_types,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name,return_of_invalid_type

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'api.pbenum.dart';

export 'api.pbenum.dart';

enum Message_Content {
  playerState, 
  sessionState, 
  notSet
}

class Message extends $pb.GeneratedMessage {
  static const $core.Map<$core.int, Message_Content> _Message_ContentByTag = {
    1 : Message_Content.playerState,
    2 : Message_Content.sessionState,
    0 : Message_Content.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('Message', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..oo(0, [1, 2])
    ..aOM<PlayerState>(1, 'playerState', subBuilder: PlayerState.create)
    ..e<SessionState>(2, 'sessionState', $pb.PbFieldType.OE, defaultOrMaker: SessionState.LOBBY, valueOf: SessionState.valueOf, enumValues: SessionState.values)
    ..hasRequiredFields = false
  ;

  Message._() : super();
  factory Message() => create();
  factory Message.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Message.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  Message clone() => Message()..mergeFromMessage(this);
  Message copyWith(void Function(Message) updates) => super.copyWith((message) => updates(message as Message));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static Message create() => Message._();
  Message createEmptyInstance() => create();
  static $pb.PbList<Message> createRepeated() => $pb.PbList<Message>();
  @$core.pragma('dart2js:noInline')
  static Message getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Message>(create);
  static Message _defaultInstance;

  Message_Content whichContent() => _Message_ContentByTag[$_whichOneof(0)];
  void clearContent() => clearField($_whichOneof(0));

  @$pb.TagNumber(1)
  PlayerState get playerState => $_getN(0);
  @$pb.TagNumber(1)
  set playerState(PlayerState v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasPlayerState() => $_has(0);
  @$pb.TagNumber(1)
  void clearPlayerState() => clearField(1);
  @$pb.TagNumber(1)
  PlayerState ensurePlayerState() => $_ensure(0);

  @$pb.TagNumber(2)
  SessionState get sessionState => $_getN(1);
  @$pb.TagNumber(2)
  set sessionState(SessionState v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasSessionState() => $_has(1);
  @$pb.TagNumber(2)
  void clearSessionState() => clearField(2);
}

class GameInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('GameInfo', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..aOS(1, 'id')
    ..hasRequiredFields = false
  ;

  GameInfo._() : super();
  factory GameInfo() => create();
  factory GameInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GameInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  GameInfo clone() => GameInfo()..mergeFromMessage(this);
  GameInfo copyWith(void Function(GameInfo) updates) => super.copyWith((message) => updates(message as GameInfo));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GameInfo create() => GameInfo._();
  GameInfo createEmptyInstance() => create();
  static $pb.PbList<GameInfo> createRepeated() => $pb.PbList<GameInfo>();
  @$core.pragma('dart2js:noInline')
  static GameInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GameInfo>(create);
  static GameInfo _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);
}

class PlayerState extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('PlayerState', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..aOS(1, 'id')
    ..aOS(2, 'leader')
    ..a<$core.int>(3, 'leaderTimeInOffice', $pb.PbFieldType.O3)
    ..a<$core.int>(4, 'health', $pb.PbFieldType.O3)
    ..a<$core.int>(5, 'wealth', $pb.PbFieldType.O3)
    ..a<$core.int>(6, 'stability', $pb.PbFieldType.O3)
    ..aOM<Card>(7, 'card', subBuilder: Card.create)
    ..hasRequiredFields = false
  ;

  PlayerState._() : super();
  factory PlayerState() => create();
  factory PlayerState.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PlayerState.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  PlayerState clone() => PlayerState()..mergeFromMessage(this);
  PlayerState copyWith(void Function(PlayerState) updates) => super.copyWith((message) => updates(message as PlayerState));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static PlayerState create() => PlayerState._();
  PlayerState createEmptyInstance() => create();
  static $pb.PbList<PlayerState> createRepeated() => $pb.PbList<PlayerState>();
  @$core.pragma('dart2js:noInline')
  static PlayerState getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PlayerState>(create);
  static PlayerState _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get leader => $_getSZ(1);
  @$pb.TagNumber(2)
  set leader($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasLeader() => $_has(1);
  @$pb.TagNumber(2)
  void clearLeader() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get leaderTimeInOffice => $_getIZ(2);
  @$pb.TagNumber(3)
  set leaderTimeInOffice($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLeaderTimeInOffice() => $_has(2);
  @$pb.TagNumber(3)
  void clearLeaderTimeInOffice() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get health => $_getIZ(3);
  @$pb.TagNumber(4)
  set health($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasHealth() => $_has(3);
  @$pb.TagNumber(4)
  void clearHealth() => clearField(4);

  @$pb.TagNumber(5)
  $core.int get wealth => $_getIZ(4);
  @$pb.TagNumber(5)
  set wealth($core.int v) { $_setSignedInt32(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasWealth() => $_has(4);
  @$pb.TagNumber(5)
  void clearWealth() => clearField(5);

  @$pb.TagNumber(6)
  $core.int get stability => $_getIZ(5);
  @$pb.TagNumber(6)
  set stability($core.int v) { $_setSignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasStability() => $_has(5);
  @$pb.TagNumber(6)
  void clearStability() => clearField(6);

  @$pb.TagNumber(7)
  Card get card => $_getN(6);
  @$pb.TagNumber(7)
  set card(Card v) { setField(7, v); }
  @$pb.TagNumber(7)
  $core.bool hasCard() => $_has(6);
  @$pb.TagNumber(7)
  void clearCard() => clearField(7);
  @$pb.TagNumber(7)
  Card ensureCard() => $_ensure(6);
}

class Card extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('Card', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..aOS(1, 'text')
    ..aOS(2, 'speaker')
    ..aOS(3, 'acceptText')
    ..aOS(4, 'rejectText')
    ..hasRequiredFields = false
  ;

  Card._() : super();
  factory Card() => create();
  factory Card.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory Card.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  Card clone() => Card()..mergeFromMessage(this);
  Card copyWith(void Function(Card) updates) => super.copyWith((message) => updates(message as Card));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static Card create() => Card._();
  Card createEmptyInstance() => create();
  static $pb.PbList<Card> createRepeated() => $pb.PbList<Card>();
  @$core.pragma('dart2js:noInline')
  static Card getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Card>(create);
  static Card _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get text => $_getSZ(0);
  @$pb.TagNumber(1)
  set text($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasText() => $_has(0);
  @$pb.TagNumber(1)
  void clearText() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get speaker => $_getSZ(1);
  @$pb.TagNumber(2)
  set speaker($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSpeaker() => $_has(1);
  @$pb.TagNumber(2)
  void clearSpeaker() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get acceptText => $_getSZ(2);
  @$pb.TagNumber(3)
  set acceptText($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAcceptText() => $_has(2);
  @$pb.TagNumber(3)
  void clearAcceptText() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get rejectText => $_getSZ(3);
  @$pb.TagNumber(4)
  set rejectText($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasRejectText() => $_has(3);
  @$pb.TagNumber(4)
  void clearRejectText() => clearField(4);
}

class CreateGameRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('CreateGameRequest', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  CreateGameRequest._() : super();
  factory CreateGameRequest() => create();
  factory CreateGameRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateGameRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  CreateGameRequest clone() => CreateGameRequest()..mergeFromMessage(this);
  CreateGameRequest copyWith(void Function(CreateGameRequest) updates) => super.copyWith((message) => updates(message as CreateGameRequest));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CreateGameRequest create() => CreateGameRequest._();
  CreateGameRequest createEmptyInstance() => create();
  static $pb.PbList<CreateGameRequest> createRepeated() => $pb.PbList<CreateGameRequest>();
  @$core.pragma('dart2js:noInline')
  static CreateGameRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateGameRequest>(create);
  static CreateGameRequest _defaultInstance;
}

class CreateGameResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('CreateGameResponse', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..aOS(1, 'gameUrl')
    ..hasRequiredFields = false
  ;

  CreateGameResponse._() : super();
  factory CreateGameResponse() => create();
  factory CreateGameResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateGameResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  CreateGameResponse clone() => CreateGameResponse()..mergeFromMessage(this);
  CreateGameResponse copyWith(void Function(CreateGameResponse) updates) => super.copyWith((message) => updates(message as CreateGameResponse));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CreateGameResponse create() => CreateGameResponse._();
  CreateGameResponse createEmptyInstance() => create();
  static $pb.PbList<CreateGameResponse> createRepeated() => $pb.PbList<CreateGameResponse>();
  @$core.pragma('dart2js:noInline')
  static CreateGameResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateGameResponse>(create);
  static CreateGameResponse _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get gameUrl => $_getSZ(0);
  @$pb.TagNumber(1)
  set gameUrl($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasGameUrl() => $_has(0);
  @$pb.TagNumber(1)
  void clearGameUrl() => clearField(1);
}

class JoinGameRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('JoinGameRequest', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..aOS(1, 'gameId')
    ..hasRequiredFields = false
  ;

  JoinGameRequest._() : super();
  factory JoinGameRequest() => create();
  factory JoinGameRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory JoinGameRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  JoinGameRequest clone() => JoinGameRequest()..mergeFromMessage(this);
  JoinGameRequest copyWith(void Function(JoinGameRequest) updates) => super.copyWith((message) => updates(message as JoinGameRequest));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static JoinGameRequest create() => JoinGameRequest._();
  JoinGameRequest createEmptyInstance() => create();
  static $pb.PbList<JoinGameRequest> createRepeated() => $pb.PbList<JoinGameRequest>();
  @$core.pragma('dart2js:noInline')
  static JoinGameRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<JoinGameRequest>(create);
  static JoinGameRequest _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get gameId => $_getSZ(0);
  @$pb.TagNumber(1)
  set gameId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasGameId() => $_has(0);
  @$pb.TagNumber(1)
  void clearGameId() => clearField(1);
}

class JoinGameResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('JoinGameResponse', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..aOS(1, 'gameUrl')
    ..hasRequiredFields = false
  ;

  JoinGameResponse._() : super();
  factory JoinGameResponse() => create();
  factory JoinGameResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory JoinGameResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  JoinGameResponse clone() => JoinGameResponse()..mergeFromMessage(this);
  JoinGameResponse copyWith(void Function(JoinGameResponse) updates) => super.copyWith((message) => updates(message as JoinGameResponse));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static JoinGameResponse create() => JoinGameResponse._();
  JoinGameResponse createEmptyInstance() => create();
  static $pb.PbList<JoinGameResponse> createRepeated() => $pb.PbList<JoinGameResponse>();
  @$core.pragma('dart2js:noInline')
  static JoinGameResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<JoinGameResponse>(create);
  static JoinGameResponse _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get gameUrl => $_getSZ(0);
  @$pb.TagNumber(1)
  set gameUrl($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasGameUrl() => $_has(0);
  @$pb.TagNumber(1)
  void clearGameUrl() => clearField(1);
}

class ListGamesRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('ListGamesRequest', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  ListGamesRequest._() : super();
  factory ListGamesRequest() => create();
  factory ListGamesRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListGamesRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  ListGamesRequest clone() => ListGamesRequest()..mergeFromMessage(this);
  ListGamesRequest copyWith(void Function(ListGamesRequest) updates) => super.copyWith((message) => updates(message as ListGamesRequest));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ListGamesRequest create() => ListGamesRequest._();
  ListGamesRequest createEmptyInstance() => create();
  static $pb.PbList<ListGamesRequest> createRepeated() => $pb.PbList<ListGamesRequest>();
  @$core.pragma('dart2js:noInline')
  static ListGamesRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListGamesRequest>(create);
  static ListGamesRequest _defaultInstance;
}

class ListGamesResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo('ListGamesResponse', package: const $pb.PackageName('couchcampaign'), createEmptyInstance: create)
    ..pc<GameInfo>(1, 'games', $pb.PbFieldType.PM, subBuilder: GameInfo.create)
    ..hasRequiredFields = false
  ;

  ListGamesResponse._() : super();
  factory ListGamesResponse() => create();
  factory ListGamesResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListGamesResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  ListGamesResponse clone() => ListGamesResponse()..mergeFromMessage(this);
  ListGamesResponse copyWith(void Function(ListGamesResponse) updates) => super.copyWith((message) => updates(message as ListGamesResponse));
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ListGamesResponse create() => ListGamesResponse._();
  ListGamesResponse createEmptyInstance() => create();
  static $pb.PbList<ListGamesResponse> createRepeated() => $pb.PbList<ListGamesResponse>();
  @$core.pragma('dart2js:noInline')
  static ListGamesResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListGamesResponse>(create);
  static ListGamesResponse _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<GameInfo> get games => $_getList(0);
}

