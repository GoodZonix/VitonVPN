import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter/services.dart';
import '../providers/auth_provider.dart';
import '../providers/vpn_provider.dart';

class SubscriptionScreen extends StatelessWidget {
  const SubscriptionScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Wallet')),
      body: Consumer2<AuthProvider, VPNProvider>(
        builder: (_, auth, vpn, __) {
          return FutureBuilder<Map<String, dynamic>>(
            future: auth.api.get('/api/wallet'),
            builder: (context, snap) {
              if (!snap.hasData) {
                return const Center(child: CircularProgressIndicator());
              }
              final data = snap.data!;
              final balance = (data['balance'] as num?)?.toDouble() ?? 0.0;
              final approxDays = (data['approx_days_left'] as num?)?.toDouble() ?? 0.0;
              return SingleChildScrollView(
                padding: const EdgeInsets.all(24),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(20),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Balance: ${balance.toStringAsFixed(2)} ₽',
                              style: Theme.of(context).textTheme.titleLarge,
                            ),
                            const SizedBox(height: 8),
                            Text('Approx. days left: ${approxDays.toStringAsFixed(1)}'),
                            const SizedBox(height: 8),
                            const Text(
                              'Rate: 100 ₽ ≈ 7 days of VPN',
                              style: TextStyle(fontSize: 12),
                            ),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 24),
                    const Text('Top up wallet', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                    const SizedBox(height: 12),
                    _AmountCard(amount: 100, subtitle: '≈ 1 week', onTap: () => _topup(context, 100)),
                    _AmountCard(amount: 300, subtitle: '≈ 3 weeks', onTap: () => _topup(context, 300)),
                    _AmountCard(amount: 500, subtitle: '≈ 5 weeks', onTap: () => _topup(context, 500)),
                    const SizedBox(height: 24),
                    const Text('Telegram', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                    const SizedBox(height: 12),
                    OutlinedButton.icon(
                      onPressed: () => _linkTelegram(context),
                      icon: const Icon(Icons.telegram),
                      label: const Text('Привязать Telegram-бота'),
                    ),
                  ],
                ),
              );
            },
          );
        },
      ),
    );
  }

  void _topup(BuildContext context, double amount) async {
    final auth = context.read<AuthProvider>();
    try {
      await auth.api.post('/api/wallet/topup', {'amount': amount});
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Wallet topped up by ${amount.toStringAsFixed(0)} ₽ (demo, YooMoney integration pending)')),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Topup failed: $e')),
      );
    }
  }

  void _linkTelegram(BuildContext context) async {
    final auth = context.read<AuthProvider>();
    try {
      final res = await auth.api.post('/api/telegram/link', {});
      final code = (res['code'] as String?) ?? '';
      final link = (res['deep_link'] as String?) ?? '';
      if (code.isEmpty) {
        throw Exception('no code');
      }
      await showDialog<void>(
        context: context,
        builder: (ctx) => AlertDialog(
          title: const Text('Привязка Telegram'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('Откройте бота @YViton_bot и отправьте команду:'),
              const SizedBox(height: 8),
              SelectableText('/start $code'),
              const SizedBox(height: 12),
              if (link.isNotEmpty) ...[
                const Text('Или перейдите по ссылке:'),
                const SizedBox(height: 8),
                SelectableText(link),
              ],
            ],
          ),
          actions: [
            TextButton(
              onPressed: () async {
                await Clipboard.setData(ClipboardData(text: '/start $code'));
                if (ctx.mounted) Navigator.pop(ctx);
                ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Скопировано'))); 
              },
              child: const Text('Скопировать'),
            ),
            TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Закрыть')),
          ],
        ),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Не удалось получить код: $e')),
      );
    }
  }
}

class _AmountCard extends StatelessWidget {
  const _AmountCard({
    required this.amount,
    this.subtitle,
    required this.onTap,
  });

  final double amount;
  final String? subtitle;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: ListTile(
        title: Text('${amount.toStringAsFixed(0)} ₽'),
        subtitle: subtitle != null ? Text(subtitle!) : null,
        trailing: const Icon(Icons.add_circle_outline),
        onTap: onTap,
      ),
    );
  }
}
