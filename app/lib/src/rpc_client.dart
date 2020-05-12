import 'package:xhttp/xhttp.dart' as http;

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
