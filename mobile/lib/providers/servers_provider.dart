import 'package:flutter/foundation.dart';
import '../core/api_client.dart';

class ServerItem {
  ServerItem({
    required this.id,
    required this.name,
    required this.region,
    required this.host,
    required this.port,
    required this.type,
    this.pingMs,
  });
  final String id;
  final String name;
  final String region;
  final String host;
  final int port;
  final String type;
  int? pingMs;
}

class ServersProvider with ChangeNotifier {
  List<ServerItem> _servers = [];
  bool _loading = false;
  String? _error;

  List<ServerItem> get servers => _servers;
  bool get loading => _loading;
  String? get error => _error;

  Future<void> load(ApiClient api) async {
    _loading = true;
    _error = null;
    notifyListeners();
    try {
      final res = await api.get('/api/servers');
      final list = (res['servers'] as List?) ?? [];
      _servers = list.map((e) {
        final m = e as Map<String, dynamic>;
        return ServerItem(
          id: m['id'] as String? ?? '',
          name: m['name'] as String? ?? '',
          region: m['region'] as String? ?? '',
          host: m['host'] as String? ?? '',
          port: m['port'] as int? ?? 443,
          type: m['type'] as String? ?? 'reality',
        );
      }).toList();
    } catch (e) {
      _error = e.toString();
    }
    _loading = false;
    notifyListeners();
  }
}
