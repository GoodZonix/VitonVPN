import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import '../providers/servers_provider.dart';
import '../providers/vpn_provider.dart';

class ServersScreen extends StatelessWidget {
  const ServersScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Select server')),
      body: Consumer3<ServersProvider, AuthProvider, VPNProvider>(
        builder: (_, servers, auth, vpn, __) {
          if (servers.loading && servers.servers.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }
          if (servers.error != null) {
            return Center(child: Text(servers.error!));
          }
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              ListTile(
                title: const Text('Auto (best server)'),
                subtitle: const Text('Recommended'),
                leading: const Icon(Icons.bolt),
                onTap: () {
                  vpn.setSelectedServer(null);
                  Navigator.pop(context);
                },
              ),
              const Divider(),
              ...servers.servers.map((s) => ListTile(
                    title: Text(s.name),
                    subtitle: Text('${s.region} • ${s.host}'),
                    trailing: s.pingMs != null ? Text('${s.pingMs} ms') : null,
                    onTap: () {
                      vpn.setSelectedServer(s.name);
                      Navigator.pop(context);
                    },
                  )),
            ],
          );
        },
      ),
    );
  }
}
