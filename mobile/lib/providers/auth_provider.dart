import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../core/api_client.dart';

class AuthProvider with ChangeNotifier {
  AuthProvider() {
    _loadStored();
  }

  final ApiClient _api = ApiClient();
  String? _token;
  String? _userId;
  String? _email;
  bool _isLoggedIn = false;

  bool get isLoggedIn => _isLoggedIn;
  String? get token => _token;
  String? get email => _email;
  String? get userId => _userId;

  Future<void> _loadStored() async {
    final prefs = await SharedPreferences.getInstance();
    _token = prefs.getString('token');
    _email = prefs.getString('email');
    _userId = prefs.getString('userId');
    _isLoggedIn = _token != null && _token!.isNotEmpty;
    _api.token = _token;
    notifyListeners();
  }

  Future<void> login(String email, String password, String deviceId, String deviceName) async {
    final res = await _api.post('/api/login', {
      'email': email,
      'password': password,
      'device_id': deviceId,
      'device_name': deviceName,
    });
    _token = res['token'] as String?;
    _email = (res['user'] as Map?)?['email'] as String?;
    _userId = (res['user'] as Map?)?['id'] as String?;
    _api.token = _token;
    _isLoggedIn = true;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('token', _token!);
    await prefs.setString('email', _email ?? '');
    if (_userId != null) await prefs.setString('userId', _userId!);
    notifyListeners();
  }

  Future<void> register(String email, String password, String deviceId, String deviceName) async {
    final res = await _api.post('/api/register', {
      'email': email,
      'password': password,
      'device_id': deviceId,
      'device_name': deviceName,
    });
    _token = res['token'] as String?;
    _email = (res['user'] as Map?)?['email'] as String?;
    _userId = (res['user'] as Map?)?['id'] as String?;
    _api.token = _token;
    _isLoggedIn = true;
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('token', _token!);
    await prefs.setString('email', _email ?? '');
    if (_userId != null) await prefs.setString('userId', _userId!);
    notifyListeners();
  }

  Future<void> logout() async {
    _token = null;
    _userId = null;
    _email = null;
    _isLoggedIn = false;
    _api.token = null;
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('token');
    await prefs.remove('email');
    await prefs.remove('userId');
    notifyListeners();
  }

  ApiClient get api => _api;
}
