import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiClient {
  ApiClient({this.baseUrl = 'http://localhost:8080', this.token});

  final String baseUrl;
  String? token;

  Map<String, String> get _headers => {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      };

  Future<Map<String, dynamic>> post(String path, Map<String, dynamic> body) async {
    final r = await http.post(
      Uri.parse('$baseUrl$path'),
      headers: _headers,
      body: jsonEncode(body),
    );
    return _decode(r);
  }

  Future<Map<String, dynamic>> get(String path) async {
    final r = await http.get(Uri.parse('$baseUrl$path'), headers: _headers);
    return _decode(r);
  }

  static Map<String, dynamic> _decode(http.Response r) {
    final decoded = jsonDecode(r.body.isEmpty ? '{}' : r.body) as Map<String, dynamic>?;
    if (r.statusCode >= 400) {
      throw ApiException(r.statusCode, decoded?['error'] as String? ?? r.body);
    }
    return decoded ?? {};
  }
}

class ApiException implements Exception {
  ApiException(this.statusCode, this.message);
  final int statusCode;
  final String? message;
  @override
  String toString() => 'ApiException($statusCode, $message)';
}
