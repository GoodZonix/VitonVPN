import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import '../providers/auth_provider.dart';
import '../providers/vpn_provider.dart';
import '../providers/servers_provider.dart';
import 'servers_screen.dart';
import 'subscription_screen.dart';
import 'settings_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  static const String appVersion = 'V 1.2.1';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.vitonBlack,
      appBar: AppBar(
        title: const Text('Viton VPN'),
        backgroundColor: AppTheme.vitonBlack,
        foregroundColor: Colors.white,
        elevation: 0,
        leading: const Padding(
          padding: EdgeInsets.only(left: 8),
          child: Icon(Icons.shield_outlined, color: Colors.white),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.account_balance_wallet_outlined),
            tooltip: 'Кошелёк',
            onPressed: () => Navigator.push(
              context,
              MaterialPageRoute(builder: (_) => const SubscriptionScreen()),
            ),
          ),
          IconButton(
            icon: const Icon(Icons.settings_outlined),
            tooltip: 'Настройки',
            onPressed: () => Navigator.push(
              context,
              MaterialPageRoute(builder: (_) => const SettingsScreen()),
            ),
          ),
        ],
      ),
      body: Consumer3<VPNProvider, AuthProvider, ServersProvider>(
        builder: (_, vpn, auth, servers, __) {
          return RefreshIndicator(
            onRefresh: () async {
              await vpn.fetchConfig(auth.api);
              await servers.load(auth.api);
            },
            color: AppTheme.vitonGreen,
            backgroundColor: Colors.grey.shade800,
            child: SingleChildScrollView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: Column(
                children: [
                  const SizedBox(height: 32),
                  _InstructionText(connected: vpn.status == VPNStatus.connected),
                  const SizedBox(height: 24),
                  _ConnectButton(
                    status: vpn.status,
                    onConnect: () => vpn.connect(),
                    onDisconnect: () => vpn.disconnect(),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    vpn.status == VPNStatus.connected ? 'VPN включен' : 'VPN выключен',
                    style: TextStyle(
                      color: vpn.status == VPNStatus.connected ? AppTheme.vitonGreen : Colors.white70,
                      fontSize: 16,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                  const SizedBox(height: 32),
                  _DataStats(receivedMb: vpn.receivedMb, sentMb: vpn.sentMb),
                  const SizedBox(height: 24),
                  _UserId(id: _shortId(auth.userId)),
                  const SizedBox(height: 16),
                  TextButton.icon(
                    onPressed: () => Navigator.push(
                      context,
                      MaterialPageRoute(builder: (_) => const ServersScreen()),
                    ).then((_) => servers.load(auth.api)),
                    icon: const Icon(Icons.public, size: 18, color: Colors.white70),
                    label: Text(
                      vpn.selectedServerName ?? 'Авто (лучший сервер)',
                      style: const TextStyle(color: Colors.white70),
                    ),
                  ),
                  const SizedBox(height: 48),
                  Align(
                    alignment: Alignment.centerRight,
                    child: Text(
                      appVersion,
                      style: const TextStyle(color: Colors.white38, fontSize: 12),
                    ),
                  ),
                  const SizedBox(height: 24),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  static String _shortId(String? userId) {
    if (userId == null || userId.isEmpty) return 'V00000000';
    final clean = userId.replaceAll('-', '').toUpperCase();
    if (clean.length >= 8) return 'V${clean.substring(0, 8)}';
    return 'V${clean.padLeft(8, '0')}';
  }
}

class _InstructionText extends StatelessWidget {
  const _InstructionText({required this.connected});

  final bool connected;

  @override
  Widget build(BuildContext context) {
    return Text(
      connected
          ? 'Нажмите на кнопку, чтобы отключить VPN'
          : 'Нажмите на кнопку, чтобы подключить VPN',
      textAlign: TextAlign.center,
      style: const TextStyle(color: Colors.white, fontSize: 16),
    );
  }
}

class _ConnectButton extends StatelessWidget {
  const _ConnectButton({
    required this.status,
    required this.onConnect,
    required this.onDisconnect,
  });

  final VPNStatus status;
  final VoidCallback onConnect;
  final VoidCallback onDisconnect;

  @override
  Widget build(BuildContext context) {
    final connected = status == VPNStatus.connected;
    final connecting = status == VPNStatus.connecting;
    return SizedBox(
      width: 200,
      height: 200,
      child: Stack(
        alignment: Alignment.center,
        children: [
          if (connecting)
            const SizedBox(
              width: 200,
              height: 200,
              child: CircularProgressIndicator(
                strokeWidth: 3,
                color: AppTheme.vitonGreen,
              ),
            ),
          Material(
            color: AppTheme.vitonGreen,
            borderRadius: BorderRadius.circular(100),
            elevation: 8,
            child: InkWell(
              onTap: connecting ? null : (connected ? onDisconnect : onConnect),
              borderRadius: BorderRadius.circular(100),
              child: Container(
                width: 200,
                height: 200,
                alignment: Alignment.center,
                child: Icon(
                  Icons.power_settings_new,
                  size: 80,
                  color: Colors.white,
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _DataStats extends StatelessWidget {
  const _DataStats({required this.receivedMb, required this.sentMb});

  final String receivedMb;
  final String sentMb;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        Column(
          children: [
            const Text('Получено', style: TextStyle(color: Colors.white70, fontSize: 14)),
            const SizedBox(height: 4),
            Text('$receivedMb MB', style: const TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.w600)),
          ],
        ),
        Column(
          children: [
            const Text('Отправлено', style: TextStyle(color: Colors.white70, fontSize: 14)),
            const SizedBox(height: 4),
            Text('$sentMb MB', style: const TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.w600)),
          ],
        ),
      ],
    );
  }
}

class _UserId extends StatelessWidget {
  const _UserId({required this.id});

  final String id;

  @override
  Widget build(BuildContext context) {
    return Text(
      'Ваш ID: $id',
      style: const TextStyle(color: Colors.white70, fontSize: 14),
    );
  }
}
