///
//  Generated code. Do not modify.
//  source: api.proto
//
// @dart = 2.3
// ignore_for_file: camel_case_types,non_constant_identifier_names,library_prefixes,unused_import,unused_shown_name,return_of_invalid_type

const SessionState$json = const {
  '1': 'SessionState',
  '2': const [
    const {'1': 'LOBBY', '2': 0},
    const {'1': 'RUNNING', '2': 1},
  ],
};

const Message$json = const {
  '1': 'Message',
  '2': const [
    const {'1': 'player_state', '3': 1, '4': 1, '5': 11, '6': '.couchcampaign.PlayerState', '9': 0, '10': 'playerState'},
    const {'1': 'session_state', '3': 2, '4': 1, '5': 14, '6': '.couchcampaign.SessionState', '9': 0, '10': 'sessionState'},
  ],
  '8': const [
    const {'1': 'content'},
  ],
};

const GameInfo$json = const {
  '1': 'GameInfo',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
  ],
};

const PlayerState$json = const {
  '1': 'PlayerState',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    const {'1': 'leader', '3': 2, '4': 1, '5': 9, '10': 'leader'},
    const {'1': 'leader_time_in_office', '3': 3, '4': 1, '5': 5, '10': 'leaderTimeInOffice'},
    const {'1': 'health', '3': 4, '4': 1, '5': 5, '10': 'health'},
    const {'1': 'wealth', '3': 5, '4': 1, '5': 5, '10': 'wealth'},
    const {'1': 'stability', '3': 6, '4': 1, '5': 5, '10': 'stability'},
    const {'1': 'card', '3': 7, '4': 1, '5': 11, '6': '.couchcampaign.Card', '10': 'card'},
  ],
};

const Card$json = const {
  '1': 'Card',
  '2': const [
    const {'1': 'text', '3': 1, '4': 1, '5': 9, '10': 'text'},
    const {'1': 'speaker', '3': 2, '4': 1, '5': 9, '10': 'speaker'},
    const {'1': 'accept_text', '3': 3, '4': 1, '5': 9, '10': 'acceptText'},
    const {'1': 'reject_text', '3': 4, '4': 1, '5': 9, '10': 'rejectText'},
  ],
};

const CreateGameRequest$json = const {
  '1': 'CreateGameRequest',
};

const CreateGameResponse$json = const {
  '1': 'CreateGameResponse',
  '2': const [
    const {'1': 'game_url', '3': 1, '4': 1, '5': 9, '10': 'gameUrl'},
  ],
};

const JoinGameRequest$json = const {
  '1': 'JoinGameRequest',
  '2': const [
    const {'1': 'game_id', '3': 1, '4': 1, '5': 9, '10': 'gameId'},
  ],
};

const JoinGameResponse$json = const {
  '1': 'JoinGameResponse',
  '2': const [
    const {'1': 'game_url', '3': 1, '4': 1, '5': 9, '10': 'gameUrl'},
  ],
};

const ListGamesRequest$json = const {
  '1': 'ListGamesRequest',
};

const ListGamesResponse$json = const {
  '1': 'ListGamesResponse',
  '2': const [
    const {'1': 'games', '3': 1, '4': 3, '5': 11, '6': '.couchcampaign.GameInfo', '10': 'games'},
  ],
};

