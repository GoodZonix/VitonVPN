import 'package:flutter/foundation.dart';
import '../core/api_client.dart';

enum VPNStatus { disconnected, connecting, connected }

class VPNProvider with ChangeNotifier {
  VPNStatus _status = VPNStatus.disconnected;
  String? _selectedServerName;
  int? _pingMs;
  int _bytesReceived = 0;
  int _bytesSent = 0;

  VPNStatus get status => _status;
  String? get selectedServerName => _selectedServerName;
  int? get pingMs => _pingMs;
  int get bytesReceived => _bytesReceived;
  int get bytesSent => _bytesSent;

  void setSelectedServer(String? name) {
    _selectedServerName = name;
    notifyListeners();
  }

  Future<void> fetchConfig(ApiClient api) async {
    try {
      await api.get('/api/config');
      notifyListeners();
    } catch (_) {
      notifyListeners();
    }
  }

  Future<void> connect() async {
    _status = VPNStatus.connecting;
    notifyListeners();
    // TODO: integrate with native VPN (Xray/V2Ray core) using config from API
    await Future.delayed(const Duration(seconds: 2));
    _status = VPNStatus.connected;
    _pingMs = 42;
    _bytesReceived = 270960000; // demo: ~270.96 MB
    _bytesSent = 43520000;     // demo: ~43.52 MB
    notifyListeners();
  }

  Future<void> disconnect() async {
    _status = VPNStatus.disconnected;
    _pingMs = null;
    _bytesReceived = 0;
    _bytesSent = 0;
    notifyListeners();
  }

  String get receivedMb => (_bytesReceived / (1024 * 1024)).toStringAsFixed(2);
  String get sentMb => (_bytesSent / (1024 * 1024)).toStringAsFixed(2);

  void setPing(int? ms) {
    _pingMs = ms;
    notifyListeners();
  }
}
