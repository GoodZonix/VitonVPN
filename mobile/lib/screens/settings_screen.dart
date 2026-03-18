import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../config/theme.dart';
import 'manual_setup_screen.dart';
import '../providers/auth_provider.dart';

class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppTheme.vitonBlack,
      appBar: AppBar(
        title: const Text('Настройки'),
        backgroundColor: AppTheme.vitonBlack,
        foregroundColor: Colors.white,
        elevation: 0,
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView(
              padding: const EdgeInsets.symmetric(vertical: 8),
              children: [
                _MenuItem(
                  icon: Icons.tune,
                  title: 'Ручная настройка',
                  onTap: () => _onManualSetup(context),
                ),
                _MenuItem(
                  icon: Icons.restart_alt,
                  title: 'Сброс конфигурации',
                  onTap: () => _onResetConfig(context),
                ),
                _MenuItem(
                  icon: Icons.support_agent,
                  title: 'Написать в поддержку',
                  onTap: () => _onContactSupport(context),
                ),
                _MenuItem(
                  icon: Icons.info_outline,
                  title: 'О приложении',
                  onTap: () => _onAbout(context),
                ),
                const Divider(color: Colors.white24, height: 32),
                ListTile(
                  leading: const Icon(Icons.logout, color: Colors.white70),
                  title: const Text('Выйти', style: TextStyle(color: Colors.white70)),
                  onTap: () async {
                    await context.read<AuthProvider>().logout();
                    if (context.mounted) Navigator.of(context).popUntil((r) => r.isFirst);
                  },
                ),
              ],
            ),
          ),
          _BottomHint(),
        ],
      ),
    );
  }

  void _onManualSetup(BuildContext context) {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (_) => const ManualSetupScreen()),
    );
  }

  void _onResetConfig(BuildContext context) {
    showDialog<void>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: const Color(0xFF1E293B),
        title: const Text('Сброс конфигурации', style: TextStyle(color: Colors.white)),
        content: const Text(
          'Сбросить сохранённую конфигурацию VPN?',
          style: TextStyle(color: Colors.white70),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Отмена'),
          ),
          TextButton(
            onPressed: () {
              Navigator.pop(ctx);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('Конфигурация сброшена')),
              );
            },
            child: const Text('Сбросить', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }

  void _onContactSupport(BuildContext context) {
    // TODO: открыть почту или чат поддержки
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Написать в поддержку — укажите email в настройках')),
    );
  }

  void _onAbout(BuildContext context) {
    showDialog<void>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: const Color(0xFF1E293B),
        title: const Text('О приложении', style: TextStyle(color: Colors.white)),
        content: const Text(
          'Viton VPN\nВерсия 1.2.1\n\nБезопасный VPN на базе VLESS (Xray).',
          style: TextStyle(color: Colors.white70),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }
}

class _MenuItem extends StatelessWidget {
  const _MenuItem({
    required this.icon,
    required this.title,
    required this.onTap,
  });

  final IconData icon;
  final String title;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Icon(icon, color: Colors.white70),
      title: Text(title, style: const TextStyle(color: Colors.white)),
      trailing: const Icon(Icons.chevron_right, color: Colors.white38),
      onTap: onTap,
    );
  }
}

class _BottomHint extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(20),
      margin: const EdgeInsets.only(top: 8),
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.06),
        border: Border(top: BorderSide(color: Colors.white12)),
      ),
      child: Text(
        'Включайте и отключайте VPN только через приложение Viton VPN. '
        'Не используйте для этого раздел «Настройки» на iPhone.',
        style: TextStyle(
          color: Colors.white.withValues(alpha: 0.7),
          fontSize: 12,
          height: 1.35,
        ),
        textAlign: TextAlign.center,
      ),
    );
  }
}
